package DTO

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLoginClient_Structure(t *testing.T) {
	t.Run("should create LoginClient with all fields", func(t *testing.T) {
		login := LoginClientUser{
			Email:    "john.doe@example.com",
			Password: "1SecurePswd!",
		}

		assert.Equal(t, "john.doe@example.com", login.Email)
		assert.Equal(t, "1SecurePswd!", login.Password)
	})

	t.Run("should handle empty LoginClient", func(t *testing.T) {
		login := LoginClientUser{}

		assert.Empty(t, login.Email)
		assert.Empty(t, login.Password)
	})
}

func TestCreateClient_Structure(t *testing.T) {
	t.Run("should create CreateClient with all fields", func(t *testing.T) {
		client := CreateClientUser{
			Name:     "John",
			Surname:  "Doe",
			Email:    "john.doe@example.com",
			Phone:    "+15555555555",
			Password: "1SecurePswd!",
		}

		assert.Equal(t, "John", client.Name)
		assert.Equal(t, "Doe", client.Surname)
		assert.Equal(t, "john.doe@example.com", client.Email)
		assert.Equal(t, "+15555555555", client.Phone)
		assert.Equal(t, "1SecurePswd!", client.Password)
	})

	t.Run("should handle partial CreateClient", func(t *testing.T) {
		client := CreateClientUser{
			Email:    "test@example.com",
			Password: "Pass123!",
		}

		assert.Equal(t, "test@example.com", client.Email)
		assert.Equal(t, "Pass123!", client.Password)
		assert.Empty(t, client.Name)
		assert.Empty(t, client.Surname)
		assert.Empty(t, client.Phone)
	})
}

func TestClient_Structure(t *testing.T) {
	t.Run("should create Client with all fields", func(t *testing.T) {
		id := uuid.New()
		client := ClientUser{
			ID:       id,
			Name:     "John",
			Surname:  "Doe",
			Email:    "john.doe@example.com",
			Phone:    "+15555555555",
			Verified: true,
		}

		assert.Equal(t, id, client.ID)
		assert.Equal(t, "John", client.Name)
		assert.Equal(t, "Doe", client.Surname)
		assert.Equal(t, "john.doe@example.com", client.Email)
		assert.Equal(t, "+15555555555", client.Phone)
		assert.True(t, client.Verified)
	})

	t.Run("should default Verified to false", func(t *testing.T) {
		client := ClientUser{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		assert.False(t, client.Verified)
	})
}

func TestClientPopulated_Structure(t *testing.T) {
	t.Run("should create ClientPopulated with all fields", func(t *testing.T) {
		id := uuid.New()
		client := ClientUserPopulated{
			ID:       id,
			Name:     "John",
			Surname:  "Doe",
			Email:    "john.doe@example.com",
			Phone:    "+1-555-555-5555",
			Verified: true,
		}

		assert.Equal(t, id, client.ID)
		assert.Equal(t, "John", client.Name)
		assert.Equal(t, "Doe", client.Surname)
		assert.Equal(t, "john.doe@example.com", client.Email)
		assert.Equal(t, "+1-555-555-5555", client.Phone)
		assert.True(t, client.Verified)
	})

	t.Run("should support different phone formats", func(t *testing.T) {
		phoneFormats := []string{
			"+1-555-555-5555",
			"+15555555555",
			"+44 20 1234 5678",
		}

		for _, phone := range phoneFormats {
			client := ClientUserPopulated{
				Phone: phone,
			}

			assert.Equal(t, phone, client.Phone)
		}
	})
}
