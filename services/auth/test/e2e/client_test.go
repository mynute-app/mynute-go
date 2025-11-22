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

	// Create test client users
	var client1Token, client2Token string
	tt.Describe("Create test client users").Test(func() error {
		// Create client 1
		createReq1 := map[string]interface{}{
			"email":    "client1@test.com",
			"password": "client1pass",
		}

		http1 := handler.NewHttpClient()
		var createResp1 map[string]interface{}
		if err := http1.
			Method("POST").
			URL("/users/client").
			ExpectedStatus(201).
			Send(createReq1).
			ParseResponse(&createResp1).
			Error; err != nil {
			return fmt.Errorf("failed to create client 1: %w", err)
		}

		// Login as client 1
		loginReq1 := map[string]interface{}{
			"email":    "client1@test.com",
			"password": "client1pass",
		}

		loginHttp1 := handler.NewHttpClient()
		if err := loginHttp1.
			Method("POST").
			URL("/auth/client/login").
			ExpectedStatus(200).
			Send(loginReq1).
			ParseResponse(&map[string]interface{}{}).
			Error; err != nil {
			return fmt.Errorf("failed to login as client 1: %w", err)
		}

		tokens1 := loginHttp1.ResHeaders["X-Auth-Token"]
		if len(tokens1) == 0 {
			return fmt.Errorf("no token returned for client 1")
		}
		client1Token = tokens1[0]

		// Create client 2
		createReq2 := map[string]interface{}{
			"email":    "client2@test.com",
			"password": "client2pass",
		}

		http2 := handler.NewHttpClient()
		var createResp2 map[string]interface{}
		if err := http2.
			Method("POST").
			URL("/users/client").
			ExpectedStatus(201).
			Send(createReq2).
			ParseResponse(&createResp2).
			Error; err != nil {
			return fmt.Errorf("failed to create client 2: %w", err)
		}

		// Login as client 2
		loginReq2 := map[string]interface{}{
			"email":    "client2@test.com",
			"password": "client2pass",
		}

		loginHttp2 := handler.NewHttpClient()
		if err := loginHttp2.
			Method("POST").
			URL("/auth/client/login").
			ExpectedStatus(200).
			Send(loginReq2).
			ParseResponse(&map[string]interface{}{}).
			Error; err != nil {
			return fmt.Errorf("failed to login as client 2: %w", err)
		}

		tokens2 := loginHttp2.ResHeaders["X-Auth-Token"]
		if len(tokens2) == 0 {
			return fmt.Errorf("no token returned for client 2")
		}
		client2Token = tokens2[0]

		return nil
	}())

	// Get client IDs from tokens (extract from JWT)
	var clientID1, clientID2 string
	tt.Describe("Extract client IDs from tokens").Test(func() error {
		// Use validation endpoint to get client info
		http1 := handler.NewHttpClient()
		http1.Header("X-Auth-Token", client1Token)

		var resp1 map[string]interface{}
		if err := http1.
			Method("POST").
			URL("/auth/validate").
			ExpectedStatus(200).
			ParseResponse(&resp1).
			Error; err != nil {
			return fmt.Errorf("failed to validate client 1 token: %w", err)
		}

		id1, ok := resp1["id"].(string)
		if !ok {
			return fmt.Errorf("client 1 ID not found in token")
		}
		clientID1 = id1

		http2 := handler.NewHttpClient()
		http2.Header("X-Auth-Token", client2Token)

		var resp2 map[string]interface{}
		if err := http2.
			Method("POST").
			URL("/auth/validate").
			ExpectedStatus(200).
			ParseResponse(&resp2).
			Error; err != nil {
			return fmt.Errorf("failed to validate client 2 token: %w", err)
		}

		id2, ok := resp2["id"].(string)
		if !ok {
			return fmt.Errorf("client 2 ID not found in token")
		}
		clientID2 = id2

		return nil
	}())

	// Test client authorization - can access own profile
	tt.Describe("Client can access own profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/client/:id",
			// subject now extracted from token
			"resource": map[string]interface{}{
				"client_id": clientID1,
			},
			"path_params": map[string]interface{}{
				"id": clientID1,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
			// subject now extracted from token
			"resource": map[string]interface{}{
				"client_id": clientID2,
			},
			"path_params": map[string]interface{}{
				"id": clientID2,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
			// subject now extracted from token
			"body": map[string]interface{}{
				"client_id": clientID1,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
			// subject now extracted from token
			"body": map[string]interface{}{
				"client_id": clientID2, // Different client!
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
			// subject now extracted from token
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
			// subject now extracted from token
			"query": map[string]interface{}{
				"client_id": clientID1,
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
			// subject now extracted from token
			"query": map[string]interface{}{
				"client_id": clientID2, // Different client!
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "", client1Token)
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
