package main

import (
	"fmt"
	"log"
	database "mynute-go/core/src/config/db"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
			BaseModel:   model.BaseModel{ID: uuid.New()},
			Name:        "superadmin",
			Description: "Full access to all tenants and resources",
		},
		{
			BaseModel:   model.BaseModel{ID: uuid.New()},
			Name:        "support",
			Description: "Customer support with read-only access to tenant data",
		},
		{
			BaseModel:   model.BaseModel{ID: uuid.New()},
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
	var existingAdmin model.Admin
	if err := db.Gorm.Where("email = ?", defaultAdminEmail).First(&existingAdmin).Error; err != nil {
		// Admin doesn't exist, create it

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		admin := model.Admin{
			BaseModel: model.BaseModel{ID: uuid.New()},
			Name:      defaultAdminName,
			Email:     defaultAdminEmail,
			Password:  string(hashedPassword),
			IsActive:  true,
		}

		// Create admin
		if err := db.Gorm.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin: %v", err)
		}

		// Assign superadmin role
		var superadminRole model.RoleAdmin
		if err := db.Gorm.Where("name = ?", "superadmin").First(&superadminRole).Error; err != nil {
			log.Fatalf("Failed to find superadmin role: %v", err)
		}

		if err := db.Gorm.Model(&admin).Association("Roles").Append(&superadminRole); err != nil {
			log.Fatalf("Failed to assign superadmin role: %v", err)
		}

		log.Println("✓ Created default admin user")
		separator := strings.Repeat("=", 60)
		fmt.Println("\n" + separator)
		fmt.Println("DEFAULT ADMIN CREDENTIALS")
		fmt.Println(separator)
		fmt.Printf("Email:    %s\n", defaultAdminEmail)
		fmt.Printf("Password: %s\n", defaultAdminPassword)
		fmt.Println(separator)
		fmt.Println("⚠️  IMPORTANT: Change the password immediately after first login!")
		fmt.Println(separator + "\n")
	} else {
		log.Printf("- Admin user '%s' already exists", defaultAdminEmail)
	}

	log.Println("\n✓ Admin seeding completed successfully!")
	fmt.Println("\nSeeded:")
	fmt.Println("  - 3 Admin Roles (superadmin, support, auditor)")
	fmt.Println("  - 1 Default Admin User")
}
