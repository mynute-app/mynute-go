package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
)

// EmployeeController embeds BaseController in order to extend it with the functions below
type EmployeeController struct {
	BaseController[models.Employee, DTO.Employee]
}

func NewEmployeeController(Req *handlers.Request, Mid middleware.IMiddleware) *EmployeeController {
	return &EmployeeController{
		BaseController: BaseController[models.Employee, DTO.Employee]{
			Request:     Req,
			Middleware:  Mid,
			Associations: []string{"Branch"},
		},
	}
}