package e2e_test

import (
	"fmt"
	"mynute-go/services/auth"
	"mynute-go/services/auth/test/src/handler"
	"testing"
)

func Test_Client(t *testing.T) {
	server := auth.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Use test UUIDs for authorization checks
	clientID1 := "650e8400-e29b-41d4-a716-446655440001"
	clientID2 := "650e8400-e29b-41d4-a716-446655440002"

	// Test client authorization - can access own profile
	tt.Describe("Client can access own profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/client/:id",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
			"resource": map[string]interface{}{
				"client_id": clientID1,
			},
			"path_params": map[string]interface{}{
				"id": clientID1,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed, got response: %v", response)
		}

		return nil
	}())

	// Test client authorization - cannot access another client's profile
	tt.Describe("Client cannot access other client profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/client/:id",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
			"resource": map[string]interface{}{
				"client_id": clientID2,
			},
			"path_params": map[string]interface{}{
				"id": clientID2,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		// Should be denied (clients can only access their own data)
		if allowed {
			return fmt.Errorf("expected authorization to be denied for accessing other client")
		}

		return nil
	}())

	// Test client authorization - can create own appointments
	tt.Describe("Client can create own appointments").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
			"body": map[string]interface{}{
				"client_id": clientID1,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for creating own appointment, got: %v", response)
		}

		return nil
	}())

	// Test client authorization - cannot create appointments for another client
	tt.Describe("Client cannot create appointments for other clients").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
			"body": map[string]interface{}{
				"client_id": clientID2, // Different client!
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		if allowed {
			return fmt.Errorf("expected authorization to be denied for creating appointment for other client")
		}

		return nil
	}())

	// Test client authorization - list available services (public endpoint)
	tt.Describe("Client can list available services").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/service",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for listing services, got: %v", response)
		}

		return nil
	}())

	// Test client can view their own appointments
	tt.Describe("Client can view own appointments").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
			"query": map[string]interface{}{
				"client_id": clientID1,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for viewing own appointments, got: %v", response)
		}

		return nil
	}())

	// Test client cannot view another client's appointments
	tt.Describe("Client cannot view other client appointments").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": clientID1,
				"role":    "client",
			},
			"query": map[string]interface{}{
				"client_id": clientID2, // Different client!
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		if allowed {
			return fmt.Errorf("expected authorization to be denied for viewing other client's appointments")
		}

		return nil
	}())
}
