package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool represents the database connection pool
var Pool *pgxpool.Pool

// Init initializes the database connection pool
func Init() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	ctx := context.Background()
	var err error
	Pool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		return err
	}

	// Test the connection
	err = Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return err
	}

	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("Database connection pool closed")
	}
}
