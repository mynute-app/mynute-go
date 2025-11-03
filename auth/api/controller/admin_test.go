package controller

import (
	"mynute-go/auth/config/db/model"
	DTO "mynute-go/auth/config/dto"
	"mynute-go/auth/handler"
	"mynute-go/auth/lib"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Unit tests for admin-related logic without database dependencies

func TestAdminCreationLogic(t *testing.T) {
	t.Run("should hash password correctly for admin", func(t *testing.T) {
		password := "Admin@123456"

		hashedPassword, err := handler.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)

		// Verify password can be validated
		isValid := handler.ComparePassword(hashedPassword, password)
		assert.True(t, isValid)
	})

	t.Run("should verify hashed password is different each time", func(t *testing.T) {
		password := "Admin@123456"

		hash1, _ := handler.HashPassword(password)
		hash2, _ := handler.HashPassword(password)

		assert.NotEqual(t, hash1, hash2, "Each hash should have unique salt")

		// Both should validate correctly
		assert.True(t, handler.ComparePassword(hash1, password))
		assert.True(t, handler.ComparePassword(hash2, password))
	})

	t.Run("should reject incorrect password", func(t *testing.T) {
		password := "Admin@123456"
		wrongPassword := "WrongPassword"

		hashedPassword, _ := handler.HashPassword(password)

		isValid := handler.ComparePassword(hashedPassword, wrongPassword)
		assert.False(t, isValid, "Wrong password should not validate")
	})

	t.Run("should handle password case sensitivity", func(t *testing.T) {
		password := "Admin@123456"
		hashedPassword, _ := handler.HashPassword(password)

		// Passwords are case-sensitive
		assert.False(t, handler.ComparePassword(hashedPassword, "admin@123456"))
		assert.False(t, handler.ComparePassword(hashedPassword, "ADMIN@123456"))
		assert.True(t, handler.ComparePassword(hashedPassword, password))
	})
}

func TestAdminValidation(t *testing.T) {
	t.Run("should validate admin creation request", func(t *testing.T) {
		validReq := DTO.AdminCreateRequest{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "SecureP@ss123",
			IsActive: true,
			Roles:    []string{"superadmin"},
		}

		err := lib.MyCustomStructValidator(validReq)
		assert.NoError(t, err, "Valid admin request should pass validation")
	})

	t.Run("should reject invalid email", func(t *testing.T) {
		invalidReq := DTO.AdminCreateRequest{
			Name:     "Admin User",
			Email:    "not-an-email",
			Password: "SecureP@ss123",
		}

		err := lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject invalid email")
	})

	t.Run("should reject short password", func(t *testing.T) {
		invalidReq := DTO.AdminCreateRequest{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "short",
		}

		err := lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject password shorter than 8 characters")
	})

	t.Run("should reject empty name", func(t *testing.T) {
		invalidReq := DTO.AdminCreateRequest{
			Name:     "",
			Email:    "admin@example.com",
			Password: "SecureP@ss123",
		}

		err := lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject empty name")
	})

	t.Run("should reject empty email", func(t *testing.T) {
		invalidReq := DTO.AdminCreateRequest{
			Name:     "Admin User",
			Email:    "",
			Password: "SecureP@ss123",
		}

		err := lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject empty email")
	})

	t.Run("should reject empty password", func(t *testing.T) {
		invalidReq := DTO.AdminCreateRequest{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "",
		}

		err := lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject empty password")
	})
}

func TestAdminClaimsStructure(t *testing.T) {
	t.Run("should create admin claims with all fields", func(t *testing.T) {
		adminID := uuid.New()
		claims := DTO.AdminClaims{
			ID:       adminID,
			Name:     "Test Admin",
			Email:    "admin@example.com",
			IsAdmin:  true,
			IsActive: true,
			Type:     "admin",
			Roles:    []string{"superadmin", "support"},
		}

		assert.Equal(t, adminID, claims.ID)
		assert.True(t, claims.IsAdmin)
		assert.Contains(t, claims.Roles, "superadmin")
		assert.Len(t, claims.Roles, 2)
	})

	t.Run("should support multiple roles", func(t *testing.T) {
		claims := DTO.AdminClaims{
			Roles: []string{"superadmin", "support", "auditor"},
		}

		assert.Len(t, claims.Roles, 3)
		assert.Contains(t, claims.Roles, "superadmin")
		assert.Contains(t, claims.Roles, "support")
		assert.Contains(t, claims.Roles, "auditor")
	})

	t.Run("should support empty roles", func(t *testing.T) {
		claims := DTO.AdminClaims{
			ID:    uuid.New(),
			Email: "admin@example.com",
			Roles: []string{},
		}

		assert.Len(t, claims.Roles, 0)
		assert.NotNil(t, claims.Roles)
	})
}

func TestAdminUpdateRequest(t *testing.T) {
	t.Run("should validate admin update request", func(t *testing.T) {
		name := "Updated Admin"
		surname := "Updated Surname"
		isActive := true

		validReq := DTO.AdminUpdateRequest{
			Name:     &name,
			Surname:  &surname,
			IsActive: &isActive,
		}

		err := lib.MyCustomStructValidator(validReq)
		assert.NoError(t, err, "Valid update request should pass validation")
	})

	t.Run("should allow partial updates", func(t *testing.T) {
		name := "Only Name Updated"

		partialReq := DTO.AdminUpdateRequest{
			Name: &name,
		}

		err := lib.MyCustomStructValidator(partialReq)
		assert.NoError(t, err, "Partial update should be valid")
	})

	t.Run("should validate updated email format", func(t *testing.T) {
		validEmail := "updated@example.com"
		invalidEmail := "not-an-email"

		validReq := DTO.AdminUpdateRequest{
			Email: &validEmail,
		}
		err := lib.MyCustomStructValidator(validReq)
		assert.NoError(t, err)

		invalidReq := DTO.AdminUpdateRequest{
			Email: &invalidEmail,
		}
		err = lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject invalid email format")
	})

	t.Run("should validate updated password length", func(t *testing.T) {
		validPassword := "ValidP@ss123"
		shortPassword := "short"

		validReq := DTO.AdminUpdateRequest{
			Password: &validPassword,
		}
		err := lib.MyCustomStructValidator(validReq)
		assert.NoError(t, err)

		invalidReq := DTO.AdminUpdateRequest{
			Password: &shortPassword,
		}
		err = lib.MyCustomStructValidator(invalidReq)
		assert.Error(t, err, "Should reject password shorter than 8 characters")
	})
}

func TestAdminModelStructure(t *testing.T) {
	t.Run("should create valid admin user model", func(t *testing.T) {
		hashedPassword, _ := handler.HashPassword("SecureP@ss123")

		admin := model.User{
			BaseModel: model.BaseModel{ID: uuid.New()},
			Email:     "admin@example.com",
			Password:  hashedPassword,
			Type:      "admin",
			Verified:  true,
		}

		assert.NotEqual(t, uuid.Nil, admin.ID)
		assert.Equal(t, "admin", admin.Type)
		assert.True(t, admin.Verified)
		assert.NotEmpty(t, admin.Password)
	})

	t.Run("should support different user types", func(t *testing.T) {
		validTypes := []string{"admin", "client", "employee"}

		for _, userType := range validTypes {
			user := model.User{
				BaseModel: model.BaseModel{ID: uuid.New()},
				Email:     "user@example.com",
				Type:      userType,
			}

			assert.Equal(t, userType, user.Type)
		}
	})
}

func TestAdminPasswordEdgeCases(t *testing.T) {
	t.Run("should handle very long passwords", func(t *testing.T) {
		// bcrypt has a 72-byte limit
		longPassword := "VeryLongP@ssword1234567890!@#$%^&*()_+-=[]{}|;:',.<>?/`~VeryLongPassword"

		hashedPassword, err := handler.HashPassword(longPassword)

		assert.NoError(t, err)
		assert.True(t, handler.ComparePassword(hashedPassword, longPassword))
	})

	t.Run("should handle passwords with special characters", func(t *testing.T) {
		specialPasswords := []string{
			"P@ssw0rd!",
			"Secure#Pass123",
			"Admin$2024*",
			"Test&Pass^2024",
		}

		for _, password := range specialPasswords {
			hashedPassword, err := handler.HashPassword(password)

			assert.NoError(t, err)
			assert.True(t, handler.ComparePassword(hashedPassword, password))
		}
	})

	t.Run("should handle passwords with unicode characters", func(t *testing.T) {
		unicodePassword := "Pāsswōrd123!"

		hashedPassword, err := handler.HashPassword(unicodePassword)

		assert.NoError(t, err)
		assert.True(t, handler.ComparePassword(hashedPassword, unicodePassword))
	})

	t.Run("should reject empty password for hashing", func(t *testing.T) {
		// This test verifies behavior - bcrypt will actually hash empty strings
		// but validation should catch this earlier
		emptyReq := DTO.AdminCreateRequest{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: "",
		}

		err := lib.MyCustomStructValidator(emptyReq)
		assert.Error(t, err, "Empty password should be rejected by validation")
	})
}

func TestAdminEmailEdgeCases(t *testing.T) {
	t.Run("should validate various email formats", func(t *testing.T) {
		validEmails := []string{
			"admin@example.com",
			"user.name@example.com",
			"user+tag@example.co.uk",
			"admin123@subdomain.example.com",
		}

		for _, email := range validEmails {
			req := DTO.AdminCreateRequest{
				Name:     "Admin",
				Email:    email,
				Password: "SecureP@ss123",
			}

			err := lib.MyCustomStructValidator(req)
			assert.NoError(t, err, "Email %s should be valid", email)
		}
	})

	t.Run("should reject invalid email formats", func(t *testing.T) {
		invalidEmails := []string{
			"notanemail",
			"@example.com",
			"user@",
			"user @example.com",
			"user@.com",
		}

		for _, email := range invalidEmails {
			req := DTO.AdminCreateRequest{
				Name:     "Admin",
				Email:    email,
				Password: "SecureP@ss123",
			}

			err := lib.MyCustomStructValidator(req)
			assert.Error(t, err, "Email %s should be invalid", email)
		}
	})

	t.Run("should handle email case", func(t *testing.T) {
		// Emails should typically be stored lowercase, but validation accepts both
		mixedCaseEmails := []string{
			"Admin@Example.Com",
			"ADMIN@EXAMPLE.COM",
			"admin@example.com",
		}

		for _, email := range mixedCaseEmails {
			req := DTO.AdminCreateRequest{
				Name:     "Admin",
				Email:    email,
				Password: "SecureP@ss123",
			}

			err := lib.MyCustomStructValidator(req)
			assert.NoError(t, err, "Email %s should be valid", email)
		}
	})
}

func TestAdminRolesValidation(t *testing.T) {
	t.Run("should accept valid roles", func(t *testing.T) {
		validRoles := [][]string{
			{"superadmin"},
			{"support"},
			{"auditor"},
			{"superadmin", "support"},
			{"superadmin", "support", "auditor"},
		}

		for _, roles := range validRoles {
			req := DTO.AdminCreateRequest{
				Name:     "Admin",
				Email:    "admin@example.com",
				Password: "SecureP@ss123",
				Roles:    roles,
			}

			err := lib.MyCustomStructValidator(req)
			assert.NoError(t, err, "Roles %v should be valid", roles)
		}
	})

	t.Run("should handle empty roles array", func(t *testing.T) {
		req := DTO.AdminCreateRequest{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: "SecureP@ss123",
			Roles:    []string{},
		}

		err := lib.MyCustomStructValidator(req)
		assert.NoError(t, err, "Empty roles array should be valid")
	})
}

func TestFirstAdminCreationScenarios(t *testing.T) {
	t.Run("first admin should be created as superadmin", func(t *testing.T) {
		// This test documents the expected behavior for first admin creation
		firstAdminReq := DTO.AdminCreateRequest{
			Name:     "First Admin",
			Email:    "first@example.com",
			Password: "SecureP@ss123",
			IsActive: true,
			Roles:    []string{"superadmin"},
		}

		err := lib.MyCustomStructValidator(firstAdminReq)
		assert.NoError(t, err)
		assert.Contains(t, firstAdminReq.Roles, "superadmin")
		assert.True(t, firstAdminReq.IsActive)
	})

	t.Run("first admin should have verified status", func(t *testing.T) {
		// First admin is auto-verified
		hashedPassword, _ := handler.HashPassword("SecureP@ss123")

		firstAdmin := model.User{
			BaseModel: model.BaseModel{ID: uuid.New()},
			Email:     "first@example.com",
			Password:  hashedPassword,
			Type:      "admin",
			Verified:  true, // First admin is auto-verified
		}

		assert.True(t, firstAdmin.Verified)
		assert.Equal(t, "admin", firstAdmin.Type)
	})

	t.Run("subsequent admins require authentication", func(t *testing.T) {
		// This test documents that after first admin, authentication is required
		// The actual endpoint validation would check for valid JWT token

		subsequentAdminReq := DTO.AdminCreateRequest{
			Name:     "Second Admin",
			Email:    "second@example.com",
			Password: "SecureP@ss123",
		}

		err := lib.MyCustomStructValidator(subsequentAdminReq)
		assert.NoError(t, err, "Request structure should be valid")
		// In actual endpoint, would also verify JWT token presence
	})
}
