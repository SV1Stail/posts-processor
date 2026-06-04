package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pool *pgxpool.Pool
}

// TODO: add config
func MustNewDB(ctx context.Context) *DB {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/posts_processor?sslmode=disable"
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
