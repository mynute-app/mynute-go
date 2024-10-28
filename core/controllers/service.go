package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
)

// serviceController embeds BaseController in order to extend it with the functions below
type serviceController struct {
	BaseController[models.Service, DTO.Service]
}

func Service(Gorm *handlers.Gorm) *serviceController {
	return &serviceController{
		BaseController: BaseController[models.Service, DTO.Service]{
			Name:         namespace.EmployeeKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Service(Gorm),
			Associations: []string{"ServiceType"},
		},
	}
}
