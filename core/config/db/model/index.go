package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;<-:create" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	} else if m.ID.Variant() != uuid.RFC4122 {
		// 2. ID is not nil. Check if it's structurally valid (conforms to RFC 4122).
		// This covers updates *and* creates where an ID might have been manually assigned.
		// uuid.Nil has Variant() == uuid.Invalid.
		// A correctly formatted UUID (v1-v5) has Variant() == uuid.RFC4122.
		// Return an error to prevent saving an invalid UUID.
		errMsg := fmt.Sprintf("BeforeSave: Invalid UUID variant for ID %s in %T", m.ID.String(), m)
		return errors.New(errMsg)
	}
	return nil
}
