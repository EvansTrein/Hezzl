package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	DB *pgxpool.Pool
}

func New(storagePath string) (*PostgresDB, error) {
	log.Println("database: connection to Postgres started")

	DB, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("database: connect to Postgres successfully")
	return &PostgresDB{DB: DB}, nil
}

func (s *PostgresDB) Close() error {
	log.Println("database: stop started")

	if s.DB == nil {
		return errors.New("database connection is already closed")
	}

	s.DB.Close()
	s.DB = nil

	log.Println("database: stop successful")
	return nil
}
