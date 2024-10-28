package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"
)

type BranchController struct {
	BaseController[models.Branch, DTO.Branch]
}

func NewBranchController(HTTP *handlers.HTTP, Mid *middleware.Registry) *BranchController {
	return &BranchController{
		BaseController: BaseController[models.Branch, DTO.Branch]{
			HTTP:         HTTP,
			Middleware:   Mid,
			Associations: []string{"Company", "Employees", "Services"},
		},
	}
}