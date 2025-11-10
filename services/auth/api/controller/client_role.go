package controller

import (
	"fmt"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// CLIENT ROLE MANAGEMENT
// =====================

// ListClientRoles retrieves all client roles with pagination
//
//	@Summary		List all client roles
//	@Description	Retrieve all client roles with pagination
//	@Tags			Client Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			limit			query	int		false	"Number of items per page (default: 10, max: 100)"
//	@Param			offset			query	int		false	"Number of items to skip (default: 0)"
//	@Produce		json
//	@Success		200	{object}	PaginatedClientRolesResponse
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/client/roles [get]
func ListClientRoles(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var roles []model.ClientRole
	limit, offset, err := List(c, &model.ClientRole{}, &roles)
	if err != nil {
		return err
	}

	return c.JSON(PaginatedClientRolesResponse{
		Data:   roles,
		Limit:  limit,
		Offset: offset,
	})
}

// GetClientRoleById retrieves a client role by ID
//
//	@Summary		Get client role by ID
//	@Description	Retrieve a client role by its ID
//	@Tags			Client Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	model.ClientRole
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/client/roles/{id} [get]
func GetClientRoleById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var role model.ClientRole
	if err := GetOneBy("id", c, &role); err != nil {
		return err
	}

	return c.JSON(role)
}

// CreateClientRole creates a new client role
//
//	@Summary		Create client role
//	@Description	Create a new client role
//	@Tags			Client Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			role	body		ClientRoleCreateRequest	true	"Role data"
//	@Success		201		{object}	model.ClientRole
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/client/roles [post]
func CreateClientRole(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req ClientRoleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	role := model.ClientRole{
		BaseModel:   model.BaseModel{ID: uuid.New()},
		Name:        req.Name,
		Description: req.Description,
	}

	if err := tx.Create(&role).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("role with name '%s' already exists", req.Name))
		}
		return lib.Error.General.CreatedError.WithError(err)
	}

	return c.Status(201).JSON(role)
}

// UpdateClientRoleById updates a client role by ID
//
//	@Summary		Update client role
//	@Description	Update a client role
//	@Tags			Client Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Role ID"
//	@Param			role	body		ClientRoleUpdateRequest	true	"Role data"
//	@Success		200		{object}	model.ClientRole
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/client/roles/{id} [patch]
func UpdateClientRoleById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	roleID := c.Params("id")
	if roleID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("role ID is required"))
	}

	id, err := uuid.Parse(roleID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid role ID format"))
	}

	var req ClientRoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var role model.ClientRole
	if err := tx.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("role not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Update fields if provided
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := tx.Save(&role).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("role with name '%s' already exists", role.Name))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(role)
}

// DeleteClientRoleById deletes a client role by ID
//
//	@Summary		Delete client role
//	@Description	Delete a client role by its ID
//	@Tags			Client Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/client/roles/{id} [delete]
func DeleteClientRoleById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	return DeleteOneById(c, &model.ClientRole{})
}

// =====================
// REQUEST TYPES
// =====================

type PaginatedClientRolesResponse struct {
	Data   []model.ClientRole `json:"data"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

type ClientRoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type ClientRoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}
