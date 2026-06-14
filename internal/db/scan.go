package db

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func scanOriginalPost(row pgx.Row) (*OriginalPost, error) {
	orPost := &OriginalPost{}
	data := make([]byte, 0)
	var linkOriginalPost, originalChannel, originalImageURL sql.NullString

	err := row.Scan(
		&orPost.ID,
		&data,
		&linkOriginalPost,
		&originalChannel,
		&originalImageURL,
		&orPost.Theme,
		&orPost.CreatedAt,
		&orPost.UpdatedAt,
		&orPost.Attempts,
	)
	if err != nil {
		log.Err(err).Msg("scan failed")
		return nil, err
	}

	if linkOriginalPost.Valid {
		orPost.LinkOriginalPost = originalImageURL.String
	}
	if originalChannel.Valid {
		orPost.OriginalChannel = originalChannel.String
	}
	if originalImageURL.Valid {
		orPost.OriginalImageURL = originalImageURL.String
	}

	orPost.Data, err = convertByteIntoOriginalPostData(data)
	if err != nil {
		log.Err(err).
			Str("original_post_id", orPost.ID).
			Str("linked_post", orPost.LinkOriginalPost).
			Msg("failed unmarshal data after scan")

		return nil, err
	}

	return orPost, nil
}

func scanPosts(rows pgx.Rows) ([]*OriginalPost, error) {
	var originalPosts []*OriginalPost
	for rows.Next() {
		orPost := &OriginalPost{}
		data := make([]byte, 0)
		var linkOriginalPost, originalChannel, originalImageURL sql.NullString

		err := rows.Scan(
			&orPost.ID,
			&data,
			&linkOriginalPost,
			&originalChannel,
			&originalImageURL,
			&orPost.Theme,
			&orPost.CreatedAt,
			&orPost.UpdatedAt,
			&orPost.Attempts,
		)
		if err != nil {
			log.Err(err).Msg("scan failed")
			return nil, err
		}

		if linkOriginalPost.Valid {
			orPost.LinkOriginalPost = linkOriginalPost.String
		}
		if originalChannel.Valid {
			orPost.OriginalChannel = originalChannel.String
		}
		if originalImageURL.Valid {
			orPost.OriginalImageURL = originalImageURL.String
		}

		orPost.Data, err = convertByteIntoOriginalPostData(data)
		if err != nil {
			log.Err(err).
				Str("original_post_id", orPost.ID).
				Str("linked_post", orPost.LinkOriginalPost).
				Msg("failed unmarshal data after scan")

			continue
		}

		originalPosts = append(originalPosts, orPost)
	}
	if err := rows.Err(); err != nil {
		log.Err(err).Msg("rows error")
		return nil, err
	}

	return originalPosts, nil
}

func convertByteIntoOriginalPostData(data []byte) (*OriginalPostData, error) {
	var text string
	err := json.Unmarshal(data, &text)
	if err != nil {
		return nil, err
	}

	return &OriginalPostData{
		Text: text,
	}, nil
}
