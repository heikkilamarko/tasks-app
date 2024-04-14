package shared

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(ctx context.Context, config *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.Shared.PostgresConnectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
