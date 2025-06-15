package model

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

type ProductUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ProjectID   int    `json:"project_id" validate:"required"`
	ID          int    `json:"id" validate:"required"`
}

type ProductRemoveResponce struct {
	ID        int  `json:"id"`
	ProjectID int  `json:"project_id"`
	Removed   bool `json:"removed"`
}

type ProductListResponce struct {
	Meta struct {
		Total   int  `json:"total"`
		Removed int `json:"removed"`
		Limit   int  `json:"limit"`
		Offset  int  `json:"offset"`
	} `json:"meta"`
	Goods []Product `json:"goods"`
}

type ProductReprioritizyRequest struct {
	NewPriority int `json:"newPriority" validate:"required"`
	ProjectID   int `json:"project_id" validate:"required"`
	ID          int `json:"id" validate:"required"`
}

type ProductReprioritizyResponce struct {
	Priorities []struct {
		ID       int `json:"id"`
		Priority int `json:"priority"`
	} `json:"priorities"`
}
