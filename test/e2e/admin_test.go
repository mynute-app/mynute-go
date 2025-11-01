package e2e_test

import (
	"fmt"
	"mynute-go/core"
	coreModel "mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"
	"testing"

	"github.com/google/uuid"
)

func Test_Admin(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Bootstrap a superadmin for testing
	// Roles are already seeded via InitialSeed()
	superadminEmail := lib.GenerateRandomEmail("superadmin")
	superadmin := &model.Admin{
		Created: &coreModel.Admin{
			BaseModel: coreModel.BaseModel{ID: uuid.New()},
			Name:      "System Administrator",
			Email:     superadminEmail,
			IsActive:  true,
		},
	}

	// Create superadmin user in database
	db := server.Db.Gorm
	lib.ChangeToPublicSchema(db)

	// Generate password (BeforeCreate hook will hash it automatically)
	password := lib.GenerateValidPassword()
	superadmin.Created.Password = password

	// Create the admin (BeforeCreate hook handles hashing)
	if err := db.Create(superadmin.Created).Error; err != nil {
		t.Fatalf("Failed to create superadmin: %v", err)
	}

	// Assign superadmin role (already seeded)
	var superadminRole coreModel.RoleAdmin
	if err := db.Where("name = ?", "superadmin").First(&superadminRole).Error; err != nil {
		t.Fatalf("Failed to find superadmin role: %v", err)
	}
	if err := db.Model(superadmin.Created).Association("Roles").Append(&superadminRole); err != nil {
		t.Fatalf("Failed to assign superadmin role: %v", err)
	}

	// Password is already stored plain text in variable for login
	// (the DB has the hashed version)

	tt.Describe("Superadmin login").Test(superadmin.Login(200, password))

	tt.Describe("Superadmin get me").Test(superadmin.GetMe(200))

	// Test admin creation (only superadmin can create admins)
	newAdmin := &model.Admin{}
	newAdmin.X_Auth_Token = superadmin.X_Auth_Token

	tt.Describe("Create support admin").Test(newAdmin.Create(201, "support"))

	// Test login with new admin
	newAdminPassword := newAdmin.Created.Password
	tt.Describe("Login with new admin").Test(newAdmin.Login(200, newAdminPassword))

	tt.Describe("Get new admin info").Test(newAdmin.GetMe(200))

	// Test that non-superadmin cannot create other admins
	anotherAdmin := &model.Admin{}
	anotherAdmin.X_Auth_Token = newAdmin.X_Auth_Token
	tt.Describe("Non-superadmin cannot create admin").Test(anotherAdmin.Create(403, "support"))

	// Test listing admins
	tt.Describe("List all admins").Test(func() error {
		admins, err := superadmin.ListAdmins(200)
		if err != nil {
			return err
		}
		if len(admins) < 2 {
			return fmt.Errorf("expected at least 2 admins, got %d", len(admins))
		}
		return nil
	}())

	//Test admin update
	// Use a separate admin object with superadmin token to avoid corrupting the superadmin object
	adminUpdater := &model.Admin{
		X_Auth_Token: superadmin.X_Auth_Token,
		Created:      &coreModel.Admin{},
	}
	tt.Describe("Update admin name").Test(adminUpdater.Update(200, newAdmin.Created.ID, map[string]any{
		"name": "Updated Support Admin",
	}))

	// Sync the updated data to newAdmin
	newAdmin.Created.Name = adminUpdater.Created.Name
	newAdmin.Created.Email = adminUpdater.Created.Email // This is critical - email must be preserved

	// Test password update
	newPassword := lib.GenerateValidPassword()
	tt.Describe("Update admin password").Test(adminUpdater.Update(200, newAdmin.Created.ID, map[string]any{
		"password": newPassword,
	}))

	// Sync the updated password to newAdmin
	newAdmin.Created.Password = newPassword

	// Test login with old password fails
	tt.Describe("Login with old password fails").Test(newAdmin.Login(401, newAdminPassword))

	// Test login with new password
	tt.Describe("Login with new password").Test(newAdmin.Login(200, newPassword))

	// Test token refresh
	oldToken := newAdmin.X_Auth_Token
	tt.Describe("Refresh admin token").Test(newAdmin.RefreshToken(200))
	if newAdmin.X_Auth_Token == oldToken {
		t.Error("Token should have been refreshed")
	}

	// Test deactivating admin
	tt.Describe("Deactivate admin").Test(adminUpdater.Update(200, newAdmin.Created.ID, map[string]any{
		"is_active": false,
	}))

	// Test that deactivated admin cannot login
	tt.Describe("Deactivated admin cannot login").Test(newAdmin.Login(401, newPassword))

	// Reactivate admin
	tt.Describe("Reactivate admin").Test(adminUpdater.Update(200, newAdmin.Created.ID, map[string]any{
		"is_active": true,
	}))

	// Test login after reactivation
	tt.Describe("Login after reactivation").Test(newAdmin.Login(200, newPassword))

	// ==================
	// ROLE MANAGEMENT
	// ==================

	// Test listing roles
	tt.Describe("List all roles").Test(func() error {
		roles, err := superadmin.ListRoles(200)
		if err != nil {
			return err
		}
		if len(roles) < 3 {
			return fmt.Errorf("expected at least 3 default roles, got %d", len(roles))
		}
		return nil
	}())

	// Test creating a new role (superadmin only)
	tt.Describe("Create new admin role").Test(func() error {
		role, err := superadmin.CreateRole(201, "developer", "Developer access for technical support")
		if err != nil {
			return err
		}
		if role.Name != "developer" {
			return fmt.Errorf("expected role name 'developer', got '%s'", role.Name)
		}
		return nil
	}())

	// Test that non-superadmin cannot create roles
	tt.Describe("Non-superadmin cannot create role").Test(func() error {
		_, err := newAdmin.CreateRole(403, "unauthorized", "Should fail")
		if err != nil {
			return fmt.Errorf("expected 403 forbidden: %w", err)
		}
		return nil
	}())

	// Test updating a role
	tt.Describe("List roles to get developer role ID").Test(func() error {
		roles, err := superadmin.ListRoles(200)
		if err != nil {
			return err
		}

		var developerRole *coreModel.RoleAdmin
		for _, role := range roles {
			if role.Name == "developer" {
				developerRole = &role
				break
			}
		}

		if developerRole == nil {
			return fmt.Errorf("developer role not found")
		}

		// Update the role
		updated, err := superadmin.UpdateRole(200, developerRole.ID, map[string]any{
			"description": "Updated developer description",
		})
		if err != nil {
			return err
		}

		if updated.Description != "Updated developer description" {
			return fmt.Errorf("role description was not updated")
		}

		return nil
	}())

	// Test creating admin with new role
	devAdmin := &model.Admin{}
	devAdmin.X_Auth_Token = superadmin.X_Auth_Token
	tt.Describe("Create admin with developer role").Test(devAdmin.Create(201, "developer"))

	tt.Describe("Login with developer admin").Test(devAdmin.Login(200, devAdmin.Created.Password))

	// Test creating admin with multiple roles
	multiRoleAdmin := &model.Admin{}
	multiRoleAdmin.X_Auth_Token = superadmin.X_Auth_Token
	tt.Describe("Create admin with multiple roles").Test(multiRoleAdmin.Create(201, "support", "auditor"))

	// Test invalid login credentials
	wrongPassword := lib.GenerateValidPassword() // Generate a valid password that's different from the correct one
	invalidLoginTest := &model.Admin{
		Created: &coreModel.Admin{
			Email: superadminEmail,
		},
	}
	tt.Describe("Login with invalid password").Test(invalidLoginTest.Login(401, wrongPassword))

	tt.Describe("Login with non-existent email").Test(func() error {
		fakeAdmin := &model.Admin{
			Created: &coreModel.Admin{
				Email: "nonexistent@example.com",
			},
		}
		return fakeAdmin.Login(401, "SomePassword123!")
	}())

	// Test admin deletion
	tt.Describe("Delete developer admin").Test(superadmin.Delete(204, devAdmin.Created.ID))

	// Test that deleted admin cannot login
	tt.Describe("Deleted admin cannot login").Test(devAdmin.Login(401, devAdmin.Created.Password))

	// Test deleting a role (get the developer role ID first)
	tt.Describe("Delete developer role").Test(func() error {
		roles, err := superadmin.ListRoles(200)
		if err != nil {
			return err
		}

		var developerRole *coreModel.RoleAdmin
		for _, role := range roles {
			if role.Name == "developer" {
				developerRole = &role
				break
			}
		}

		if developerRole == nil {
			return fmt.Errorf("developer role not found")
		}

		return superadmin.DeleteRole(204, developerRole.ID)
	}())

	// Verify role was deleted
	tt.Describe("Verify developer role was deleted").Test(func() error {
		roles, err := superadmin.ListRoles(200)
		if err != nil {
			return err
		}

		for _, role := range roles {
			if role.Name == "developer" {
				return fmt.Errorf("developer role should have been deleted")
			}
		}
		return nil
	}())

	// Test that admin can access tenant endpoints (bypass check)
	// This would require setting up a company and testing access
	// For now, we've tested the core admin functionality

	// Use the original password and email for final login test
	finalLoginTest := &model.Admin{
		Created: &coreModel.Admin{
			Email: superadminEmail,
		},
	}
	tt.Describe("Final superadmin login").Test(finalLoginTest.Login(200, password))
}
