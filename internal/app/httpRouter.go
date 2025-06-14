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

	return engine
}
