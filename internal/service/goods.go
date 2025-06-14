package service

import (
	"context"
	models "hezzl/internal/model"
	"log/slog"
)

type IGoodsRepo interface {
	Create(ctx context.Context, data models.ProductCreateRequest) (*models.Product, error)
}

type Goods struct {
	log  *slog.Logger
	repo IGoodsRepo
}

type GoodsDeps struct {
	*slog.Logger
	IGoodsRepo
}

func NewGoods(deps *GoodsDeps) *Goods {
	return &Goods{
		log:  deps.Logger,
		repo: deps.IGoodsRepo,
	}
}

func (s *Goods) Create(ctx context.Context, data models.ProductCreateRequest) (*models.Product, error) {

	result, err := s.repo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}
