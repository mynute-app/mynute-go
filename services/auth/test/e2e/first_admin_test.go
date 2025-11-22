package e2e_test

import (
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_FirstSuperAdmin(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Test: Check no superadmin exists initially
	tt.Describe("Check no superadmin exists initially").Test(func() error {
		http := handler.NewHttpClient()

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/users/admin/are_there_any_superadmin").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		hasSuperAdmin, ok := response["has_superadmin"].(bool)
		if !ok {
			return fmt.Errorf("expected 'has_superadmin' field in response")
		}

		if hasSuperAdmin {
			return fmt.Errorf("expected no superadmin to exist, but found one")
		}

		return nil
	}())

	// Test: Create first superadmin successfully
	var firstAdminEmail string
	tt.Describe("Create first superadmin successfully").Test(func() error {
		http := handler.NewHttpClient()

		firstAdminEmail = "first-admin@mynute.test"
		createData := map[string]interface{}{
			"name":     "First",
			"surname":  "Admin",
			"email":    firstAdminEmail,
			"password": "FirstAdmin123!",
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/users/admin/first_superadmin").
			Send(createData).
			ExpectedStatus(201).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		// Verify email matches
		email, ok := response["email"].(string)
		if !ok || email != firstAdminEmail {
			return fmt.Errorf("expected email %s, got %v", firstAdminEmail, email)
		}

		// Verify roles include superadmin
		roles, ok := response["roles"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'roles' array in response")
		}

		hasSuperAdminRole := false
		for _, role := range roles {
			roleMap, ok := role.(map[string]interface{})
			if !ok {
				continue
			}
			if roleName, ok := roleMap["name"].(string); ok && roleName == "superadmin" {
				hasSuperAdminRole = true
				break
			}
		}

		if !hasSuperAdminRole {
			return fmt.Errorf("expected first admin to have superadmin role")
		}

		return nil
	}())

	// Test: Verify superadmin now exists
	tt.Describe("Verify superadmin now exists").Test(func() error {
		http := handler.NewHttpClient()

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/users/admin/are_there_any_superadmin").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		hasSuperAdmin, ok := response["has_superadmin"].(bool)
		if !ok {
			return fmt.Errorf("expected 'has_superadmin' field in response")
		}

		if !hasSuperAdmin {
			return fmt.Errorf("expected superadmin to exist after creation")
		}

		return nil
	}())

	// Test: Attempt to create another first superadmin should fail
	tt.Describe("Attempt to create another first superadmin should fail").Test(func() error {
		http := handler.NewHttpClient()

		createData := map[string]interface{}{
			"name":     "Second",
			"surname":  "Admin",
			"email":    "second-admin@mynute.test",
			"password": "SecondAdmin123!",
		}

		var response map[string]interface{}
		// This should fail with 403 (Forbidden) because admin already exists
		if err := http.
			Method("POST").
			URL("/users/admin/first_superadmin").
			Send(createData).
			ExpectedStatus(403).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		// Should have an error message
		if response["error"] == nil {
			return fmt.Errorf("expected error message when trying to create first admin again")
		}

		return nil
	}())

	// Test: Login as the first superadmin
	var token string
	tt.Describe("Login as the first superadmin").Test(func() error {
		http := handler.NewHttpClient()

		loginData := map[string]interface{}{
			"email":    firstAdminEmail,
			"password": "FirstAdmin123!",
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/auth/admin/login").
			Send(loginData).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		// Get token from response header
		tokenHeader := http.ResHeaders["X-Auth-Token"]
		if len(tokenHeader) == 0 {
			return fmt.Errorf("no token returned for first admin")
		}

		token = tokenHeader[0]
		if token == "" {
			return fmt.Errorf("token is empty")
		}

		return nil
	}())

	// Test: First superadmin can create other admins
	tt.Describe("First superadmin can create other admins").Test(func() error {
		http := handler.NewHttpClient()
		http.Header("X-Auth-Token", token)

		createData := map[string]interface{}{
			"name":     "Support",
			"surname":  "Admin",
			"email":    "support-admin@mynute.test",
			"password": "SupportAdmin123!",
			"roles":    []string{"support"},
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/users/admin").
			Send(createData).
			ExpectedStatus(201).
			ParseResponse(&response).
			Error; err != nil {
			return fmt.Errorf("first superadmin should be able to create other admins: %w", err)
		}

		// Verify the new admin has support role
		roles, ok := response["roles"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'roles' array in response")
		}

		hasSupportRole := false
		for _, role := range roles {
			roleMap, ok := role.(map[string]interface{})
			if !ok {
				continue
			}
			if roleName, ok := roleMap["name"].(string); ok && roleName == "support" {
				hasSupportRole = true
				break
			}
		}

		if !hasSupportRole {
			return fmt.Errorf("expected new admin to have support role")
		}

		return nil
	}())
}
