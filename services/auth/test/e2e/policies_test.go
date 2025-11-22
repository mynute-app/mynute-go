package e2e_test

import (
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_AdminPolicies(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Authenticate as superadmin
	http, err := handler.WithSuperAdminAuth()
	if err != nil {
		t.Fatalf("Failed to authenticate as superadmin: %v", err)
	}

	var createdPolicyID string
	var endpointID string

	// Test: List admin policies (should have seeded policies)
	tt.Describe("List admin policies").Test(func() error {
		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/policies/admin").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		policies, ok := response["data"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'data' array in response, got: %v", response)
		}

		// Should have at least the seeded policies (or empty array if none)
		// This is okay - we just check the endpoint works
		_ = policies

		return nil
	}())

	// Test: Get an existing endpoint for policy creation (from seed data)
	tt.Describe("Get existing endpoint for policy creation").Test(func() error {
		// Get endpoints
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var endpoints []map[string]interface{}
		if err := http.Method("GET").URL("/endpoints").Send(nil).ExpectedStatus(200).ParseResponse(&endpoints).Error; err != nil {
			return err
		}

		// Look for /admin/users endpoint (GET method)
		for _, endpoint := range endpoints {
			method, _ := endpoint["method"].(string)
			path, _ := endpoint["path"].(string)
			if method == "GET" && path == "/admin/users" {
				id, ok := endpoint["id"].(string)
				if !ok {
					return fmt.Errorf("endpoint ID is not a string: %v", endpoint["id"])
				}
				endpointID = id
				return nil
			}
		}

		return fmt.Errorf("GET /admin/users endpoint not found in seeded data (found %d endpoints)", len(endpoints))
	}())

	// Test: Create admin policy with meaningful conditions
	tt.Describe("Create admin policy").Test(func() error {
		if endpointID == "" {
			return fmt.Errorf("no endpoint ID available from previous test")
		}

		policyData := map[string]interface{}{
			"name":         "test-admin-policy",
			"description":  "Allow admins with auditor role to view admin users",
			"effect":       "Allow",
			"end_point_id": endpointID,
			"conditions": map[string]interface{}{
				"logic_type": "AND",
				"children": []map[string]interface{}{
					{
						"leaf": map[string]interface{}{
							"attribute": "subject.roles",
							"operator":  "Contains",
							"value":     "auditor",
						},
					},
					{
						"leaf": map[string]interface{}{
							"attribute": "method",
							"operator":  "Equals",
							"value":     "GET",
						},
					},
				},
			},
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/policies/admin").
			ExpectedStatus(201).
			Send(policyData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdPolicyID = id

		name, _ := response["name"].(string)
		if name != "test-admin-policy" {
			return fmt.Errorf("expected policy name 'test-admin-policy', got: %s", name)
		}

		effect, _ := response["effect"].(string)
		if effect != "Allow" {
			return fmt.Errorf("expected effect 'Allow', got: %s", effect)
		}

		return nil
	}())

	// Test: Get admin policy by ID
	tt.Describe("Get admin policy by ID").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL(fmt.Sprintf("/policies/admin/%s", createdPolicyID)).
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok || id != createdPolicyID {
			return fmt.Errorf("expected policy ID %s, got: %v", createdPolicyID, id)
		}

		name, _ := response["name"].(string)
		if name != "test-admin-policy" {
			return fmt.Errorf("expected policy name 'test-admin-policy', got: %s", name)
		}

		return nil
	}())

	// Test: Update admin policy
	tt.Describe("Update admin policy").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		updateData := map[string]interface{}{
			"name":        "test-admin-policy-updated",
			"description": "Updated test admin policy",
			"effect":      "Deny",
		}

		var response map[string]interface{}
		if err := http.
			Method("PATCH").
			URL(fmt.Sprintf("/policies/admin/%s", createdPolicyID)).
			ExpectedStatus(200).
			Send(updateData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		name, _ := response["name"].(string)
		if name != "test-admin-policy-updated" {
			return fmt.Errorf("expected updated name 'test-admin-policy-updated', got: %s", name)
		}

		effect, _ := response["effect"].(string)
		if effect != "Deny" {
			return fmt.Errorf("expected updated effect 'Deny', got: %s", effect)
		}

		return nil
	}())

	// Test: Delete admin policy
	tt.Describe("Delete admin policy").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http.
			Method("DELETE").
			URL(fmt.Sprintf("/policies/admin/%s", createdPolicyID)).
			Send(nil).
			ExpectedStatus(200).
			Error; err != nil {
			return err
		}

		// Verify policy is deleted by trying to get it (should return 404)
		http2, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http2.
			Method("GET").
			URL(fmt.Sprintf("/policies/admin/%s", createdPolicyID)).
			Send(nil).
			ExpectedStatus(404).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		return nil
	}())
}

func Test_ClientPolicies(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	var createdPolicyID string
	var endpointID string

	// Test: Get existing endpoint for client policy creation
	tt.Describe("Get existing endpoint for client policy creation").Test(func() error {
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var endpoints []map[string]interface{}
		if err := http.Method("GET").URL("/endpoints").Send(nil).ExpectedStatus(200).ParseResponse(&endpoints).Error; err != nil {
			return err
		}

		// Look for POST /appointment endpoint
		for _, endpoint := range endpoints {
			method, _ := endpoint["method"].(string)
			path, _ := endpoint["path"].(string)
			if method == "POST" && path == "/appointment" {
				id, ok := endpoint["id"].(string)
				if !ok {
					return fmt.Errorf("endpoint ID is not a string: %v", endpoint["id"])
				}
				endpointID = id
				return nil
			}
		}

		return fmt.Errorf("POST /appointment endpoint not found in seeded data (found %d endpoints)", len(endpoints))
	}())

	// Test: Create client policy
	tt.Describe("Create client policy").Test(func() error {
		if endpointID == "" {
			return fmt.Errorf("no endpoint ID available from previous test")
		}

		policyData := map[string]interface{}{
			"name":         "test-client-policy",
			"description":  "Test client policy for E2E tests",
			"effect":       "Allow",
			"end_point_id": endpointID,
			"conditions": map[string]interface{}{
				"logic_type": "AND",
				"children": []map[string]interface{}{
					{
						"leaf": map[string]interface{}{
							"attribute":          "subject.user_id",
							"operator":           "Equals",
							"resource_attribute": "resource.client_id",
						},
					},
				},
			},
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/policies/client").
			ExpectedStatus(201).
			Send(policyData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdPolicyID = id

		name, _ := response["name"].(string)
		if name != "test-client-policy" {
			return fmt.Errorf("expected policy name 'test-client-policy', got: %s", name)
		}

		return nil
	}())

	// Test: List client policies
	tt.Describe("List client policies").Test(func() error {
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/policies/client").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		policies, ok := response["data"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'data' array in response, got: %v", response)
		}

		// Should include our created policy plus seeded ones
		found := false
		for _, p := range policies {
			policy := p.(map[string]interface{})
			if policy["id"].(string) == createdPolicyID {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("created policy not found in list")
		}

		return nil
	}())

	// Test: Delete client policy
	tt.Describe("Delete client policy").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http.
			Method("DELETE").
			URL(fmt.Sprintf("/policies/client/%s", createdPolicyID)).
			Send(nil).
			ExpectedStatus(200).
			Error; err != nil {
			return err
		}

		return nil
	}())
}

func Test_TenantPolicies(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	var createdPolicyID string
	var endpointID string

	// Test: Get existing endpoint for tenant policy creation
	tt.Describe("Get existing endpoint for tenant policy creation").Test(func() error {
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var endpoints []map[string]interface{}
		if err := http.Method("GET").URL("/endpoints").Send(nil).ExpectedStatus(200).ParseResponse(&endpoints).Error; err != nil {
			return err
		}

		// Look for GET /employee/:id endpoint
		for _, endpoint := range endpoints {
			method, _ := endpoint["method"].(string)
			path, _ := endpoint["path"].(string)
			if method == "GET" && path == "/employee/:id" {
				id, ok := endpoint["id"].(string)
				if !ok {
					return fmt.Errorf("endpoint ID is not a string: %v", endpoint["id"])
				}
				endpointID = id
				return nil
			}
		}

		return fmt.Errorf("GET /employee/:id endpoint not found in seeded data (found %d endpoints)", len(endpoints))
	}())

	// Test: Create tenant policy
	tt.Describe("Create tenant policy").Test(func() error {
		if endpointID == "" {
			return fmt.Errorf("no endpoint ID available from previous test")
		}

		policyData := map[string]interface{}{
			"name":         "test-tenant-policy",
			"description":  "Test tenant policy for E2E tests",
			"effect":       "Allow",
			"end_point_id": endpointID,
			"conditions": map[string]interface{}{
				"logic_type": "AND",
				"children": []map[string]interface{}{
					{
						"leaf": map[string]interface{}{
							"attribute":          "subject.company_id",
							"operator":           "Equals",
							"resource_attribute": "resource.company_id",
						},
					},
				},
			},
		}

		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		// Add X-Company-ID header required for tenant policies
		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/policies/tenant").
			Header("X-Company-ID", "00000000-0000-0000-0000-000000000001").
			ExpectedStatus(201).
			Send(policyData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdPolicyID = id

		name, _ := response["name"].(string)
		if name != "test-tenant-policy" {
			return fmt.Errorf("expected policy name 'test-tenant-policy', got: %s", name)
		}

		return nil
	}())

	// Test: Get tenant policy by ID
	tt.Describe("Get tenant policy by ID").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL(fmt.Sprintf("/policies/tenant/%s", createdPolicyID)).
			Header("X-Company-ID", "00000000-0000-0000-0000-000000000001").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok || id != createdPolicyID {
			return fmt.Errorf("expected policy ID %s, got: %v", createdPolicyID, id)
		}

		return nil
	}())

	// Test: Update tenant policy
	tt.Describe("Update tenant policy").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		updateData := map[string]interface{}{
			"description": "Updated tenant policy description",
		}

		var response map[string]interface{}
		if err := http.
			Method("PATCH").
			URL(fmt.Sprintf("/policies/tenant/%s", createdPolicyID)).
			Header("X-Company-ID", "00000000-0000-0000-0000-000000000001").
			ExpectedStatus(200).
			Send(updateData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		description, _ := response["description"].(string)
		if description != "Updated tenant policy description" {
			return fmt.Errorf("expected updated description, got: %s", description)
		}

		return nil
	}())

	// Test: Delete tenant policy
	tt.Describe("Delete tenant policy").Test(func() error {
		if createdPolicyID == "" {
			return fmt.Errorf("no policy ID available from previous test")
		}

		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http.
			Method("DELETE").
			URL(fmt.Sprintf("/policies/tenant/%s", createdPolicyID)).
			Header("X-Company-ID", "00000000-0000-0000-0000-000000000001").
			Send(nil).
			ExpectedStatus(200).
			Error; err != nil {
			return err
		}

		return nil
	}())
}
