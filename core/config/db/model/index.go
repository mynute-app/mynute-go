package model

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type GeneralResourceInfo struct { // size=88 (0x58)
	gorm.Model
	Permissions map[string][]int `json:"permissions" gorm:"type:jsonb"`
}

var (
	validate           = validator.New()
	passwordRegex      = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{6,16}$`)
	ErrInvalidPassword = errors.New("password must be 6-16 chars, containing uppercase, lowercase, number, special character")
)
