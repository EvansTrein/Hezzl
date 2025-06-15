package model

import (
	"errors"
)

type apiError struct {
	Error string `json:"error"`
}

func NewApiError(msg string) apiError {
	return apiError{msg}
}

type Custom404 struct {
	Message string   `json:"message"`
	Code    int      `json:"code"`
	Details struct{} `json:"details"`
}

var (
	ErrValidate   = errors.New("data validate error")
	ErrQueryParam = errors.New("invalid query parameters")
	ErrNotFound   = errors.New("not found")
)
