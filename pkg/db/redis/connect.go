package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client   *redis.Client
	TTLKeys  time.Duration
	numberDb int
}

func New(host, port, password string, ttlKeys time.Duration, numberDb int) (*RedisDB, error) {
	log.Println("Redis: connection to Redis started")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       numberDb,
	})

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Redis: connect to Redis successfully")
	return &RedisDB{Client: client, TTLKeys: ttlKeys, numberDb: numberDb}, nil
}

func (r *RedisDB) Close() error {
	log.Println("Redis: stop started")

	if r.Client == nil {
		return errors.New("redis connection is already closed")
	}

	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	r.Client = nil

	log.Println("Redis: stop successful")
	return nil
}
