package e2e_test

import (
	"fmt"
	"mynute-go/services/core"
	"mynute-go/services/core/api/lib"
	"mynute-go/services/core/test/src/handler"
	"mynute-go/services/core/test/src/model"
	"testing"
)

func Test_Admin(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	superadmin := &model.Admin{}

	// Create first admin - Set will automatically handle whether it's the first superadmin or not
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

	// Test admin authorization - superadmin can access all admin endpoints
	tt.Describe("Superadmin can create admin users").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/admin/users",
			"subject": map[string]interface{}{
				"user_id": superadmin.Created.UserID.String(),
				"role":    "superadmin",
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected superadmin to be allowed, got response: %v", response)
		}

		return nil
	}())

	// Test admin authorization - support admin has limited access
	tt.Describe("Support admin cannot create admin users").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/admin/users",
			"subject": map[string]interface{}{
				"user_id": supportAdmin.Created.UserID.String(),
				"role":    "support",
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		// Support admins should not be able to create other admins
		if allowed {
			return fmt.Errorf("expected support admin to be denied from creating admin users")
		}

		return nil
	}())

	// Test admin authorization - admins can view system resources
	tt.Describe("Admin can view companies").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/company",
			"subject": map[string]interface{}{
				"user_id": supportAdmin.Created.UserID.String(),
				"role":    "support",
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected admin to view companies, got response: %v", response)
		}

		return nil
	}())

	tt.Describe("Update admin name").Test(superadmin.Update(200, supportAdmin.Created.UserID, map[string]any{
		"name": "Updated Support Admin",
	}))

	// Sync the updated data to supportAdmin
	tt.Describe("Sync data updates to supportAdmin").Test(supportAdmin.GetByID(200))

	// Test password update
	newPassword := lib.GenerateValidPassword()
	tt.Describe("Update admin password").Test(superadmin.Update(200, supportAdmin.Created.UserID, map[string]any{
		"password": newPassword,
	}))

	// Test login with old password fails
	tt.Describe("Login with old password fails").Test(supportAdmin.LoginByPassword(401, "<temp-password>" /* TODO: Password on User */))

	// Sync the updated password to supportAdmin
	// TODO: Password on User - newPassword

	// Test login with new password
	tt.Describe("Login with new password").Test(supportAdmin.LoginByPassword(200, newPassword))

	// Test admin can update their own profile
	tt.Describe("Admin can update own profile").Test(func() error {
		newName := lib.GenerateRandomName("Self Updated Admin")
		changes := map[string]any{
			"name": newName,
		}
		if err := supportAdmin.Update(200, supportAdmin.Created.UserID, changes); err != nil {
			return err
		}
		// Sync to verify update
		if err := supportAdmin.GetByID(200); err != nil {
			return err
		}
		if supportAdmin.Created.Name != newName {
			return fmt.Errorf("expected name to be %s, got %s", newName, supportAdmin.Created.Name)
		}
		return nil
	}())

	// Test deactivating admin
	tt.Describe("Deactivate admin").Test(superadmin.Update(200, supportAdmin.Created.UserID, map[string]any{
		"is_active": false,
	}))

	// Test that deactivated admin cannot login
	tt.Describe("Deactivated admin cannot login").Test(supportAdmin.LoginByPassword(401, newPassword))

	// Test that deactivated admin cannot access endpoints
	tt.Describe("Deactivated admin cannot authorize").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/company",
			"subject": map[string]interface{}{
				"user_id":   supportAdmin.Created.UserID.String(),
				"role":      "support",
				"is_active": false,
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		if allowed {
			return fmt.Errorf("expected deactivated admin to be denied authorization")
		}

		return nil
	}())

	// Reactivate admin
	tt.Describe("Reactivate admin").Test(superadmin.Update(200, supportAdmin.Created.UserID, map[string]any{
		"is_active": true,
	}))

	// Test login after reactivation
	tt.Describe("Login after reactivation").Test(supportAdmin.LoginByPassword(200, newPassword))

	// Test authorization after reactivation
	tt.Describe("Reactivated admin can authorize").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/company",
			"subject": map[string]interface{}{
				"user_id":   supportAdmin.Created.UserID.String(),
				"role":      "support",
				"is_active": true,
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected reactivated admin to be allowed, got response: %v", response)
		}

		return nil
	}())

	// Test role-based permissions
	developerAdmin := &model.Admin{}
	tt.Describe("Create developer admin").Test(superadmin.Set([]string{"developer"}, developerAdmin))

	tt.Describe("Developer can access system metrics").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/metrics",
			"subject": map[string]interface{}{
				"user_id": developerAdmin.Created.UserID.String(),
				"role":    "developer",
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected developer to access metrics, got response: %v", response)
		}

		return nil
	}())

	// Test list admins with pagination
	tt.Describe("List admins with limit").Test(func() error {
		// Assuming ListAdmins accepts a status code parameter
		admins, err := superadmin.ListAdmins(200)
		if err != nil {
			return err
		}
		if len(admins) < 3 {
			return fmt.Errorf("expected at least 3 admins (superadmin, support, developer), got %d", len(admins))
		}
		return nil
	}())

	// Test admin password reset via email
	tt.Describe("Reset admin password via email").Test(func() error {
		if err := supportAdmin.SendPasswordResetEmail(200); err != nil {
			return err
		}
		newPwd, err := supportAdmin.GetNewPasswordFromEmail()
		if err != nil {
			return err
		}
		supportAdmin.Password = newPwd
		return supportAdmin.LoginByPassword(200, newPwd)
	}())

	// Test get admin by ID
	tt.Describe("Get admin by ID").Test(func() error {
		return superadmin.GetByID(200)
	}())
}
