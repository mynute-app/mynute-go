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

	// Use test company ID (in real scenario, this would come from the business service)
	companyID := "550e8400-e29b-41d4-a716-446655440001"
	companyID2 := "550e8400-e29b-41d4-a716-446655440004"

	// Create test tenant users
	var employee1Token, employee2Token string
	tt.Describe("Create test tenant users").Test(func() error {
		// Create employee 1
		createReq1 := map[string]interface{}{
			"email":     "employee1@test.com",
			"password":  "employee1pass",
			"tenant_id": companyID,
		}

		http1 := handler.NewHttpClient()
		var createResp1 map[string]interface{}
		if err := http1.
			Method("POST").
			URL("/users/tenant").
			ExpectedStatus(201).
			Send(createReq1).
			ParseResponse(&createResp1).
			Error; err != nil {
			return fmt.Errorf("failed to create employee 1: %w", err)
		}

		// Login as employee 1
		loginReq1 := map[string]interface{}{
			"email":    "employee1@test.com",
			"password": "employee1pass",
		}

		loginHttp1 := handler.NewHttpClient()
		loginHttp1.Header("X-Company-ID", companyID)
		if err := loginHttp1.
			Method("POST").
			URL("/auth/tenant/login").
			ExpectedStatus(200).
			Send(loginReq1).
			ParseResponse(&map[string]interface{}{}).
			Error; err != nil {
			return fmt.Errorf("failed to login as employee 1: %w", err)
		}

		tokens1 := loginHttp1.ResHeaders["X-Auth-Token"]
		if len(tokens1) == 0 {
			return fmt.Errorf("no token returned for employee 1")
		}
		employee1Token = tokens1[0]

		// Create employee 2
		createReq2 := map[string]interface{}{
			"email":     "employee2@test.com",
			"password":  "employee2pass",
			"tenant_id": companyID,
		}

		http2 := handler.NewHttpClient()
		var createResp2 map[string]interface{}
		if err := http2.
			Method("POST").
			URL("/users/tenant").
			ExpectedStatus(201).
			Send(createReq2).
			ParseResponse(&createResp2).
			Error; err != nil {
			return fmt.Errorf("failed to create employee 2: %w", err)
		}

		// Login as employee 2
		loginReq2 := map[string]interface{}{
			"email":    "employee2@test.com",
			"password": "employee2pass",
		}

		loginHttp2 := handler.NewHttpClient()
		loginHttp2.Header("X-Company-ID", companyID)
		if err := loginHttp2.
			Method("POST").
			URL("/auth/tenant/login").
			ExpectedStatus(200).
			Send(loginReq2).
			ParseResponse(&map[string]interface{}{}).
			Error; err != nil {
			return fmt.Errorf("failed to login as employee 2: %w", err)
		}

		tokens2 := loginHttp2.ResHeaders["X-Auth-Token"]
		if len(tokens2) == 0 {
			return fmt.Errorf("no token returned for employee 2")
		}
		employee2Token = tokens2[0]

		return nil
	}())

	// Get employee IDs from tokens
	var employeeID1, employeeID2 string
	tt.Describe("Extract employee IDs from tokens").Test(func() error {
		// Use validation endpoint to get employee info
		http1 := handler.NewHttpClient()
		http1.Header("X-Auth-Token", employee1Token)

		var resp1 map[string]interface{}
		if err := http1.
			Method("POST").
			URL("/auth/validate").
			ExpectedStatus(200).
			ParseResponse(&resp1).
			Error; err != nil {
			return fmt.Errorf("failed to validate employee 1 token: %w", err)
		}

		id1, ok := resp1["id"].(string)
		if !ok {
			return fmt.Errorf("employee 1 ID not found in token")
		}
		employeeID1 = id1

		http2 := handler.NewHttpClient()
		http2.Header("X-Auth-Token", employee2Token)

		var resp2 map[string]interface{}
		if err := http2.
			Method("POST").
			URL("/auth/validate").
			ExpectedStatus(200).
			ParseResponse(&resp2).
			Error; err != nil {
			return fmt.Errorf("failed to validate employee 2 token: %w", err)
		}

		id2, ok := resp2["id"].(string)
		if !ok {
			return fmt.Errorf("employee 2 ID not found in token")
		}
		employeeID2 = id2

		return nil
	}())

	// Test tenant authorization - employee can access their own profile
	tt.Describe("Tenant can access own employee profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/employee/:id",
			// subject now extracted from token
			"resource": map[string]interface{}{
				"employee_id": employeeID1,
			},
			"path_params": map[string]interface{}{
				"id": employeeID1,
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID, employee1Token)
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
			// subject now extracted from token
			"resource": map[string]interface{}{
				"employee_id": employeeID2,
			},
			"path_params": map[string]interface{}{
				"id": employeeID2,
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID, employee1Token)
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
			// subject now extracted from token
			"query": map[string]interface{}{
				"company_id": companyID,
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID, employee1Token)
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
	tt.Describe("Employee from company1 cannot access company2 resources").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/service",
			// subject now extracted from token
			"query": map[string]interface{}{
				"company_id": companyID2, // Different company!
			},
		}

		response, err := makeAuthorizationRequest("tenant", authReq, companyID, employee1Token)
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
