package controller

import (
	"errors"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// ADMIN MANAGEMENT
// =====================

// ListAdmins returns all admin users
// @Summary List all admins
// @Description Get list of all admin users with their roles
// @Tags Admin Management
// @Produce json
// @Success 200 {array} DTO.AdminDetail
// @Failure 401 {object} lib.ErrorResponse
// @Security BearerAuth
// @Router /admin/list [get]
func ListAdmins(c *fiber.Ctx) error {
	// Verify admin authentication
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

	// Convert to DTO
	adminList := make([]DTO.AdminDetail, len(admins))
	for i, admin := range admins {
		roleNames := make([]string, len(admin.Roles))
		for j, role := range admin.Roles {
			roleNames[j] = role.Name
		}
		adminList[i] = DTO.AdminDetail{
			ID:       admin.ID,
			Name:     admin.Name,
			Email:    admin.Email,
			IsActive: admin.IsActive,
			Roles:    roleNames,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    adminList,
	})
}

// CreateAdmin creates a new admin user
// @Summary Create admin
// @Description Create a new admin user with specified roles
// @Tags Admin Management
// @Accept json
// @Produce json
// @Param admin body DTO.AdminCreateRequest true "Admin creation data"
// @Success 201 {object} DTO.AdminDetail
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Security BearerAuth
// @Router /admin/create [post]
func CreateAdmin(c *fiber.Ctx) error {
	// Verify admin authentication (only superadmin can create admins)
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req DTO.AdminCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid request body: %w", err))
	}

	if err := lib.MyCustomStructValidator(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Check if email already exists
	var existingAdmin model.Admin
	if err := tx.Where("email = ?", req.Email).First(&existingAdmin).Error; err == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("admin with email %s already exists", req.Email))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Create admin
	admin := model.Admin{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		IsActive: req.IsActive,
	}

	// Start transaction
	if err := tx.Transaction(func(tx *gorm.DB) error {
		// Create admin
		if err := tx.Create(&admin).Error; err != nil {
			return err
		}

		// Assign roles if provided
		if len(req.Roles) > 0 {
			var roles []model.RoleAdmin
			if err := tx.Where("name IN ?", req.Roles).Find(&roles).Error; err != nil {
				return err
			}

			if len(roles) != len(req.Roles) {
				return fmt.Errorf("some roles not found")
			}

			if err := tx.Model(&admin).Association("Roles").Append(&roles); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Reload admin with roles
	if err := tx.Preload("Roles").First(&admin, admin.ID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	roleNames := make([]string, len(admin.Roles))
	for i, role := range admin.Roles {
		roleNames[i] = role.Name
	}

	response := DTO.AdminDetail{
		ID:       admin.ID,
		Name:     admin.Name,
		Email:    admin.Email,
		IsActive: admin.IsActive,
		Roles:    roleNames,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateAdmin updates an existing admin
// @Summary Update admin
// @Description Update admin user information
// @Tags Admin Management
// @Accept json
// @Produce json
// @Param id path string true "Admin ID"
// @Param admin body DTO.AdminUpdateRequest true "Admin update data"
// @Success 200 {object} DTO.AdminDetail
// @Failure 400 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Security BearerAuth
// @Router /admin/{id} [patch]
func UpdateAdmin(c *fiber.Ctx) error {
	// Verify admin authentication
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	adminID := c.Params("id")
	if adminID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("admin ID required"))
	}

	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid admin ID"))
	}

	var req DTO.AdminUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid request body: %w", err))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Find admin
	var admin model.Admin
	if err := tx.Preload("Roles").First(&admin, adminUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("admin not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Update fields
	if req.Name != nil {
		admin.Name = *req.Name
	}
	if req.Email != nil {
		admin.Email = *req.Email
	}
	if req.Password != nil {
		admin.Password = *req.Password
	}
	if req.IsActive != nil {
		admin.IsActive = *req.IsActive
	}

	// Start transaction
	if err := tx.Transaction(func(tx *gorm.DB) error {
		// Save admin
		if err := tx.Save(&admin).Error; err != nil {
			return err
		}

		// Update roles if provided
		if len(req.Roles) > 0 {
			var roles []model.RoleAdmin
			if err := tx.Where("name IN ?", req.Roles).Find(&roles).Error; err != nil {
				return err
			}

			if err := tx.Model(&admin).Association("Roles").Replace(&roles); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Reload admin with roles
	if err := tx.Preload("Roles").First(&admin, admin.ID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	roleNames := make([]string, len(admin.Roles))
	for i, role := range admin.Roles {
		roleNames[i] = role.Name
	}

	response := DTO.AdminDetail{
		ID:       admin.ID,
		Name:     admin.Name,
		Email:    admin.Email,
		IsActive: admin.IsActive,
		Roles:    roleNames,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteAdmin soft deletes an admin
// @Summary Delete admin
// @Description Soft delete an admin user
// @Tags Admin Management
// @Param id path string true "Admin ID"
// @Success 204
// @Failure 404 {object} lib.ErrorResponse
// @Security BearerAuth
// @Router /admin/{id} [delete]
func DeleteAdmin(c *fiber.Ctx) error {
	// Verify admin authentication
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	adminID := c.Params("id")
	if adminID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("admin ID required"))
	}

	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid admin ID"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Soft delete
	result := tx.Delete(&model.Admin{}, adminUUID)
	if result.Error != nil {
		return lib.Error.General.InternalError.WithError(result.Error)
	}
	if result.RowsAffected == 0 {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("admin not found"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// =====================
// ROLE MANAGEMENT
// =====================

// ListRoles returns all admin roles
// @Summary List admin roles
// @Description Get list of all admin roles
// @Tags Admin Roles
// @Produce json
// @Success 200 {array} DTO.RoleAdminDetail
// @Security BearerAuth
// @Router /admin/roles [get]
func ListRoles(c *fiber.Ctx) error {
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

	roleList := make([]DTO.RoleAdminDetail, len(roles))
	for i, role := range roles {
		roleList[i] = DTO.RoleAdminDetail{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    roleList,
	})
}

// CreateRole creates a new admin role
// @Summary Create admin role
// @Description Create a new admin role
// @Tags Admin Roles
// @Accept json
// @Produce json
// @Param role body DTO.RoleAdminCreateRequest true "Role data"
// @Success 201 {object} DTO.RoleAdminDetail
// @Security BearerAuth
// @Router /admin/roles [post]
func CreateRole(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req DTO.RoleAdminCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid request body: %w", err))
	}

	if err := lib.MyCustomStructValidator(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	role := model.RoleAdmin{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := tx.Create(&role).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	response := DTO.RoleAdminDetail{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateRole updates an admin role
// @Summary Update admin role
// @Description Update an existing admin role
// @Tags Admin Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body DTO.RoleAdminUpdateRequest true "Role update data"
// @Success 200 {object} DTO.RoleAdminDetail
// @Security BearerAuth
// @Router /admin/roles/{id} [patch]
func UpdateRole(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	roleID := c.Params("id")
	if roleID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("role ID required"))
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid role ID"))
	}

	var req DTO.RoleAdminUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid request body: %w", err))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var role model.RoleAdmin
	if err := tx.First(&role, roleUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("role not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := tx.Save(&role).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	response := DTO.RoleAdminDetail{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteRole soft deletes an admin role
// @Summary Delete admin role
// @Description Soft delete an admin role
// @Tags Admin Roles
// @Param id path string true "Role ID"
// @Success 204
// @Security BearerAuth
// @Router /admin/roles/{id} [delete]
func DeleteRole(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	roleID := c.Params("id")
	if roleID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("role ID required"))
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid role ID"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	result := tx.Delete(&model.RoleAdmin{}, roleUUID)
	if result.Error != nil {
		return lib.Error.General.InternalError.WithError(result.Error)
	}
	if result.RowsAffected == 0 {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("role not found"))
	}

	return c.SendStatus(fiber.StatusNoContent)
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
		ListAdmins,
		CreateAdmin,
		UpdateAdmin,
		DeleteAdmin,
		ListRoles,
		CreateRole,
		UpdateRole,
		DeleteRole,
	})
}
