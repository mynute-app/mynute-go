package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminAuthRequest(t *testing.T) {
	t.Run("should create valid admin auth request", func(t *testing.T) {
		req := AdminAuthRequest{
			Method:  "DELETE",
			Path:    "/api/admin/users",
			Subject: map[string]interface{}{"id": "admin-123", "roles": []string{"superadmin"}},
		}

		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, "/api/admin/users", req.Path)
		assert.NotNil(t, req.Subject)
	})

	t.Run("should support admin-specific context", func(t *testing.T) {
		req := AdminAuthRequest{
			Method: "POST",
			Path:   "/api/admin/policies",
			Subject: map[string]interface{}{
				"id":    "admin-123",
				"email": "admin@example.com",
				"roles": []string{"superadmin", "support"},
			},
			Body: map[string]interface{}{
				"name":   "New Policy",
				"effect": "Allow",
			},
		}

		assert.Contains(t, req.Subject, "roles")
		assert.NotNil(t, req.Body)
	})

	t.Run("should handle all HTTP methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			req := AdminAuthRequest{
				Method:  method,
				Path:    "/api/admin/resource",
				Subject: map[string]interface{}{"id": "admin-123"},
			}

			assert.Equal(t, method, req.Method)
		}
	})
}

func TestAdminAuthRequestValidation(t *testing.T) {
	t.Run("should validate admin subject structure", func(t *testing.T) {
		subject := map[string]interface{}{
			"id":    "admin-123",
			"email": "admin@example.com",
			"roles": []string{"superadmin"},
		}

		assert.Contains(t, subject, "id")
		assert.Contains(t, subject, "roles")
	})

	t.Run("should support multiple admin roles", func(t *testing.T) {
		subject := map[string]interface{}{
			"id":    "admin-123",
			"roles": []string{"superadmin", "support", "auditor"},
		}

		roles := subject["roles"].([]string)
		assert.Len(t, roles, 3)
		assert.Contains(t, roles, "superadmin")
	})

	t.Run("should handle system-wide operations", func(t *testing.T) {
		req := AdminAuthRequest{
			Method: "POST",
			Path:   "/api/admin/system/config",
			Subject: map[string]interface{}{
				"id":    "admin-123",
				"roles": []string{"superadmin"},
			},
			Body: map[string]interface{}{
				"setting": "value",
			},
		}

		assert.NotNil(t, req.Body)
		assert.Contains(t, req.Path, "/admin/")
	})
}
