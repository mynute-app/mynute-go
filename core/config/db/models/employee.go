package models

import (
	"time"

	"gorm.io/gorm"
)

// Define a custom TimeRange struct for start and end times
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Fourth step: Choosing the employee.
type Employee struct {
	gorm.Model
	Name           string      `json:"name"`
	Role           string      `json:"role"`
	Branches       []Branch    `gorm:"many2many:branch_employees;"`      // Many-to-many relation
	Services       []Service   `gorm:"many2many:employee_services;"`     // Many-to-many relation
	AvailableSlots []TimeRange `gorm:"type:json" json:"available_slots"` // Store availability as JSON in the database
}

// Check if the employee is available for a given service at a specific time.
func (e *Employee) IsAvailable(service Service, requestedTime time.Time) bool {
	serviceEnd := requestedTime.Add(time.Duration(service.Duration) * time.Minute)

	for _, slot := range e.AvailableSlots {
		// Check if the requested time fits within any available slot
		if requestedTime.After(slot.Start) || requestedTime.Equal(slot.Start) {
			if serviceEnd.Before(slot.End) || serviceEnd.Equal(slot.End) {
				return true
			}
		}
	}
	return false
}