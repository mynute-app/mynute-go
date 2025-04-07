package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	} else if m.ID.Variant() != uuid.RFC4122 {
		// 2. ID is not nil. Check if it's structurally valid (conforms to RFC 4122).
		// This covers updates *and* creates where an ID might have been manually assigned.
		// uuid.Nil has Variant() == uuid.Invalid.
		// A correctly formatted UUID (v1-v5) has Variant() == uuid.RFC4122.
		errMsg := fmt.Sprintf("BeforeSave: Invalid UUID variant for ID %s in %T", m.ID.String(), m)
		// Return an error to prevent saving an invalid UUID.
		return errors.New(errMsg)
	}
	log.Println("BeforeCreate called for", m)
	return nil
}
