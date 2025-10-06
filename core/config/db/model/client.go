package model

import (
	"fmt"
	mJSON "mynute-go/core/config/db/model/json"
	"mynute-go/core/lib"
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
	Name     string             `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"`
	Surname  string             `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"surname"`
	Email    string             `gorm:"type:varchar(100);uniqueIndex" validate:"required,email" json:"email"`
	Phone    string             `gorm:"type:varchar(20);uniqueIndex" validate:"required,e164" json:"phone"`
	Password string             `gorm:"type:varchar(255)" validate:"required,myPasswordValidation" json:"password"`
	Verified bool               `gorm:"default:false" json:"verified"`
	Design   mJSON.DesignConfig `gorm:"type:jsonb" json:"design"`
}

// Updated Client model
type Client struct {
	ClientMeta
}

func (Client) TableName() string  { return "public.clients" }
func (Client) SchemaType() string { return "public" }

func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if err := lib.MyCustomStructValidator(c); err != nil {
		return err
	}
	if err := c.HashPassword(); err != nil {
		return err
	}
	return nil
}

func (c *Client) BeforeUpdate(tx *gorm.DB) (err error) {
	if c.Password != "" {
		var dbClient Client
		tx.First(&dbClient, "id = ?", c.ID)
		if c.Password == dbClient.Password || c.MatchPassword(dbClient.Password) {
			return nil
		}
		if err := lib.ValidatorV10.Var(c.Password, "myPasswordValidation"); err != nil {
			if _, ok := err.(validator.ValidationErrors); ok {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("password invalid"))
			} else {
				return lib.Error.General.InternalError.WithError(err)
			}
		}
		return c.HashPassword()
	}
	return nil
}

func (c *Client) MatchPassword(hashedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(c.Password))
	return err == nil
}

// Method to set hashed password:
func (c *Client) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hash)
	return nil
}

func (c *Client) GetFullClient(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}
	cID := c.ID.String()
	if err := tx.First(c, "id = ?", cID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func (c *Client) AddAppointment(a *Appointment, tx *gorm.DB) error {
	appointment := ClientAppointment{
		AppointmentID: a.ID,
		ClientID:      c.ID,
		CompanyID:     a.CompanyID,
		StartTime:     a.StartTime,
		TimeZone:      a.TimeZone,
	}
	tx.Create(&appointment)
	if err := tx.Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}
