package e2e_test

import (
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_Admin(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Authenticate as superadmin to get token
	var token string
	tt.Describe("Login as superadmin").Test(func() error {
		var err error
		token, err = handler.AuthenticateAsSuperAdmin()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		if token == "" {
			return fmt.Errorf("token is empty")
		}
		return nil
	}())

	// Test basic authorization endpoint
	tt.Describe("Admin authorization endpoint exists").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/admin/users",
			// subject is now extracted from JWT token
		}

		response, err := makeAuthorizationRequest("admin", authReq, "", token)
		if err != nil {
			return err
		}

		_, ok := response["allowed"]
		if !ok {
			return fmt.Errorf("expected 'allowed' field in response, got: %v", response)
		}

		return nil
	}())

	// Test admin authorization - superadmin can access all admin endpoints
	tt.Describe("Superadmin can create admin users").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/admin/users",
			// subject is now extracted from JWT token
		}

		response, err := makeAuthorizationRequest("admin", authReq, "", token)
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected superadmin to be allowed, got response: %v", response)
		}

		return nil
	}())

	// Test admin authorization - admins can view system resources
	tt.Describe("Admin can view companies").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/company",
			// subject is now extracted from JWT token
		}

		response, err := makeAuthorizationRequest("admin", authReq, "", token)
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected admin to view companies, got response: %v", response)
		}

		return nil
	}())

	// Test creating a support admin (with limited permissions)
	tt.Describe("Create support admin user").Test(func() error {
		http := handler.NewHttpClient()
		http.Header("X-Auth-Token", token)

		createReq := map[string]interface{}{
			"name":     "Support Admin",
			"surname":  "User",
			"email":    "support@mynute.test",
			"password": "support123456",
			"roles":    []string{"support"},
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/users/admin").
			ExpectedStatus(201).
			Send(createReq).
			ParseResponse(&response).
			Error; err != nil {
			return fmt.Errorf("failed to create support admin: %w", err)
		}

		return nil
	}())

	// Test support admin login and authorization
	tt.Describe("Support admin cannot create admin users").Test(func() error {
		// Login as support admin
		loginReq := map[string]interface{}{
			"email":    "support@mynute.test",
			"password": "support123456",
		}

		http := handler.NewHttpClient()
		var loginResponse map[string]interface{}
		if err := http.
			Method("POST").
			URL("/auth/admin/login").
			ExpectedStatus(200).
			Send(loginReq).
			ParseResponse(&loginResponse).
			Error; err != nil {
			return fmt.Errorf("failed to login as support admin: %w", err)
		}

		// Get token from response header
		supportToken := http.ResHeaders["X-Auth-Token"]
		if len(supportToken) == 0 {
			return fmt.Errorf("no token returned for support admin")
		}

		// Try to create admin user (should be denied)
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/admin/users",
			// subject is now extracted from JWT token
		}

		response, err := makeAuthorizationRequest("admin", authReq, "", supportToken[0])
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
}

// Helper function to make authorization requests
func makeAuthorizationRequest(userType string, authReq map[string]interface{}, companyID string, token string) (map[string]interface{}, error) {
	http := handler.NewHttpClient()

	// Set the authentication token
	if token != "" {
		http.Header("X-Auth-Token", token)
	}

	endpoint := ""
	switch userType {
	case "tenant":
		endpoint = "/authorize/tenant"
		http.Header("X-Company-ID", companyID)
	case "client":
		endpoint = "/authorize/client"
	case "admin":
		endpoint = "/authorize/admin"
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
