# Production Seeding Guide

## Overview

In production environments, automatic migrations and seeding are disabled for safety. This guide explains how to seed system data (endpoints, policies, roles, resources) in production.

## Why Seeding is Disabled in Production

```go
// core/server.go
if app_env == "dev" || app_env == "test" {
    db.Migrate(model.GeneralModels)
    db.InitialSeed()  // Only runs in dev/test
} else {
    log.Println("Production environment detected - skipping automatic migrations and seeding")
}
```

This prevents accidental data changes during application restarts and ensures controlled deployment processes.

## What Gets Seeded

The seeding process updates/creates:

1. **System Resources** - Table configurations for authorization
2. **System Roles** - Owner, General Manager, Branch Manager, etc.
3. **API Endpoints** - All routes with their permission settings
4. **Access Policies** - RBAC/ABAC rules for each endpoint

## Seeding Methods

### Method 1: Using the Seed Command (Recommended)

#### For Development/Testing
```bash
# Using Make
make seed

# Using PowerShell script
.\scripts\seed.ps1

# Direct execution
go run cmd/seed/main.go
```

#### For Production Deployment

1. **Build the seed binary:**
   ```bash
   make seed-build
   # or
   go build -o bin/seed cmd/seed/main.go
   ```

2. **Deploy to production server:**
   - Copy `bin/seed` (or `bin/seed.exe` on Windows) to your server
   - Ensure `.env` file has correct production database credentials

3. **Run seeding:**
   ```bash
   # Linux/Mac
   ./bin/seed
   
   # Windows
   .\bin\seed.exe
   ```

### Method 2: Using SQL Migrations (Limited)

The SQL migration file `migrations/20251026195226_seed_system_data.up.sql` seeds roles and resources but **NOT** endpoints and policies (they're too complex for static SQL).

```bash
# Run migrations (includes basic seeding)
make migrate-up
```

**Note:** This method won't update endpoints/policies when you change code.

## When to Run Seeding

### Required
- Initial production deployment
- After deploying endpoint changes (new routes, changed permissions)
- After deploying policy changes (RBAC/ABAC rules)
- After system role modifications

### Not Required
- Regular application restarts
- Code changes that don't affect routes/permissions
- Database record updates (users, companies, etc.)

## Seeding is Idempotent

The seed command is **safe to run multiple times**. It:
- Updates existing records with new values
- Creates missing records
- Won't duplicate data
- Won't delete existing data

Example from `database.go`:
```go
// For each model, it either:
if recordExists {
    // Update with new values (including boolean false!)
    tx.Model(oldModel).Select(fieldsToUpdate).Updates(newModel)
} else {
    // Create new record
    tx.Create(newModel)
}
```

## Production Deployment Workflow

### Standard Deployment
```bash
# 1. Run migrations first
make migrate-up

# 2. Build and deploy seed binary
make seed-build
# Copy bin/seed to server

# 3. On production server
./bin/seed

# 4. Start/restart application
./mynute-go
```

### Docker Deployment
Add to your deployment script or `docker-entrypoint.sh`:
```bash
#!/bin/bash
set -e

# Run migrations
./migrate -path ./migrations -database "$DATABASE_URL" up

# Run seeding
./seed

# Start application
./mynute-go
```

### Kubernetes Deployment
Use an init container or Job:
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: mynute-seed
spec:
  template:
    spec:
      containers:
      - name: seed
        image: mynute-go:latest
        command: ["/app/seed"]
        env:
          - name: APP_ENV
            value: "prod"
          - name: POSTGRES_HOST
            valueFrom:
              secretKeyRef:
                name: db-secret
                key: host
      restartPolicy: OnFailure
```

## Updating Endpoint Permissions in Production

When you change an endpoint's configuration (e.g., `DenyUnauthorized: false`):

1. **Update the code** in `core/src/config/db/model/endpoint.go`:
   ```go
   var GetAppointmentByID = &EndPoint{
       Path:             "/appointment/:id",
       Method:           "GET",
       DenyUnauthorized: false,  // Changed from true
       // ...
   }
   ```

2. **Deploy the change:**
   ```bash
   # Build new version
   go build -o bin/mynute-go main.go
   go build -o bin/seed cmd/seed/main.go
   
   # Deploy to production
   scp bin/mynute-go bin/seed production:/app/
   
   # On production server
   ssh production
   cd /app
   ./seed              # Update endpoint configurations
   ./mynute-go         # Restart app with new code
   ```

3. **Verify the change:**
   ```bash
   # Check the database
   psql -d mynute_prod -c "SELECT path, method, deny_unauthorized FROM endpoints WHERE path = '/appointment/:id';"
   ```

## Troubleshooting

### Seeding Fails with Foreign Key Errors
This was fixed in the recent update. If you still see this:
- Ensure you're using the latest code with `LoadEndpointIDs()`
- Check that resources and roles are seeded before endpoints/policies

### Endpoints Not Updating
- Verify the seed command completed successfully
- Check the query keys match (method + path for endpoints)
- Ensure `getUpdateableFields()` is being used (recent fix)

### Boolean Fields Not Updating
Fixed in the recent update. The `getUpdateableFields()` function now always includes boolean fields to allow `true → false` updates.

## Help Commands

```bash
# Show migration help
make migrate-help

# Show seeding help
make seed-help

# Show all available commands
make help
```

## See Also

- [MIGRATIONS.md](./MIGRATIONS.md) - Database migration guide
- [MIGRATIONS_QUICKSTART.md](../MIGRATIONS_QUICKSTART.md) - Quick reference
- `core/src/config/db/database.go` - Seeding implementation
