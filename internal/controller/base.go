package controller

import (
	"context"
	"encoding/json"
	"errors"
	models "hezzl/internal/model"
	"log/slog"
	"net/http"
)

// Universal structure for sending responses
type BaseControllerResponce struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

type BaseController struct {
	Log *slog.Logger
}

type BaseControllerDeps struct {
	*slog.Logger
}

func NewBaseController(deps *BaseControllerDeps) *BaseController {
	return &BaseController{Log: deps.Logger}
}

func (h *BaseController) SendJsonResp(w http.ResponseWriter, status int, data any) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		h.Log.Error("failed to marshal JSON", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(jsonResponse); err != nil {
		h.Log.Error("!!ATTENTION!! failed to write JSON response", "error", err)
	}
}

func (b *BaseController) SendJsonError(w http.ResponseWriter, mess string, err error) {
	switch {
	case errors.Is(err, models.ErrValidate):
		b.Log.Error(mess, "error", err)
		b.SendJsonResp(w, 400, &BaseControllerResponce{
			Status:  http.StatusBadRequest,
			Message: mess,
			Error:   err.Error(),
		})
	case errors.Is(err, models.ErrQueryParam):
		b.Log.Error(mess, "error", err)
		b.SendJsonResp(w, 400, &BaseControllerResponce{
			Status:  http.StatusBadRequest,
			Message: mess,
			Error:   err.Error(),
		})
	case errors.Is(err, context.DeadlineExceeded):
		b.Log.Error("request processing exceeded the allowed time limit", "error", err)
		b.SendJsonResp(w, 504, &BaseControllerResponce{
			Status:  http.StatusGatewayTimeout,
			Message: "request processing exceeded the allowed time limit",
			Error:   err.Error(),
		})
	default:
		b.Log.Error("internal server error", "error", err)
		b.SendJsonResp(w, 500, &BaseControllerResponce{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
			Error:   err.Error(),
		})
	}
}
