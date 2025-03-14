package model

import (
	"gorm.io/gorm"
)

type GeneralResourceInfo struct { // size=88 (0x58)
	gorm.Model
	Permissions map[string][]int `json:"permissions" gorm:"type:jsonb"`
}
