package utils

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(pool *pgxpool.Pool) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db := stdlib.OpenDBFromPool(pool)
	return Retry(context.Background(), 3, time.Second, func() error {
		return goose.Up(db, "./migrations")
	})
}
