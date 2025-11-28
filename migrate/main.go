package main

// Run all pending migrations
// go run migrate/main.go -action=up

// Check current migration version
// go run migrate/main.go -action=version

// Rollback last migration
// go run migrate/main.go -action=down -steps=1

// Create a new migration file
// go run migrate/main.go -action=create your_migration_name

import (
	"flag"
	"fmt"
	"log"
	"mynute-go/core/src/lib"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	lib.LoadEnv()

	var (
		action         string
		steps          int
		version        int
		migrationsPath string
	)

	flag.StringVar(&action, "action", "", "Migration action: up, down, version, force, create")
	flag.IntVar(&steps, "steps", 1, "Number of steps to migrate down")
	flag.IntVar(&version, "version", -1, "Version to force migration to")
	flag.StringVar(&migrationsPath, "path", "./migrations", "Path to migrations directory")
	flag.Parse()

	if action == "" {
		log.Fatal("Please specify an action: up, down, version, force, create")
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	switch action {
	case "up":
		log.Println("Running migrations...")
		if err := lib.RunMigrations(absPath); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations completed successfully!")

	case "down":
		log.Printf("Rolling back %d migration(s)...\n", steps)
		if err := lib.RollbackMigration(absPath, steps); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}

	case "version":
		version, dirty, err := lib.MigrationVersion(absPath)
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current version: %d (dirty: %t)\n", version, dirty)

	case "force":
		if version < 0 {
			log.Fatal("Please specify a version using -version flag")
		}
		log.Printf("Forcing migration version to %d...\n", version)
		if err := lib.ForceMigrationVersion(absPath, version); err != nil {
			log.Fatalf("Force failed: %v", err)
		}

	case "create":
		migrationName := "migration"
		if len(flag.Args()) > 0 {
			migrationName = flag.Args()[0]
		}
		if err := runSmartMigration(migrationName); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}

	default:
		log.Fatalf("Unknown action: %s", action)
	}
}

func runSmartMigration(name string) error {
	log.Println("Running smart-migration tool to generate migration...")

	// Run the smart-migration tool
	cmd := exec.Command("go", "run", "tools/smart-migration/main.go", "-name="+name, "-models=all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run smart-migration tool: %w", err)
	}

	log.Println("Migration files created successfully!")
	return nil
}