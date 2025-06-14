package app

import (
	"hezzl/internal/controller"
	"hezzl/pkg/logger"
	"hezzl/pkg/middleware"
	"net/http"
)

type activeHandlers struct {
	*controller.Goods
}

type activeHandlersDeps struct {
	*controller.Goods
}

func NewActiveHandlers(deps *activeHandlersDeps) *activeHandlers {
	return &activeHandlers{
		Goods: deps.Goods,
	}
}

func (h *activeHandlers) InitRouters() *http.ServeMux {
	engine := &http.ServeMux{}

	engine.Handle("POST /good/create", middleware.ChainMiddleware(
		middleware.HandlerLog(logger.GetLogger()),
	)(h.Goods.Create()))

	engine.Handle("PATCH /good/update", middleware.ChainMiddleware(
		middleware.HandlerLog(logger.GetLogger()),
	)(h.Goods.Update()))

	engine.Handle("DELETE /good/remove", middleware.ChainMiddleware(
		middleware.HandlerLog(logger.GetLogger()),
	)(h.Goods.Remove()))

	engine.Handle("GET /goods/list", middleware.ChainMiddleware(
		middleware.HandlerLog(logger.GetLogger()),
	)(h.Goods.List()))

	engine.Handle("PATCH /good/reprioritizy", middleware.ChainMiddleware(
		middleware.HandlerLog(logger.GetLogger()),
	)(h.Goods.Reprioritizy()))

	return engine
}
