package model

import mJSON "mynute-go/services/auth/config/db/json"

type User struct {
	BaseModel
	Email    string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string         `gorm:"type:varchar(255);not null" json:"-"`
	Verified bool           `gorm:"default:false" json:"verified"`
	Meta     mJSON.UserMeta `gorm:"type:jsonb" json:"meta"`
}
