package e2e_test

import (
	"fmt"
	"mynute-go/core"
	"mynute-go/core/src/lib"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"
	"testing"
)

func Test_Admin(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	var superadmin *model.Admin

	// The first admin created in the system should get superadmin privileges automatically
	tt.Describe("Create first admin").Test(superadmin.Set([]string{}, nil))

	// Test admin creation (only superadmin can create admins)
	supportAdmin := &model.Admin{}
	tt.Describe("Create support admin").Test(superadmin.Set([]string{"support"}, supportAdmin))

	// Test that non-superadmin cannot create other admins

	tt.Describe("Non-superadmin cannot create admin").Test(func() error {
		newFailedToCreateAdmin := &model.Admin{}
		err := supportAdmin.Set([]string{"support"}, newFailedToCreateAdmin)
		if err == nil {
			return fmt.Errorf("expected error when non-superadmin tries to create admin")
		}
		return nil
	}())

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

	tt.Describe("Update admin name").Test(superadmin.Update(200, supportAdmin.Created.ID, map[string]any{
		"name": "Updated Support Admin",
	}))

	// Sync the updated data to supportAdmin
	tt.Describe("Sync data updates to supportAdmin").Test(supportAdmin.GetByID(200))

	// Test password update
	newPassword := lib.GenerateValidPassword()
	tt.Describe("Update admin password").Test(superadmin.Update(200, supportAdmin.Created.ID, map[string]any{
		"password": newPassword,
	}))

	// Test login with old password fails
	tt.Describe("Login with old password fails").Test(supportAdmin.LoginByPassword(401, supportAdmin.Created.Password))

	// Sync the updated password to supportAdmin
	supportAdmin.Created.Password = newPassword

	// Test login with new password
	tt.Describe("Login with new password").Test(supportAdmin.LoginByPassword(200, newPassword))

	// Test deactivating admin
	tt.Describe("Deactivate admin").Test(superadmin.Update(200, supportAdmin.Created.ID, map[string]any{
		"is_active": false,
	}))

	// Test that deactivated admin cannot login
	tt.Describe("Deactivated admin cannot login").Test(supportAdmin.LoginByPassword(401, newPassword))

	// Reactivate admin
	tt.Describe("Reactivate admin").Test(superadmin.Update(200, supportAdmin.Created.ID, map[string]any{
		"is_active": true,
	}))

	// Test login after reactivation
	tt.Describe("Login after reactivation").Test(supportAdmin.LoginByPassword(200, newPassword))
}
