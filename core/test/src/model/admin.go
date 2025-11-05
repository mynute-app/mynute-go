package model

import (
	"fmt"
	DTO "mynute-go/core/src/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/lib/emailclient"
	"mynute-go/core/test/src/handler"
	"net/url"

	"github.com/google/uuid"
)

type Admin struct {
	Created      *model.Admin
	X_Auth_Token string
	Email        string // Email is stored in auth.User, cached here for tests
	Password     string // Password is stored in auth.User, cached here for tests
}

// Set creates and configures a complete admin user for testing
func (A *Admin) Set(roles []string, newAdmin *Admin) error {
	if len(roles) == 0 {
		roles = []string{"superadmin"}
	}
	var a *Admin
	if newAdmin == nil {
		a = A
	} else {
		a = newAdmin
	}
	a, err := a.Create(200, roles...)
	if err != nil {
		return err
	}
	if err := a.VerifyEmail(200); err != nil {
		return err
	}
	// TODO: Password is on User model in auth service, not on Admin
	// This test needs architectural update to work with User/Admin split
	if err := a.LoginByPassword(200, "<temp-password>"); err != nil {
		return err
	}
	if err := a.ResetPasswordByEmail(200); err != nil {
		return err
	}
	return nil
}

func (a *Admin) SendPasswordResetEmail(s int) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/admin/reset-password/%s?lang=en", url.PathEscape(a.Created.UserID.String() /* Email on User */))).
		ExpectedStatus(s).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send password reset email to client: %w", err)
	}
	return nil
}

func (a *Admin) GetNewPasswordFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog := emailclient.NewMailHogClient()

	// Get the latest email sent to the admin
	message, err := mailhog.FindMessageByRecipient(a.Created.UserID.String() /* Email on User */)
	if err != nil {
		return "", err
	}

	// Verify the email has a subject
	if len(message.Content.Headers["Subject"]) == 0 {
		return "", fmt.Errorf("email subject is empty")
	}

	// Extract the new password from the email
	password, err := message.ExtractPassword()
	if err != nil {
		return "", fmt.Errorf("failed to extract password: %w", err)
	}

	return password, nil
}

func (a *Admin) ResetPasswordByEmail(s int) error {
	if err := a.SendPasswordResetEmail(s); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	newPassword, err := a.GetNewPasswordFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get new password from email: %w", err)
	}

	// Update the password in test wrapper
	a.Password = newPassword

	// Try to login with the new password
	if err := a.LoginByPassword(200, newPassword); err != nil {
		return fmt.Errorf("failed to login with new password: %w", err)
	}

	return nil
}

func (a *Admin) SendVerificationEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/admin/send-verification-code/email/%s?language=en", url.PathEscape(a.Created.UserID.String() /* Email on User */))).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	return nil
}

func (a *Admin) GetVerificationCodeFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog := emailclient.NewMailHogClient()

	// Get the latest email sent to the client
	message, err := mailhog.FindMessageByRecipient(a.Created.UserID.String() /* Email on User */)
	if err != nil {
		return "", err
	}

	// Verify the email has a subject
	if message.GetSubject() == "" {
		return "", fmt.Errorf("email subject is empty")
	}

	// Extract the verification code from the email
	code, err := message.ExtractValidationCode()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (a *Admin) VerifyEmailByCode(s int, code string) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("GET").
		URL(fmt.Sprintf("/admin/verify-email/%s/%s", url.PathEscape(a.Created.UserID.String() /* Email on User */), code)).
		ExpectedStatus(s).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to verify admin email: %w", err)
	}

	// Note: Admin data will be refreshed after login, no need to fetch it here
	// since we don't have an auth token yet
	return nil
}

func (a *Admin) VerifyEmail(s int) error {
	if err := a.SendVerificationEmail(s); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	code, err := a.GetVerificationCodeFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get verification code from email: %w", err)
	}

	if err := a.VerifyEmailByCode(s, code); err != nil {
		return fmt.Errorf("failed to verify email with code: %w", err)
	}

	return nil
}

// checkSuperAdminExists checks if there are any superadmin users in the system
// This is a private helper to avoid requiring authentication
func checkSuperAdminExists() (bool, error) {
	var response map[string]bool
	// Try without authentication - endpoint may require auth
	err := handler.NewHttpClient().
		Method("GET").
		URL("/admin/are_there_any_superadmin").
		Send(nil).
		ParseResponse(&response).
		Error

	if err != nil {
		return false, err
	}

	hasSuperAdmin, ok := response["has_superadmin"]
	if !ok {
		// Invalid response format, assume no superadmin exists
		return false, nil
	}

	return hasSuperAdmin, nil
}

// createAdminFromDTO converts DTO.Admin to model.Admin with the provided password
func createAdminFromDTO(dtoAdmin DTO.Admin, password string) *Admin {
	admin := &Admin{
		Created: &model.Admin{
			UserID: dtoAdmin.ID, // Admin uses UserID as primary key
			Name:   dtoAdmin.Name,
			// Email:    dtoAdmin.Email, // TODO: Email is on User model
			// Password: password, // TODO: Password is on User model
			IsActive: dtoAdmin.IsActive,
			Roles:    make([]model.RoleAdmin, len(dtoAdmin.Roles)),
		},
		Email:    dtoAdmin.Email, // Cache email in test wrapper
		Password: password,       // Cache password in test wrapper
	}

	// Convert roles from DTO to model
	for i, dtoRole := range dtoAdmin.Roles {
		admin.Created.Roles[i] = model.RoleAdmin{
			ID:          dtoRole.ID, // RoleAdmin has direct ID field
			Name:        dtoRole.Name,
			Description: dtoRole.Description,
		}
	}

	return admin
}

// Create creates a new admin user, automatically detecting if this should be
// the first superadmin or a regular admin creation
func (a *Admin) Create(s int, roles ...string) (*Admin, error) {
	// Default to superadmin role if no roles specified
	if len(roles) == 0 {
		roles = []string{"superadmin"}
	}

	// Check if we should try creating the first superadmin
	hasSuperAdmin, err := checkSuperAdminExists()

	if err != nil {
		return nil, err
	}

	if !hasSuperAdmin && len(roles) == 1 && roles[0] == "superadmin" {
		// Try creating as first superadmin (no auth required)
		admin, err := a.createFirstSuperAdmin(s, roles)
		if err == nil {
			return admin, nil
		}
		// If it fails, fall through to regular creation
	}

	// Create using regular endpoint (requires authentication)
	return a.createRegularAdmin(s, roles)
}

// createFirstSuperAdmin creates the first superadmin using the special endpoint
func (a *Admin) createFirstSuperAdmin(s int, roles []string) (*Admin, error) {
	pswd := lib.GenerateValidPassword()

	createReq := DTO.AdminCreateRequest{
		Name:     lib.GenerateRandomName("Admin Name"),
		Surname:  lib.GenerateRandomName("Admin Surname"),
		Email:    lib.GenerateRandomEmail("admin"),
		Password: pswd,
		Roles:    roles,
	}

	var dtoAdmin DTO.Admin
	if err := handler.NewHttpClient().
		Method("POST").
		ExpectedStatus(s).
		URL("/admin/first_superadmin").
		Send(createReq).
		ParseResponse(&dtoAdmin).
		Error; err != nil {
		return nil, err
	}

	return createAdminFromDTO(dtoAdmin, pswd), nil
}

// createRegularAdmin creates an admin using the regular endpoint (requires auth)
func (a *Admin) createRegularAdmin(s int, roles []string) (*Admin, error) {
	pswd := lib.GenerateValidPassword()

	createReq := DTO.AdminCreateRequest{
		Name:     lib.GenerateRandomName("Admin Name"),
		Surname:  lib.GenerateRandomName("Admin Surname"),
		Email:    lib.GenerateRandomEmail("admin"),
		Password: pswd,
		IsActive: true,
		Roles:    roles,
	}

	client := handler.NewHttpClient().
		Method("POST").
		URL("/admin").
		ExpectedStatus(s)

	// Only add auth header if token is not empty
	if a.X_Auth_Token != "" {
		client.Header(namespace.HeadersKey.Auth, a.X_Auth_Token)
	}

	// Parse response as DTO first
	var dtoAdmin DTO.Admin
	if err := client.Send(createReq).
		ParseResponse(&dtoAdmin).
		Error; err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	// Create admin from DTO and preserve auth token
	newAdmin := createAdminFromDTO(dtoAdmin, pswd)
	newAdmin.X_Auth_Token = a.X_Auth_Token

	return newAdmin, nil
}

// CreateSuperAdmin creates a new superadmin user
func (a *Admin) CreateSuperAdmin(s int) (*Admin, error) {
	return a.Create(s, "superadmin")
}

// Login authenticates an admin user
func (a *Admin) LoginByPassword(s int, password string) error {
	loginReq := DTO.AdminLoginRequest{
		Email:    a.Created.UserID.String(), /* Email on User */
		Password: password,
	}

	client := handler.NewHttpClient().
		Method("POST").
		URL("/admin/auth/login").
		ExpectedStatus(s).
		Send(loginReq)

	if err := client.Error; err != nil {
		return fmt.Errorf("failed to login admin: %w", err)
	}

	if s == 200 {
		// Extract token from response headers
		a.X_Auth_Token = client.ResHeaders[namespace.HeadersKey.Auth][0]
	}

	return nil
}

// GetByID retrieves the admin's information by ID
func (a *Admin) GetByID(s int) error {
	var adminDetail DTO.Admin
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/admin/%d", a.Created.UserID)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		ParseResponse(&adminDetail).
		Error; err != nil {
		return fmt.Errorf("failed to get admin by id: %w", err)
	}

	if s == 200 {
		a.Created.Name = adminDetail.Name
		// TODO: Email on User - adminDetail.Email
		a.Created.IsActive = adminDetail.IsActive
	}

	return nil
}

// Update updates admin information
func (a *Admin) Update(s int, adminID uuid.UUID, changes map[string]any) error {
	if len(changes) == 0 {
		return fmt.Errorf("no changes provided")
	}

	var response struct {
		Data DTO.Admin `json:"data"`
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/admin/%s", adminID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(changes).
		ParseResponse(&response).
		Error; err != nil {
		return fmt.Errorf("failed to update admin: %w", err)
	}

	if s >= 200 && s < 300 {
		// Map DTO to model.Admin for validation
		updated := model.Admin{
			UserID: response.Data.ID, // Admin uses UserID
			Name:   response.Data.Name,
			// Email:     response.Data.Email, // Email is on User model
			IsActive: response.Data.IsActive,
		}

		// Skip password validation since passwords are not returned in responses (security)
		changesToValidate := make(map[string]any)
		for k, v := range changes {
			if k != "password" {
				changesToValidate[k] = v
			}
		}

		if len(changesToValidate) > 0 {
			if err := ValidateUpdateChanges("Admin", &updated, changesToValidate); err != nil {
				return err
			}
		}

		// Update the Created field with the latest data from server
		a.Created.UserID = response.Data.ID
		a.Created.Name = response.Data.Name
		// TODO: Email on User - response.Data.Email
		a.Created.IsActive = response.Data.IsActive

		// Convert roles from DTO to model
		a.Created.Roles = make([]model.RoleAdmin, len(response.Data.Roles))
		for i, dtoRole := range response.Data.Roles {
			a.Created.Roles[i] = model.RoleAdmin{
				ID:          dtoRole.ID, // RoleAdmin has direct ID field
				Name:        dtoRole.Name,
				Description: dtoRole.Description,
			}
		}

		// If password was updated, store the new plain password for future login tests
		if newPass, ok := changes["password"].(string); ok {
			a.Password = newPass // Cache in test wrapper
		}
	}

	return nil
}

// Delete deletes an admin user
func (a *Admin) Delete(s int, adminID uuid.UUID) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/admin/%s", adminID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}
	return nil
}

// ListAdmins retrieves all admin users
func (a *Admin) ListAdmins(s int) ([]model.Admin, error) {
	var response struct {
		Data []DTO.Admin `json:"data"`
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/admin").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		ParseResponse(&response).
		Error; err != nil {
		return nil, fmt.Errorf("failed to list admins: %w", err)
	}

	// Convert DTO to model.Admin
	admins := make([]model.Admin, len(response.Data))
	for i, dto := range response.Data {
		admins[i] = model.Admin{
			UserID: dto.ID, // Admin uses UserID
			Name:   dto.Name,
			// Email:     dto.Email, // Email is on User model
			IsActive: dto.IsActive,
		}
	}
	return admins, nil
}

// CreateRole creates a new admin role
func (a *Admin) CreateRole(s int, name, description string) (*model.RoleAdmin, error) {
	roleReq := DTO.RoleAdminCreateRequest{
		Name:        name,
		Description: description,
	}

	var response struct {
		Success bool          `json:"success"`
		Data    DTO.AdminRole `json:"data"`
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/admin/role").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(roleReq).
		ParseResponse(&response).
		Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Convert DTO to model.RoleAdmin
	role := &model.RoleAdmin{
		ID:          response.Data.ID, // RoleAdmin has direct ID field
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}
	return role, nil
}

// ListRoles retrieves all admin roles
func (a *Admin) ListRoles(s int) ([]model.RoleAdmin, error) {
	var response struct {
		Success bool            `json:"success"`
		Data    []DTO.AdminRole `json:"data"`
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/admin/role").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		ParseResponse(&response).
		Error; err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	// Convert DTO to model.RoleAdmin
	roles := make([]model.RoleAdmin, len(response.Data))
	for i, dto := range response.Data {
		roles[i] = model.RoleAdmin{
			ID:          dto.ID, // RoleAdmin has direct ID field
			Name:        dto.Name,
			Description: dto.Description,
		}
	}
	return roles, nil
}

// UpdateRole updates an admin role
func (a *Admin) UpdateRole(s int, roleID uuid.UUID, changes map[string]any) (*model.RoleAdmin, error) {
	if len(changes) == 0 {
		return nil, fmt.Errorf("no changes provided")
	}

	var response struct {
		Success bool          `json:"success"`
		Data    DTO.AdminRole `json:"data"`
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/admin/role/%s", roleID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(changes).
		ParseResponse(&response).
		Error; err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Convert DTO to model.RoleAdmin
	updated := &model.RoleAdmin{
		ID:          response.Data.ID, // RoleAdmin has direct ID field
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}
	return updated, nil
}

// DeleteRole deletes an admin role
func (a *Admin) DeleteRole(s int, roleID uuid.UUID) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/admin/role/%s", roleID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
