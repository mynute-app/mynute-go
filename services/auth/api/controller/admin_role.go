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
// ADMIN ROLE MANAGEMENT
// =====================

// ListAdminRoles retrieves all admin roles with pagination
//
//	@Summary		List all admin roles
//	@Description	Retrieve all admin roles with pagination
//	@Tags			Admin Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			limit			query	int		false	"Number of items per page (default: 10, max: 100)"
//	@Param			offset			query	int		false	"Number of items to skip (default: 0)"
//	@Produce		json
//	@Success		200	{object}	PaginatedAdminRolesResponse
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/roles [get]
func ListAdminRoles(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var roles []model.AdminRole
	limit, offset, err := List(c, &model.AdminRole{}, &roles)
	if err != nil {
		return err
	}

	return c.JSON(PaginatedAdminRolesResponse{
		Data:   roles,
		Limit:  limit,
		Offset: offset,
	})
}

// GetAdminRoleById retrieves an admin role by ID
//
//	@Summary		Get admin role by ID
//	@Description	Retrieve an admin role by its ID
//	@Tags			Admin Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	model.AdminRole
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/roles/{id} [get]
func GetAdminRoleById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var role model.AdminRole
	if err := GetOneBy("id", c, &role); err != nil {
		return err
	}

	return c.JSON(role)
}

// CreateAdminRole creates a new admin role
//
//	@Summary		Create admin role
//	@Description	Create a new admin role
//	@Tags			Admin Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			role	body		AdminRoleCreateRequest	true	"Role data"
//	@Success		201		{object}	model.AdminRole
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/roles [post]
func CreateAdminRole(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req AdminRoleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	role := model.AdminRole{
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

// UpdateAdminRoleById updates an admin role by ID
//
//	@Summary		Update admin role
//	@Description	Update an admin role
//	@Tags			Admin Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Role ID"
//	@Param			role	body		AdminRoleUpdateRequest	true	"Role data"
//	@Success		200		{object}	model.AdminRole
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/roles/{id} [patch]
func UpdateAdminRoleById(c *fiber.Ctx) error {
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

	var req AdminRoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var role model.AdminRole
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

// DeleteAdminRoleById deletes an admin role by ID
//
//	@Summary		Delete admin role
//	@Description	Delete an admin role by its ID
//	@Tags			Admin Roles
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/roles/{id} [delete]
func DeleteAdminRoleById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	return DeleteOneById(c, &model.AdminRole{})
}

// =====================
// REQUEST TYPES
// =====================

type PaginatedAdminRolesResponse struct {
	Data   []model.AdminRole `json:"data"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

type AdminRoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type AdminRoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}
