package e2e_test

import (
	"encoding/json"
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_Endpoints(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	var createdEndpointID string

	// Test: List endpoints
	tt.Describe("List endpoints").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		// Get raw response body and parse manually since response is an array
		var rawBody []byte
		if err := http.
			Method("GET").
			URL("/endpoints").
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&rawBody).
			Error; err != nil {
			return err
		}

		// Parse the array response
		var endpoints []interface{}
		if err := json.Unmarshal(rawBody, &endpoints); err != nil {
			return fmt.Errorf("failed to parse endpoints array: %w", err)
		}

		// Should have at least the seeded endpoints
		if len(endpoints) < 5 {
			return fmt.Errorf("expected at least 5 endpoints, got %d", len(endpoints))
		}

		return nil
	}())

	// Test: Create endpoint
	tt.Describe("Create endpoint").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		endpointData := map[string]interface{}{
			"method":            "GET",
			"path":              "/test/endpoint",
			"description":       "Test endpoint for e2e testing",
			"controller_name":   "TestController",
			"deny_unauthorized": true,
		}

		var response map[string]interface{}
		if err := http.
			Method("POST").
			URL("/endpoints").
			ExpectedStatus(201).
			Send(endpointData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		id, ok := response["id"].(string)
		if !ok {
			return fmt.Errorf("expected 'id' field in response, got: %v", response)
		}

		createdEndpointID = id

		path, _ := response["path"].(string)
		if path != "/test/endpoint" {
			return fmt.Errorf("expected path '/test/endpoint', got: %s", path)
		}

		method, _ := response["method"].(string)
		if method != "GET" {
			return fmt.Errorf("expected method 'GET', got: %s", method)
		}

		return nil
	}())

	// Test: Get endpoint by ID
	tt.Describe("Get endpoint by ID").Test(func() error {
		if createdEndpointID == "" {
			return fmt.Errorf("no endpoint ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		var response map[string]interface{}
		if err := http.
			Method("GET").
			URL("/endpoints/" + createdEndpointID).
			Send(nil).
			ExpectedStatus(200).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		path, _ := response["path"].(string)
		if path != "/test/endpoint" {
			return fmt.Errorf("expected path '/test/endpoint', got: %s", path)
		}

		return nil
	}())

	// Test: Update endpoint
	tt.Describe("Update endpoint").Test(func() error {
		if createdEndpointID == "" {
			return fmt.Errorf("no endpoint ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		updateData := map[string]interface{}{
			"description": "Updated test endpoint description",
		}

		var response map[string]interface{}
		if err := http.
			Method("PATCH").
			URL("/endpoints/" + createdEndpointID).
			ExpectedStatus(200).
			Send(updateData).
			ParseResponse(&response).
			Error; err != nil {
			return err
		}

		description, _ := response["description"].(string)
		if description != "Updated test endpoint description" {
			return fmt.Errorf("expected updated description, got: %s", description)
		}

		return nil
	}())

	// Test: Delete endpoint
	tt.Describe("Delete endpoint").Test(func() error {
		if createdEndpointID == "" {
			return fmt.Errorf("no endpoint ID available from previous test")
		}

		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := http.
			Method("DELETE").
			URL("/endpoints/" + createdEndpointID).
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
			URL("/endpoints/" + createdEndpointID).
			Send(nil).
			ExpectedStatus(404).
			Error; err != nil {
			return err
		}

		return nil
	}())

	// Test: Cannot create duplicate endpoint (same method + path)
	tt.Describe("Cannot create duplicate endpoint").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		endpointData := map[string]interface{}{
			"method":      "GET",
			"path":        "/admin/users",
			"description": "Duplicate endpoint",
		}

		if err := http.
			Method("POST").
			URL("/endpoints").
			ExpectedStatus(500). // Server returns 500 for database constraint violations
			Send(endpointData).
			Error; err != nil {
			return err
		}

		return nil
	}())

	// Test: Validate HTTP method
	tt.Describe("Reject invalid HTTP method").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		endpointData := map[string]interface{}{
			"method":      "INVALID",
			"path":        "/test/invalid",
			"description": "Invalid method endpoint",
		}

		if err := http.
			Method("POST").
			URL("/endpoints").
			ExpectedStatus(500). // Server returns 500 for validation errors
			Send(endpointData).
			Error; err != nil {
			return err
		}

		return nil
	}())

	// Test: Create endpoint with different methods for same path
	tt.Describe("Create endpoints with different methods for same path").Test(func() error {
		// Create new authenticated client for this test
		http, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		// Create POST endpoint
		postEndpointData := map[string]interface{}{
			"method":      "POST",
			"path":        "/test/multi-method",
			"description": "POST test endpoint",
		}

		var postResponse map[string]interface{}
		if err := http.
			Method("POST").
			URL("/endpoints").
			ExpectedStatus(201).
			Send(postEndpointData).
			ParseResponse(&postResponse).
			Error; err != nil {
			return err
		}

		// Create PUT endpoint with same path (should succeed) - new client
		http2, err := handler.WithSuperAdminAuth()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		putEndpointData := map[string]interface{}{
			"method":      "PUT",
			"path":        "/test/multi-method",
			"description": "PUT test endpoint",
		}

		var putResponse map[string]interface{}
		if err := http2.
			Method("POST").
			URL("/endpoints").
			ExpectedStatus(201).
			Send(putEndpointData).
			ParseResponse(&putResponse).
			Error; err != nil {
			return err
		}

		// Cleanup
		postID, _ := postResponse["id"].(string)
		putID, _ := putResponse["id"].(string)

		if postID != "" {
			httpCleanup1, _ := handler.WithSuperAdminAuth()
			_ = httpCleanup1.Method("DELETE").URL("/endpoints/" + postID).Send(nil).ExpectedStatus(204).Error
		}
		if putID != "" {
			httpCleanup2, _ := handler.WithSuperAdminAuth()
			_ = httpCleanup2.Method("DELETE").URL("/endpoints/" + putID).Send(nil).ExpectedStatus(204).Error
		}

		return nil
	}())
}
