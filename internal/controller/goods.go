package controller

import (
	"context"
	"encoding/json"
	models "hezzl/internal/model"
	"hezzl/pkg/validate"
	"net/http"
	"strconv"
)

type IGoodsService interface {
	Create(ctx context.Context, data models.ProductCreateRequest) (*models.Product, error)
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
		projectIdStr := r.URL.Query().Get("projectId")
		projectId, err := strconv.Atoi(projectIdStr)
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

		if err := validate.IsValid(&reqData); err != nil {
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
