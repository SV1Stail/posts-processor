package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	post_processor_pb "github.com/SV1Stail/tg-project-protos/gen/go/posts_processor"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	insertSystemPromptQuery = `
	INSERT INTO system_prompts (owner_name, title_name, message, created, updated) 
	VALUES ($1, $2, $3::jsonb, NOW(), NOW())
	`
	updateSystemPromptQuery = `
    UPDATE system_prompts 
    SET message = $3::jsonb, updated = NOW() 
    WHERE owner_name = $1 AND title_name = $2
	`
)

type DB struct {
	pool *pgxpool.Pool
}

type OriginalPost struct {
	ID               string            `db:"id"`
	Data             *OriginalPostData `db:"data"`
	LinkOriginalPost *string           `db:"link_original_post"`
	OriginalChannel  *string           `db:"original_channel"`
	OriginalImageURL *string           `db:"original_image_url"`
	Theme            string            `db:"theme"`
	CreatedAt        time.Time         `db:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at"`
	Attempts         int               `db:"attempts"`
}

// TODO: заполнить
type OriginalPostData struct {
	Text string
}

// TODO: add config
func MustNewDB(ctx context.Context) *DB {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/posts_processor?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	return &DB{pool: pool}
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) WrapWithTransAction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("tx begin failed")
		return err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Err(err).Ctx(ctx).Msg("roll back failed")
		}
	}()

	err = fn(tx)
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("fn failed")
		return err
	}

	return tx.Commit(ctx)
}

type GetOriginalPosts struct {
	ID              string
	Theme           string
	OriginalChannel string
}

func (db *DB) GetOriginalPosts(ctx context.Context, in *GetOriginalPosts) ([]*OriginalPost, error) {
	filter := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	counter := 0
	if in.ID != "" {
		counter++
		filter = append(filter, fmt.Sprintf("ID = $%d", counter))
		args = append(args, in.ID)
	}
	if in.Theme != "" {
		counter++
		filter = append(filter, fmt.Sprintf("theme = $%d", counter))
		args = append(args, in.Theme)
	}
	if in.OriginalChannel != "" {
		counter++
		filter = append(filter, fmt.Sprintf("original_channel = $%d", counter))
		args = append(args, in.OriginalChannel)
	}
	baseQuery := `
    SELECT 
        id, data, link_original_post,
        original_channel, original_image_url,
        theme, created_at, updated_at, attempts
    FROM original_posts
	`
	if len(filter) > 0 {
		baseQuery += " WHERE " + strings.Join(filter, " AND ")
	} else {
		baseQuery += " WHERE 1=0"
	}
	baseQuery += " ORDER BY created_at DESC LIMIT 10"

	rows, err := db.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("fn failed")

	}
	defer rows.Close()

	posts, err := scanPosts(rows)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

type UpsertSystemPrompt struct {
	OwnerName     string
	TitleName     string
	systemMessage string
}

func (db *DB) InsertSystemPrompt(ctx context.Context, in *UpsertSystemPrompt) error {
	systemMessage, err := json.Marshal(in.systemMessage)
	if err != nil {
		log.Err(err).Msg("failed marshal system chat message")

		return err
	}

	_, err = db.pool.Exec(ctx, insertSystemPromptQuery, in.OwnerName,
		in.TitleName, systemMessage)

	return nil
}

func (db *DB) UpdateSystemPrompt(ctx context.Context, in *UpsertSystemPrompt) error {
	systemMessage, err := json.Marshal(in.systemMessage)
	if err != nil {
		log.Err(err).Msg("failed marshal system chat message")

		return err
	}

	_, err = db.pool.Exec(ctx, updateSystemPromptQuery, systemMessage,
		in.OwnerName, in.TitleName)
	if err != nil {
		log.Err(err).Msg("failed update system prompt")

		return err
	}

	return nil
}

func (db *DB) CreateOriginalPost(ctx context.Context, in *post_processor_pb.CreateOriginalPostRequest) error {
	query := `
	INSERT INTO original_posts (
	"id", "data", "link_original_post",
    "original_channel", "theme",
    "created_at", "updated_at") 
	VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	generatedID := uuid.NewString()
	_, err := db.pool.Exec(ctx, query, generatedID, []byte(in.GetText()),
		in.GetLinkOriginalPost(), in.GetOriginalChannel(), in.GetTheme())
	if err != nil {
		log.Err(err).Msg("failed insert original post")

		return err
	}
	log.Info().Ctx(ctx).Msg("success insert original post")

	return err
}

func (db *DB) DeleteOriginalPostsByChannel(ctx context.Context, in *post_processor_pb.DeleteOriginalPostsByChannelRequest) error {
	query := `
	DELETE FROM original_posts 
	WHERE original_channel = $1
	`

	_, err := db.pool.Exec(ctx, query, in.GetOriginalChannel())
	if err != nil {
		log.Err(err).Msg("failed delete original post")

		return err
	}
	log.Info().Ctx(ctx).Msg("success delete original post")

	return err
}

func (db *DB) GetOriginalPostsFromDB(ctx context.Context, in *post_processor_pb.GetOriginalPostsFromDBRequest) (*post_processor_pb.GetOriginalPostsFromDBResponse, error) {
	query := `
	SELECT 
		id, data, link_original_post,
		original_channel, original_image_url,
		theme, created_at, updated_at 
	FROM original_posts 
	ORDER BY created_at DESC
	OFFSET $1
	LIMIT $2
	`

	rows, err := db.pool.Query(ctx, query, in.GetOffset(), in.GetLimit())
	if err != nil {
		log.Err(err).Msg("failed delete original post")

		return nil, err
	}

	originalPosts := make([]*post_processor_pb.OriginalPost, 0)
	for rows.Next() {
		orPost := post_processor_pb.OriginalPost{}
		data := make([]byte, 0)
		err := rows.Scan(
			&orPost.Id,
			&data,
			&orPost.LinkOriginalPost,
			&orPost.OriginalChannel,
			&orPost.OriginalImageUrl,
			&orPost.Theme,
			&orPost.CreatedAt,
			&orPost.UpdatedAt,
		)
		if err != nil {
			log.Err(err).Msg("scan failed")
			return nil, err
		}

		d, err := convertByteIntoOriginalPostData(data)
		if err != nil {
			log.Err(err).
				Str("original_post_id", orPost.GetId()).
				Str("linked_post", orPost.GetLinkOriginalPost()).
				Msg("failed unmarshal data after scan")

			continue
		}
		orPost.Data = d.Text

		originalPosts = append(originalPosts, &orPost)
	}

	log.Info().Ctx(ctx).Msg("success GetOriginalPostsFromDB original post")

	return &post_processor_pb.GetOriginalPostsFromDBResponse{
		Posts: originalPosts,
	}, err
}
