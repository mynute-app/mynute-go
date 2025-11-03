package handler

import (
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "MyP@ssw0rd123"

		hashedPassword, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)
		assert.Greater(t, len(hashedPassword), len(password))
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		password := "MyP@ssw0rd123"

		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "bcrypt should generate unique salts")
	})

	t.Run("should handle empty password", func(t *testing.T) {
		password := ""

		hashedPassword, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("should return true for matching password", func(t *testing.T) {
		password := "MyP@ssw0rd123"
		hashedPassword, _ := HashPassword(password)

		result := ComparePassword(hashedPassword, password)

		assert.True(t, result)
	})

	t.Run("should return false for non-matching password", func(t *testing.T) {
		password := "MyP@ssw0rd123"
		wrongPassword := "WrongP@ss456"
		hashedPassword, _ := HashPassword(password)

		result := ComparePassword(hashedPassword, wrongPassword)

		assert.False(t, result)
	})

	t.Run("should return false for invalid hash", func(t *testing.T) {
		password := "MyP@ssw0rd123"
		invalidHash := "not-a-valid-bcrypt-hash"

		result := ComparePassword(invalidHash, password)

		assert.False(t, result)
	})

	t.Run("should handle empty password comparison", func(t *testing.T) {
		password := ""
		hashedPassword, _ := HashPassword(password)

		result := ComparePassword(hashedPassword, password)

		assert.True(t, result)
	})

	t.Run("should be case sensitive", func(t *testing.T) {
		password := "MyP@ssw0rd123"
		hashedPassword, _ := HashPassword(password)

		result := ComparePassword(hashedPassword, "myp@ssw0rd123")

		assert.False(t, result)
	})
}

func TestNewCookieStore(t *testing.T) {
	t.Run("should create cookie store with correct options", func(t *testing.T) {
		opts := SessionsOptions{
			CookiesKey: "test-secret-key-12345",
			MaxAge:     3600,
			HttpOnly:   true,
			Secure:     false,
		}

		store := NewCookieStore(opts)

		assert.NotNil(t, store)
		assert.IsType(t, &sessions.CookieStore{}, store)
	})

	t.Run("should set correct path", func(t *testing.T) {
		opts := SessionsOptions{
			CookiesKey: "test-secret-key-12345",
			MaxAge:     3600,
			HttpOnly:   true,
			Secure:     false,
		}

		store := NewCookieStore(opts)

		assert.NotNil(t, store)
		assert.Equal(t, "/", store.Options.Path)
	})

	t.Run("should respect HttpOnly option", func(t *testing.T) {
		opts := SessionsOptions{
			CookiesKey: "test-secret-key-12345",
			MaxAge:     3600,
			HttpOnly:   true,
			Secure:     false,
		}

		store := NewCookieStore(opts)

		assert.True(t, store.Options.HttpOnly)
	})

	t.Run("should respect Secure option", func(t *testing.T) {
		opts := SessionsOptions{
			CookiesKey: "test-secret-key-12345",
			MaxAge:     3600,
			HttpOnly:   true,
			Secure:     true,
		}

		store := NewCookieStore(opts)

		assert.True(t, store.Options.Secure)
	})
}
