package e2e_test

import (
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_AdminRoles(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Authenticate as superadmin
	http, err := handler.WithSuperAdminAuth()
	if err != nil {
		t.Fatalf("Failed to authenticate as superadmin: %v", err)
	}

	var createdRoleID string

	// Test: List admin roles
	tt.Describe("List admin roles").Test(func() error {
		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/roles/admin").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		roles, ok := response["data"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'data' array in response, got: %v", response)
		}

		// Should have at least the seeded roles (superadmin, support)
		if len(roles) < 2 {
			return fmt.Errorf("expected at least 2 roles, got %d", len(roles))
		}

		return nil
	}())

	// Test: Create admin role
	tt.Describe("Create admin role").Test(func() error {
		roleData := map[string]interface{}{
			"name":        "auditor",
			"description": "Audit administrator role",
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/roles/admin").
			ExpectedStatus(201).
			Send(roleData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdRoleID = id

		name, _ := response["name"].(string)
		if name != "auditor" {
			return fmt.Errorf("expected role name 'auditor', got: %s", name)
		}

		return nil
	}())

	// Test: Get admin role by ID
	tt.Describe("Get admin role by ID").Test(func() error {
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/roles/admin/" + createdRoleID).
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		name, _ := response["name"].(string)
		if name != "auditor" {
			return fmt.Errorf("expected role name 'auditor', got: %s", name)
		}

		return nil
	}())

	// Test: Update admin role
	tt.Describe("Update admin role").Test(func() error {
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		updateData := map[string]interface{}{
			"description": "Updated audit administrator role",
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("PATCH").
			URL("/roles/admin/" + createdRoleID).
			ExpectedStatus(200).
			Send(updateData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		description, _ := response["description"].(string)
		if description != "Updated audit administrator role" {
			return fmt.Errorf("expected updated description, got: %s", description)
		}

		return nil
	}())

	// Test: Delete admin role
	tt.Describe("Delete admin role").Test(func() error {
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http.
			Method("DELETE").
			URL("/roles/admin/" + createdRoleID).
			Send(nil).
			ExpectedStatus(204).
			Error; err != nil {
			return err
		}

		// Verify it's deleted - create new client for verification
		http2, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http2.
			Method("GET").
			URL("/roles/admin/" + createdRoleID).
			Send(nil).
			ExpectedStatus(404).
			Error; err != nil {
			return err
		}

		return nil
	}())

	// Test: Cannot create duplicate role
	tt.Describe("Cannot create duplicate admin role").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		roleData := map[string]interface{}{
			"name":        "superadmin",
			"description": "Duplicate superadmin",
		}

		if err := http.
			Method("POST").
			URL("/roles/admin").
			ExpectedStatus(500). // Server returns 500 for database constraint violations
			Send(roleData).
			Error; err != nil {
			return err
		}

		return nil
	}())
}

func Test_TenantRoles(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	companyID := "550e8400-e29b-41d4-a716-446655440001"
	var createdRoleID string

	// Test: Create tenant role
	tt.Describe("Create tenant role").Test(func() error {
		http := handler.NewHttpClient()
		roleData := map[string]interface{}{
			"tenant_id":   companyID,
			"name":        "manager",
			"description": "Branch manager role",
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/roles/tenant").
			Header("X-Company-ID", companyID).
			ExpectedStatus(201).
			Send(roleData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdRoleID = id

		name, _ := response["name"].(string)
		if name != "manager" {
			return fmt.Errorf("expected role name 'manager', got: %s", name)
		}

		return nil
	}())

	// Test: List tenant roles
	tt.Describe("List tenant roles").Test(func() error {
		http := handler.NewHttpClient()
		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/roles/tenant").
			Header("X-Company-ID", companyID).
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		data, ok := response["data"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'data' array in response, got: %v", response)
		}

		if len(data) == 0 {
			return fmt.Errorf("expected at least 1 role, got %d", len(data))
		}

		return nil
	}())

	// Test: Get tenant role by ID
	tt.Describe("Get tenant role by ID").Test(func() error {
		http := handler.NewHttpClient()
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/roles/tenant/"+createdRoleID).
			Header("X-Company-ID", companyID).
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		name, _ := response["name"].(string)
		if name != "manager" {
			return fmt.Errorf("expected role name 'manager', got: %s", name)
		}

		return nil
	}())

	// Test: Update tenant role
	tt.Describe("Update tenant role").Test(func() error {
		http := handler.NewHttpClient()
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		updateData := map[string]interface{}{
			"description": "Senior branch manager role",
		}

		var response map[string]interface{}
		if err := http.
			Method("PATCH").
			URL("/roles/tenant/"+createdRoleID).
			Header("X-Company-ID", companyID).
			ExpectedStatus(200).
			Send(updateData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		description, _ := response["description"].(string)
		if description != "Senior branch manager role" {
			return fmt.Errorf("expected updated description, got: %s", description)
		}

		return nil
	}())

	// Test: Delete tenant role
	tt.Describe("Delete tenant role").Test(func() error {
		http := handler.NewHttpClient()
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		if err := http.
			Method("DELETE").
			URL("/roles/tenant/"+createdRoleID).
			Header("X-Company-ID", companyID).
			Send(nil).
			ExpectedStatus(204).
			Error; err != nil {
			return err
		}

		return nil
	}())
}

func Test_ClientRoles(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	var createdRoleID string

	// Test: Create client role
	tt.Describe("Create client role").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		roleData := map[string]interface{}{
			"name":        "premium_client",
			"description": "Premium client with extra features",
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/roles/client").
			ExpectedStatus(201).
			Send(roleData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdRoleID = id

		name, _ := response["name"].(string)
		if name != "premium_client" {
			return fmt.Errorf("expected role name 'premium_client', got: %s", name)
		}

		return nil
	}())

	// Test: List client roles
	tt.Describe("List client roles").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/roles/client").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		roles, ok := response["data"].([]interface{})
		if !ok {
			return fmt.Errorf("expected 'data' array in response, got: %v", response)
		}

		if len(roles) == 0 {
			return fmt.Errorf("expected at least 1 role, got %d", len(roles))
		}

		return nil
	}())

	// Test: Get client role by ID
	tt.Describe("Get client role by ID").Test(func() error {
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/roles/client/" + createdRoleID).
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		name, _ := response["name"].(string)
		if name != "premium_client" {
			return fmt.Errorf("expected role name 'premium_client', got: %s", name)
		}

		return nil
	}())

	// Test: Update client role
	tt.Describe("Update client role").Test(func() error {
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		updateData := map[string]interface{}{
			"description": "VIP premium client with all features",
		}

		var response map[string]interface{}
		if err := http.
			Method("PATCH").
			URL("/roles/client/" + createdRoleID).
			ExpectedStatus(200).
			Send(updateData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		description, _ := response["description"].(string)
		if description != "VIP premium client with all features" {
			return fmt.Errorf("expected updated description, got: %s", description)
		}

		return nil
	}())

	// Test: Delete client role
	tt.Describe("Delete client role").Test(func() error {
		if createdRoleID == "" {
			return fmt.Errorf("no role ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http.
			Method("DELETE").
			URL("/roles/client/" + createdRoleID).
			Send(nil).
			ExpectedStatus(204).
			Error; err != nil {
			return err
		}

		return nil
	}())
}
