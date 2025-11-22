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
// ADMIN AUTHORIZATION
// =====================

// AuthorizeAdmin evaluates admin-specific policies for access control
//
//	@Summary		Authorize admin access
//	@Description	Evaluate if an admin can access a resource based on admin policies
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AdminAuthRequest	true	"Authorization request"
//	@Success		200		{object}	AuthorizationResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/admin/authorize [post]
func AuthorizeAdmin(c *fiber.Ctx) error {
	var req AdminAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Extract admin claims from JWT token
	adminClaims, err := handler.JWT(c).WhoAreYouAdmin()
	if err != nil {
		return lib.Error.Auth.InvalidToken.WithError(err)
	}
	if adminClaims == nil {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin token required"))
	}

	// Build subject from token claims
	subject := map[string]interface{}{
		"user_id": adminClaims.ID.String(),
		"email":   adminClaims.Email,
		"type":    adminClaims.Type,
	}
	// Add roles to subject - check if user has specific role
	if len(adminClaims.Roles) > 0 {
		// For authorization, we check the primary role
		subject["role"] = adminClaims.Roles[0]
		subject["roles"] = adminClaims.Roles
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

	// Get all admin policies for this endpoint
	var policies []model.AdminPolicy
	if err := tx.Where("end_point_id = ?", endpoint.ID).Find(&policies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if len(policies) == 0 {
		return c.JSON(AuthorizationResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("no admin policies defined for endpoint: %s %s", req.Method, req.Path),
		})
	}

	// Create access controller
	accessCtrl := handler.NewAccessController(tx)

	// Separate allow and deny policies
	var allowPolicies []model.AdminPolicy
	var denyPolicies []model.AdminPolicy

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
			subject,
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
			subject,
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

type AdminAuthRequest struct {
	Method string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path   string `json:"path" validate:"required"`
	// Subject is now extracted from JWT token, not from request body
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}
