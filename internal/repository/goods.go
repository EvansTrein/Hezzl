package repository

import (
	"context"
	models "hezzl/internal/model"
	"hezzl/pkg/db/postgres"
	"log/slog"
)

type goodsRepo struct {
	log *slog.Logger
	*postgres.PostgresDB
}

type GoodsRepoDeps struct {
	*slog.Logger
	*postgres.PostgresDB
}

func NewGoodsRepo(deps *GoodsRepoDeps) *goodsRepo {
	return &goodsRepo{
		log:        deps.Logger,
		PostgresDB: deps.PostgresDB,
	}
}

func (r *goodsRepo) Create(ctx context.Context, data models.ProductCreateRequest) (*models.Product, error) {
	product := &models.Product{}

	return product, nil
}
