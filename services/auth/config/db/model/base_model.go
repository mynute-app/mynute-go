package model

import (
	"fmt"
	"mynute-go/services/auth/api/lib"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid();<-:create" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	if m.ID != uuid.Nil && m.ID.Variant() != uuid.RFC4122 {
		errMsg := fmt.Errorf("BeforeSave: Invalid UUID variant for ID %s in %T", m.ID.String(), m)
		return lib.Error.General.UpdatedError.WithError(errMsg)
	}
	return nil
}
