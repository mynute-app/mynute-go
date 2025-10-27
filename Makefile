# Makefile for database migrations
# Requires: golang-migrate CLI tool installed

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Migration commands
.PHONY: migrate-up migrate-down migrate-create migrate-version migrate-force migrate-help

# Run all pending migrations
migrate-up:
	@echo "Running database migrations..."
	@go run migrate/main.go -action=up -path=./migrations

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	@go run migrate/main.go -action=down -steps=1 -path=./migrations

# Rollback N migrations
migrate-down-n:
	@echo "Rolling back $(STEPS) migration(s)..."
	@go run migrate/main.go -action=down -steps=$(STEPS) -path=./migrations

# Check current migration version
migrate-version:
	@echo "Checking migration version..."
	@go run migrate/main.go -action=version -path=./migrations

# Force migration to specific version (use with caution!)
migrate-force:
	@echo "Forcing migration to version $(VERSION)..."
	@go run migrate/main.go -action=force -version=$(VERSION) -path=./migrations

# Create a new migration file (manual)
migrate-create:
ifndef NAME
	@echo "Error: NAME is required. Usage: make migrate-create NAME=your_migration_name"
	@exit 1
endif
	@echo "Creating new migration: $(NAME)"
	@go run migrate/main.go -action=create $(NAME) -path=./migrations

# Auto-generate migration with template for specific models
migrate-generate:
ifndef NAME
	@echo "Error: NAME is required. Usage: make migrate-generate NAME=migration_name [MODELS=Employee,Branch]"
	@exit 1
endif
ifdef MODELS
	@echo "Generating migration for models: $(MODELS)"
	@go run tools/generate-migration/main.go -name=$(NAME) -models=$(MODELS)
else
	@echo "Generating migration template for all models"
	@go run tools/generate-migration/main.go -name=$(NAME) -models=all
endif
	@echo ""
	@echo "⚠️  IMPORTANT: Review and edit the generated SQL files before applying!"

# Smart migration - Auto-detect schema changes
migrate-smart:
ifndef NAME
	@echo "Error: NAME is required. Usage: make migrate-smart NAME=migration_name MODELS=Employee,Branch"
	@exit 1
endif
ifndef MODELS
	@echo "Error: MODELS is required. Usage: make migrate-smart NAME=migration_name MODELS=Employee,Branch"
	@exit 1
endif
	@echo "Analyzing schema changes for models: $(MODELS)"
	@go run tools/smart-migration/main.go -name=$(NAME) -models=$(MODELS)
	@echo ""
	@echo "💡 SQL generated based on detected changes!"

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
