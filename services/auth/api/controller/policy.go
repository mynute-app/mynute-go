package controller

import (
	"encoding/json"
	"fmt"
	"mynute-go/services/auth/api/lib"
	authModel "mynute-go/services/auth/config/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// POLICY MANAGEMENT
// =====================

// ListPolicies retrieves all policy rules
//
//	@Summary		List all policies
//	@Description	Retrieve all policy rules
//	@Tags			Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Produce		json
//	@Success		200	{array}		authModel.PolicyRule
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/policies [get]
func ListPolicies(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var policies []authModel.PolicyRule
	if err := tx.Preload("EndPoint").Find(&policies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(policies)
}

// GetPolicyById retrieves a policy rule by ID
//
//	@Summary		Get policy by ID
//	@Description	Retrieve a policy rule by its ID
//	@Tags			Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			id				path		string	true	"Policy ID"
//	@Produce		json
//	@Success		200	{object}	authModel.PolicyRule
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/policies/{id} [get]
func GetPolicyById(c *fiber.Ctx) error {
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

	var policy authModel.PolicyRule
	if err := tx.Preload("EndPoint").Where("id = ?", id).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("policy not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(policy)
}

// CreatePolicy creates a new policy rule
//
//	@Summary		Create policy
//	@Description	Create a new policy rule
//	@Tags			Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string					true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			policy			body		PolicyCreateRequest		true	"Policy data"
//	@Success		201				{object}	authModel.PolicyRule
//	@Failure		400				{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/policies [post]
func CreatePolicy(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req PolicyCreateRequest
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
	var conditionsNode authModel.ConditionNode
	if err := json.Unmarshal(req.Conditions, &conditionsNode); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid conditions JSON: %w", err))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Check if endpoint exists
	var endpoint authModel.EndPoint
	if err := tx.Where("id = ?", endpointID).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	policy := authModel.PolicyRule{
		BaseModel:   authModel.BaseModel{ID: uuid.New()},
		Name:        req.Name,
		Description: req.Description,
		Effect:      req.Effect,
		EndPointID:  endpointID,
		Conditions:  req.Conditions,
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

// UpdatePolicyById updates a policy rule by ID
//
//	@Summary		Update policy
//	@Description	Update a policy rule
//	@Tags			Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string					true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string					true	"Policy ID"
//	@Param			policy			body		PolicyUpdateRequest		true	"Policy data"
//	@Success		200				{object}	authModel.PolicyRule
//	@Failure		400				{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404				{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/policies/{id} [patch]
func UpdatePolicyById(c *fiber.Ctx) error {
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

	var req PolicyUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var policy authModel.PolicyRule
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
		var endpoint authModel.EndPoint
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
		var conditionsNode authModel.ConditionNode
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

// DeletePolicyById deletes a policy rule by ID
//
//	@Summary		Delete policy
//	@Description	Delete a policy rule by its ID
//	@Tags			Policies
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			id				path		string	true	"Policy ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Failure		404	{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/policies/{id} [delete]
func DeletePolicyById(c *fiber.Ctx) error {
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

	var policy authModel.PolicyRule
	if err := tx.Where("id = ?", id).First(&policy).Error; err != nil {
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

type PolicyCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description"`
	Effect      string          `json:"effect" validate:"required,oneof=Allow Deny"`
	EndPointID  string          `json:"end_point_id" validate:"required,uuid"`
	Conditions  json.RawMessage `json:"conditions" validate:"required"`
}

type PolicyUpdateRequest struct {
	Name        *string         `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty" validate:"omitempty,oneof=Allow Deny"`
	EndPointID  *string         `json:"end_point_id,omitempty" validate:"omitempty,uuid"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
}
