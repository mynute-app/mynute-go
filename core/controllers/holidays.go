package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
)

type HolidaysController struct {
	BaseController[models.Holidays, DTO.Holidays]
}

func Holidays(Gorm *handlers.Gorm) *HolidaysController {

	return &HolidaysController{
		BaseController: BaseController[models.Holidays, DTO.Holidays]{
			Name:         namespace.HolidaysKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Holidays(Gorm),
			Associations: []string{},
		},
	}
}
