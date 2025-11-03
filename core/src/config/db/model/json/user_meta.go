package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// LoginConfig contains authentication-related metadata
// TODO: This should be moved to auth service's UserMeta
type LoginConfig struct {
	NewPassword               *string    `json:"new_password,omitempty"`
	NewPasswordExpiry         *time.Time `json:"new_password_expiry,omitempty"`
	NewPasswordRequestedAt    *time.Time `json:"new_password_requested_at,omitempty"`
	NewPasswordRequestsCount  *int       `json:"new_password_requests_count,omitempty"`
	ValidationCode            *string    `json:"validation_code,omitempty"`
	ValidationExpiry          *time.Time `json:"validation_expiry,omitempty"`
	ValidationRequestedAt     *time.Time `json:"validation_requested_at,omitempty"`
	ValidationRequestsCount   *int       `json:"validation_requests_count,omitempty"`
	VerificationCode          *string    `json:"verification_code,omitempty"`
	VerificationExpiry        *time.Time `json:"verification_expiry,omitempty"`
	VerificationRequestedAt   *time.Time `json:"verification_requested_at,omitempty"`
	VerificationRequestsCount *int       `json:"verification_requests_count,omitempty"`
}

type UserMeta struct {
	Design DesignConfig `json:"design" gorm:"-"`
	Login  LoginConfig  `json:"login" gorm:"-"` // TODO: Move to auth service
}

// Scan implements the sql.Scanner interface for UserMeta
func (m *UserMeta) Scan(value any) error {
	if value == nil {
		*m = UserMeta{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &m)
}

// Value implements the driver.Valuer interface for UserMeta
func (m UserMeta) Value() (driver.Value, error) {
	return json.Marshal(m)
}
