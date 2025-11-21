package controller

import (
	"encoding/json"
	"mynute-go/services/auth/api/lib"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminPolicyCreateRequest(t *testing.T) {
	t.Run("should create valid admin policy request", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		req := AdminPolicyCreateRequest{
			Name:        "Allow System Access",
			Description: "Allows admin to access system resources",
			Effect:      "Allow",
			EndPointID:  uuid.New().String(),
			Conditions:  conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate effect values", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		validReq := AdminPolicyCreateRequest{
			Name:       "Test Policy",
			Effect:     "Deny",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}
		err := lib.MyCustomStructValidator(validReq)
		assert.NoError(t, err)

		invalidReq := AdminPolicyCreateRequest{
			Name:       "Test Policy",
			Effect:     "Unknown",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}
		err = lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err)
	})

	t.Run("should support admin role-based conditions", func(t *testing.T) {
		conditions := json.RawMessage(`{
			"logic_type": "OR",
			"children": [
				{
					"leaf": {
						"attribute": "subject.roles",
						"operator": "Contains",
						"value": "superadmin"
					}
				},
				{
					"leaf": {
						"attribute": "subject.roles",
						"operator": "Contains",
						"value": "support"
					}
				}
			]
		}`)

		req := AdminPolicyCreateRequest{
			Name:       "Admin Role Check",
			Effect:     "Allow",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})
}

func TestAdminPolicyUpdateRequest(t *testing.T) {
	t.Run("should support partial updates", func(t *testing.T) {
		description := "Updated admin policy description"

		req := AdminPolicyUpdateRequest{
			Description: &description,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate optional fields", func(t *testing.T) {
		name := "New Admin Policy"
		effect := "Allow"
		endpointID := uuid.New().String()

		req := AdminPolicyUpdateRequest{
			Name:       &name,
			Effect:     &effect,
			EndPointID: &endpointID,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})
}

