package DTO

import (
	"mynute-go/services/auth/config/db/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminClaims_Structure(t *testing.T) {
	t.Run("should create AdminClaims with all fields", func(t *testing.T) {
		id := uuid.New()
		claims := AdminUserClaims{
			ID:       id,
			Name:     "Admin User",
			Email:    "admin@example.com",
			IsAdmin:  true,
			IsActive: true,
			Type:     "admin",
			Roles:    []string{"superadmin", "moderator"},
		}

		assert.Equal(t, id, claims.ID)
		assert.Equal(t, "Admin User", claims.Name)
		assert.Equal(t, "admin@example.com", claims.Email)
		assert.True(t, claims.IsAdmin)
		assert.True(t, claims.IsActive)
		assert.Equal(t, "admin", claims.Type)
		assert.Len(t, claims.Roles, 2)
		assert.Contains(t, claims.Roles, "superadmin")
	})

	t.Run("should handle admin claims without roles", func(t *testing.T) {
		claims := AdminUserClaims{
			ID:       uuid.New(),
			Email:    "admin@example.com",
			IsAdmin:  true,
			IsActive: true,
		}

		assert.True(t, claims.IsAdmin)
		assert.Empty(t, claims.Roles)
	})
}

func TestAdminLoginRequest_Structure(t *testing.T) {
	t.Run("should create AdminLoginRequest with all fields", func(t *testing.T) {
		login := AdminUserLoginRequest{
			Email:    "admin@example.com",
			Password: "StrongPassword123!",
		}

		assert.Equal(t, "admin@example.com", login.Email)
		assert.Equal(t, "StrongPassword123!", login.Password)
	})

	t.Run("should handle empty AdminLoginRequest", func(t *testing.T) {
		login := AdminUserLoginRequest{}

		assert.Empty(t, login.Email)
		assert.Empty(t, login.Password)
	})
}

func TestAdmin_Structure(t *testing.T) {
	t.Run("should create Admin with all fields", func(t *testing.T) {
		id := uuid.New()
		admin := AdminUser{
			ID:       id,
			Name:     "John",
			Surname:  "Doe",
			Email:    "admin@example.com",
			IsActive: true,
			Roles: []model.AdminRole{
				{Role: model.Role{Name: "superadmin"}},
			},
		}

		assert.Equal(t, id, admin.ID)
		assert.Equal(t, "John", admin.Name)
		assert.Equal(t, "Doe", admin.Surname)
		assert.Equal(t, "admin@example.com", admin.Email)
		assert.True(t, admin.IsActive)
		assert.Len(t, admin.Roles, 1)
	})

	t.Run("should support inactive admins", func(t *testing.T) {
		admin := AdminUser{
			ID:       uuid.New(),
			Email:    "inactive@example.com",
			IsActive: false,
		}

		assert.False(t, admin.IsActive)
	})
}

func TestAdminUserList_Structure(t *testing.T) {
	t.Run("should create AdminUserList with admins", func(t *testing.T) {
		admins := []AdminUser{
			{
				ID:       uuid.New(),
				Name:     "Admin 1",
				Email:    "admin1@example.com",
				IsActive: true,
			},
			{
				ID:       uuid.New(),
				Name:     "Admin 2",
				Email:    "admin2@example.com",
				IsActive: true,
			},
		}

		adminList := AdminUserList{
			Admins: admins,
		}

		assert.Len(t, adminList.Admins, 2)
	})

	t.Run("should handle empty admin list", func(t *testing.T) {
		adminList := AdminUserList{
			Admins: []AdminUser{},
		}

		assert.Empty(t, adminList.Admins)
	})
}

func TestAdminCreateRequest_Structure(t *testing.T) {
	t.Run("should create AdminCreateRequest with all fields", func(t *testing.T) {
		request := AdminUserCreateRequest{
			Name:     "Admin User",
			Surname:  "Doe",
			Email:    "admin@example.com",
			Password: "StrongPassword123!",
			IsActive: true,
			Roles:    []string{"superadmin"},
		}

		assert.Equal(t, "Admin User", request.Name)
		assert.Equal(t, "Doe", request.Surname)
		assert.Equal(t, "admin@example.com", request.Email)
		assert.Equal(t, "StrongPassword123!", request.Password)
		assert.True(t, request.IsActive)
		assert.Len(t, request.Roles, 1)
	})

	t.Run("should allow creating admin without surname", func(t *testing.T) {
		request := AdminUserCreateRequest{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "StrongPassword123!",
			IsActive: true,
		}

		assert.Equal(t, "Admin User", request.Name)
		assert.Empty(t, request.Surname)
	})
}

func TestAdminUpdateRequest_Structure(t *testing.T) {
	t.Run("should create AdminUpdateRequest with all fields", func(t *testing.T) {
		name := "Updated Name"
		surname := "Updated Surname"

		request := AdminUserUpdateRequest{
			Name:    &name,
			Surname: &surname,
		}

		assert.NotNil(t, request.Name)
		assert.Equal(t, "Updated Name", *request.Name)
		assert.NotNil(t, request.Surname)
		assert.Equal(t, "Updated Surname", *request.Surname)
	})

	t.Run("should handle nil fields in AdminUpdateRequest", func(t *testing.T) {
		request := AdminUserUpdateRequest{}

		assert.Nil(t, request.Name)
		assert.Nil(t, request.Surname)
	})

	t.Run("should allow partial updates", func(t *testing.T) {
		name := "Just Name Update"

		request := AdminUserUpdateRequest{
			Name: &name,
		}

		assert.NotNil(t, request.Name)
		assert.Nil(t, request.Surname)
	})
}
