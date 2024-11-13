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
type User struct {
	gorm.Model
	Name         string        `gorm:"not null" json:"name"`
	Surname      string        `json:"surname"`
	Role         string        `json:"role"`
	Email        string        `gorm:"not null;unique" json:"email"`
	Phone        string        `gorm:"unique" json:"phone"`
	Password     string        `gorm:"not null" json:"password"`
	Appointments []Appointment `gorm:"foreignKey:EmployeeID"` // One-to-many relation
	EmployeeInfo
}

type EmployeeInfo struct {
	CompanyID      uint        `json:"company_id"`                       // Foreign key to Company
	Company        Company     `gorm:"constraint:OnDelete:CASCADE;"`     // Relation to Company
	Branches       []Branch    `gorm:"many2many:employee_branches;"`     // Many-to-many relation with Branch
	Services       []Service   `gorm:"many2many:employee_services;"`     // Many-to-many relation with Service
	AvailableSlots []TimeRange `gorm:"type:json" json:"available_slots"` // Store availability as JSON in the database
}

// Check if the employee is available for a given service at a specific time.
func (e *User) IsAvailable(service Service, requestedTime time.Time) bool {
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
