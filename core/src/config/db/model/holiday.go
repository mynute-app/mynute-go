package model

import (
	"time"
)

type Holiday struct {
	BaseModel
	Name        string    `gorm:"not null;index" json:"name"`
	Date        time.Time `gorm:"not null" json:"date"`
	Type        string    `gorm:"not null" json:"type"`
	Description string    `gorm:"not null" json:"description"`
	Recurrent   bool      `gorm:"not null;index" json:"recurrent"`
	DayMonth    string    `gorm:"not null" json:"dayMonth"`
}

const HolidayTableName = "public.holidays"

func (Holiday) TableName() string { return HolidayTableName }
func (Holiday) SchemaType() string { return "public" }