package controller

import (
	"fmt"
	"mynute-go/auth"
	authModel "mynute-go/auth/model"
	"mynute-go/core/src/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// AUTHORIZATION / ACCESS CONTROL
// =====================

// CheckAccess evaluates if a subject can perform an action based on policies
//
//	@Summary		Check access authorization
//	@Description	Evaluate if a subject can access a resource based on endpoint policies
//	@Tags			Authorization
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AccessCheckRequest	true	"Access check request"
//	@Success		200		{object}	AccessCheckResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/authorize/check [post]
func CheckAccess(c *fiber.Ctx) error {
	var req AccessCheckRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Find the endpoint
	var endpoint authModel.EndPoint
	if err := tx.Where("method = ? AND path = ?", req.Method, req.Path).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(AccessCheckResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("endpoint not found: %s %s", req.Method, req.Path),
			})
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Get all policies for this endpoint
	var policies []authModel.PolicyRule
	if err := tx.Where("end_point_id = ?", endpoint.ID).Find(&policies).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if len(policies) == 0 {
		return c.JSON(AccessCheckResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("no policies defined for endpoint: %s %s", req.Method, req.Path),
		})
	}

	// Create access controller
	accessCtrl := auth.NewAccessController(tx)

	// Evaluate each policy
	var allowPolicies []authModel.PolicyRule
	var denyPolicies []authModel.PolicyRule

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
			return c.Status(500).JSON(AccessCheckResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("policy evaluation error: %v", decision.Error),
				Error:   decision.Error.Error(),
			})
		}

		// If deny policy says "not allowed", that means the deny condition was met
		if !decision.Allowed {
			return c.JSON(AccessCheckResponse{
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
			return c.Status(500).JSON(AccessCheckResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("policy evaluation error: %v", decision.Error),
				Error:   decision.Error.Error(),
			})
		}

		// If allow policy says "allowed", grant access
		if decision.Allowed {
			return c.JSON(AccessCheckResponse{
				Allowed:    true,
				Reason:     "Access granted",
				PolicyID:   policy.ID.String(),
				PolicyName: policy.Name,
				Effect:     "Allow",
			})
		}
	}

	// No allow policy matched
	return c.JSON(AccessCheckResponse{
		Allowed: false,
		Reason:  "no allow policy matched the request",
	})
}

// EvaluatePolicy evaluates a single policy against a request
//
//	@Summary		Evaluate single policy
//	@Description	Test a specific policy against request context
//	@Tags			Authorization
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Policy ID"
//	@Param			request	body		PolicyEvaluateRequest	true	"Policy evaluation request"
//	@Success		200		{object}	AccessCheckResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/authorize/policy/{id}/evaluate [post]
func EvaluatePolicy(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	policyIDStr := c.Params("id")
	if policyIDStr == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("policy ID is required"))
	}

	policyID, err := uuid.Parse(policyIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid policy ID format"))
	}

	var req PolicyEvaluateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Find the policy
	var policy authModel.PolicyRule
	if err := tx.Preload("EndPoint").Where("id = ?", policyID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(AccessCheckResponse{
				Allowed: false,
				Reason:  "policy not found",
			})
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Create access controller
	accessCtrl := auth.NewAccessController(tx)

	// Evaluate the policy
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
		return c.Status(500).JSON(AccessCheckResponse{
			Allowed:    false,
			Reason:     fmt.Sprintf("policy evaluation error: %v", decision.Error),
			Error:      decision.Error.Error(),
			PolicyID:   policy.ID.String(),
			PolicyName: policy.Name,
			Effect:     policy.Effect,
		})
	}

	reason := decision.Reason
	if decision.Allowed {
		reason = "Policy conditions met"
	}

	return c.JSON(AccessCheckResponse{
		Allowed:    decision.Allowed,
		Reason:     reason,
		PolicyID:   policy.ID.String(),
		PolicyName: policy.Name,
		Effect:     policy.Effect,
	})
}

// =====================
// REQUEST/RESPONSE TYPES
// =====================

type AccessCheckRequest struct {
	Method     string                 `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path       string                 `json:"path" validate:"required"`
	Subject    map[string]interface{} `json:"subject" validate:"required"`
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}

type AccessCheckByIdRequest struct {
	Subject    map[string]interface{} `json:"subject" validate:"required"`
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}

type PolicyEvaluateRequest struct {
	Subject    map[string]interface{} `json:"subject" validate:"required"`
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}

type AccessCheckResponse struct {
	Allowed    bool   `json:"allowed"`
	Reason     string `json:"reason"`
	PolicyID   string `json:"policy_id,omitempty"`
	PolicyName string `json:"policy_name,omitempty"`
	Effect     string `json:"effect,omitempty"`
	Error      string `json:"error,omitempty"`
}
