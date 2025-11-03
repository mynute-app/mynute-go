package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UserMeta contains only authentication-related metadata
// Business-specific fields (design, appointments, etc.) belong in core service
type UserMeta struct {
	Login LoginConfig `json:"login" gorm:"-"`
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
