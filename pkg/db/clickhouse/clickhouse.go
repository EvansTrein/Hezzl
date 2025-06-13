package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickhouseDB struct {
	DB *sql.DB
}

func New(host, port, dbName, userName, password string) (*ClickhouseDB, error) {
	log.Println("database: connection to ClickHouse started")

	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: dbName,
			Username: userName,
			Password: password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	log.Println("database: connect to ClickHouse successfully")
	return &ClickhouseDB{DB: db}, nil
}

func (c *ClickhouseDB) Close() error {
	log.Println("database: ClickHouse stop started")

	if c.DB == nil {
		return errors.New("database connection is already closed")
	}

	if err := c.DB.Close(); err != nil {
		return fmt.Errorf("failed to close ClickHouse connection: %w", err)
	}

	c.DB = nil

	log.Println("database: ClickHouse stop successful")
	return nil
}
