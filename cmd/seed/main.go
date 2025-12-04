package main

import (
	"fmt"
	"log"
	database "mynute-go/core/src/config/db"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"os"
)

// Seed command for production environments
// Usage: go run cmd/seed/main.go
// Or build: go build -o bin/seed cmd/seed/main.go && ./bin/seed
//
// IMPORTANT: This command uses POSTGRES_DB_PROD environment variable
// to determine which database to seed (same as migrations).
// Set POSTGRES_DB_PROD=maindb for production, or POSTGRES_DB_PROD=devdb for dev.
func main() {
	log.Println("Starting seeding process...")

	// Load environment variables
	lib.LoadEnv()

	// Verify POSTGRES_DB_PROD is set
	dbName := os.Getenv("POSTGRES_DB_PROD")
	if dbName == "" {
		log.Fatal("Error: POSTGRES_DB_PROD environment variable is required for seeding")
	}

	log.Printf("Target database: %s\n", dbName)
	log.Println("⚠️  WARNING: Seeding will modify the database specified by POSTGRES_DB_PROD")
	log.Println("")

	// Connect to database using production configuration
	db := database.ConnectForTools()
	defer db.Disconnect()

	// Start transaction for seeding
	tx, end, err := database.Transaction(db.Gorm)
	defer end(nil)
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	// Seed system roles (if not already seeded by migrations)
	model.AllowSystemRoleCreation = true
	defer func() { model.AllowSystemRoleCreation = false }()

	if err := db.
		Seed("Resources", model.Resources, `"table" = ?`, []string{"Table"}).
		Seed("Roles", model.Roles, "name = ? AND company_id IS NULL", []string{"Name"}).
		Error; err != nil {
		log.Fatalf("Failed to seed roles/resources: %v", err)
	}

	// Load system role IDs
	if err := model.LoadSystemRoleIDs(tx); err != nil {
		log.Fatalf("Failed to load system role IDs: %v", err)
	}

	// Seed endpoints
	endpoints, deferEndpoint, err := model.EndPoints(&model.EndpointCfg{AllowCreation: true}, tx)
	if err != nil {
		log.Fatalf("Failed to prepare endpoints: %v", err)
	}
	defer deferEndpoint()

	if err := db.
		Seed("Endpoints", endpoints, "method = ? AND path = ?", []string{"Method", "Path"}).
		Error; err != nil {
		log.Fatalf("Failed to seed endpoints: %v", err)
	}

	// Load endpoint IDs from database so policies can reference them
	if err := model.LoadEndpointIDs(tx); err != nil {
		log.Fatalf("Failed to load endpoint IDs: %v", err)
	}

	// Seed policies
	policies, deferPolicies := model.Policies(&model.PolicyCfg{AllowNilCompanyID: true, AllowNilCreatedBy: true})
	defer deferPolicies()

	if err := db.
		Seed("Policies", policies, "name = ?", []string{"Name"}).
		Error; err != nil {
		log.Fatalf("Failed to seed policies: %v", err)
	}

	log.Println("✓ Seeding completed successfully!")
	fmt.Println("\nSeeded:")
	fmt.Println("  - Resources")
	fmt.Println("  - System Roles")
	fmt.Printf("  - %d Endpoints\n", len(endpoints))
	fmt.Printf("  - %d Policies\n", len(policies))
}
