package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Default for local development
		dbUrl = "postgres://user:password@localhost:5432/biohorror"
	}
	var err error
	Pool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	return nil
}
