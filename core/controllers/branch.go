package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type branchController struct {
	service.Base[models.Branch, DTO.Branch]
}

// CreateBranch creates a branch
//
//	@Summary		Create branch
//	@Description	Create a branch
//	@Tags			Branch
//	@Accept			json
//	@Produce		json
//	@Param			branch	body		DTO.CreateBranch	true	"Branch"
//	@Success		200		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch [post]
func (cc *branchController) CreateBranch(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetBranchById retrieves a branch by ID
//
//	@Summary		Get branch by ID
//	@Description	Retrieve a branch by its ID
//	@Tags			Branch
//	@Param			id	path	string	true	"Branch ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [get]
func (cc *branchController) GetBranchById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// UpdateBranch updates a branch by ID
//
//	@Summary		Update branch
//	@Description	Update a branch
//	@Tags			Branch
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Branch ID"
//	@Param			branch	body		DTO.UpdateBranch	true	"Branch"
//	@Success		200		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [patch]
func (cc *branchController) UpdateBranchById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteBranchById deletes a branch by ID
//
//	@Summary		Delete branch by ID
//	@Description	Delete a branch by its ID
//	@Tags			Branch
//	@Param			id	path	string	true	"Branch ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [delete]
func (cc *branchController) DeleteBranchById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// CreateBranch creates a branch
func Branch(Gorm *handlers.Gorm) *branchController {
	return &branchController{
		Base: service.Base[models.Branch, DTO.Branch]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Branch(Gorm),
			Associations: []string{"Employees", "Services"},
		},
	}
}
