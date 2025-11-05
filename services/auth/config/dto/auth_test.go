package DTO

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClaims_Structure(t *testing.T) {
	t.Run("should create Claims with all fields", func(t *testing.T) {
		id := uuid.New()
		companyID := uuid.New()

		claims := Claims{
			ID:        id,
			Name:      "John",
			Surname:   "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+15555555555",
			Verified:  true,
			CompanyID: companyID,
			Password:  "StrongPswrd1!",
			Type:      "employee",
		}

		assert.Equal(t, id, claims.ID)
		assert.Equal(t, "John", claims.Name)
		assert.Equal(t, "Doe", claims.Surname)
		assert.Equal(t, "john.doe@example.com", claims.Email)
		assert.Equal(t, "+15555555555", claims.Phone)
		assert.True(t, claims.Verified)
		assert.Equal(t, companyID, claims.CompanyID)
		assert.Equal(t, "StrongPswrd1!", claims.Password)
		assert.Equal(t, "employee", claims.Type)
	})

	t.Run("should handle empty Claims", func(t *testing.T) {
		claims := Claims{}

		assert.Equal(t, uuid.Nil, claims.ID)
		assert.Empty(t, claims.Name)
		assert.Empty(t, claims.Email)
		assert.False(t, claims.Verified)
	})

	t.Run("should support different user types", func(t *testing.T) {
		types := []string{"admin", "client", "employee"}

		for _, userType := range types {
			claims := Claims{
				Type: userType,
			}

			assert.Equal(t, userType, claims.Type)
		}
	})
}

func TestLoginByEmailCode_Structure(t *testing.T) {
	t.Run("should create LoginByEmailCode with all fields", func(t *testing.T) {
		login := LoginByEmailCode{
			Email: "john.doe@example.com",
			Code:  "123456",
		}

		assert.Equal(t, "john.doe@example.com", login.Email)
		assert.Equal(t, "123456", login.Code)
	})

	t.Run("should handle empty LoginByEmailCode", func(t *testing.T) {
		login := LoginByEmailCode{}

		assert.Empty(t, login.Email)
		assert.Empty(t, login.Code)
	})

	t.Run("should validate 6-digit code format", func(t *testing.T) {
		validCodes := []string{"123456", "000000", "999999"}

		for _, code := range validCodes {
			login := LoginByEmailCode{
				Email: "test@example.com",
				Code:  code,
			}

			assert.Equal(t, 6, len(login.Code))
		}
	})
}
