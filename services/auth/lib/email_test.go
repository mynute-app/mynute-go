package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareEmail(t *testing.T) {
	t.Run("should validate and decode valid email", func(t *testing.T) {
		email := "test@example.com"

		result, err := PrepareEmail(email)

		assert.NoError(t, err)
		assert.Equal(t, email, result)
	})

	t.Run("should decode URL-encoded email", func(t *testing.T) {
		encodedEmail := "test%2Buser@example.com"
		expectedEmail := "test+user@example.com"

		result, err := PrepareEmail(encodedEmail)

		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, result)
	})

	t.Run("should return error for empty email", func(t *testing.T) {
		email := ""

		result, err := PrepareEmail(email)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "email parameter is empty")
	})

	t.Run("should return error for invalid email format", func(t *testing.T) {
		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"test@",
			"test..email@example.com",
			"test@example",
		}

		for _, email := range invalidEmails {
			result, err := PrepareEmail(email)

			assert.Error(t, err, "Email '%s' should be invalid", email)
			assert.Empty(t, result)
		}
	})

	t.Run("should accept various valid email formats", func(t *testing.T) {
		validEmails := []string{
			"simple@example.com",
			"user.name@example.com",
			"test123@test-domain.com",
			"a@b.c",
		}

		for _, email := range validEmails {
			result, err := PrepareEmail(email)

			assert.NoError(t, err, "Email '%s' should be valid", email)
			assert.Equal(t, email, result)
		}
	})

	t.Run("should handle email with special characters that need decoding", func(t *testing.T) {
		encodedEmail := "user%40company@example.com"

		result, err := PrepareEmail(encodedEmail)

		// After decoding, user%40company becomes user@company, which is invalid
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
