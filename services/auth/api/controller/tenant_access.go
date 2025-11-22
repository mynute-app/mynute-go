package controller

import (
	"fmt"
	"mynute-go/services/auth/api/handler"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// TENANT AUTHORIZATION
// =====================

// AuthorizeTenant evaluates tenant-specific policies for access control
//
//	@Summary		Authorize tenant access
//	@Description	Evaluate if a tenant can access a resource based on tenant policies
//	@Tags			Tenant
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			request	body		TenantAuthRequest	true	"Authorization request"
//	@Success		200		{object}	AuthorizationResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/tenant/authorize [post]
func AuthorizeTenant(c *fiber.Ctx) error {
	var req TenantAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Extract tenant claims from JWT token
	claims, err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return lib.Error.Auth.InvalidToken.WithError(err)
	}
	if claims == nil {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("tenant token required"))
	}

	// Build subject from token claims
	subject := map[string]interface{}{
		"user_id":    claims.ID.String(),
		"email":      claims.Email,
		"type":       claims.Type,
		"company_id": claims.CompanyID.String(),
	}
	// Add role from claims (use first role or default)
	if len(claims.Roles) > 0 {
		subject["role"] = claims.Roles[0]
		subject["roles"] = claims.Roles
	} else {
		subject["role"] = "employee"
		subject["roles"] = []string{"employee"}
	}

	// Get tenant ID from header
	tenantIDStr := c.Get("X-Company-ID")
	if tenantIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("X-Company-ID header is required"))
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid X-Company-ID format"))
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

	// Get all tenant policies for this endpoint and tenant
	var policies []model.TenantPolicy
	if err := tx.Where("end_point_id = ? AND tenant_id = ?", endpoint.ID, tenantID).Find(&policies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Also get general tenant policies that apply to all tenants
	var generalPolicies []model.TenantGeneralPolicy
	if err := tx.Where("end_point_id = ?", endpoint.ID).Find(&generalPolicies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if len(policies) == 0 && len(generalPolicies) == 0 {
		return c.JSON(AuthorizationResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("no tenant policies defined for endpoint: %s %s", req.Method, req.Path),
		})
	}

	// Create access controller
	accessCtrl := handler.NewAccessController(tx)

	// Combine policies into a unified list for processing
	// We'll check tenant-specific policies first, then general policies
	type policyItem struct {
		policy model.PolicyInterface
		effect string
	}

	var allPolicies []policyItem

	// Add tenant-specific policies
	for _, policy := range policies {
		p := policy // Create a copy to avoid pointer issues
		allPolicies = append(allPolicies, policyItem{policy: &p, effect: p.Effect})
	}

	// Add general policies
	for _, policy := range generalPolicies {
		p := policy // Create a copy to avoid pointer issues
		allPolicies = append(allPolicies, policyItem{policy: &p, effect: p.Effect})
	}

	// Separate allow and deny policies
	var allowPolicies []policyItem
	var denyPolicies []policyItem

	for _, item := range allPolicies {
		switch item.effect {
		case "Allow":
			allowPolicies = append(allowPolicies, item)
		case "Deny":
			denyPolicies = append(denyPolicies, item)
		}
	}

	// Check deny policies first (explicit deny takes precedence)
	for _, item := range denyPolicies {
		decision := accessCtrl.Validate(
			subject,
			req.Resource,
			req.PathParams,
			req.Body,
			req.Query,
			req.Headers,
			item.policy,
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
				PolicyID:   item.policy.GetID().String(),
				PolicyName: item.policy.GetName(),
				Effect:     "Deny",
			})
		}
	}

	// Check allow policies
	for _, item := range allowPolicies {
		decision := accessCtrl.Validate(
			subject,
			req.Resource,
			req.PathParams,
			req.Body,
			req.Query,
			req.Headers,
			item.policy,
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
				PolicyID:   item.policy.GetID().String(),
				PolicyName: item.policy.GetName(),
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

type TenantAuthRequest struct {
	Method string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path   string `json:"path" validate:"required"`
	// Subject is now extracted from JWT token, not from request body
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}

type AuthorizationResponse struct {
	Allowed    bool   `json:"allowed"`
	Reason     string `json:"reason"`
	PolicyID   string `json:"policy_id,omitempty"`
	PolicyName string `json:"policy_name,omitempty"`
	Effect     string `json:"effect,omitempty"`
	Error      string `json:"error,omitempty"`
}
