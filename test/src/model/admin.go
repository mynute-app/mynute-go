package model

import (
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"mynute-go/test/src/handler"

	"github.com/google/uuid"
)

type Admin struct {
	Created      *model.Admin
	Roles        []model.RoleAdmin
	X_Auth_Token string
}

// Set creates and configures a complete admin user for testing
func (a *Admin) Set() error {
	if err := a.Create(200); err != nil {
		return err
	}
	if err := a.Login(200, a.Created.Password); err != nil {
		return err
	}
	if err := a.GetMe(200); err != nil {
		return err
	}
	return nil
}

// Create creates a new admin user
func (a *Admin) Create(s int, roles ...string) error {
	pswd := lib.GenerateValidPassword()

	// Default to support role if no roles specified
	if len(roles) == 0 {
		roles = []string{"support"}
	}

	createReq := DTO.AdminCreateRequest{
		Name:     lib.GenerateRandomName("Admin"),
		Email:    lib.GenerateRandomEmail("admin"),
		Password: pswd,
		IsActive: true,
		Roles:    roles,
	}

	var response struct {
		Data DTO.AdminDetail `json:"data"`
	}

	if err := handler.NewHttpClient().
		Method("POST").
		URL("/admin").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(createReq).
		ParseResponse(&response).
		Error; err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	// Map the response to the admin model
	if s >= 200 && s < 300 {
		a.Created = &model.Admin{
			BaseModel: model.BaseModel{ID: response.Data.ID},
			Name:      response.Data.Name,
			Email:     response.Data.Email,
			Password:  pswd, // Store the plain password for future login tests
			IsActive:  response.Data.IsActive,
		}
	}

	return nil
}

// CreateSuperAdmin creates a new superadmin user
func (a *Admin) CreateSuperAdmin(s int) error {
	return a.Create(s, "superadmin")
}

// Login authenticates an admin user
func (a *Admin) Login(s int, password string) error {
	loginReq := DTO.AdminLoginRequest{
		Email:    a.Created.Email,
		Password: password,
	}

	var loginResp DTO.AdminLoginResponse
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/admin/auth/login").
		ExpectedStatus(s).
		Send(loginReq).
		ParseResponse(&loginResp).
		Error; err != nil {
		return fmt.Errorf("failed to login admin: %w", err)
	}

	if s == 200 {
		a.X_Auth_Token = loginResp.Token
		if loginResp.Admin != nil {
			a.Created.Name = loginResp.Admin.Name
			a.Created.Email = loginResp.Admin.Email
			a.Created.IsActive = loginResp.Admin.IsActive
		}
	}

	return nil
}

// GetMe retrieves the current admin's information
func (a *Admin) GetMe(s int) error {
	var adminDetail DTO.AdminDetail
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/admin/auth/me").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		ParseResponse(&adminDetail).
		Error; err != nil {
		return fmt.Errorf("failed to get admin me: %w", err)
	}

	if s == 200 {
		a.Created.Name = adminDetail.Name
		a.Created.Email = adminDetail.Email
		a.Created.IsActive = adminDetail.IsActive
	}

	return nil
}

// RefreshToken refreshes the admin's JWT token
func (a *Admin) RefreshToken(s int) error {
	var tokenResp struct {
		Token string `json:"token"`
	}

	if err := handler.NewHttpClient().
		Method("POST").
		URL("/admin/auth/refresh").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, a.X_Auth_Token).
		Send(nil).
		ParseResponse(&tokenResp).
		Error; err != nil {
		return fmt.Errorf("failed to refresh admin token: %w", err)
	}

	if s == 200 {
		a.X_Auth_Token = tokenResp.Token
	}

	return nil
}

// Update updates admin information
func (a *Admin) Update(s int, adminID uuid.UUID, changes map[string]any) error {
	if len(changes) == 0 {
		return fmt.Errorf("no changes provided")
	}

	var response struct {
		Data DTO.AdminDetail `json:"data"`
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
			BaseModel: model.BaseModel{ID: response.Data.ID},
			Name:      response.Data.Name,
			Email:     response.Data.Email,
			IsActive:  response.Data.IsActive,
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
		a.Created.ID = response.Data.ID
		a.Created.Name = response.Data.Name
		a.Created.Email = response.Data.Email
		a.Created.IsActive = response.Data.IsActive
		// If password was updated, store the new plain password for future login tests
		if newPass, ok := changes["password"].(string); ok {
			a.Created.Password = newPass
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
		Data []DTO.AdminDetail `json:"data"`
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
			BaseModel: model.BaseModel{ID: dto.ID},
			Name:      dto.Name,
			Email:     dto.Email,
			IsActive:  dto.IsActive,
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
		Success bool                `json:"success"`
		Data    DTO.RoleAdminDetail `json:"data"`
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
		BaseModel:   model.BaseModel{ID: response.Data.ID},
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}
	return role, nil
}

// ListRoles retrieves all admin roles
func (a *Admin) ListRoles(s int) ([]model.RoleAdmin, error) {
	var response struct {
		Success bool                  `json:"success"`
		Data    []DTO.RoleAdminDetail `json:"data"`
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
			BaseModel:   model.BaseModel{ID: dto.ID},
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
		Success bool                `json:"success"`
		Data    DTO.RoleAdminDetail `json:"data"`
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
		BaseModel:   model.BaseModel{ID: response.Data.ID},
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
