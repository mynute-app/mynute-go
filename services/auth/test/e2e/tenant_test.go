package e2e_test

import (
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_Tenant(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Use test UUIDs for authorization checks
	companyID := "550e8400-e29b-41d4-a716-446655440001"
	employeeID1 := "550e8400-e29b-41d4-a716-446655440002"
	employeeID2 := "550e8400-e29b-41d4-a716-446655440003"

	// Test tenant authorization - employee can access their own profile
	tt.Describe("Tenant can access own employee profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/employee/:id",
			"subject": map[string]interface{}{
				"user_id":    employeeID1,
				"company_id": companyID,
				"role":       "employee",
			},
			"resource": map[string]interface{}{
				"employee_id": employeeID1,
			},
			"path_params": map[string]interface{}{
				"id": employeeID1,
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID)
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
				"user_id":    employeeID1,
				"company_id": companyID,
				"role":       "employee",
			},
			"resource": map[string]interface{}{
				"employee_id": employeeID2,
			},
			"path_params": map[string]interface{}{
				"id": employeeID2,
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID)
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
				"user_id":    employeeID1,
				"company_id": companyID,
				"role":       "employee",
			},
			"query": map[string]interface{}{
				"company_id": companyID,
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID)
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for company services, got: %v", response)
		}

		return nil
	}())

	// Test cross-company isolation - employee from one company cannot access another company
	companyID2 := "550e8400-e29b-41d4-a716-446655440004"

	tt.Describe("Employee from company1 cannot access company2 resources").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/service",
			"subject": map[string]interface{}{
				"user_id":    employeeID1,
				"company_id": companyID,
				"role":       "employee",
			},
			"query": map[string]interface{}{
				"company_id": companyID2, // Different company!
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID)
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
