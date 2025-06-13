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
	ErrNotFound = errors.New("not found")
)
