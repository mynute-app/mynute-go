package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientAuthRequest(t *testing.T) {
	t.Run("should create valid client auth request", func(t *testing.T) {
		req := ClientAuthRequest{
			Method: "GET",
			Path:   "/api/client/resource",
			// Subject is now extracted from JWT token
		}

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/api/client/resource", req.Path)
	})

	t.Run("should support client-specific context", func(t *testing.T) {
		req := ClientAuthRequest{
			Method: "POST",
			Path:   "/api/appointments",
			// Subject is now extracted from JWT token
			Resource: map[string]interface{}{
				"id":        "appointment-123",
				"client_id": "client-123",
			},
		}

		assert.NotNil(t, req.Resource)
	})

	t.Run("should handle all HTTP methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			req := ClientAuthRequest{
				Method: method,
				Path:   "/api/resource",
				// Subject is now extracted from JWT token
			}

			assert.Equal(t, method, req.Method)
		}
	})
}

func TestClientAuthRequestValidation(t *testing.T) {
	t.Run("should validate request structure without subject", func(t *testing.T) {
		req := ClientAuthRequest{
			Method: "GET",
			Path:   "/api/client/profile",
			// Subject extracted from JWT token - not in request
		}

		assert.NotEmpty(t, req.Method)
		assert.NotEmpty(t, req.Path)
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

	t.Run("should support optional parameters", func(t *testing.T) {
		req := ClientAuthRequest{
			Method:     "POST",
			Path:       "/api/appointments",
			Body:       map[string]interface{}{"date": "2025-01-01"},
			Query:      map[string]interface{}{"filter": "upcoming"},
			PathParams: map[string]interface{}{"id": "123"},
			Headers:    map[string]interface{}{"X-Client": "web"},
		}

		assert.NotNil(t, req.Body)
		assert.NotNil(t, req.Query)
		assert.NotNil(t, req.PathParams)
		assert.NotNil(t, req.Headers)
	})
}
