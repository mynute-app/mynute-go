package controller

import (
	DTO "mynute-go/services/auth/config/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantAuthRequest(t *testing.T) {
	t.Run("should create valid tenant auth request", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "GET",
			Path:   "/api/resource",
			// Subject is now extracted from JWT token
		}

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/api/resource", req.Path)
	})

	t.Run("should support all HTTP methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			req := DTO.AuthRequest{
				Method: method,
				Path:   "/api/resource",
				// Subject is now extracted from JWT token
			}

			assert.Equal(t, method, req.Method)
		}
	})

	t.Run("should support optional context parameters", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "POST",
			Path:   "/api/resource",
			// Subject is now extracted from JWT token
			Resource:   map[string]interface{}{"id": "resource-123"},
			PathParams: map[string]interface{}{"id": "123"},
			Body:       map[string]interface{}{"name": "test"},
			Query:      map[string]interface{}{"filter": "active"},
			Headers:    map[string]interface{}{"X-Custom": "value"},
		}

		assert.NotNil(t, req.Resource)
		assert.NotNil(t, req.PathParams)
		assert.NotNil(t, req.Body)
		assert.NotNil(t, req.Query)
		assert.NotNil(t, req.Headers)
	})
}

func TestAuthorizationResponse(t *testing.T) {
	t.Run("should create allowed response", func(t *testing.T) {
		resp := AuthorizationResponse{
			Allowed:    true,
			Reason:     "Access granted",
			PolicyID:   "policy-123",
			PolicyName: "Allow Read",
			Effect:     "Allow",
		}

		assert.True(t, resp.Allowed)
		assert.Equal(t, "Access granted", resp.Reason)
		assert.Equal(t, "Allow", resp.Effect)
	})

	t.Run("should create denied response", func(t *testing.T) {
		resp := AuthorizationResponse{
			Allowed:    false,
			Reason:     "Access denied",
			PolicyID:   "policy-456",
			PolicyName: "Deny Write",
			Effect:     "Deny",
		}

		assert.False(t, resp.Allowed)
		assert.Equal(t, "Access denied", resp.Reason)
		assert.Equal(t, "Deny", resp.Effect)
	})

	t.Run("should include error information", func(t *testing.T) {
		resp := AuthorizationResponse{
			Allowed: false,
			Reason:  "Policy evaluation failed",
			Error:   "invalid policy structure",
		}

		assert.False(t, resp.Allowed)
		assert.NotEmpty(t, resp.Error)
	})
}

func TestTenantAuthRequestValidation(t *testing.T) {
	t.Run("should validate request structure without subject", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "GET",
			Path:   "/api/resource",
			// Subject extracted from JWT token - not in request
		}

		assert.NotEmpty(t, req.Method)
		assert.NotEmpty(t, req.Path)
	})

	t.Run("should validate resource context structure", func(t *testing.T) {
		resource := map[string]interface{}{
			"id":        "resource-123",
			"tenant_id": "tenant-456",
			"owner_id":  "user-123",
		}

		assert.Contains(t, resource, "id")
		assert.Contains(t, resource, "tenant_id")
	})

	t.Run("should support optional parameters", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method:     "PUT",
			Path:       "/api/resource/:id",
			PathParams: map[string]interface{}{"id": "123"},
			Body:       map[string]interface{}{"status": "active"},
			Query:      map[string]interface{}{"notify": "true"},
			Headers:    map[string]interface{}{"X-Tenant": "tenant-456"},
		}

		assert.NotNil(t, req.PathParams)
		assert.NotNil(t, req.Body)
		assert.NotNil(t, req.Query)
		assert.NotNil(t, req.Headers)
	})
}
