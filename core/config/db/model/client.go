package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Custom TimeRange struct for start and end times
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type ClientMeta struct {
	BaseModel
	Name             string   `gorm:"type:varchar(100);not null" json:"name"`
	Surname          string   `gorm:"type:varchar(100)" json:"surname"`
	Email            string   `gorm:"type:varchar(100);not null;uniqueIndex" json:"email" validate:"required,email"`
	Phone            string   `gorm:"type:varchar(20);not null;uniqueIndex" json:"phone" validate:"required,e164"`
	Tags             []string `gorm:"type:json" json:"tags"`
	Password         string   `gorm:"type:varchar(255);not null" json:"password" validate:"required,myPasswordValidation"`
	ChangePassword   bool     `gorm:"default:false;not null" json:"change_password"`
	VerificationCode string   `gorm:"type:varchar(100)" json:"verification_code"`
	Verified         bool     `gorm:"default:false;not null" json:"verified"`
}

// Updated Client model
type ClientFull struct {
	ClientMeta
	Appointments mJSON.ClientAppointments `gorm:"type:jsonb" json:"appointments"`
}

func (ClientFull) TableName() string { return "public.clients" }
func (ClientFull) SchemaType() string { return "public" }

func (c *ClientFull) BeforeCreate(tx *gorm.DB) (err error) {
	if err := lib.ValidatorV10.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			BadReq := lib.Error.General.BadRequest
			for _, fieldErr := range validationErrors {
				// You can customize the message
				BadReq.WithError(
					fmt.Errorf("field '%s' failed on the '%s' rule", fieldErr.Field(), fieldErr.Tag()),
				)
			}
			return BadReq
		} else {
			return lib.Error.General.InternalError.WithError(err)
		}
	}
	if err := c.HashPassword(); err != nil {
		return err
	}
	return nil
}

// Method to set hashed password:
func (c *ClientFull) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hash)
	return nil
}
