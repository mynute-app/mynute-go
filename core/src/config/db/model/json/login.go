package mJSON

import "time"

type LoginConfig struct {
	NewPassword               *string    `json:"new_password,omitempty"`
	NewPasswordExpiry         *time.Time `json:"new_password_expiry,omitempty"`
	NewPasswordRequestedAt    *time.Time `json:"new_password_requested_at,omitempty"`
	NewPasswordRequestsCount  *int       `json:"new_password_requests_count,omitempty"` // This is always reset when a new password is set successfully
	ValidationCode            *string    `json:"validation_code,omitempty"`
	ValidationExpiry          *time.Time `json:"validation_expiry,omitempty"`
	ValidationRequestedAt     *time.Time `json:"validation_requested_at,omitempty"`
	ValidationRequestsCount   *int       `json:"validation_requests_count,omitempty"` // This is always reset when a code is used successfully
	VerificationCode          *string    `json:"verification_code,omitempty"`
	VerificationExpiry        *time.Time `json:"verification_expiry,omitempty"`
	VerificationRequestedAt   *time.Time `json:"verification_requested_at,omitempty"`
	VerificationRequestsCount *int       `json:"verification_requests_count,omitempty"` // This is always reset when a code is used successfully
}
