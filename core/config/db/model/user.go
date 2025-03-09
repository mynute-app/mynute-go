package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Define a custom TimeRange struct for start and end times
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CreateUser struct {
	Name     string `gorm:"not null" json:"name" example:"John"`
	Surname  string `json:"surname" example:"Doe"`
	Role     string `json:"role" example:"user"`
	Email    string `gorm:"not null;unique" json:"email" example:"john.doe@example.com"`
	Phone    string `gorm:"not null;unique" json:"phone" example:"+15555555555"`
	Password string `gorm:"not null" json:"password" example:"1VerySecurePassword!"`
}

// Fourth step: Choosing the employee.
type User struct {
	gorm.Model
	CreateUser
	Tags             []string      `gorm:"type:json" json:"tags"` // Tags for the user
	VerificationCode string        `json:"verification_code"`
	Verified         bool          `gorm:"not null" json:"verified"`
	AvailableSlots   []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments     []Appointment `gorm:"foreignKey:UserID"` // One-to-many relation
	EmployeeID       uint          `json:"employee_id"`
	CompanyID        uint          `json:"company_id"`
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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Phone == "" {
		return errors.New("phone is required")
	} else if u.Password == "" {
		return errors.New("password is required")
	} else if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
