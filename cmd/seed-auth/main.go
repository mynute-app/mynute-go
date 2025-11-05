package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	authModel "mynute-go/services/auth/config/db/model"
	database "mynute-go/services/core/src/config/db"
	coreModel "mynute-go/services/core/src/config/db/model"
	endpointSeed "mynute-go/services/core/src/config/db/seed/endpoint"
	resourceSeed "mynute-go/services/core/src/config/db/seed/resource"
	"mynute-go/services/core/src/lib"
)

// SeedAuth sends all endpoints, resources, and policies to the auth service
// Usage: go run cmd/seed-auth/main.go
func main() {
	log.Println("Starting auth service seeding...")

	// Load environment variables
	lib.LoadEnv()

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:4001"
	}

	log.Printf("Auth service URL: %s\n", authServiceURL)

	// Connect to database to load endpoint IDs if needed
	db := database.Connect()
	defer db.Disconnect()

	tx, end, err := database.Transaction(db.Gorm)
	defer end(nil)
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Step 1: Seed Resources
	log.Println("\n=== Seeding Resources ===")
	resources := resourceSeed.Resources
	for _, resource := range resources {
		if err := seedResource(client, authServiceURL, resource); err != nil {
			log.Printf("⚠️  Warning: Failed to seed resource '%s': %v\n", resource.Name, err)
		} else {
			log.Printf("✓ Seeded resource: %s\n", resource.Name)
		}
	}

	// Step 2: Seed Endpoints
	log.Println("\n=== Seeding Endpoints ===")

	// Enable endpoint creation temporarily
	coreModel.AllowSystemRoleCreation = true
	authModel.AllowEndpointCreation = true
	defer func() {
		coreModel.AllowSystemRoleCreation = false
		authModel.AllowEndpointCreation = false
	}()

	endpoints := endpointSeed.GetAllEndpoints()
	endpoints, deferEndpoints, err := authModel.EndPoints(endpoints, &authModel.EndpointCfg{AllowCreation: true}, tx)
	if err != nil {
		log.Fatalf("Failed to prepare endpoints: %v", err)
	}
	defer deferEndpoints()

	successCount := 0
	for _, endpoint := range endpoints {
		if err := seedEndpoint(client, authServiceURL, endpoint); err != nil {
			log.Printf("⚠️  Warning: Failed to seed endpoint '%s %s': %v\n", endpoint.Method, endpoint.Path, err)
		} else {
			log.Printf("✓ Seeded endpoint: %s %s\n", endpoint.Method, endpoint.Path)
			successCount++
		}
	}

	log.Printf("\n✓ Seeding completed! %d/%d endpoints seeded successfully\n", successCount, len(endpoints))
	fmt.Println("\nNote: Policies should be seeded separately after reviewing the policy definitions")
	fmt.Println("Run: go run cmd/seed-auth-policies/main.go")
}

// seedResource sends a resource to the auth service
func seedResource(client *http.Client, baseURL string, resource *authModel.Resource) error {
	// For now, we'll skip resources since they might need to be created via a dedicated endpoint
	// Or we can add them directly to the auth database
	// This is a placeholder for future implementation
	return nil
}

// seedEndpoint sends an endpoint to the auth service
func seedEndpoint(client *http.Client, baseURL string, endpoint *authModel.EndPoint) error {
	// Prepare endpoint data
	data := map[string]interface{}{
		"method":      endpoint.Method,
		"path":        endpoint.Path,
		"description": endpoint.Description,
	}

	// Add optional fields if present
	if endpoint.ControllerName != "" {
		data["controller_name"] = endpoint.ControllerName
	}
	if endpoint.DenyUnauthorized {
		data["deny_unauthorized"] = endpoint.DenyUnauthorized
	}
	if endpoint.Resource != nil {
		data["resource_id"] = endpoint.Resource.ID
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoint: %w", err)
	}

	// Send POST request to create endpoint
	req, err := http.NewRequest("POST", baseURL+"/endpoints", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if endpoint already exists (409) or was created successfully (200/201)
	if resp.StatusCode == http.StatusConflict {
		// Endpoint already exists, this is fine
		return nil
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf("unexpected status %d: %v", resp.StatusCode, errorResp)
	}

	return nil
}

