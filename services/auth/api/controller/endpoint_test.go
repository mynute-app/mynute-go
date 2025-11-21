package controller

import (
	"mynute-go/services/auth/api/lib"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEndpointCreateRequest(t *testing.T) {
	t.Run("should create valid endpoint request", func(t *testing.T) {
		req := EndpointCreateRequest{
			ControllerName:   "UserController",
			Description:      "Get user by ID",
			Method:           "GET",
			Path:             "/api/users/:id",
			DenyUnauthorized: true,
			NeedsCompanyId:   false,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate HTTP methods", func(t *testing.T) {
		validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range validMethods {
			req := EndpointCreateRequest{
				ControllerName: "TestController",
				Method:         method,
				Path:           "/api/test",
			}

			err := lib.MyCustomStructValidator(req)
			assert.NoError(t, err)
		}
	})

	t.Run("should reject invalid HTTP methods", func(t *testing.T) {
		req := EndpointCreateRequest{
			ControllerName: "TestController",
			Method:         "INVALID",
			Path:           "/api/test",
		}

		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})

	t.Run("should validate controller name length", func(t *testing.T) {
		req := EndpointCreateRequest{
			ControllerName: "AB", // Too short
			Method:         "GET",
			Path:           "/api/test",
		}

		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)

		req.ControllerName = "ValidControllerName"
		err = lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should accept optional resource ID", func(t *testing.T) {
		resourceID := uuid.New().String()
		req := EndpointCreateRequest{
			ControllerName: "UserController",
			Method:         "GET",
			Path:           "/api/users",
			ResourceID:     &resourceID,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate resource ID format", func(t *testing.T) {
		invalidID := "not-a-uuid"
		req := EndpointCreateRequest{
			ControllerName: "UserController",
			Method:         "GET",
			Path:           "/api/users",
			ResourceID:     &invalidID,
		}

		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})

	t.Run("should handle boolean flags", func(t *testing.T) {
		req := EndpointCreateRequest{
			ControllerName:   "AuthController",
			Method:           "POST",
			Path:             "/api/auth/login",
			DenyUnauthorized: false,
			NeedsCompanyId:   false,
		}

		assert.False(t, req.DenyUnauthorized)
		assert.False(t, req.NeedsCompanyId)

		req.DenyUnauthorized = true
		req.NeedsCompanyId = true

		assert.True(t, req.DenyUnauthorized)
		assert.True(t, req.NeedsCompanyId)
	})
}

func TestEndpointUpdateRequest(t *testing.T) {
	t.Run("should support partial updates", func(t *testing.T) {
		description := "Updated description"

		req := EndpointUpdateRequest{
			Description: &description,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate optional method", func(t *testing.T) {
		validMethod := "POST"
		invalidMethod := "INVALID"

		req := EndpointUpdateRequest{
			Method: &validMethod,
		}
		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)

		req.Method = &invalidMethod
		err = lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})

	t.Run("should validate optional controller name length", func(t *testing.T) {
		shortName := "AB"
		validName := "ValidControllerName"

		req := EndpointUpdateRequest{
			ControllerName: &shortName,
		}
		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)

		req.ControllerName = &validName
		err = lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should update all fields", func(t *testing.T) {
		controllerName := "UpdatedController"
		description := "Updated description"
		method := "PUT"
		path := "/api/updated"
		denyUnauth := true
		needsCompany := false
		resourceID := uuid.New().String()

		req := EndpointUpdateRequest{
			ControllerName:   &controllerName,
			Description:      &description,
			Method:           &method,
			Path:             &path,
			DenyUnauthorized: &denyUnauth,
			NeedsCompanyId:   &needsCompany,
			ResourceID:       &resourceID,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should clear resource ID with empty string", func(t *testing.T) {
		emptyID := ""

		req := EndpointUpdateRequest{
			ResourceID: &emptyID,
		}

		assert.NotNil(t, req.ResourceID)
		assert.Equal(t, "", *req.ResourceID)
	})
}

func TestEndpointPathValidation(t *testing.T) {
	t.Run("should accept various path formats", func(t *testing.T) {
		validPaths := []string{
			"/api/users",
			"/api/users/:id",
			"/api/v1/users",
			"/api/companies/:companyId/users/:userId",
			"/auth/login",
		}

		for _, path := range validPaths {
			req := EndpointCreateRequest{
				ControllerName: "TestController",
				Method:         "GET",
				Path:           path,
			}

			err := lib.MyCustomStructValidator(req)
			assert.NoError(t, err, "Path %s should be valid", path)
		}
	})

	t.Run("should require path to be non-empty", func(t *testing.T) {
		req := EndpointCreateRequest{
			ControllerName: "TestController",
			Method:         "GET",
			Path:           "",
		}

		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})
}

