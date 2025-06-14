package controller

import (
	"context"
	"encoding/json"
	models "hezzl/internal/model"
	"hezzl/pkg/validate"
	"net/http"
)

type IGoodsService interface {
	Create(ctx context.Context, data models.ProductCreateRequest) (*models.Product, error)
	Update(ctx context.Context, data models.ProductUpdateRequest) (*models.Product, error)
	Remove(ctx context.Context, id, projectId int) (*models.ProductRemoveResponce, error)
	List(ctx context.Context, offset, limit int) (*models.ProductListResponce, error)
	Reprioritizy(ctx context.Context, data models.ProductReprioritizyRequest) (*models.ProductReprioritizyResponce, error)
}

type Goods struct {
	base    *BaseController
	service IGoodsService
}

type GoodsDeps struct {
	*BaseController
	IGoodsService
}

func NewGoods(deps *GoodsDeps) *Goods {
	return &Goods{
		base:    deps.BaseController,
		service: deps.IGoodsService,
	}
}

func (g *Goods) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		reqData := models.ProductCreateRequest{
			ProjectID: projectId,
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		if err := validate.IsValid(reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrValidate)
			return
		}

		resp, err := g.service.Create(r.Context(), reqData)
		if err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		g.base.SendJsonResp(w, 201, resp)
	}
}

func (g *Goods) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := g.base.GetIntQueryParam(r, "id")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		reqData := models.ProductUpdateRequest{
			ID:        id,
			ProjectID: projectId,
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		if err := validate.IsValid(reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrValidate)
			return
		}

		resp, err := g.service.Update(r.Context(), reqData)
		if err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		g.base.SendJsonResp(w, 200, resp)
	}
}

func (g *Goods) Remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := g.base.GetIntQueryParam(r, "id")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		resp, err := g.service.Remove(r.Context(), id, projectId)
		if err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		g.base.SendJsonResp(w, 200, resp)
	}
}

func (g *Goods) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset, err := g.base.GetIntQueryParam(r, "offset")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		limit, err := g.base.GetIntQueryParam(r, "limit")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		resp, err := g.service.List(r.Context(), offset, limit)
		if err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		g.base.SendJsonResp(w, 200, resp)
	}
}

func (g *Goods) Reprioritizy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := g.base.GetIntQueryParam(r, "id")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrQueryParam)
			return
		}

		reqData := models.ProductReprioritizyRequest{
			ID:        id,
			ProjectID: projectId,
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		if err := validate.IsValid(reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), models.ErrValidate)
			return
		}

		resp, err := g.service.Reprioritizy(r.Context(), reqData)
		if err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		g.base.SendJsonResp(w, 200, resp)
	}
}
