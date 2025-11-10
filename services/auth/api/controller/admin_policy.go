package controller

import (
	"encoding/json"
	"fmt"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// ADMIN POLICY MANAGEMENT
// =====================

// ListAdminPolicies retrieves all admin policy rules with pagination
//
//	@Summary		List all admin policies
//	@Description	Retrieve all admin policy rules with pagination
//	@Tags			Admin Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			limit			query	int		false	"Number of items per page (default: 10, max: 100)"
//	@Param			offset			query	int		false	"Number of items to skip (default: 0)"
//	@Produce		json
//	@Success		200	{object}	PaginatedAdminPoliciesResponse
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/policies [get]
func ListAdminPolicies(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var policies []model.AdminPolicy
	limit, offset, err := List(c, &model.AdminPolicy{}, &policies)
	if err != nil {
		return err
	}

	return c.JSON(PaginatedAdminPoliciesResponse{
		Data:   policies,
		Limit:  limit,
		Offset: offset,
	})
}

// GetAdminPolicyById retrieves an admin policy rule by ID
//
//	@Summary		Get admin policy by ID
//	@Description	Retrieve an admin policy rule by its ID
//	@Tags			Admin Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Policy ID"
//	@Produce		json
//	@Success		200	{object}	model.AdminPolicy
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/policies/{id} [get]
func GetAdminPolicyById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
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

	var policy model.AdminPolicy
	if err := tx.Preload("EndPoint").Where("id = ?", id).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("policy not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(policy)
}

// CreateAdminPolicy creates a new admin policy rule
//
//	@Summary		Create admin policy
//	@Description	Create a new admin policy rule
//	@Tags			Admin Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			policy	body		AdminPolicyCreateRequest	true	"Policy data"
//	@Success		201		{object}	model.AdminPolicy
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/policies [post]
func CreateAdminPolicy(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req AdminPolicyCreateRequest
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

	policy := model.AdminPolicy{
		Policy: model.Policy{
			BaseModel:   model.BaseModel{ID: uuid.New()},
			Name:        req.Name,
			Description: req.Description,
			Effect:      req.Effect,
			EndPointID:  endpointID,
			Conditions:  req.Conditions,
		},
	}

	if err := tx.Create(&policy).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	// Load the endpoint relation
	if err := tx.Preload("EndPoint").Where("id = ?", policy.ID).First(&policy).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.Status(201).JSON(policy)
}

// UpdateAdminPolicyById updates an admin policy rule by ID
//
//	@Summary		Update admin policy
//	@Description	Update an admin policy rule
//	@Tags			Admin Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Policy ID"
//	@Param			policy	body		AdminPolicyUpdateRequest	true	"Policy data"
//	@Success		200		{object}	model.AdminPolicy
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/policies/{id} [patch]
func UpdateAdminPolicyById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	policyID := c.Params("id")
	if policyID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy ID is required"))
	}

	id, err := uuid.Parse(policyID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid policy ID format"))
	}

	var req AdminPolicyUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var policy model.AdminPolicy
	if err := tx.Where("id = ?", id).First(&policy).Error; err != nil {
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
		return lib.Error.General.InternalError.WithError(err)
	}

	// Load the endpoint relation
	if err := tx.Preload("EndPoint").Where("id = ?", policy.ID).First(&policy).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(policy)
}

// DeleteAdminPolicyById deletes an admin policy rule by ID
//
//	@Summary		Delete admin policy
//	@Description	Delete an admin policy rule by its ID
//	@Tags			Admin Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Policy ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/admin/policies/{id} [delete]
func DeleteAdminPolicyById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	return DeleteOneById(c, &model.AdminPolicy{})
}

// =====================
// REQUEST TYPES
// =====================

type PaginatedAdminPoliciesResponse struct {
	Data   []model.AdminPolicy `json:"data"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

type AdminPolicyCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description"`
	Effect      string          `json:"effect" validate:"required,oneof=Allow Deny"`
	EndPointID  string          `json:"end_point_id" validate:"required,uuid"`
	Conditions  json.RawMessage `json:"conditions" validate:"required"`
}

type AdminPolicyUpdateRequest struct {
	Name        *string         `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty" validate:"omitempty,oneof=Allow Deny"`
	EndPointID  *string         `json:"end_point_id,omitempty" validate:"omitempty,uuid"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
}
