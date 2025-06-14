package models

import (
	"errors"
)

type apiError struct {
	Error string `json:"error"`
}

func NewApiError(msg string) apiError {
	return apiError{msg}
}

var (
	ErrValidate  = errors.New("data validate error")
	ErrQueryParam = errors.New("invalid query parameters")
)
