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
	Appointments *mJSON.ClientAppointments `gorm:"type:jsonb" json:"appointments"`
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

func (c *Client) AddAppointment(a *Appointment, s *Service, company *Company, b *Branch, e *Employee) {
	if c.Appointments == nil {
		c.Appointments = &mJSON.ClientAppointments{}
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
	}

	if a.Payment != nil {
		if a.Payment.Price != 0 && a.Payment.Currency != "" {
			ca.Price = &a.Payment.Price
			ca.Currency = &a.Payment.Currency
		}
	}

	c.Appointments.Add(ca)
}
