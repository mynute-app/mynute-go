package model

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type GeneralResourceInfo struct { // size=88 (0x58)
	gorm.Model
	Permissions map[string][]int `json:"permissions" gorm:"type:jsonb"`
}

var validator_v10 = validator.New()