package controller

import (
	"encoding/json"
	"fmt"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// TENANT POLICY MANAGEMENT
// =====================

// ListTenantPolicies retrieves all tenant policy rules with pagination
//
//	@Summary		List all tenant policies
//	@Description	Retrieve all tenant policy rules with pagination
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			limit			query	int		false	"Number of items per page (default: 10, max: 100)"
//	@Param			offset			query	int		false	"Number of items to skip (default: 0)"
//	@Produce		json
//	@Success		200	{object}	PaginatedTenantPoliciesResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/tenant/policies [get]
func ListTenantPolicies(c *fiber.Ctx) error {
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

	var policies []model.TenantPolicy
	if err := tx.Model(&model.TenantPolicy{}).
		Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Find(&policies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(PaginatedTenantPoliciesResponse{
		Data:   policies,
		Limit:  limit,
		Offset: offset,
	})
}

// GetTenantPolicyById retrieves a tenant policy rule by ID
//
//	@Summary		Get tenant policy by ID
//	@Description	Retrieve a tenant policy rule by its ID
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Policy ID"
//	@Produce		json
//	@Success		200	{object}	model.TenantPolicy
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/tenant/policies/{id} [get]
func GetTenantPolicyById(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	policyID := c.Params("id")
	if policyID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy ID is required"))
	}

	id, err := uuid.Parse(policyID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid policy ID format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var policy model.TenantPolicy
	if err := tx.Preload("EndPoint").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("policy not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(policy)
}

// CreateTenantPolicy creates a new tenant policy rule
//
//	@Summary		Create tenant policy
//	@Description	Create a new tenant policy rule
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			policy	body		TenantPolicyCreateRequest	true	"Policy data"
//	@Success		201		{object}	model.TenantPolicy
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/tenant/policies [post]
func CreateTenantPolicy(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	var req TenantPolicyCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Validate endpoint ID
	endpointID, err := uuid.Parse(req.EndPointID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid endpoint ID"))
	}

	// Validate effect
	if req.Effect != "Allow" && req.Effect != "Deny" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("effect must be 'Allow' or 'Deny'"))
	}

	// Validate conditions JSON
	var conditionsNode model.ConditionNode
	if err := json.Unmarshal(req.Conditions, &conditionsNode); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid conditions JSON: %w", err))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Check if endpoint exists
	var endpoint model.EndPoint
	if err := tx.Where("id = ?", endpointID).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	policy := model.TenantPolicy{
		BaseModel:   model.BaseModel{ID: uuid.New()},
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Effect:      req.Effect,
		EndPointID:  endpointID,
		Conditions:  req.Conditions,
	}

	if err := tx.Create(&policy).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "idx_tenant_policy_name") || strings.Contains(err.Error(), "duplicate key") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy with name '%s' already exists for tenant %s", req.Name, tenantID))
		}
		return lib.Error.General.CreatedError.WithError(err)
	}

	// Load the endpoint relation
	if err := tx.Preload("EndPoint").Where("id = ?", policy.ID).First(&policy).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.Status(201).JSON(policy)
}

// UpdateTenantPolicyById updates a tenant policy rule by ID
//
//	@Summary		Update tenant policy
//	@Description	Update a tenant policy rule
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Policy ID"
//	@Param			policy	body		TenantPolicyUpdateRequest	true	"Policy data"
//	@Success		200		{object}	model.TenantPolicy
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/tenant/policies/{id} [patch]
func UpdateTenantPolicyById(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	policyID := c.Params("id")
	if policyID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy ID is required"))
	}

	id, err := uuid.Parse(policyID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid policy ID format"))
	}

	var req TenantPolicyUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var policy model.TenantPolicy
	if err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("policy not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Update fields if provided
	if req.Name != nil {
		policy.Name = *req.Name
	}
	if req.Description != nil {
		policy.Description = *req.Description
	}
	if req.Effect != nil {
		if *req.Effect != "Allow" && *req.Effect != "Deny" {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("effect must be 'Allow' or 'Deny'"))
		}
		policy.Effect = *req.Effect
	}
	if req.EndPointID != nil {
		endpointID, err := uuid.Parse(*req.EndPointID)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid endpoint ID"))
		}

		// Check if endpoint exists
		var endpoint model.EndPoint
		if err := tx.Where("id = ?", endpointID).First(&endpoint).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint not found"))
			}
			return lib.Error.General.InternalError.WithError(err)
		}

		policy.EndPointID = endpointID
	}
	if req.Conditions != nil {
		// Validate conditions JSON
		var conditionsNode model.ConditionNode
		if err := json.Unmarshal(req.Conditions, &conditionsNode); err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid conditions JSON: %w", err))
		}
		policy.Conditions = req.Conditions
	}

	if err := tx.Save(&policy).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "idx_tenant_policy_name") || strings.Contains(err.Error(), "duplicate key") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy with name '%s' already exists for tenant %s", policy.Name, tenantID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Load the endpoint relation
	if err := tx.Preload("EndPoint").Where("id = ?", policy.ID).First(&policy).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(policy)
}

// DeleteTenantPolicyById deletes a tenant policy rule by ID
//
//	@Summary		Delete tenant policy
//	@Description	Delete a tenant policy rule by its ID
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Policy ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/tenant/policies/{id} [delete]
func DeleteTenantPolicyById(c *fiber.Ctx) error {
	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
	}

	policyID := c.Params("id")
	if policyID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy ID is required"))
	}

	id, err := uuid.Parse(policyID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid policy ID format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var policy model.TenantPolicy
	if err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("policy not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Delete(&policy).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(fiber.Map{
		"message": "Policy deleted successfully",
	})
}

// =====================
// REQUEST TYPES
// =====================

type PaginatedTenantPoliciesResponse struct {
	Data   []model.TenantPolicy `json:"data"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

type TenantPolicyCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description"`
	Effect      string          `json:"effect" validate:"required,oneof=Allow Deny"`
	EndPointID  string          `json:"end_point_id" validate:"required,uuid"`
	Conditions  json.RawMessage `json:"conditions" validate:"required" swaggertype:"string"`
}

type TenantPolicyUpdateRequest struct {
	Name        *string         `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty" validate:"omitempty,oneof=Allow Deny"`
	EndPointID  *string         `json:"end_point_id,omitempty" validate:"omitempty,uuid"`
	Conditions  json.RawMessage `json:"conditions,omitempty" swaggertype:"string"`
}
