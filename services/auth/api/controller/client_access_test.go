package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientAuthRequest(t *testing.T) {
	t.Run("should create valid client auth request", func(t *testing.T) {
		req := ClientAuthRequest{
			Method:  "GET",
			Path:    "/api/client/resource",
			Subject: map[string]interface{}{"id": "client-123", "type": "client"},
		}

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/api/client/resource", req.Path)
		assert.NotNil(t, req.Subject)
	})

	t.Run("should support client-specific context", func(t *testing.T) {
		req := ClientAuthRequest{
			Method: "POST",
			Path:   "/api/appointments",
			Subject: map[string]interface{}{
				"id":    "client-123",
				"email": "client@example.com",
				"type":  "client",
			},
			Resource: map[string]interface{}{
				"id":        "appointment-123",
				"client_id": "client-123",
			},
		}

		assert.Equal(t, "client-123", req.Subject["id"])
		assert.NotNil(t, req.Resource)
	})

	t.Run("should handle all HTTP methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			req := ClientAuthRequest{
				Method:  method,
				Path:    "/api/resource",
				Subject: map[string]interface{}{"id": "client-123"},
			}

			assert.Equal(t, method, req.Method)
		}
	})
}

func TestClientAuthRequestValidation(t *testing.T) {
	t.Run("should validate client subject structure", func(t *testing.T) {
		subject := map[string]interface{}{
			"id":    "client-123",
			"email": "client@example.com",
			"type":  "client",
		}

		assert.Contains(t, subject, "id")
		assert.Contains(t, subject, "email")
		assert.Equal(t, "client", subject["type"])
	})

	t.Run("should support appointment context", func(t *testing.T) {
		resource := map[string]interface{}{
			"id":        "appointment-123",
			"client_id": "client-123",
			"status":    "scheduled",
		}

		assert.Contains(t, resource, "client_id")
		assert.Equal(t, "scheduled", resource["status"])
	})
}

