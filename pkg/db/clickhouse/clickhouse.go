package clickhouse

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Clickhouse struct {
	DB *sql.DB
}

func New(addr, dbName, userName, password string) (*Clickhouse, error) {
	log.Println("database: connection to ClickHouse started")

	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: dbName,
			Username: userName,
			Password: password,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
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
	return &Clickhouse{DB: db}, nil
}

func (c *Clickhouse) Close() error {
	log.Println("database: ClickHouse stop started")

	if c.DB == nil {
		return errors.New("database connection is already closed")
	}

	c.DB.Close()
	c.DB = nil

	log.Println("database: ClickHouse stop successful")
	return nil
}
