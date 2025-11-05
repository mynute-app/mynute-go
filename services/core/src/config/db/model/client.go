package model

import (
	"fmt"
	mJSON "mynute-go/services/core/src/config/db/model/json"
	"mynute-go/services/core/src/lib"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Custom TimeRange struct for start and end times
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Client represents a customer/client in the business service
// Authentication is handled by auth service - UserID links to auth.users.id
type Client struct {
	UserID  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"user_id"` // Primary key, references auth.users.id
	Name    string         `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"`
	Surname string         `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"surname"`
	Phone   string         `gorm:"type:varchar(20);uniqueIndex" validate:"required,e164" json:"phone"`
	Meta    mJSON.UserMeta `gorm:"type:jsonb" json:"meta"` // Business-specific metadata
}

const ClientTableName = "public.clients"

func (Client) TableName() string  { return ClientTableName }
func (Client) SchemaType() string { return "public" }

func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if c.UserID == uuid.Nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("user_id is required"))
	}
	if err := lib.MyCustomStructValidator(c); err != nil {
		return err
	}
	return nil
}

func (c *Client) BeforeUpdate(tx *gorm.DB) (err error) {
	return lib.MyCustomStructValidator(c)
}

func (c *Client) GetFullClient(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}
	if err := tx.First(c, "user_id = ?", c.UserID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (c *Client) AddAppointment(a *Appointment, tx *gorm.DB) error {
	appointment := ClientAppointment{
		AppointmentID: a.ID,
		ClientID:      c.UserID,
		CompanyID:     a.CompanyID,
		StartTime:     a.StartTime,
		TimeZone:      a.TimeZone,
	}
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to public schema: %w", err))
	}
	tx.Create(&appointment)
	if err := tx.Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	companySchema := fmt.Sprintf("company_%s", a.CompanyID.String())
	if err := lib.ChangeToCompanySchema(tx, companySchema); err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to company schema: %w", err))
	}
	return nil
}
