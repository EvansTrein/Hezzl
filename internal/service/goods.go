package service

import (
	"context"
	models "hezzl/internal/model"
	"log/slog"
)

type IGoodsRepo interface {
	Create(ctx context.Context, data models.ProductCreateRequest) (*models.Product, error)
	Update(ctx context.Context, data models.ProductUpdateRequest) (*models.Product, error)
	Remove(ctx context.Context, id, projectId int) (*models.ProductRemoveResponce, error)
	List(ctx context.Context, offset, limit int) (*models.ProductListResponce, error)
	Reprioritizy(ctx context.Context, data models.ProductReprioritizyRequest) (*models.ProductReprioritizyResponce, error)
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

func (s *Goods) Update(ctx context.Context, data models.ProductUpdateRequest) (*models.Product, error) {

	return nil, nil
}

func (s *Goods) Remove(ctx context.Context, id, projectId int) (*models.ProductRemoveResponce, error) {

	return nil, nil
}

func (s *Goods) List(ctx context.Context, offset, limit int) (*models.ProductListResponce, error) {

	return nil, nil
}

func (s *Goods) Reprioritizy(ctx context.Context, data models.ProductReprioritizyRequest) (*models.ProductReprioritizyResponce, error) {

	return nil, nil
}
