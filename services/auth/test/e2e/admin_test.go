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

	// Test basic authorization endpoint
	tt.Describe("Admin authorization endpoint exists").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/admin/users",
			"subject": map[string]interface{}{
				"user_id": "550e8400-e29b-41d4-a716-446655440000",
				"role":    "superadmin",
			},
		}

		response, err := makeAuthorizationRequest("admin", authReq, "")
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
			"subject": map[string]interface{}{
				"user_id": "550e8400-e29b-41d4-a716-446655440000",
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
				"user_id": "550e8400-e29b-41d4-a716-446655440001",
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
				"user_id": "550e8400-e29b-41d4-a716-446655440001",
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
}

// Helper function to make authorization requests
func makeAuthorizationRequest(userType string, authReq map[string]interface{}, companyID string) (map[string]interface{}, error) {
	http := handler.NewHttpClient()

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
