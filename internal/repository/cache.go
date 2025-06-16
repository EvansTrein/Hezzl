package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"hezzl/internal/model"
	"hezzl/pkg/db/redis"
	"log/slog"
	"strings"
	"time"
)

const (
	cacheName   = "goodsList:"
	methodTimer = time.Second * 5
)

type cacheRepo struct {
	log *slog.Logger
	*redis.RedisDB
}

type CacheRepoDeps struct {
	*slog.Logger
	*redis.RedisDB
}

func NewCacheRepo(deps *CacheRepoDeps) *cacheRepo {
	return &cacheRepo{
		log:     deps.Logger,
		RedisDB: deps.RedisDB,
	}
}

func (r *cacheRepo) AddGoodsList(data *model.ProductListResponce) {
	op := "cache repository: creating"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func AddGoodsList", "data", data)

	ctx, cancel := context.WithTimeout(context.Background(), methodTimer)
	defer cancel()

	key := fmt.Sprintf("%soffset=%d:limit=%d", cacheName, data.Meta.Offset, data.Meta.Limit)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("failed to marshal json", "error", err)
		return
	}

	if res := r.Client.Set(ctx, key, jsonData, r.TTLKeys); res.Err() != nil {
		log.Error("failed to add a record", "error", err)
		return
	}

	log.Info("successfully added")
}

func (r *cacheRepo) GetGoodsList(offset, limit int) *model.ProductListResponce {
	op := "cache repository: retrieving"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func GetGoodsList", "offset", offset, "limit", limit)

	ctx, cancel := context.WithTimeout(context.Background(), methodTimer)
	defer cancel()

	key := fmt.Sprintf("%soffset=%d:limit=%d", cacheName, offset, limit)

	result, err := r.Client.Get(ctx, key).Result()
	switch {
	case err == nil:
		var goodsList model.ProductListResponce
		if unmarshalErr := json.Unmarshal([]byte(result), &goodsList); unmarshalErr != nil {
			log.Error("failed to unmarshal data", "error", unmarshalErr)
			return nil
		}
		log.Info("successfully retrieved")
		return &goodsList

	case strings.Contains(err.Error(), "redis: nil"):
		log.Warn("data not found")
		return nil

	default:
		log.Error("failed to get data from redis", "error", err)
		return nil
	}
}

func (r *cacheRepo) InvalidateGoods() {
	op := "cache repository: invalidating"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func InvalidateGoods")

	ctx, cancel := context.WithTimeout(context.Background(), methodTimer)
	defer cancel()

	pattern := cacheName + "*"

	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()
	var deleted int

	for iter.Next(ctx) {
		key := iter.Val()
		if res := r.Client.Del(ctx, key); res.Err() != nil {
			log.Error("failed to delete key", "key", key, "error", res.Err())
			continue
		}
		deleted++
	}

	if err := iter.Err(); err != nil {
		log.Error("error during key scanning", "error", err)
		return
	}

	log.Info("successfully invalidated cache", "deletedKeys", deleted)
}
