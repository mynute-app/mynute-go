# Makefile for database migrations
# Requires: golang-migrate CLI tool installed

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Migration commands (Atlas-based)
.PHONY: migrate-up migrate-down migrate-diff migrate-status migrate-help

# Seeding commands
.PHONY: seed seed-help

# Run all pending migrations
migrate-up:
	@echo "Applying database migrations..."
	@go run migrate/main.go -action=up -env=dev

# Run migrations in production
migrate-up-prod:
	@echo "Applying database migrations (PRODUCTION)..."
	@go run migrate/main.go -action=up -env=prod

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	@go run migrate/main.go -action=down -env=dev

# Generate a new migration by comparing models with database
migrate-diff:
ifndef NAME
	@echo "Error: NAME is required. Usage: make migrate-diff NAME=your_migration_name"
	@exit 1
endif
	@echo "Generating migration: $(NAME)"
	@go run migrate/main.go -action=diff -name=$(NAME) -env=dev

# Check migration status
migrate-status:
	@echo "Checking migration status..."
	@go run migrate/main.go -action=status -env=dev
migrate-generate:
ifndef NAME
	@echo "Error: NAME is required. Usage: make migrate-generate NAME=migration_name [MODELS=Employee,Branch]"
	@exit 1
endif
# Show help
migrate-help:
	@echo "Database Migration Commands:"
	@echo ""
	@echo "Basic Commands:"
	@echo "  make migrate-up              - Run all pending migrations"
	@echo "  make migrate-down            - Rollback the last migration"
	@echo "  make migrate-down-n STEPS=N  - Rollback N migrations"
	@echo "  make migrate-version         - Show current migration version"
	@echo ""
	@echo "Creating Migrations:"
	@echo "  make migrate-create NAME=x              - Create empty migration files"
	@echo "  make migrate-generate NAME=x MODELS=Y   - Generate template (no detection)"
	@echo "  make migrate-smart NAME=x MODELS=Y      - Auto-detect changes (smart!)"
	@echo "  make migrate-smart NAME=x MODELS=all    - Check ALL models for changes"
	@echo ""
	@echo "Testing:"
	@echo "  make test-migrate            - Auto-test migration (up->down->up)"
	@echo "  make test-migrate-interactive - Interactive migration test"
	@echo ""
	@echo "Advanced:"
	@echo "  make migrate-force VERSION=N - Force migration version (dangerous!)"
	@echo "  make db-reset                - Reset database (down all + up all)"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-smart NAME=add_bio MODELS=Employee"
	@echo "  make migrate-smart NAME=check_all MODELS=all"
	@echo "  make migrate-generate NAME=add_fields MODELS=Branch,Service"
	@echo "  make test-migrate"

# Development helpers
.PHONY: db-reset db-fresh

# Reset database (down all + up all)
db-reset:
	@echo "Resetting database..."
	@go run migrate/main.go -action=down -steps=100 -path=./migrations || true
	@go run migrate/main.go -action=up -path=./migrations

# Fresh database (for development only)
db-fresh:
	@echo "⚠️  WARNING: This will DROP and recreate the database!"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	@$(MAKE) db-reset

# Test helpers
.PHONY: test-migrate

# Automated migration testing (up -> down -> up)
test-migrate:
	@echo "Running automated migration tests..."
	@pwsh -File scripts/test-migration.ps1 -SkipConfirmation

# Interactive migration testing
test-migrate-interactive:
	@echo "Running interactive migration tests..."
	@pwsh -File scripts/test-migration.ps1

# ============================================
# SEEDING COMMANDS
# ============================================
# NOTE: Seeding uses POSTGRES_DB_PROD (same as migrations)
# Set POSTGRES_DB_PROD=maindb for production, or POSTGRES_DB_PROD=devdb for dev

# Run seeding (endpoints, policies, roles, resources)
seed:
	@echo "Running database seeding..."
	@go run cmd/seed/main.go

# Build seed binary for production
seed-build:
	@echo "Building seed binary..."
	@go build -o bin/seed cmd/seed/main.go
	@echo "✓ Binary created at: bin/seed"

# Show seeding help
seed-help:
	@echo "Database Seeding Commands:"
	@echo ""
	@echo "Commands:"
	@echo "  make seed              - Run seeding (endpoints, policies, roles, resources)"
	@echo "  make seed-build        - Build seed binary for production deployment"
	@echo ""
	@echo "IMPORTANT: Seeding uses POSTGRES_DB_PROD environment variable"
	@echo "  - Same as migrations for consistency"
	@echo "  - Set POSTGRES_DB_PROD=maindb for production"
	@echo "  - Set POSTGRES_DB_PROD=devdb for development seeding"
	@echo ""
	@echo "What gets seeded:"
	@echo "  - System Resources (tables configuration)"
	@echo "  - System Roles (Owner, General Manager, etc.)"
	@echo "  - API Endpoints (all routes with permissions)"
	@echo "  - Access Policies (RBAC/ABAC rules)"
	@echo ""
	@echo "Usage in production:"
	@echo "  1. Set: POSTGRES_DB_PROD=maindb in .env"
	@echo "  2. Build: make seed-build"
	@echo "  3. Deploy bin/seed to production server"
	@echo "  4. Run: ./bin/seed (or seed.exe on Windows)"
	@echo ""
	@echo "Note: Seeding is idempotent - safe to run multiple times"
	@echo "      Updates existing records, creates new ones"
