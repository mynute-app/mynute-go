package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
)

// ServiceController embeds BaseController in order to extend it with the functions below
type ServiceController struct {
	BaseController[models.Service, DTO.Service]
}

func NewServiceController(Req *handlers.Request, Mid middleware.IMiddleware) *ServiceController {
	return &ServiceController{
		BaseController: BaseController[models.Service, DTO.Service]{
			Request:     Req,
			Middleware:  Mid,
			Associations: []string{"Branch"},
		},
	}
}