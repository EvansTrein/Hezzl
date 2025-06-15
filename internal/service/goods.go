package service

import (
	"context"
	"hezzl/internal/model"
	"log/slog"
)

type IGoodsRepo interface {
	Create(ctx context.Context, data model.ProductCreateRequest) (*model.Product, error)
	Update(ctx context.Context, data model.ProductUpdateRequest) (*model.Product, error)
	Remove(ctx context.Context, id, projectId int) (*model.ProductRemoveResponce, error)
	List(ctx context.Context, offset, limit int) (*model.ProductListResponce, error)
	Reprioritizy(ctx context.Context, data model.ProductReprioritizyRequest) (*model.ProductReprioritizyResponce, error)
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

func (s *Goods) Create(ctx context.Context, data model.ProductCreateRequest) (*model.Product, error) {
	op := "goods service: creating"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Call func Create", "data", data)

	result, err := s.repo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	log.Info("successfully created")
	return result, nil
}

func (s *Goods) Update(ctx context.Context, data model.ProductUpdateRequest) (*model.Product, error) {
	op := "goods service: updating"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Call func Update", "data", data)

	result, err := s.repo.Update(ctx, data)
	if err != nil {
		return nil, err
	}

	log.Info("successfully updated")
	return result, nil
}

func (s *Goods) Remove(ctx context.Context, id, projectId int) (*model.ProductRemoveResponce, error) {
	op := "goods service: removing"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Call func Remove", "id", id, "projectId", projectId)

	result, err := s.repo.Remove(ctx, id, projectId)
	if err != nil {
		return nil, err
	}

	log.Info("successfully removed")
	return result, nil
}

func (s *Goods) List(ctx context.Context, offset, limit int) (*model.ProductListResponce, error) {
	op := "goods service: goods list retrieval"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Call func List", "offset", offset, "limit", limit)

	result, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	log.Info("successful search")
	return result, nil
}

func (s *Goods) Reprioritizy(ctx context.Context, data model.ProductReprioritizyRequest) (*model.ProductReprioritizyResponce, error) {
	op := "goods service: reprioritizing"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Call func Reprioritizy", "data", data)

	result, err := s.repo.Reprioritizy(ctx, data)
	if err != nil {
		return nil, err
	}

	log.Info("successfully reprioritized")
	return result, nil
}
