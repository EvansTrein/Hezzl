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

type ICacheRepo interface {
	AddGoodsList(data *model.ProductListResponce)
	GetGoodsList(offset, limit int) *model.ProductListResponce
	InvalidateGoods()
}

type IEventManager interface {
	SendToBroker(data *model.Product)
}

type Goods struct {
	log   *slog.Logger
	repo  IGoodsRepo
	cache ICacheRepo
	event IEventManager
}

type GoodsDeps struct {
	*slog.Logger
	IGoodsRepo
	ICacheRepo
	IEventManager
}

func NewGoods(deps *GoodsDeps) *Goods {
	return &Goods{
		log:   deps.Logger,
		repo:  deps.IGoodsRepo,
		cache: deps.ICacheRepo,
		event: deps.IEventManager,
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

	go s.cache.InvalidateGoods()
	go s.event.SendToBroker(result)

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

	go s.cache.InvalidateGoods()
	go s.event.SendToBroker(&model.Product{
		ID:        result.ID,
		ProjectID: result.ProjectID,
		Removed:   result.Removed,
	})

	log.Info("successfully removed")
	return result, nil
}

func (s *Goods) List(ctx context.Context, offset, limit int) (*model.ProductListResponce, error) {
	op := "goods service: goods list retrieval"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Call func List", "offset", offset, "limit", limit)

	if cacheResult := s.cache.GetGoodsList(offset, limit); cacheResult != nil {
		log.Debug("data was retrieved from the cache")
		return cacheResult, nil
	}

	result, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	go s.cache.AddGoodsList(result)

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

	go s.cache.InvalidateGoods()
	go func() {
		for _, el := range result.Priorities {
			s.event.SendToBroker(&model.Product{
				ID:        el.ID,
				ProjectID: el.ProjectID,
				Priority:  el.Priority,
			})
		}
	}()

	log.Info("successfully reprioritized")
	return result, nil
}
