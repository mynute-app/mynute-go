package controller

import (
	"fmt"
	"mynute-go/services/auth/api/handler"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =====================
// CLIENT AUTHORIZATION
// =====================

// AuthorizeClient evaluates client-specific policies for access control
//
//	@Summary		Authorize client access
//	@Description	Evaluate if a client can access a resource based on client policies
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ClientAuthRequest	true	"Authorization request"
//	@Success		200		{object}	AuthorizationResponse
//	@Failure		400		{object}	mynute-go_auth_config_dto.ErrorResponse
//	@Router			/client/authorize [post]
func AuthorizeClient(c *fiber.Ctx) error {
	var req ClientAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Find the endpoint
	var endpoint model.EndPoint
	if err := tx.Where("method = ? AND path = ?", req.Method, req.Path).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(AuthorizationResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("endpoint not found: %s %s", req.Method, req.Path),
			})
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Get all client policies for this endpoint
	var policies []model.ClientPolicy
	if err := tx.Where("end_point_id = ?", endpoint.ID).Find(&policies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if len(policies) == 0 {
		return c.JSON(AuthorizationResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("no client policies defined for endpoint: %s %s", req.Method, req.Path),
		})
	}

	// Create access controller
	accessCtrl := handler.NewAccessController(tx)

	// Separate allow and deny policies
	var allowPolicies []model.ClientPolicy
	var denyPolicies []model.ClientPolicy

	for _, policy := range policies {
		switch policy.Effect {
		case "Allow":
			allowPolicies = append(allowPolicies, policy)
		case "Deny":
			denyPolicies = append(denyPolicies, policy)
		}
	}

	// Check deny policies first (explicit deny takes precedence)
	for _, policy := range denyPolicies {
		decision := accessCtrl.Validate(
			req.Subject,
			req.Resource,
			req.PathParams,
			req.Body,
			req.Query,
			req.Headers,
			&policy,
		)

		if decision.Error != nil {
			return c.Status(500).JSON(AuthorizationResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("policy evaluation error: %v", decision.Error),
				Error:   decision.Error.Error(),
			})
		}

		// If deny policy is triggered (not allowed), deny access
		if !decision.Allowed {
			return c.JSON(AuthorizationResponse{
				Allowed:    false,
				Reason:     decision.Reason,
				PolicyID:   policy.ID.String(),
				PolicyName: policy.Name,
				Effect:     "Deny",
			})
		}
	}

	// Check allow policies
	for _, policy := range allowPolicies {
		decision := accessCtrl.Validate(
			req.Subject,
			req.Resource,
			req.PathParams,
			req.Body,
			req.Query,
			req.Headers,
			&policy,
		)

		if decision.Error != nil {
			return c.Status(500).JSON(AuthorizationResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("policy evaluation error: %v", decision.Error),
				Error:   decision.Error.Error(),
			})
		}

		// If allow policy matches, grant access
		if decision.Allowed {
			return c.JSON(AuthorizationResponse{
				Allowed:    true,
				Reason:     "Access granted",
				PolicyID:   policy.ID.String(),
				PolicyName: policy.Name,
				Effect:     "Allow",
			})
		}
	}

	// No allow policy matched
	return c.JSON(AuthorizationResponse{
		Allowed: false,
		Reason:  "no allow policy matched the request",
	})
}

// =====================
// REQUEST/RESPONSE TYPES
// =====================

type ClientAuthRequest struct {
	Method     string                 `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path       string                 `json:"path" validate:"required"`
	Subject    map[string]interface{} `json:"subject" validate:"required"`
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}
