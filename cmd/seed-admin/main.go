package main

import (
	"log"
	database "mynute-go/services/core/src/config/db"
	"mynute-go/services/core/src/config/db/model"
	"mynute-go/services/core/src/lib"
	"os"

	"github.com/google/uuid"
)

// Admin seeder - creates default admin user and roles
// Usage: go run cmd/seed-admin/main.go
func main() {
	log.Println("Starting admin seeding process...")

	// Load environment variables
	lib.LoadEnv()

	app_env := os.Getenv("APP_ENV")
	log.Printf("Environment: %s\n", app_env)

	// Connect to database
	db := database.Connect()
	defer db.Disconnect()

	// Ensure we're working in public schema
	if err := lib.ChangeToPublicSchema(db.Gorm); err != nil {
		log.Fatalf("Failed to switch to public schema: %v", err)
	}

	// Create admin roles
	adminRoles := []model.RoleAdmin{
		{
			ID:          uuid.New(),
			Name:        "superadmin",
			Description: "Full access to all tenants and resources",
		},
		{
			ID:          uuid.New(),
			Name:        "support",
			Description: "Customer support with read-only access to tenant data",
		},
		{
			ID:          uuid.New(),
			Name:        "auditor",
			Description: "Read-only access for compliance and auditing",
		},
	}

	for _, role := range adminRoles {
		var existingRole model.RoleAdmin
		if err := db.Gorm.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			if err := db.Gorm.Create(&role).Error; err != nil {
				log.Fatalf("Failed to create role '%s': %v", role.Name, err)
			}
			log.Printf("✓ Created admin role: %s", role.Name)
		} else {
			log.Printf("- Admin role '%s' already exists", role.Name)
		}
	}

	// Get default admin credentials from environment or use defaults
	defaultAdminEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	if defaultAdminEmail == "" {
		defaultAdminEmail = "admin@mynute.com"
	}

	defaultAdminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if defaultAdminPassword == "" {
		defaultAdminPassword = "Admin@123456"
		log.Printf("⚠️  WARNING: Using default password. Set DEFAULT_ADMIN_PASSWORD in .env for production!")
	}

	defaultAdminName := os.Getenv("DEFAULT_ADMIN_NAME")
	if defaultAdminName == "" {
		defaultAdminName = "System Administrator"
	}

	// Check if admin already exists
	// NOTE: This seed file needs to be updated for the new User/Admin architecture
	// where User records are created in the auth service first, then Admin records
	// reference them via UserID. For now, this is disabled.

	log.Println("⚠️  WARNING: Admin seeding temporarily disabled - needs update for new User/Admin architecture")
	log.Println("Please use the admin registration endpoint to create the first admin user.")

	log.Println("\n✓ Admin role seeding completed successfully!")
	log.Println("\nSeeded:")
	log.Println("  - 3 Admin Roles (superadmin, support, auditor)")
}
