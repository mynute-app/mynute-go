package main

// Atlas-based migration tool
// Run all pending migrations: go run migrate/main.go -action=up
// Generate new migration: go run migrate/main.go -action=diff -name=your_migration_name
// Check current status: go run migrate/main.go -action=status

import (
	"flag"
	"log"
	"mynute-go/core/src/lib"
	"os"
	"os/exec"
)

func main() {
	lib.LoadEnv()

	var (
		action string
		name   string
		env    string
	)

	flag.StringVar(&action, "action", "", "Migration action: up, diff, status, down")
	flag.StringVar(&name, "name", "", "Migration name (required for diff)")
	flag.StringVar(&env, "env", "dev", "Environment: dev or prod")
	flag.Parse()

	if action == "" {
		log.Fatal("Please specify an action: up, diff, status, down")
	}

	switch action {
	case "up":
		log.Println("Applying migrations...")
		if err := runAtlasCommand(env, "migrate", "apply", "--url", getDatabaseURL()); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("✅ Migrations applied successfully!")

	case "diff":
		if name == "" {
			log.Fatal("Please specify a migration name using -name flag")
		}
		log.Printf("Generating migration: %s...\n", name)
		if err := runAtlasCommand(env, "migrate", "diff", name, "--env", env); err != nil {
			log.Fatalf("Failed to generate migration: %v", err)
		}
		log.Println("✅ Migration generated successfully!")

	case "status":
		log.Println("Checking migration status...")
		if err := runAtlasCommand(env, "migrate", "status", "--url", getDatabaseURL()); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}

	case "down":
		log.Println("Rolling back last migration...")
		if err := runAtlasCommand(env, "migrate", "down", "--url", getDatabaseURL()); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("✅ Rollback completed successfully!")

	default:
		log.Fatalf("Unknown action: %s. Use: up, diff, status, down", action)
	}
}

func runAtlasCommand(env string, args ...string) error {
	cmd := exec.Command("atlas", args...)
	cmd.Env = append(os.Environ(),
		"ATLAS_ENV="+env,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getDatabaseURL() string {
	config := lib.GetMigrationConfig()
	return config.GetDatabaseURL()
}
