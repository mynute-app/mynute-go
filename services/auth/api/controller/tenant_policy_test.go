package controller

import (
	"encoding/json"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTenantPolicyCreateRequest(t *testing.T) {
	t.Run("should create valid tenant policy request", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		req := TenantPolicyCreateRequest{
			Name:        "Allow Read",
			Description: "Allows read access",
			Effect:      "Allow",
			EndPointID:  uuid.New().String(),
			Conditions:  conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
		assert.Equal(t, "Allow", req.Effect)
	})

	t.Run("should validate effect values", func(t *testing.T) {
		validEffects := []string{"Allow", "Deny"}

		for _, effect := range validEffects {
			conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

			req := TenantPolicyCreateRequest{
				Name:        "Test Policy",
				Description: "Test",
				Effect:      effect,
				EndPointID:  uuid.New().String(),
				Conditions:  conditions,
			}

			err := lib.MyCustomStructValidator(req)
			assert.NoError(t, err)
		}
	})

	t.Run("should reject invalid effect values", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		req := TenantPolicyCreateRequest{
			Name:        "Test Policy",
			Description: "Test",
			Effect:      "Invalid",
			EndPointID:  uuid.New().String(),
			Conditions:  conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})

	t.Run("should validate required fields", func(t *testing.T) {
		req := TenantPolicyCreateRequest{
			Name:   "",
			Effect: "Allow",
		}

		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err, "Should reject empty name")
	})

	t.Run("should validate name length", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		// Too short
		req := TenantPolicyCreateRequest{
			Name:       "AB",
			Effect:     "Allow",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}
		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)

		// Valid length
		req.Name = "Valid Policy Name"
		err = lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should accept complex conditions", func(t *testing.T) {
		conditions := json.RawMessage(`{
			"logic_type": "AND",
			"children": [
				{
					"leaf": {
						"attribute": "subject.tenant_id",
						"operator": "Equals",
						"resource_attribute": "resource.tenant_id"
					}
				}
			]
		}`)

		req := TenantPolicyCreateRequest{
			Name:       "Complex Policy",
			Effect:     "Allow",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})
}

func TestTenantPolicyUpdateRequest(t *testing.T) {
	t.Run("should support partial updates", func(t *testing.T) {
		name := "Updated Name"

		req := TenantPolicyUpdateRequest{
			Name: &name,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate optional effect", func(t *testing.T) {
		validEffect := "Allow"
		invalidEffect := "Invalid"

		req := TenantPolicyUpdateRequest{
			Effect: &validEffect,
		}
		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)

		req.Effect = &invalidEffect
		err = lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})

	t.Run("should validate optional name length", func(t *testing.T) {
		shortName := "AB"
		validName := "Valid Policy Name"

		req := TenantPolicyUpdateRequest{
			Name: &shortName,
		}
		err := lib.MyCustomStructValidator(req)
		assert.Error(t, err)

		req.Name = &validName
		err = lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should update conditions", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"OR","children":[]}`)

		req := TenantPolicyUpdateRequest{
			Conditions: conditions,
		}

		assert.NotNil(t, req.Conditions)
	})
}

func TestPaginatedTenantPoliciesResponse(t *testing.T) {
	t.Run("should create paginated response", func(t *testing.T) {
		resp := PaginatedTenantPoliciesResponse{
			Data:   []model.TenantPolicy{},
			Limit:  10,
			Offset: 0,
		}

		assert.Equal(t, 10, resp.Limit)
		assert.Equal(t, 0, resp.Offset)
		assert.NotNil(t, resp.Data)
	})
}

