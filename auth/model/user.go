package model

import (
	"fmt"
	"mynute-go/auth/lib"
	mJSON "mynute-go/auth/model/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a unified authentication user
// Type determines if user is admin, client, employee, etc.
type User struct {
	BaseModel
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	Type      string         `gorm:"type:varchar(50);not null;index" json:"type"` // "admin", "client", "employee"
	Meta      mJSON.UserMeta `gorm:"type:jsonb" json:"meta"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to set UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Validate type
	validTypes := map[string]bool{
		"admin":    true,
		"client":   true,
		"employee": true,
	}
	if !validTypes[u.Type] {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid user type: %s", u.Type))
	}

	return nil
}

// BeforeUpdate hook for validation
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Validate type if being updated
	if u.Type != "" {
		validTypes := map[string]bool{
			"admin":    true,
			"client":   true,
			"employee": true,
		}
		if !validTypes[u.Type] {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid user type: %s", u.Type))
		}
	}

	return nil
}
