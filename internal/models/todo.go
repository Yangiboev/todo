package models

import (
	"time"

	"github.com/google/uuid"
)

// ToDo Swagger model
type ToDoSwagger struct {
	Title string `json:"title" db:"title" validate:"required,gte=3"`
}

// ToDo model
type ToDo struct {
	ToDoID    uuid.UUID `json:"id" db:"id" validate:"omitempty,uuid"`
	Title     string    `json:"title" db:"title" validate:"required,gte=3"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// All ToDo response
type ToDosList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	ToDos      []*ToDo `json:"todos"`
}
