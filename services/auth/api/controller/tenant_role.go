package controller

import (
	"fmt"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"
	DTO "mynute-go/services/auth/config/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// TENANT ROLE MANAGEMENT
// =====================

// ListTenantRoles retrieves all tenant roles with pagination
//
//	@Summary		List all tenant roles
//	@Description	Retrieve all tenant roles with pagination
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			limit			query	int		false	"Number of items per page (default: 10, max: 100)"
//	@Param			offset			query	int		false	"Number of items to skip (default: 0)"
//	@Produce		json
//	@Success		200	{object}	PaginatedTenantRolesResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/tenant/roles [get]
func ListTenantRoles(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	// Parse pagination parameters
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var roles []model.TenantRole
	if err := tx.Model(&model.TenantRole{}).
		Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Find(&roles).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var roleBases []DTO.TenantRoleBase
	for _, role := range roles {
		roleBases = append(roleBases, DTO.TenantRoleBase{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return c.JSON(PaginatedTenantRolesResponse{
		Data:   roleBases,
		Limit:  limit,
		Offset: offset,
	})
}

// GetTenantRoleById retrieves a tenant role by ID
//
//	@Summary		Get tenant role by ID
//	@Description	Retrieve a tenant role by its ID
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	model.TenantRole
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/tenant/roles/{id} [get]
func GetTenantRoleById(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	roleID := c.Params("id")
	if roleID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("role ID is required"))
	}

	id, err := uuid.Parse(roleID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid role ID format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var role model.TenantRole
	if err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("role not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(role)
}

// CreateTenantRole creates a new tenant role
//
//	@Summary		Create tenant role
//	@Description	Create a new tenant role
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			role	body		TenantRoleCreateRequest	true	"Role data"
//	@Success		201		{object}	model.TenantRole
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/tenant/roles [post]
func CreateTenantRole(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	var req TenantRoleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	role := model.TenantRole{
		BaseModel:   model.BaseModel{ID: uuid.New()},
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := tx.Create(&role).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "idx_tenant_role_name") || strings.Contains(err.Error(), "duplicate key") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("role with name '%s' already exists for tenant %s", req.Name, tenantID))
		}
		return lib.Error.General.CreatedError.WithError(err)
	}

	return c.Status(201).JSON(role)
}

// UpdateTenantRoleById updates a tenant role by ID
//
//	@Summary		Update tenant role
//	@Description	Update a tenant role
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Role ID"
//	@Param			role	body		TenantRoleUpdateRequest	true	"Role data"
//	@Success		200		{object}	model.TenantRole
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/tenant/roles/{id} [patch]
func UpdateTenantRoleById(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	roleID := c.Params("id")
	if roleID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("role ID is required"))
	}

	id, err := uuid.Parse(roleID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid role ID format"))
	}

	var req TenantRoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var role model.TenantRole
	if err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&role).Error; err != nil {
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
		if strings.Contains(err.Error(), "idx_tenant_role_name") || strings.Contains(err.Error(), "duplicate key") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("role with name '%s' already exists for tenant %s", role.Name, tenantID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(role)
}

// DeleteTenantRoleById deletes a tenant role by ID
//
//	@Summary		Delete tenant role
//	@Description	Delete a tenant role by its ID
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/tenant/roles/{id} [delete]
func DeleteTenantRoleById(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	roleID := c.Params("id")
	if roleID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("role ID is required"))
	}

	id, err := uuid.Parse(roleID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid role ID format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var role model.TenantRole
	if err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("role not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Delete(&role).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}

// =====================
// REQUEST TYPES
// =====================

type PaginatedTenantRolesResponse struct {
	Data   []DTO.TenantRoleBase `json:"data"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

type TenantRoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type TenantRoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

