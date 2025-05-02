package model

import (
	"agenda-kaki-go/core/lib"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Custom TimeRange struct for start and end times
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Updated Client model
type Client struct {
	BaseModel
	Name             string        `gorm:"type:varchar(100);not null" json:"name"`
	Surname          string        `gorm:"type:varchar(100)" json:"surname"`
	Email            string        `gorm:"type:varchar(100);not null;uniqueIndex" json:"email" validate:"required,email"`
	Phone            string        `gorm:"type:varchar(20);not null;uniqueIndex" json:"phone" validate:"required,e164"`
	Tags             []string      `gorm:"type:json" json:"tags"`
	Password         string        `gorm:"type:varchar(255);not null" json:"password" validate:"required,myPasswordValidation"`
	ChangePassword   bool          `gorm:"default:false;not null" json:"change_password"`
	VerificationCode string        `gorm:"type:varchar(100)" json:"verification_code"`
	Verified         bool          `gorm:"default:false;not null" json:"verified"`
	AvailableSlots   []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments     []Appointment `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"appointments"`
}

func (Client) TableName() string {
	return "public.clients"
}

func (u *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if err := lib.ValidatorV10.Struct(u); err != nil {
		return err
	}
	if err := u.HashPassword(); err != nil {
		return err
	}
	return nil
}

// Method to set hashed password:
func (u *Client) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *Client) CheckAvailability(service Service, requestedTime time.Time) bool {
	serviceEnd := requestedTime.Add(time.Duration(service.Duration) * time.Minute)

	for _, slot := range u.AvailableSlots {
		if requestedTime.After(slot.Start) || requestedTime.Equal(slot.Start) {
			if serviceEnd.Before(slot.End) || serviceEnd.Equal(slot.End) {
				return true
			}
		}
	}
	return false
}
