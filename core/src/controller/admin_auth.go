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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AdminLogin handles admin authentication
// @Summary Admin login
// @Description Authenticate admin user and return JWT token
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param login body DTO.AdminLoginRequest true "Admin login credentials"
// @Success 200 {object} DTO.AdminLoginResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 400 {object} lib.ErrorResponse
// @Router /admin/auth/login [post]
func AdminLogin(c *fiber.Ctx) error {
	var loginReq DTO.AdminLoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid request body: %w", err))
	}

	// Validate request
	if err := lib.MyCustomStructValidator(&loginReq); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Ensure we're in public schema
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Find admin by email
	var admin model.Admin
	if err := tx.Where("email = ?", loginReq.Email).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, description")
		}).
		First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("invalid credentials"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Check if admin is active
	if !admin.IsActive {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin account is inactive"))
	}

	// Verify password
	if !admin.MatchPassword(admin.Password) {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("invalid credentials"))
	}

	// Extract role names
	roleNames := make([]string, len(admin.Roles))
	for i, role := range admin.Roles {
		roleNames[i] = role.Name
	}

	// Create admin claims
	adminClaims := DTO.AdminClaims{
		ID:       admin.ID,
		Name:     admin.Name,
		Email:    admin.Email,
		Password: admin.Password, // Store hashed password for token validation
		IsAdmin:  true,
		IsActive: admin.IsActive,
		Roles:    roleNames,
		Type:     namespace.AdminKey.Name,
	}

	// Generate JWT token
	token, err := handler.JWT(c).Encode(adminClaims)
	if err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to generate token: %w", err))
	}

	// Prepare response
	response := DTO.AdminLoginResponse{
		Token: token,
		Admin: &DTO.AdminDetail{
			ID:       admin.ID,
			Name:     admin.Name,
			Email:    admin.Email,
			IsActive: admin.IsActive,
			Roles:    roleNames,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// AdminMe returns current admin information from JWT token
// @Summary Get current admin info
// @Description Returns the currently authenticated admin's information
// @Tags Admin Auth
// @Produce json
// @Success 200 {object} DTO.AdminDetail
// @Failure 401 {object} lib.ErrorResponse
// @Security BearerAuth
// @Router /admin/auth/me [get]
func AdminMe(c *fiber.Ctx) error {
	// Get admin claims from context (set by WhoAreYou middleware)
	adminClaims := c.Locals(namespace.RequestKey.Auth_Claims + "_admin")
	claim, ok := adminClaims.(*DTO.AdminClaims)
	if !ok || claim == nil || !claim.IsAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin authentication required"))
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Ensure we're in public schema
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Fetch fresh admin data with roles
	var admin model.Admin
	if err := tx.Where("id = ?", claim.ID).
		Preload("Roles").
		First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("admin not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Verify password hasn't changed (token invalidation check)
	if admin.Password != claim.Password {
		return lib.Error.Auth.InvalidToken.WithError(fmt.Errorf("token is outdated, please login again"))
	}

	// Extract role names
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

// AdminRefreshToken refreshes the admin JWT token
// @Summary Refresh admin token
// @Description Generates a new JWT token for the authenticated admin
// @Tags Admin Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} lib.ErrorResponse
// @Security BearerAuth
// @Router /admin/auth/refresh [post]
func AdminRefreshToken(c *fiber.Ctx) error {
	// Get admin claims from context
	adminClaims := c.Locals(namespace.RequestKey.Auth_Claims + "_admin")
	claim, ok := adminClaims.(*DTO.AdminClaims)
	if !ok || claim == nil || !claim.IsAdmin {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin authentication required"))
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Ensure we're in public schema
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Fetch fresh admin data
	var admin model.Admin
	if err := tx.Where("id = ?", claim.ID).
		Preload(clause.Associations).
		First(&admin).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Check if admin is still active
	if !admin.IsActive {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("admin account is inactive"))
	}

	// Extract role names
	roleNames := make([]string, len(admin.Roles))
	for i, role := range admin.Roles {
		roleNames[i] = role.Name
	}

	// Create fresh admin claims
	freshClaims := DTO.AdminClaims{
		ID:       admin.ID,
		Name:     admin.Name,
		Email:    admin.Email,
		Password: admin.Password,
		IsAdmin:  true,
		IsActive: admin.IsActive,
		Roles:    roleNames,
		Type:     namespace.AdminKey.Name,
	}

	// Generate new JWT token
	token, err := handler.JWT(c).Encode(freshClaims)
	if err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to generate token: %w", err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": token,
		},
	})
}

// AdminAuth registers all admin authentication route handlers
func AdminAuth(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		AdminLogin,
		AdminMe,
		AdminRefreshToken,
	})
}
