package handler

import (
	"fmt"
)

// AuthenticateAsSuperAdmin logs in as the test superadmin and returns the JWT token
// This token can be used in the X-Auth-Token header for authenticated requests
func AuthenticateAsSuperAdmin() (string, error) {
	http := NewHttpClient()

	// First, check if any superadmin exists
	var checkResponse map[string]interface{}
	if err := http.
		Method("GET").
		URL("/users/admin/are_there_any_superadmin").
		Send(nil).
		ExpectedStatus(200).
		ParseResponse(&checkResponse).
		Error; err != nil {
		return "", fmt.Errorf("failed to check for superadmin: %w", err)
	}

	hasSuperAdmin, _ := checkResponse["has_superadmin"].(bool)

	// If no superadmin exists, create the first one
	if !hasSuperAdmin {
		createData := map[string]interface{}{
			"name":     "Test",
			"surname":  "Superadmin",
			"email":    "test-superadmin@mynute.local",
			"password": "test123456",
		}

		var createResponse map[string]interface{}
		if err := http.
			Method("POST").
			URL("/users/admin/first_superadmin").
			ExpectedStatus(201).
			Send(createData).
			ParseResponse(&createResponse).
			Error; err != nil {
			return "", fmt.Errorf("failed to create first superadmin: %w", err)
		}
	}

	// Login credentials for test superadmin
	loginData := map[string]interface{}{
		"email":    "test-superadmin@mynute.local",
		"password": "test123456",
	}

	var response map[string]interface{}
	if err := http.
		Method("POST").
		URL("/auth/admin/login").
		ExpectedStatus(200).
		Send(loginData).
		ParseResponse(&response).
		Error; err != nil {
		return "", fmt.Errorf("failed to authenticate as superadmin: %w", err)
	}

	// Token is returned in X-Auth-Token header
	token := http.ResHeaders["X-Auth-Token"]
	if len(token) == 0 {
		return "", fmt.Errorf("no token returned in X-Auth-Token header")
	}

	return token[0], nil
}

// WithSuperAdminAuth creates a new HTTP client with superadmin authentication
func WithSuperAdminAuth() (*httpActions, error) {
	token, err := AuthenticateAsSuperAdmin()
	if err != nil {
		return nil, err
	}

	http := NewHttpClient()
	http.Header("X-Auth-Token", token)
	return http, nil
}
