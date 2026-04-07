package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	log.Println("Successfully connected to the database via pgxpool.")
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
