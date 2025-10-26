package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type UserMeta struct {
	Login  LoginConfig  `json:"login"`
	Design DesignConfig `json:"design"`
}

// Scan implements the sql.Scanner interface for UserMeta
func (m *UserMeta) Scan(value interface{}) error {
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
