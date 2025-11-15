package e2e_test

import (
	"fmt"
	"mynute-go/services/core"
	"mynute-go/services/core/api/lib"
	"mynute-go/services/core/test/src/handler"
	"mynute-go/services/core/test/src/model"
	"testing"
)

func Test_Tenant(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Setup: Create admin and company
	superadmin := &model.Admin{}
	tt.Describe("Create superadmin").Test(superadmin.Set([]string{}, nil))

	company := &model.Company{}
	tt.Describe("Create company").Test(company.Create(200))

	// Create first employee (tenant user)
	employee1 := &model.Employee{Company: company}
	tt.Describe("Create first tenant employee").Test(employee1.Create(200, &superadmin.X_Auth_Token, nil))

	// Test tenant user verification
	tt.Describe("Send verification email to tenant").Test(employee1.SendVerificationEmail(200, nil))

	tt.Describe("Verify tenant email with code").Test(func() error {
		code, err := employee1.GetVerificationCodeFromEmail()
		if err != nil {
			return fmt.Errorf("failed to get verification code: %w", err)
		}
		return employee1.VerifyEmailByCode(200, code, nil)
	}())

	// Test tenant login with password
	tt.Describe("Login tenant with password").Test(employee1.LoginByPassword(200, employee1.Password, nil))

	// Test invalid password login
	tt.Describe("Login tenant with wrong password fails").Test(employee1.LoginByPassword(401, "wrong-password", nil))

	// Test login with email code
	tt.Describe("Send login code to tenant email").Test(employee1.SendLoginCode(200, nil))

	tt.Describe("Login tenant with email code").Test(func() error {
		code, err := employee1.GetLoginCodeFromEmail()
		if err != nil {
			return fmt.Errorf("failed to get login code: %w", err)
		}
		return employee1.LoginByEmailCode(200, code, nil)
	}())

	// Create second employee to test authorization
	employee2 := &model.Employee{Company: company}
	tt.Describe("Create second tenant employee").Test(employee2.Create(200, &superadmin.X_Auth_Token, nil))

	tt.Describe("Verify second employee email").Test(func() error {
		if err := employee2.SendVerificationEmail(200, nil); err != nil {
			return err
		}
		code, err := employee2.GetVerificationCodeFromEmail()
		if err != nil {
			return err
		}
		return employee2.VerifyEmailByCode(200, code, nil)
	}())

	tt.Describe("Login second employee").Test(employee2.LoginByPassword(200, employee2.Password, nil))

	// Test tenant authorization - employee can access their own profile
	tt.Describe("Tenant can access own employee profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/employee/:id",
			"subject": map[string]interface{}{
				"user_id":    employee1.Created.UserID.String(),
				"company_id": company.Created.ID.String(),
				"role":       "employee",
			},
			"resource": map[string]interface{}{
				"employee_id": employee1.Created.UserID.String(),
			},
			"path_params": map[string]interface{}{
				"id": employee1.Created.UserID.String(),
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, company.Created.ID.String())
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed, got response: %v", response)
		}

		return nil
	}())

	// Test tenant authorization - employee cannot access another employee's private data without permission
	tt.Describe("Tenant cannot access other employee without permission").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/employee/:id",
			"subject": map[string]interface{}{
				"user_id":    employee1.Created.UserID.String(),
				"company_id": company.Created.ID.String(),
				"role":       "employee",
			},
			"resource": map[string]interface{}{
				"employee_id": employee2.Created.UserID.String(),
			},
			"path_params": map[string]interface{}{
				"id": employee2.Created.UserID.String(),
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, company.Created.ID.String())
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		// Should be denied (employees can only access their own data by default)
		if allowed {
			return fmt.Errorf("expected authorization to be denied for accessing other employee")
		}

		return nil
	}())

	// Test tenant authorization - company-wide resource access
	tt.Describe("Tenant can list company services").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/service",
			"subject": map[string]interface{}{
				"user_id":    employee1.Created.UserID.String(),
				"company_id": company.Created.ID.String(),
				"role":       "employee",
			},
			"query": map[string]interface{}{
				"company_id": company.Created.ID.String(),
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, company.Created.ID.String())
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for company services, got: %v", response)
		}

		return nil
	}())

	// Test password reset
	tt.Describe("Reset tenant password by email").Test(func() error {
		if err := employee1.SendPasswordResetEmail(200, nil); err != nil {
			return err
		}
		newPassword, err := employee1.GetNewPasswordFromEmail()
		if err != nil {
			return err
		}
		employee1.Password = newPassword
		return employee1.LoginByPassword(200, newPassword, nil)
	}())

	// Test tenant user deactivation
	tt.Describe("Deactivate tenant user").Test(func() error {
		changes := map[string]any{
			"is_active": false,
		}
		return employee1.Update(200, changes, &superadmin.X_Auth_Token, nil)
	}())

	// Test deactivated tenant cannot login
	tt.Describe("Deactivated tenant cannot login").Test(employee1.LoginByPassword(401, employee1.Password, nil))

	// Reactivate tenant
	tt.Describe("Reactivate tenant user").Test(func() error {
		changes := map[string]any{
			"is_active": true,
		}
		return employee1.Update(200, changes, &superadmin.X_Auth_Token, nil)
	}())

	// Test login after reactivation
	tt.Describe("Login after reactivation").Test(employee1.LoginByPassword(200, employee1.Password, nil))

	// Test tenant update
	tt.Describe("Update tenant user details").Test(func() error {
		newName := lib.GenerateRandomName("Updated Employee")
		changes := map[string]any{
			"name": newName,
		}
		if err := employee1.Update(200, changes, &employee1.X_Auth_Token, nil); err != nil {
			return err
		}
		if employee1.Created.Name != newName {
			return fmt.Errorf("expected name to be updated to %s, got %s", newName, employee1.Created.Name)
		}
		return nil
	}())

	// Test cross-company isolation - employee from one company cannot access another company
	company2 := &model.Company{}
	tt.Describe("Create second company").Test(company2.Create(200))

	employee3 := &model.Employee{Company: company2}
	tt.Describe("Create employee in second company").Test(employee3.Create(200, &superadmin.X_Auth_Token, nil))

	tt.Describe("Employee from company1 cannot access company2 resources").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/service",
			"subject": map[string]interface{}{
				"user_id":    employee1.Created.UserID.String(),
				"company_id": company.Created.ID.String(),
				"role":       "employee",
			},
			"query": map[string]interface{}{
				"company_id": company2.Created.ID.String(), // Different company!
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, company.Created.ID.String())
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		if allowed {
			return fmt.Errorf("expected cross-company access to be denied")
		}

		return nil
	}())
}

// Helper function to make authorization requests
func makeAuthorizationRequest(userType string, authReq map[string]interface{}, companyID string) (map[string]interface{}, error) {
	http := handler.NewHttpClient()

	endpoint := ""
	switch userType {
	case "tenant":
		endpoint = "/auth/tenant/authorize"
		http.Header("X-Company-ID", companyID)
	case "client":
		endpoint = "/auth/client/authorize"
	case "admin":
		endpoint = "/auth/admin/authorize"
	default:
		return nil, fmt.Errorf("invalid user type: %s", userType)
	}

	var response map[string]interface{}
	if err := http.
		Method("POST").
		URL(endpoint).
		ExpectedStatus(200).
		Send(authReq).
		ParseResponse(&response).
		Error; err != nil {
		return nil, fmt.Errorf("failed to make authorization request: %w", err)
	}

	return response, nil
}
