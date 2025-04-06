package DTO

import "github.com/google/uuid"

type Sector struct {
	ID          uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name        string    `json:"name" example:"Your Company Sector Name"`
	Description string    `json:"description" example:"The Company Sector Description"`
}
