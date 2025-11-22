package model

import (
	mJSON "mynute-go/services/auth/config/db/json"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type TenantUser struct {
	BaseModel
	TenantID uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_tenant_email;not null" json:"tenant_id"`
	Email    string         `gorm:"type:varchar(255);uniqueIndex:idx_tenant_email;not null" json:"email"`
	Password string         `gorm:"type:varchar(255);not null" json:"-"`
	Verified bool           `gorm:"default:false" json:"verified"`
	Meta     mJSON.UserMeta `gorm:"type:jsonb" json:"meta"`
	Roles    []TenantRole   `gorm:"many2many:tenant_user_roles;" json:"roles"`
}

func (u *TenantUser) BeforeCreate(tx *gorm.DB) error {
	// Hash the password before creating the user
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}
