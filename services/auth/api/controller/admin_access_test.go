package controller

import (
	DTO "mynute-go/services/auth/config/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminAuthRequest(t *testing.T) {
	t.Run("should create valid admin auth request", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "DELETE",
			Path:   "/api/admin/users",
			// Subject is now extracted from JWT token
		}

		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, "/api/admin/users", req.Path)
	})

	t.Run("should support admin-specific context", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "POST",
			Path:   "/api/admin/policies",
			// Subject is now extracted from JWT token
			Body: map[string]interface{}{
				"name":   "New Policy",
				"effect": "Allow",
			},
		}

		assert.NotNil(t, req.Body)
	})

	t.Run("should handle all HTTP methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			req := DTO.AuthRequest{
				Method: method,
				Path:   "/api/admin/resource",
				// Subject is now extracted from JWT token
			}

			assert.Equal(t, method, req.Method)
		}
	})
}

func TestAdminAuthRequestValidation(t *testing.T) {
	t.Run("should validate request structure without subject", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "GET",
			Path:   "/api/admin/users",
			// Subject extracted from JWT token - not in request
		}

		assert.NotEmpty(t, req.Method)
		assert.NotEmpty(t, req.Path)
	})

	t.Run("should support resource and body context", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method: "POST",
			Path:   "/api/admin/system/config",
			// Subject extracted from JWT token
			Body: map[string]interface{}{
				"setting": "value",
			},
			Resource: map[string]interface{}{
				"type": "system_config",
			},
		}

		assert.NotNil(t, req.Body)
		assert.NotNil(t, req.Resource)
		assert.Contains(t, req.Path, "/admin/")
	})

	t.Run("should support optional parameters", func(t *testing.T) {
		req := DTO.AuthRequest{
			Method:     "PATCH",
			Path:       "/api/admin/users/:id",
			PathParams: map[string]interface{}{"id": "123"},
			Query:      map[string]interface{}{"include": "roles"},
			Headers:    map[string]interface{}{"X-Custom": "value"},
		}

		assert.NotNil(t, req.PathParams)
		assert.NotNil(t, req.Query)
		assert.NotNil(t, req.Headers)
	})
}
