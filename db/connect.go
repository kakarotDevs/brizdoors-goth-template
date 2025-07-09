package db

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var DB *bun.DB

// InitDB initializes Bun using database/sql and pgx
func InitDB() error {
	dsn := os.Getenv("DEV_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
	}
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	DB = bun.NewDB(sqldb, pgdialect.New())
	return DB.PingContext(context.Background())
}

// CloseDB cleans up resources.
func CloseDB() error {
	if DB != nil {
		DB.Close()
	}
	return nil
}
