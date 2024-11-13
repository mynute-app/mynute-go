package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
)

type branchController struct {
	BaseController[models.Branch, DTO.Branch]
}

func Branch(Gorm *handlers.Gorm) *branchController {
	return &branchController{
		BaseController: BaseController[models.Branch, DTO.Branch]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Branch(Gorm),
			Associations: []string{"Company", "Employees", "Services"},
		},
	}
}
