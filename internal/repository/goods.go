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

func (r *goodsRepo) Update(ctx context.Context, data models.ProductUpdateRequest) (*models.Product, error) {

	return nil, nil
}

func (r *goodsRepo) Remove(ctx context.Context, id, projectId int) (*models.ProductRemoveResponce, error) {

	return nil, nil
}

func (r *goodsRepo)List(ctx context.Context, offset, limit int) (*models.ProductListResponce, error) {

	return nil, nil
}

func (r *goodsRepo) Reprioritizy(ctx context.Context, data models.ProductReprioritizyRequest) (*models.ProductReprioritizyResponce, error) {

	return nil, nil
}
