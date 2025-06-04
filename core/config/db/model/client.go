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
type Client struct {
	ClientMeta
	Appointments mJSON.ClientAppointments `gorm:"type:jsonb" json:"appointments"`
}

func (Client) TableName() string  { return "public.clients" }
func (Client) SchemaType() string { return "public" }

func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if err := lib.ValidatorV10.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			BadReq := lib.Error.General.BadRequest
			for _, fieldErr := range validationErrors {
				// You can customize the message
				BadReq = BadReq.WithError(
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

func (c *Client) AddAppointment(a *Appointment, s *Service, company *Company, b *Branch, e *Employee) {
	if c.Appointments == nil {
		c.Appointments = mJSON.ClientAppointments{}
	}

	ca := &mJSON.ClientAppointment{
		AppointmentID:    a.ID,
		ServiceName:      s.Name,
		ServicePrice:     s.Price,
		ServiceID:        s.ID,
		CompanyTradeName: company.TradeName,
		CompanyLegalName: company.LegalName,
		CompanyID:        a.CompanyID,
		BranchAddress:    a.Branch.GetAddress(),
		BranchID:         a.Branch.ID,
		EmployeeName:     a.Employee.Name,
		EmployeeID:       a.Employee.ID,
		IsCancelled:      a.IsCancelled,
		StartTime:        a.StartTime,
		Price:            &a.Payment.Price,
		Currency:         &a.Payment.Currency,
	}

	c.Appointments.Add(ca)
}
