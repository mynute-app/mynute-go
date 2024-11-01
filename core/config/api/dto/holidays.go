package DTO

import (
	"time"
)

type Holidays struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Recurrent   bool      `json:"recurrent"`
	DayMonth    string    `json:"dayMonth"`
}
