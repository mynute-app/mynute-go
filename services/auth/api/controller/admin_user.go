package controller

import (
	"fmt"
	"mynute-go/services/auth/api/handler"
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"
	DTO "mynute-go/services/auth/config/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// ADMIN MANAGEMENT
// =====================

// AreThereAnyAdmin checks if there are any superadmin users in the system
//
//	@Summary		Check for superadmin existence
//	@Description	Check if there are any superadmin users in the system
//	@Tags			Admin
//	@Produce		json
//	@Success		200	{object}	map[string]bool
//	@Router			/admin/users/are_there_any_superadmin [get]
func AreThereAnyAdmin(c *fiber.Ctx) error {
	hasAdmin, err := areThereAnySuperAdmin(c)
	if err != nil {
		return err
	}
	return lib.ResponseFactory(c).Send(200, map[string]bool{
		"has_superadmin": hasAdmin,
	})
}

// CreateFirstAdmin creates the first admin user in the system
//
//	@Summary		Create first admin
//	@Description	Create the first admin user in the system (only works when no admins exist)
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		DTO.AdminUserCreateRequest	true	"Admin creation data"
//	@Success		201		{object}	DTO.AdminUser
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		403		{object}	DTO.ErrorResponse
//	@Router			/admin/users/first_superadmin [post]
func CreateFirstAdmin(c *fiber.Ctx) error {
	// Check if any admin users already exist
	hasAdmin, err := areThereAnySuperAdmin(c)
	if err != nil {
		return err
	}
	if hasAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("cannot create first admin: admin users already exist"))
	}

	// Parse request body
	var req DTO.AdminUserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Validate request
	if err := lib.MyCustomStructValidator(req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Check if user with same email already exists
	var existingUser model.AdminUser
	if err := tx.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("user with email %s already exists", req.Email))
	} else if err != gorm.ErrRecordNotFound {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Hash the password
	hashedPassword, err := handler.HashPassword(req.Password)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Create first admin user
	user := model.AdminUser{
		BaseModel: model.BaseModel{ID: uuid.New()},
		Email:     req.Email,
		Password:  hashedPassword,
		Verified:  true, // Admins are verified by default
	}

	if err := tx.Create(&user).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	// Always assign superadmin role to the first admin
	var superadminRole model.AdminRole
	if err := tx.Where("name = ?", "superadmin").First(&superadminRole).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("superadmin role not found: %w", err))
	}
	if err := tx.Model(&user).Association("Roles").Append(&superadminRole); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Reload user with roles for response
	if err := tx.Preload("Roles").Where("id = ?", user.ID).First(&user).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Return user (without password)
	return lib.ResponseFactory(c).SendDTO(201, &user, &DTO.AdminUser{})
}

// CreateAdmin creates a new admin user
//
//	@Summary		Create admin
//	@Description	Create a new admin user
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		DTO.AdminUserCreateRequest	true	"Admin creation data"
//	@Success		201		{object}	DTO.AdminUser
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/admin/users [post]
func CreateAdmin(c *fiber.Ctx) error {
	// Verify admin authentication (only superadmin can create admins)
	hasSuperAdmin, err := areThereAnySuperAdmin(c)
	if err != nil {
		return err
	} else if hasSuperAdmin {
		if err := requireSuperAdmin(c); err != nil {
			return err
		}
	}

	// Parse request body
	var req DTO.AdminUserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Validate request
	if err := lib.MyCustomStructValidator(req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Check if user with same email already exists
	var existingUser model.AdminUser
	if err := tx.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("user with email %s already exists", req.Email))
	} else if err != gorm.ErrRecordNotFound {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Hash the password
	hashedPassword, err := handler.HashPassword(req.Password)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Create admin user
	user := model.AdminUser{
		BaseModel: model.BaseModel{ID: uuid.New()},
		Email:     req.Email,
		Password:  hashedPassword,
		Verified:  true, // Admins are verified by default
	}

	if err := tx.Create(&user).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	// Assign roles to the admin user
	if len(req.Roles) > 0 {
		// Find and assign the specified roles
		var roles []model.AdminRole
		for _, roleName := range req.Roles {
			var role model.AdminRole
			if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return lib.Error.General.BadRequest.WithError(fmt.Errorf("role '%s' not found", roleName))
				}
				return lib.Error.General.InternalError.WithError(err)
			}
			roles = append(roles, role)
		}
		if err := tx.Model(&user).Association("Roles").Append(roles); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	} else if !hasSuperAdmin {
		// If this is the first admin and no roles specified, assign superadmin role
		var superadminRole model.AdminRole
		if err := tx.Where("name = ?", "superadmin").First(&superadminRole).Error; err != nil {
			return lib.Error.General.InternalError.WithError(fmt.Errorf("superadmin role not found"))
		}
		if err := tx.Model(&user).Association("Roles").Append(&superadminRole); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	// Reload user with roles for response
	if err := tx.Preload("Roles").Where("id = ?", user.ID).First(&user).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Return user (without password)
	return lib.ResponseFactory(c).SendDTO(201, &user, &DTO.AdminUser{})
}

// GetAdminById retrieves an admin by ID
//
//	@Summary		Get admin by ID
//	@Description	Retrieve an admin by its ID
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Admin ID"
//	@Produce		json
//	@Success		200	{object}	DTO.AdminUser
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/admin/users/{id} [get]
func GetAdminById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var user model.AdminUser
	if err := GetOneBy("id", c, &user); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.AdminUser{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateAdminById updates an admin by ID
//
//	@Summary		Update admin
//	@Description	Update an admin
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Admin ID"
//	@Param			admin	body		DTO.AdminUserUpdateRequest	true	"Admin"
//	@Success		200		{object}	DTO.AdminUser
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/admin/users/{id} [patch]
func UpdateAdminById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var user model.AdminUser
	if err := UpdateOneById(c, &user); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.AdminUser{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// DeleteAdminById deletes an admin by ID
//
//	@Summary		Delete admin by ID
//	@Description	Delete an admin by its ID
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Admin ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/admin/users/{id} [delete]
func DeleteAdminById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}
	return DeleteOneById(c, &model.AdminUser{})
}

// ListAdmins retrieves all admins
//
//	@Summary		List all admins
//	@Description	Retrieve all admin users
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Produce		json
//	@Success		200	{array}		DTO.AdminUser
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/admin/users [get]
func ListAdmins(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var users []model.AdminUser
	if err := tx.Preload("Roles").Find(&users).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(200, &users, &[]DTO.AdminUser{})
}

// =====================
// HELPER FUNCTIONS
// =====================

// areThereAnySuperAdmin checks if there are any admin users in the system
func areThereAnySuperAdmin(c *fiber.Ctx) (bool, error) {
	tx, err := lib.Session(c)
	if err != nil {
		return false, err
	}

	var count int64
	if err := tx.Model(&model.AdminUser{}).Count(&count).Error; err != nil {
		return false, lib.Error.General.InternalError.WithError(err)
	}

	return count > 0, nil
}

// requireSuperAdmin checks if the current user is a superadmin
func requireSuperAdmin(c *fiber.Ctx) error {
	claims, err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return lib.Error.Auth.InvalidToken.WithError(err)
	}
	if claims == nil {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("authentication required"))
	}

	// Verify this is an admin token
	if !claims.IsAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin token required"))
	}

	// Check if admin has superadmin role
	hasSuperAdmin := false
	for _, role := range claims.Roles {
		if role == "superadmin" {
			hasSuperAdmin = true
			break
		}
	}
	if !hasSuperAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("superadmin privileges required"))
	}
	return nil
}
