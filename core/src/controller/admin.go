package controller

import (
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/middleware"

	"github.com/gofiber/fiber/v2"
)

// =====================
// ADMIN AUTH
// =====================

// AdminLoginByPassword handles admin authentication
//
//	@Summary		Admin login
//	@Description	Authenticate admin user and return JWT token in X-Auth-Token header
//	@Tags			Admin Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body		DTO.AdminLoginRequest	true	"Admin login credentials"
//	@Success		200		"Token returned in X-Auth-Token header"
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/admin/login [post]
func AdminLoginByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.AdminKey.Name, &model.Admin{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

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
//	@Router			/admin/are_there_any_superadmin [get]
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
//	@Description	Create the first admin user in the system
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		DTO.Admin	true	"Admin creation data"
//	@Success		201		{object}	DTO.Admin
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/admin/first_superadmin [post]
func CreateFirstAdmin(c *fiber.Ctx) error {
	hasSuperAdmin, err := areThereAnySuperAdmin(c)
	if err != nil {
		return err
	} else if hasSuperAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("first admin already exists"))
	}

	return CreateAdmin(c)
}

// CreateAdmin creates a new admin user
//
//	@Summary		Create admin
//	@Description	Create a new admin user
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string					true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			admin			body		DTO.AdminCreateRequest	true	"Admin creation data"
//	@Success		201				{object}	DTO.Admin
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/admin [post]
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
	var req DTO.AdminCreateRequest
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
		return lib.Error.General.InternalError.WithError(err)
	}

	// Create admin model
	admin := model.Admin{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		IsActive: req.IsActive,
	}

	// Create admin in database (BeforeCreate hook will validate and hash password)
	if err := tx.Create(&admin).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	// Handle roles
	var rolesToAssign []string
	if !hasSuperAdmin {
		// If this is the first admin, force superadmin role
		rolesToAssign = []string{"superadmin"}
	} else {
		// Use roles from request
		rolesToAssign = req.Roles
	}

	// Find and assign roles
	if len(rolesToAssign) > 0 {
		var roles []model.RoleAdmin
		if err := tx.Where("name IN ?", rolesToAssign).Find(&roles).Error; err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}

		if len(roles) != len(rolesToAssign) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("one or more roles not found"))
		}

		// Associate roles with admin
		if err := tx.Model(&admin).Association("Roles").Append(&roles); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	// Reload admin with roles preloaded
	if err := tx.Preload("Roles").First(&admin, admin.ID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(201, &admin, &DTO.Admin{})
}

// GetAdminByID retrieves an admin by its ID
//
//	@Summary		Get admin by ID
//	@Description	Retrieve an admin by its ID
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Admin ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Admin
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/admin/{id} [get]
func GetAdminByID(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	var admin model.Admin
	if err := GetOneBy("id", c, &admin, nil, nil); err != nil {
		return err
	}
	return lib.ResponseFactory(c).SendDTO(200, &admin, &DTO.Admin{})
}

// GetAdminByEmail retrieves an admin by email
//
//	@Summary		Get admin by email
//	@Description	Retrieve an admin by email address
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Admin email"
//	@Produce		json
//	@Success		200	{object}	DTO.Admin
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/admin/email/{email} [get]
func GetAdminByEmail(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	var admin model.Admin
	if err := GetOneBy("email", c, &admin, nil, nil); err != nil {
		return err
	}
	return lib.ResponseFactory(c).SendDTO(200, &admin, &DTO.Admin{})
}

// ListAdmins returns all admin users
//
//	@Summary		List all admins
//	@Description	Get list of all admin users with their roles
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Produce		json
//	@Success		200	{array}	DTO.AdminList
//	@Router			/admin [get]
func ListAdmins(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var admins []model.Admin
	if err := tx.Preload("Roles").Find(&admins).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	adminList := model.AdminList{
		Admins: admins,
		Total:  len(admins),
	}

	return lib.ResponseFactory(c).SendDTO(200, &adminList, &DTO.AdminList{})
}

// UpdateAdminByID updates an existing admin
//
//	@Summary		Update admin
//	@Description	Update admin user information
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Admin ID"
//	@Param			admin	body		DTO.Admin	true	"Admin update data"
//	@Success		200		{object}	DTO.Admin
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/admin/{id} [patch]
func UpdateAdminByID(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var admin model.Admin
	if err := UpdateOneById(c, &admin, nil); err != nil {
		return err
	}
	return lib.ResponseFactory(c).SendDTO(200, &admin, &DTO.Admin{})
}

// DeleteAdminByID soft deletes an admin
//
//	@Summary		Delete admin
//	@Description	Soft delete an admin user
//	@Tags			Admin
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Admin ID"
//	@Success		204
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/admin/{id} [delete]
func DeleteAdminByID(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	if err := DeleteOneById(c, &model.Admin{}); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// =====================
// ADMIN ROLE MANAGEMENT
// =====================

// CreateAdminRole creates a new admin role
//
//	@Summary		Create admin role
//	@Description	Create a new admin role
//	@Tags			Admin Role
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			role	body		DTO.RoleAdminCreateRequest	true	"Role data"
//	@Success		201		{object}	DTO.AdminRole
//	@Router			/admin/role [post]
func CreateAdminRole(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var role model.RoleAdmin
	if err := Create(c, &role); err != nil {
		return err
	}
	return lib.ResponseFactory(c).SendDTO(201, &role, &DTO.AdminRole{})
}

// GetAdminRoleByID retrieves an admin role by its ID
//
//	@Summary		Get admin role by ID
//	@Description	Retrieve an admin role by its ID
//	@Tags			Admin Role
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Role ID"
//	@Produce		json
//	@Success		200	{object}	DTO.AdminRole
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/admin/role/{id} [get]
func GetAdminRoleByID(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	var role model.RoleAdmin
	if err := GetOneBy("id", c, &role, nil, nil); err != nil {
		return err
	}
	return lib.ResponseFactory(c).SendDTO(200, &role, &DTO.AdminRole{})
}

// ListAdminRoles returns all admin roles
//
//	@Summary		List admin roles
//	@Description	Get list of all admin roles
//	@Tags			Admin Role
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Produce		json
//	@Success		200	{array}	DTO.AdminRole
//	@Router			/admin/role [get]
func ListAdminRoles(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var roles []model.RoleAdmin
	if err := tx.Find(&roles).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	roleList := make([]DTO.AdminRole, len(roles))
	for i, role := range roles {
		roleList[i] = DTO.AdminRole{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return lib.ResponseFactory(c).Send(200, roleList)
}

// UpdateAdminRoleByID updates an admin role
//
//	@Summary		Update admin role
//	@Description	Update an existing admin role
//	@Tags			Admin Role
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Role ID"
//	@Param			role	body		DTO.RoleAdminUpdateRequest	true	"Role update data"
//	@Success		200		{object}	DTO.AdminRole
//	@Router			/admin/role/{id} [patch]
func UpdateAdminRoleByID(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var role model.RoleAdmin
	if err := UpdateOneById(c, &role, nil); err != nil {
		return err
	}
	return lib.ResponseFactory(c).SendDTO(200, &role, &DTO.AdminRole{})
}

// DeleteAdminRoleByID soft deletes an admin role
//
//	@Summary		Delete admin role by ID
//	@Description	Soft delete an admin role
//	@Tags			Admin Role
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Role ID"
//	@Success		204
//	@Router			/admin/role/{id} [delete]
func DeleteAdminRoleByID(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	return DeleteOneById(c, &model.RoleAdmin{})
}

// =====================
// HELPER FUNCTIONS
// =====================

// requireAdmin checks if the request is from an authenticated admin
func requireAdmin(c *fiber.Ctx) error {
	adminClaims := c.Locals(namespace.RequestKey.Auth_Claims + "_admin")
	claim, ok := adminClaims.(*DTO.AdminClaims)
	if !ok || claim == nil || !claim.IsAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin authentication required"))
	}
	if !claim.IsActive {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin account is inactive"))
	}
	return nil
}

func areThereAnySuperAdmin(c *fiber.Ctx) (bool, error) {
	tx, err := lib.Session(c)
	if err != nil {
		return false, lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return false, lib.Error.General.InternalError.WithError(err)
	}

	var count int64
	if err := tx.Model(&model.Admin{}).
		Joins("JOIN admin_role_admins ON admin_role_admins.admin_id = admins.id").
		Joins("JOIN role_admins ON role_admins.id = admin_role_admins.role_admin_id").
		Where("role_admins.name = ?", "superadmin").
		Count(&count).Error; err != nil {
		return false, lib.Error.General.InternalError.WithError(err)
	}
	return count > 0, nil
}

// requireSuperAdmin checks if the request is from a superadmin
func requireSuperAdmin(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	adminClaims := c.Locals(namespace.RequestKey.Auth_Claims + "_admin")
	claim := adminClaims.(*DTO.AdminClaims)

	for _, role := range claim.Roles {
		if role == "superadmin" {
			return nil
		}
	}

	return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("superadmin role required"))
}

// Admin registers all admin management route handlers
func Admin(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		AdminLoginByPassword,
		AreThereAnyAdmin,
		CreateFirstAdmin,
		GetAdminByID,
		GetAdminByEmail,
		ListAdmins,
		CreateAdmin,
		UpdateAdminByID,
		DeleteAdminByID,
		ListAdminRoles,
		CreateAdminRole,
		GetAdminRoleByID,
		UpdateAdminRoleByID,
		DeleteAdminRoleByID,
	})
}
