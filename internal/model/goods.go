package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductCreateRequest struct {
	Name      string `json:"name" validate:"required"`
	ProjectID int    `json:"project_id" validate:"required"`
}
