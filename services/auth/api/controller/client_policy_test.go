package controller

import (
	"encoding/json"
	"mynute-go/services/auth/api/lib"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClientPolicyCreateRequest(t *testing.T) {
	t.Run("should create valid client policy request", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		req := ClientPolicyCreateRequest{
			Name:        "Allow Client Access",
			Description: "Allows clients to access their data",
			Effect:      "Allow",
			EndPointID:  uuid.New().String(),
			Conditions:  conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate effect values", func(t *testing.T) {
		conditions := json.RawMessage(`{"logic_type":"AND","children":[]}`)

		validReq := ClientPolicyCreateRequest{
			Name:       "Test Policy",
			Effect:     "Allow",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}
		err := lib.MyCustomStructValidator(validReq)
		assert.NoError(t, err)

		invalidReq := ClientPolicyCreateRequest{
			Name:       "Test Policy",
			Effect:     "Maybe",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}
		err = lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err)
	})

	t.Run("should support client-specific conditions", func(t *testing.T) {
		conditions := json.RawMessage(`{
			"logic_type": "AND",
			"children": [
				{
					"leaf": {
						"attribute": "subject.id",
						"operator": "Equals",
						"resource_attribute": "resource.client_id"
					}
				}
			]
		}`)

		req := ClientPolicyCreateRequest{
			Name:       "Client Self Access",
			Effect:     "Allow",
			EndPointID: uuid.New().String(),
			Conditions: conditions,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})
}

func TestClientPolicyUpdateRequest(t *testing.T) {
	t.Run("should support partial updates", func(t *testing.T) {
		name := "Updated Client Policy"

		req := ClientPolicyUpdateRequest{
			Name: &name,
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)
	})

	t.Run("should validate updated endpoint ID", func(t *testing.T) {
		validID := uuid.New().String()
		invalidID := "not-a-uuid"

		req := ClientPolicyUpdateRequest{
			EndPointID: &validID,
		}
		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err)

		req.EndPointID = &invalidID
		err = lib.MyCustomStructValidator(req)
		assert.Error(t, err)
	})
}
