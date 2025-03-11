package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type branch_controller struct {
	service.Base[model.Branch, DTO.Branch]
}

// CreateBranch creates a branch
//
//	@Summary		Create branch
//	@Description	Create a branch
//	@Tags			Branch
//	@Accept			json
//	@Produce		json
//	@Param			branch	body		DTO.CreateBranch	true	"Branch"
//	@Success		201		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch [post]
func (cc *branch_controller) CreateBranch(c *fiber.Ctx) error {
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
func (cc *branch_controller) GetBranchById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// GetBranchByName retrieves a branch by name
//
//	@Summary		Get branch by name
//	@Description	Retrieve a branch by its name
//	@Tags			Branch
//	@Param			name	path	string	true	"Branch Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/name/{name} [get]
func (cc *branch_controller) GetBranchByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// UpdateBranch updates a branch by ID
//
//	@Summary		Update branch
//	@Description	Update a branch
//	@Tags			Branch
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Branch ID"
//	@Param			branch	body		DTO.UpdateBranch	true	"Branch"
//	@Success		200		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [patch]
func (cc *branch_controller) UpdateBranchById(c *fiber.Ctx) error {
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
func (cc *branch_controller) DeleteBranchById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// CreateBranch creates a branch
func Branch(Gorm *handler.Gorm) *branch_controller {
	return &branch_controller{
		Base: service.Base[model.Branch, DTO.Branch]{
			Name:         namespace.UserKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Employees", "Services"},
		},
	}
}
