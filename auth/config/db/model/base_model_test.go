package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBaseModel_UUID(t *testing.T) {
	t.Run("should create valid UUID", func(t *testing.T) {
		id := uuid.New()
		baseModel := BaseModel{
			ID: id,
		}

		assert.NotEqual(t, uuid.Nil, baseModel.ID)
		assert.Equal(t, uuid.RFC4122, baseModel.ID.Variant())
	})

	t.Run("should have zero UUID by default", func(t *testing.T) {
		baseModel := BaseModel{}

		assert.Equal(t, uuid.Nil, baseModel.ID)
	})
}

func TestBaseModel_BeforeSave(t *testing.T) {
	t.Run("should accept valid RFC4122 UUID", func(t *testing.T) {
		baseModel := BaseModel{
			ID: uuid.New(),
		}

		err := baseModel.BeforeSave(nil)

		assert.NoError(t, err)
	})

	t.Run("should accept nil UUID", func(t *testing.T) {
		baseModel := BaseModel{
			ID: uuid.Nil,
		}

		err := baseModel.BeforeSave(nil)

		assert.NoError(t, err)
	})

	t.Run("should validate UUID variant", func(t *testing.T) {
		validID := uuid.New()
		baseModel := BaseModel{
			ID: validID,
		}

		err := baseModel.BeforeSave(nil)

		assert.NoError(t, err)
		assert.Equal(t, uuid.RFC4122, baseModel.ID.Variant())
	})
}

func TestBaseModel_Fields(t *testing.T) {
	t.Run("should have all required fields", func(t *testing.T) {
		baseModel := BaseModel{
			ID: uuid.New(),
		}

		assert.NotNil(t, baseModel.ID)
		// CreatedAt and UpdatedAt are zero by default until saved to DB
		assert.True(t, baseModel.CreatedAt.IsZero())
		assert.True(t, baseModel.UpdatedAt.IsZero())
	})
}
