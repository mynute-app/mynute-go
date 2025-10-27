package main

import (
	"flag"
	"fmt"
	"log"
	"mynute-go/core/src/lib"
	"os"
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
		if len(flag.Args()) == 0 {
			log.Fatal("Please provide a migration name: go run migrate/main.go -action=create <migration_name>")
		}
		migrationName := flag.Args()[0]
		if err := createMigration(absPath, migrationName); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}

	default:
		log.Fatalf("Unknown action: %s", action)
	}
}

func createMigration(migrationsPath, name string) error {
	// Generate timestamp-based version
	timestamp := lib.GetTimestampVersion()
	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	// Create up migration file
	upContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Add your UP migration here\n", name, timestamp)
	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Add your DOWN migration here\n", name, timestamp)
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	log.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
	return nil
}
