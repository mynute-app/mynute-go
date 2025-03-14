package model

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name             string        `gorm:"type:varchar(100);not null" json:"name" example:"John"`
	Surname          string        `gorm:"type:varchar(100)" json:"surname" example:"Doe"`
	Role             string        `gorm:"type:varchar(50);default:user;not null" json:"role" example:"user"`
	Email            string        `gorm:"type:varchar(100);not null;uniqueIndex" json:"email" example:"john.doe@example.com"`
	Phone            string        `gorm:"type:varchar(20);not null;uniqueIndex" json:"phone" example:"+15555555555"`
	Tags             []string      `gorm:"type:json" json:"tags"`
	Password         string        `gorm:"type:varchar(255);not null" json:"-"`
	VerificationCode string        `gorm:"type:varchar(100)" json:"-"`
	Verified         bool          `gorm:"default:false;not null" json:"verified"`
	AvailableSlots   []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments     []Appointment `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;"`
	CompanyID        uint          `gorm:"not null;index" json:"company_id"`
	Company          Company       `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;"`
	Branches         []*Branch     `gorm:"many2many:employee_branches;"`
	Services         []*Service    `gorm:"many2many:employee_services;"`
}

func (e *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	if err := e.validate_password(); err != nil {
		return err
	}
	if err := e.hash_password(); err != nil {
		return err
	}
	if err := e.validate_email(); err != nil {
		return err
	}
	return nil
}

func (e *Employee) validate_email() error {
	if e.Email == "" {
		return errors.New("email is required")
	}

	var validate = validator.New()

	if err := validate.Var(e.Email, "email"); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}


func (e *Employee) validate_password() error {
	// Password needs at least:
	// 1 uppercase letter
	// 1 lowercase letter
	// 1 number
	// 1 special character
	// min of 6 characters
	// max of 16 characters
	pswrdRegexp := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{6,16}$`)
	if !pswrdRegexp.MatchString(e.Password) {
		return errors.New("password must be 6-16 chars, containing uppercase, lowercase, number, special character")
	}
	return nil
}

// Method to set hashed password:
func (e *Employee) hash_password() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hash)
	return nil
}

// Method to verify password:
func (e *Employee) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password)) == nil
}

func (e *Employee) CheckAvailability(service Service, requestedTime time.Time) bool {
	serviceEnd := requestedTime.Add(time.Duration(service.Duration) * time.Minute)

	for _, slot := range e.AvailableSlots {
		if requestedTime.After(slot.Start) || requestedTime.Equal(slot.Start) {
			if serviceEnd.Before(slot.End) || serviceEnd.Equal(slot.End) {
				return true
			}
		}
	}
	return false
}

