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

func NewBranchController(Req *handlers.Request, Mid middleware.IMiddleware) *BranchController {
	return &BranchController{
		BaseController: BaseController[models.Branch, DTO.Branch]{
			Request:     Req,
			Middleware:  Mid,
			Associations: []string{"Company"},
		},
	}
}