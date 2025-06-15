package controller

import (
	"context"
	"encoding/json"
	"hezzl/internal/model"
	"hezzl/pkg/validate"
	"net/http"
)

type IGoodsService interface {
	Create(ctx context.Context, data model.ProductCreateRequest) (*model.Product, error)
	Update(ctx context.Context, data model.ProductUpdateRequest) (*model.Product, error)
	Remove(ctx context.Context, id, projectId int) (*model.ProductRemoveResponce, error)
	List(ctx context.Context, offset, limit int) (*model.ProductListResponce, error)
	Reprioritizy(ctx context.Context, data model.ProductReprioritizyRequest) (*model.ProductReprioritizyResponce, error)
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
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
			return
		}

		reqData := model.ProductCreateRequest{
			ProjectID: projectId,
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		if err := validate.IsValid(reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), model.ErrValidate)
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
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
			return
		}

		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
			return
		}

		reqData := model.ProductUpdateRequest{
			ID:        id,
			ProjectID: projectId,
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		if err := validate.IsValid(reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), model.ErrValidate)
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
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
			return
		}

		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
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
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		var offset, limit int

		switch offsetStr {
		case "":
			offset = 1
		default:
			var err error
			offset, err = g.base.GetIntQueryParam(r, "offset")
			if err != nil {
				g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
				return
			}
		}

		switch limitStr {
		case "":
			limit = 10
		default:
			var err error
			limit, err = g.base.GetIntQueryParam(r, "limit")
			if err != nil {
				g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
				return
			}
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
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
			return
		}

		projectId, err := g.base.GetIntQueryParam(r, "projectId")
		if err != nil {
			g.base.SendJsonError(w, err.Error(), model.ErrQueryParam)
			return
		}

		reqData := model.ProductReprioritizyRequest{
			ID:        id,
			ProjectID: projectId,
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), err)
			return
		}

		if err := validate.IsValid(reqData); err != nil {
			g.base.SendJsonError(w, err.Error(), model.ErrValidate)
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
