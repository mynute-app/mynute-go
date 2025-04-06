package DTO

import (
	"time"

	"github.com/google/uuid"
)

type Holidays struct {
	ID          uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name        string    `json:"name" example:"New Year's Day"`
	Date        time.Time `json:"date" example:"2025-01-01T00:00:00Z"`
	Type        string    `json:"type" example:"Public"`
	Description string    `json:"description" example:"Celebration of the first day of the new year"`
	Recurrent   bool      `json:"recurrent" example:"true"`
	DayMonth    string    `json:"dayMonth" example:"01-01"`
}
