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

// Updated User model
type User struct {
	gorm.Model
	Name             string        `gorm:"type:varchar(100);not null" json:"name"`
	Surname          string        `gorm:"type:varchar(100)" json:"surname"`
	Role             string        `gorm:"type:varchar(50);default:user;not null" json:"role"`
	Email            string        `gorm:"type:varchar(100);not null;uniqueIndex" json:"email" validate:"required,email"`
	Phone            string        `gorm:"type:varchar(20);not null;uniqueIndex" json:"phone" validate:"required,e164"`
	Tags             []string      `gorm:"type:json" json:"tags"`
	Password         string        `gorm:"type:varchar(255);not null" json:"password" validate:"required,myPasswordValidation"`
	ChangePassword   bool          `gorm:"default:false;not null" json:"change_password"`
	VerificationCode string        `gorm:"type:varchar(100)" json:"verification_code"`
	Verified         bool          `gorm:"default:false;not null" json:"verified"`
	AvailableSlots   []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments     []Appointment `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"appointments"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err := lib.ValidatorV10.Struct(u); err != nil {
		return err
	}
	if err := u.HashPassword(); err != nil {
		return err
	}
	return nil
}

// func (u *User) ValidatePhone() error {
// 	if u.Phone == "" {
// 		return errors.New("phone is required")
// 	}

// 	if err := validator_v10.Var(u.Phone, "e164"); err != nil {
// 		return errors.New("invalid phone format")
// 	}

// 	return nil
// }

// func (u *User) ValidateEmail() error {
// 	if u.Email == "" {
// 		return errors.New("email is required")
// 	}

// 	if err := validator_v10.Var(u.Email, "email"); err != nil {
// 		return errors.New("invalid email format")
// 	}

// 	return nil
// }

// func (u *User) ValidatePassword() error {
// 	// Password needs at least:
// 	// 1 uppercase letter
// 	// 1 lowercase letter
// 	// 1 number
// 	// 1 special character
// 	// min of 6 characters
// 	// max of 16 characters
// 	pswd := u.Password

// 	if len(pswd) < 6 || len(pswd) > 16 {
// 		return errors.New("password must be between 6 and 16 characters")
// 	}

// 	var (
// 		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
// 		hasLower   = regexp.MustCompile(`[a-z]`).MatchString
// 		hasDigit   = regexp.MustCompile(`\d`).MatchString
// 		hasSpecial = regexp.MustCompile(`[!@#$%^&*]`).MatchString
// 	)

// 	if !hasUpper(pswd) {
// 		return errors.New("password must contain at least one uppercase letter")
// 	} else if !hasLower(pswd) {
// 		return errors.New("password must contain at least one lowercase letter")
// 	} else if !hasDigit(pswd) {
// 		return errors.New("password must contain at least one number")
// 	} else if !hasSpecial(pswd) {
// 		return errors.New("password must contain at least one special character")
// 	}

// 	return nil
// }

// Method to set hashed password:
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckAvailability(service Service, requestedTime time.Time) bool {
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
