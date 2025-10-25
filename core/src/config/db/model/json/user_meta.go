package mJSON

import "time"

type UserMeta struct {
	LoginValidationCode   *string      `json:"login_validation_code,omitempty"`
	LoginValidationExpiry *time.Time   `json:"login_validation_expiry,omitempty"`
	Design                DesignConfig `json:"design"`
}
