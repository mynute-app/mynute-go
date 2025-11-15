package main

import (
	"fmt"
	"log"
	"mynute-go/services/core/api/lib"
	database "mynute-go/services/core/config/db"
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
	tenantPolicySeed "mynute-go/services/core/config/db/seed/policy"
	resourceSeed "mynute-go/services/core/config/db/seed/resource"
	"os"
)

// Seed command for production environments
// Usage: go run cmd/seed/main.go
// Or build: go build -o bin/seed cmd/seed/main.go && ./bin/seed
func main() {
	log.Println("Starting seeding process...")

	// Load environment variables
	lib.LoadEnv()

	app_env := os.Getenv("APP_ENV")
	log.Printf("Environment: %s\n", app_env)

	// Connect to database
	db := database.Connect()
	defer db.Disconnect()

	// Start transaction for seeding
	tx, end, err := database.Transaction(db.Gorm)
	defer end(nil)
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	// Seed system roles (if not already seeded by migrations)
	model.AllowSystemRoleCreation = true
	defer func() {
		model.AllowSystemRoleCreation = false
	}()

	if err := db.
		Seed("Resources", resourceSeed.Resources, `"table" = ?`, []string{"Table"}).
		Seed("Roles", model.Roles, "name = ? AND company_id IS NULL", []string{"Name"}).
		Error; err != nil {
		log.Fatalf("Failed to seed roles/resources: %v", err)
	}

	// Load system role IDs
	if err := model.LoadSystemRoleIDs(tx); err != nil {
		log.Fatalf("Failed to load system role IDs: %v", err)
	}

	// Seed endpoints
	endpoints := endpointSeed.GetAllEndpoints()

	if err := db.
		Seed("Endpoints", endpoints, "method = ? AND path = ?", []string{"Method", "Path"}).
		Error; err != nil {
		log.Fatalf("Failed to seed endpoints: %v", err)
	}

	// Load endpoint IDs from database for reference
	for i, endpoint := range endpoints {
		var dbEndpoint model.EndPoint
		if err := tx.Where("method = ? AND path = ?", endpoint.Method, endpoint.Path).First(&dbEndpoint).Error; err != nil {
			log.Printf("⚠️  Warning: Failed to load endpoint ID for '%s %s': %v\n", endpoint.Method, endpoint.Path, err)
			continue
		}
		endpoints[i] = &dbEndpoint
	}

	// Seed tenant policies
	tenantPolicies := tenantPolicySeed.GetAllTenantPolicies()
	if err := db.Seed("TenantPolicies", tenantPolicies, "name = ?", []string{"Name"}).Error; err != nil {
		log.Fatalf("Failed to seed tenant policies: %v", err)
	}

	// Seed client policies
	clientPolicies := tenantPolicySeed.GetAllClientPolicies()
	if err := db.Seed("ClientPolicies", clientPolicies, "name = ?", []string{"Name"}).Error; err != nil {
		log.Fatalf("Failed to seed client policies: %v", err)
	}

	log.Println("✓ Seeding completed successfully!")
	fmt.Println("\nSeeded:")
	fmt.Println("  - Resources")
	fmt.Println("  - System Roles")
	fmt.Printf("  - %d Endpoints\n", len(endpoints))
	fmt.Printf("  - %d Tenant Policies\n", len(tenantPolicies))
	fmt.Printf("  - %d Client Policies\n", len(clientPolicies))
}
