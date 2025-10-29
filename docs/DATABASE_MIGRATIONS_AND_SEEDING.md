# Database Migrations and Seeding - Complete Guide

> **Last Updated:** October 28, 2025  
> **Status:** Production Ready ✅

This is your comprehensive guide to database migrations and seeding in the mynute-go project. Read this first before using any migration or seeding tools.

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Understanding the System](#understanding-the-system)
3. [Migrations](#migrations)
4. [Seeding](#seeding)
5. [Production Workflow](#production-workflow)
6. [Development Workflow](#development-workflow)
7. [Troubleshooting](#troubleshooting)
8. [Tools Reference](#tools-reference)

---

## Quick Start

### First Time Setup

```bash
# 1. Set your database target (CRITICAL!)
# Edit .env file:
POSTGRES_DB=devdb          # Your target database

# 2. Run initial schema migration
go run migrate/main.go up

# 3. Run initial seeding
go run cmd/seed/main.go

# 4. Start your app (dev auto-migrates and seeds)
go run main.go
```

### Daily Development

```bash
# Just start your app - it auto-migrates and auto-seeds
go run main.go
```

### Production Deployment

```bash
# 1. Run migrations manually
POSTGRES_DB=proddb go run migrate/main.go up

# 2. Run seeding manually
POSTGRES_DB=proddb go run cmd/seed/main.go

# 3. Start your app (no auto-migration in prod)
APP_ENV=prod go run main.go
```

---

## Understanding the System

### The Two Pillars

```
┌─────────────────────────────────────────────────────┐
│                  DATABASE SETUP                      │
├──────────────────────┬──────────────────────────────┤
│    MIGRATIONS        │         SEEDING              │
│  (Schema Changes)    │      (Initial Data)          │
├──────────────────────┼──────────────────────────────┤
│ • CREATE TABLE       │ • System Roles               │
│ • ALTER TABLE        │ • API Endpoints              │
│ • DROP TABLE         │ • Access Policies            │
│ • Add/Remove Columns │ • Resources                  │
│ • Indexes            │                              │
│ • Constraints        │ Runs: AFTER migrations       │
│                      │ Why: Needs tables to exist   │
│ Runs: FIRST          │                              │
│ Why: Creates tables  │                              │
└──────────────────────┴──────────────────────────────┘
```

### Multi-Tenant Architecture

Your database has **three schema types**:

1. **`public` schema** - System-wide tables (companies, roles, endpoints)
2. **`company_*` schemas** - Per-tenant tables (branches, employees, appointments)
3. **`tenant` type models** - Special tables that should be in public but reference tenants

```sql
-- Public schema (system-wide)
public.companies
public.roles
public.endpoints
public.policy_rules

-- Company schemas (per-tenant)
company_<uuid>.employees
company_<uuid>.branches
company_<uuid>.appointments
company_<uuid>.services
```

---

## Migrations

### What Are Migrations?

Migrations are **SQL files** that define your database schema. They come in pairs:

- **`.up.sql`** - Apply changes (CREATE, ALTER, etc.)
- **`.down.sql`** - Rollback changes (DROP, revert ALTER, etc.)

### Migration Files Location

```
migrations/
├── 20251028112254_initial_schema.up.sql       ← Creates all tables
├── 20251028112254_initial_schema.down.sql     ← Drops all tables
├── 20251026201359_add_employee_bio.up.sql     ← Adds bio column
└── 20251026201359_add_employee_bio.down.sql   ← Removes bio column
```

### Creating Migrations

#### Option 1: Generate Initial Schema from Go Models

```bash
# Generate CREATE TABLE statements from your GORM models
go run tools/generate-schema/main.go -name initial_schema

# Output:
# migrations/20251028HHMMSS_initial_schema.up.sql
# migrations/20251028HHMMSS_initial_schema.down.sql
```

**Use this for:** First-time schema creation from scratch.

#### Option 2: Detect Schema Drift (Smart Migration)

```bash
# Compare your Go models against current database
# Generates ALTER TABLE statements for differences
go run tools/smart-migration/main.go

# Check specific models only
go run tools/smart-migration/main.go -models Employee,Branch

# Check all models at once
go run tools/smart-migration/main.go -models all
```

**Use this for:** Finding differences between code and database.

⚠️ **Important:** Smart migration compares against the **current database state**, not migration history!

#### Option 3: Manual Migration Templates

```bash
# Generate empty migration template for manual editing
go run tools/generate-migration/main.go -name add_new_field

# Edit the generated files to add your changes
```

**Use this for:** Custom schema changes you want to write yourself.

### Running Migrations

```bash
# Apply all pending migrations
go run migrate/main.go up

# Rollback last migration
go run migrate/main.go down

# Check migration status
go run migrate/main.go version

# Force to specific version
go run migrate/main.go force <version>
```

### Migration Best Practices

✅ **DO:**
- Always use `IF NOT EXISTS` for CREATE statements
- Always use `IF EXISTS` for DROP statements
- Test migrations on a copy of production data first
- Review generated SQL before applying
- Keep migrations small and focused
- Use `ON CONFLICT` for idempotency

❌ **DON'T:**
- Never edit already-applied migrations
- Never run migrations directly in production without testing
- Don't mix schema changes and data changes in one migration
- Don't forget the DOWN migration

---

## Seeding

### What Is Seeding?

Seeding is **populating your database with initial required data**:
- System roles (Owner, Manager, etc.)
- API endpoints (all routes)
- Access policies (RBAC rules)
- Resources (table configurations)

### How Seeding Works

**Seeding is IDEMPOTENT** - Safe to run multiple times:
- Existing records are **updated**
- New records are **inserted**
- No duplicates are created

### Automatic vs Manual Seeding

```go
// In core/server.go
app_env := os.Getenv("APP_ENV")
if app_env == "dev" || app_env == "test" {
    db.Migrate(model.GeneralModels)  // Auto-migrate
    db.InitialSeed()                 // Auto-seed
} else {
    log.Println("Production - run migrations/seeding manually")
}
```

| Environment | Migrations | Seeding |
|-------------|------------|---------|
| `dev`       | ✅ Auto    | ✅ Auto |
| `test`      | ✅ Auto    | ✅ Auto |
| `prod`      | ❌ Manual  | ❌ Manual |

### Running Seeds

```bash
# Development - automatic on app start
go run main.go

# Production - manual execution
go run cmd/seed/main.go

# Build seed binary for deployment
go build -o bin/seed cmd/seed/main.go
./bin/seed  # Run in production
```

### What Gets Seeded

#### 1. Resources
Tables that can be managed via RBAC:
```go
// From core/src/config/db/model/resource.go
var Resources = []*Resource{
    {Table: "appointments"},
    {Table: "branches"},
    {Table: "employees"},
    // ... etc
}
```

#### 2. Roles
System-wide roles (company_id IS NULL):
```go
// From core/src/config/db/model/role.go
var Roles = []*Role{
    {Name: "Owner", Description: "Company Owner..."},
    {Name: "General Manager", Description: "..."},
    {Name: "Branch Manager", Description: "..."},
    // ... etc
}
```

#### 3. Endpoints
All API routes with permissions:
```go
// From core/src/config/db/model/endpoint.go
var CreateAppointment = &EndPoint{
    Path:             "/appointment",
    Method:           "POST",
    DenyUnauthorized: false,
    NeedsCompanyId:   true,
    Resource:         BranchResource,
}
```

#### 4. Policies
RBAC/ABAC access rules:
```go
// From core/src/config/db/model/policy.go
var Policies = []*PolicyRule{
    {
        Role:     OwnerRole,
        EndPoint: CreateAppointment,
        Effect:   "allow",
    },
    // ... etc
}
```

### Updating Seeds

When you add new endpoints, roles, or policies:

1. **Add to model files** (`core/src/config/db/model/`)
2. **Dev:** Restart app → auto-seeds
3. **Prod:** Run `go run cmd/seed/main.go`

**Example - Adding a new endpoint:**

```go
// In core/src/config/db/model/endpoint.go

// 1. Define the endpoint
var MyNewEndpoint = &EndPoint{
    Path:             "/my-new-route",
    Method:           "POST",
    ControllerName:   "MyController",
    Description:      "My new feature",
    DenyUnauthorized: true,
    Resource:         SomeResource,
}

// 2. Add to endpoints list
func EndPoints(cfg *EndpointCfg, db *gorm.DB) ([]*EndPoint, func(), error) {
    // ... existing code ...
    endpoints = append(endpoints, MyNewEndpoint)
    // ...
}

// 3. Dev: Restart app (auto-seeds)
// 4. Prod: go run cmd/seed/main.go
```

---

## Production Workflow

### Pre-Deployment Checklist

- [ ] All migrations tested on staging database
- [ ] Migration DOWN scripts tested (rollback plan)
- [ ] Seeding tested with production-like data
- [ ] Database backup created
- [ ] `POSTGRES_DB` environment variable set correctly
- [ ] No `APP_ENV=dev` in production!

### Deployment Steps

```bash
# 1. Set environment variables
export APP_ENV=prod
export POSTGRES_DB=proddb
export POSTGRES_USER=produser
export POSTGRES_PASSWORD=***

# 2. Backup database
pg_dump -h localhost -U produser proddb > backup_$(date +%Y%m%d).sql

# 3. Run migrations
go run migrate/main.go up

# 4. Verify migrations
go run migrate/main.go version

# 5. Run seeding
go run cmd/seed/main.go

# 6. Verify seeding
psql -h localhost -U produser -d proddb -c "SELECT COUNT(*) FROM public.endpoints;"
psql -h localhost -U produser -d proddb -c "SELECT COUNT(*) FROM public.roles WHERE company_id IS NULL;"

# 7. Start application
go run main.go
```

### Rollback Plan

```bash
# If something goes wrong:

# 1. Stop application
kill <app_pid>

# 2. Rollback last migration
go run migrate/main.go down

# 3. Restore database backup (if needed)
psql -h localhost -U produser -d proddb < backup_YYYYMMDD.sql

# 4. Investigate issue before retrying
```

---

## Development Workflow

### Daily Development

```bash
# Just start your app - everything is automatic!
go run main.go

# What happens:
# 1. Connects to POSTGRES_DB (from .env)
# 2. Auto-runs GORM migrations for GeneralModels
# 3. Auto-seeds roles, endpoints, policies
# 4. Starts server
```

### Adding a New Model

```bash
# 1. Create your model
# core/src/config/db/model/mymodel.go
type MyModel struct {
    BaseModel
    Name string `gorm:"type:varchar(100)"`
}

# 2. Add to GeneralModels or TenantModels
# core/src/config/db/model/general.go
var GeneralModels = []any{
    // ... existing models ...
    &MyModel{},
}

# 3. Restart app (auto-migrates)
go run main.go

# 4. Check for drift (optional)
go run tools/smart-migration/main.go -models MyModel

# 5. Generate proper migration for production
go run tools/generate-schema/main.go -name add_mymodel
# Then manually extract just the MyModel CREATE statement
```

### Checking Schema Drift

```bash
# Check if code differs from database
go run tools/smart-migration/main.go

# Output tells you:
# - What's NEW (in code, not in DB)
# - What's MODIFIED (different between code and DB)
# - Generates SQL to sync database with code
```

### Testing Migrations

```bash
# 1. Create test database
createdb testdb

# 2. Point to test database
export POSTGRES_DB=testdb

# 3. Run migrations
go run migrate/main.go up

# 4. Verify schema
psql -d testdb -c "\dt public.*"
psql -d testdb -c "\dt company_*.*"

# 5. Test rollback
go run migrate/main.go down
```

---

## Troubleshooting

### Migration Errors

#### Error: "Dirty database version"

```bash
# Check current state
go run migrate/main.go version

# Force clean state (use with caution!)
go run migrate/main.go force <version>

# Then retry
go run migrate/main.go up
```

#### Error: "Duplicate key violation"

Your migration isn't idempotent. Fix with:

```sql
-- Instead of:
CREATE TABLE users (...);

-- Use:
CREATE TABLE IF NOT EXISTS users (...);

-- Instead of:
INSERT INTO roles (name) VALUES ('Admin');

-- Use:
INSERT INTO roles (name) VALUES ('Admin')
ON CONFLICT (name) DO NOTHING;
```

#### Error: "Column already exists"

```sql
-- Instead of:
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- Use:
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'email'
    ) THEN
        ALTER TABLE users ADD COLUMN email VARCHAR(255);
    END IF;
END $$;
```

### Seeding Errors

#### Error: "Foreign key violation"

Seeding runs in order:
1. Resources (no dependencies)
2. Roles (no dependencies)
3. Endpoints (depends on Resources)
4. Policies (depends on Roles + Endpoints)

If you get FK errors, check the seeding order in `core/src/config/db/database.go`.

#### Error: "Endpoint not found"

You forgot to call `LoadEndpointIDs()` before seeding policies:

```go
// In database.go InitialSeed()
db.Seed("Endpoints", endpoints, ...)

// MUST call this before policies!
model.LoadEndpointIDs(tx)

db.Seed("Policies", policies, ...)
```

#### Seeding Not Updating Records

Check your match fields:

```go
// This matches on "name" field
db.Seed("Roles", model.Roles, "name = ? AND company_id IS NULL", []string{"Name"})

// If name hasn't changed, record won't update!
// Seeding compares based on match fields, updates all other fields
```

### Database Targeting Issues

#### ⚠️ CRITICAL: Always Verify Your Target Database

```bash
# Check what database you're targeting
echo $POSTGRES_DB

# If wrong, set it:
export POSTGRES_DB=devdb

# Never rely on APP_ENV to determine database!
# Migration tools now explicitly use POSTGRES_DB only
```

See `docs/MIGRATION_DATABASE_CONFIG.md` for the critical POSTGRES_DB change.

---

## Tools Reference

### Migration Tools

| Tool | Command | Purpose |
|------|---------|---------|
| **migrate** | `go run migrate/main.go` | Run migrations up/down |
| **generate-schema** | `go run tools/generate-schema/main.go -name <name>` | Generate CREATE TABLE from models |
| **smart-migration** | `go run tools/smart-migration/main.go [-models <models>]` | Detect schema drift |
| **generate-migration** | `go run tools/generate-migration/main.go -name <name>` | Generate ALTER TABLE template |

### Seeding Tools

| Tool | Command | Purpose |
|------|---------|---------|
| **seed** | `go run cmd/seed/main.go` | Run seeding (prod) |
| **InitialSeed** | Auto in dev/test | Auto-seed on app start |

### Environment Variables

| Variable | Required | Purpose | Example |
|----------|----------|---------|---------|
| `POSTGRES_DB` | ✅ Yes | Target database for migrations | `devdb` |
| `POSTGRES_USER` | ✅ Yes | Database user | `postgres` |
| `POSTGRES_PASSWORD` | ✅ Yes | Database password | `***` |
| `POSTGRES_HOST` | ✅ Yes | Database host | `localhost` |
| `POSTGRES_PORT` | ✅ Yes | Database port | `5432` |
| `APP_ENV` | ✅ Yes | Environment (dev/test/prod) | `dev` |

### Configuration Files

| File | Purpose |
|------|---------|
| `.env` | Local environment variables |
| `.env.example` | Template with defaults |
| `go.mod` | Go dependencies (includes golang-migrate) |

---

## Best Practices Summary

### Migrations

✅ **Always:**
- Use `POSTGRES_DB` to explicitly target database
- Test on staging before production
- Create both UP and DOWN migrations
- Use `IF NOT EXISTS` / `IF EXISTS`
- Keep migrations atomic and small
- Version control all migrations

❌ **Never:**
- Edit applied migrations
- Run migrations in production without testing
- Use `APP_ENV` to determine database target
- Mix schema and data changes
- Delete migration files

### Seeding

✅ **Always:**
- Keep seed data in Go models
- Make seeding idempotent
- Run seeds AFTER migrations
- Test seed updates don't break existing data
- Use transactions for seeding

❌ **Never:**
- Hard-code IDs in seed data
- Assume seeding order
- Skip LoadEndpointIDs() before policies
- Mix system and tenant data

---

## Related Documentation

- **[MIGRATION_DATABASE_CONFIG.md](./MIGRATION_DATABASE_CONFIG.md)** - CRITICAL: POSTGRES_DB change
- **[SEEDING.md](./SEEDING.md)** - Detailed seeding guide
- **[MIGRATIONS.md](./MIGRATIONS.md)** - Legacy migration docs
- **[MIGRATION_WORKFLOW.md](./MIGRATION_WORKFLOW.md)** - Workflow diagrams

---

## Quick Reference Card

```bash
# ==========================================
# DAILY DEVELOPMENT
# ==========================================
go run main.go                              # Auto-migrates + seeds + runs

# ==========================================
# SCHEMA CHANGES
# ==========================================
# Detect drift
go run tools/smart-migration/main.go

# Generate initial schema
go run tools/generate-schema/main.go -name initial_schema

# Generate empty template
go run tools/generate-migration/main.go -name my_change

# ==========================================
# RUNNING MIGRATIONS
# ==========================================
go run migrate/main.go up                   # Apply all pending
go run migrate/main.go down                 # Rollback last
go run migrate/main.go version              # Check status

# ==========================================
# SEEDING
# ==========================================
# Development (automatic)
go run main.go

# Production (manual)
go run cmd/seed/main.go

# Build for deployment
go build -o bin/seed cmd/seed/main.go

# ==========================================
# DATABASE TARGETING
# ==========================================
export POSTGRES_DB=devdb                    # Set target
echo $POSTGRES_DB                           # Verify target

# ==========================================
# PRODUCTION DEPLOYMENT
# ==========================================
export POSTGRES_DB=proddb
go run migrate/main.go up
go run cmd/seed/main.go
APP_ENV=prod go run main.go
```

---

**Questions or Issues?**
1. Check [Troubleshooting](#troubleshooting) section
2. Review related documentation
3. Check migration tool output for hints
4. Verify `POSTGRES_DB` is set correctly

**Last Updated:** October 28, 2025  
**Maintainer:** mynute-app team
