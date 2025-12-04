# Production Database Seeding Guide

> **Last Updated:** December 4, 2025  
> **Status:** Production Ready ✅

This is your complete guide to seeding system data in **production environments**. This guide covers:
1. **First-time production setup** - Initial seeding when deploying for the first time
2. **Ongoing seeding** - Updating system data after changes to endpoints, roles, or policies

---

## Table of Contents

1. [Understanding Database Seeding](#understanding-database-seeding)
2. [Why Seeding Matters in Production](#why-seeding-matters-in-production)
3. [First-Time Production Seeding](#first-time-production-seeding)
4. [Ongoing Production Seeding](#ongoing-production-seeding)
5. [What Gets Seeded](#what-gets-seeded)
6. [Seeding Commands](#seeding-commands)
7. [Common Scenarios](#common-scenarios)
8. [Troubleshooting](#troubleshooting)

---

## Understanding Database Seeding

### Critical Configuration

**IMPORTANT:** Both seeding AND migration tools ALWAYS use `POSTGRES_DB_PROD` environment variable.

This ensures:
- ✅ **Explicit targeting** - You always know which database you're affecting
- ✅ **Consistency** - Same variable for both migrations and seeding  
- ✅ **Production safety** - Can't accidentally seed wrong database

```powershell
# Before running seeding, ALWAYS verify:
echo $env:POSTGRES_DB_PROD  # Windows
echo $POSTGRES_DB_PROD      # Linux/Mac
```

### What is Database Seeding?

Seeding is the process of **populating your database with required system data** that your application needs to function:

- **System Roles** (Owner, Manager, Employee, etc.)
- **API Endpoints** (all routes with their permissions)
- **Access Policies** (RBAC/ABAC authorization rules)
- **Resources** (table configurations for authorization)

### Seeding vs Migrations

| Aspect | Migrations | Seeding |
|--------|-----------|---------|
| **Purpose** | Schema structure (CREATE TABLE, ALTER, etc.) | Initial system data |
| **When** | Before seeding | After migrations |
| **Changes** | Table structure, columns, indexes | System records (roles, endpoints, policies) |
| **Example** | `CREATE TABLE roles (...)` | `INSERT INTO roles VALUES (...)` |

**Critical Order:** Always run migrations **BEFORE** seeding (seeding needs tables to exist).

### Environment Behavior

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

| Environment | Auto-Seeding |
|-------------|--------------|
| `dev`       | ✅ Yes       |
| `test`      | ✅ Yes       |
| **`prod`**  | **❌ No - Manual Required** |

---

## Why Seeding Matters in Production

### The Problem

Without proper seeding in production:
- ❌ API endpoints won't have proper authorization rules
- ❌ Users can't be assigned roles (roles don't exist)
- ❌ Access control policies won't work
- ❌ Changes to endpoint permissions won't take effect

### The Solution

Run seeding manually in production:
- ✅ Creates all required system data
- ✅ Updates endpoint configurations when code changes
- ✅ Safe to run multiple times (idempotent)
- ✅ Won't duplicate or delete existing data

### Idempotency (Safe to Re-run)

The seeding process is **idempotent** - you can run it multiple times safely:

```go
// For each record:
if recordExists {
    // Update with new values from code
    db.Updates(newData)
} else {
    // Create new record
    db.Create(newData)
}
```

**What this means:**
- Existing records are **updated** with latest values
- Missing records are **created**
- No duplicates are created
- No records are deleted

### How Seeding Targets Database

**Important:** The seed command uses `POSTGRES_DB_PROD` (just like migrations):

```bash
# Seeding will affect whatever database POSTGRES_DB_PROD points to
POSTGRES_DB_PROD=maindb ./bin/seed     # Seeds maindb
POSTGRES_DB_PROD=devdb ./bin/seed      # Seeds devdb

# This is consistent with migrations
POSTGRES_DB_PROD=maindb make migrate-up  # Migrates maindb
POSTGRES_DB_PROD=maindb make seed        # Seeds maindb
```

**Why this matters:**
- Both migrations and seeding use the same targeting mechanism
- You explicitly control which database is affected
- Safer than relying on `APP_ENV` alone

---

## First-Time Production Seeding

### Scenario: Setting Up Production Database for the First Time

After you've run migrations and have an empty database schema, you need to populate it with system data.

### Prerequisites

✅ Migrations have been run successfully (see [MIGRATIONS.md](./MIGRATIONS.md))
✅ Database is accessible
✅ Environment variables are configured

### Step-by-Step Process

#### 1. Verify Environment Configuration

```powershell
# Check your environment variables
echo $env:APP_ENV          # Should be "prod"
echo $env:POSTGRES_DB_PROD  # Your production database name
```

Your production `.env` should have:
```env
APP_ENV=prod

# Database connection
POSTGRES_HOST=your-prod-db-host.com
POSTGRES_PORT=5432
POSTGRES_USER=prod_user
POSTGRES_PASSWORD=your_secure_password

# CRITICAL: Migration AND seeding tools use this variable
POSTGRES_DB_PROD=maindb

# These are used by application runtime based on APP_ENV
POSTGRES_DB_DEV=devdb    # App uses when APP_ENV=dev
POSTGRES_DB_TEST=testdb  # App uses when APP_ENV=test
```

#### 2. Build the Seed Binary

**Option A: Build on your development machine**

```powershell
# Build for Linux server (from Windows)
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build -o bin/seed cmd/seed/main.go

# Or use Make
make seed-build
```

**Option B: Build directly on production server**

```bash
# SSH into production server
ssh production

# Build the seed binary
go build -o bin/seed cmd/seed/main.go
```

#### 3. Deploy Seed Binary to Production

```powershell
# Copy to production server
scp bin/seed production:/app/bin/

# Or if using Windows binary locally
.\bin\seed.exe
```

#### 4. Run Seeding

On your production server:

```bash
# Navigate to application directory
cd /app

# Run the seed binary
./bin/seed
```

**Expected output:**
```
Connecting to database...
Connected successfully

Seeding Resources...
✅ Seeded 15 resources

Seeding System Roles...
✅ Seeded 5 roles (Owner, General Manager, Branch Manager, Employee, Other)

Seeding Endpoints...
✅ Seeded 47 endpoints

Seeding Access Policies...
✅ Seeded 47 policies

Seeding completed successfully! 🎉
```

#### 5. Verify Seeding Success

```bash
# Check that data was created
psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB_PROD

# In psql:
\dt                                    -- List tables
SELECT COUNT(*) FROM roles;            -- Should show 5+ roles
SELECT COUNT(*) FROM endpoints;        -- Should show 47+ endpoints
SELECT COUNT(*) FROM policy_rules;     -- Should show 47+ policies
SELECT COUNT(*) FROM resources;        -- Should show 15+ resources
```

#### 6. Start Your Application

Now your database has all required system data:

```bash
./mynute-go
```

The application will:
- ✅ Connect to database
- ✅ Find all required system data
- ✅ Start serving requests with proper authorization

### First-Time Setup Checklist

- [ ] Migrations completed successfully
- [ ] Environment variables configured for production
- [ ] Seed binary built and deployed
- [ ] Seeding executed successfully
- [ ] Verified data in database (roles, endpoints, policies, resources)
- [ ] Application started successfully
- [ ] Tested API endpoints with proper authorization

---

## Ongoing Production Seeding

### Scenario: Your Production is Running, Now You Need to Update System Data

This is when you've made changes to:
- Endpoint configurations (permissions, authorization rules)
- System roles
- Access policies
- API routes

### When to Run Seeding

✅ **Run seeding after these changes:**

1. **Adding new API endpoints**
   ```go
   // Added new endpoint in core/src/config/db/model/endpoint.go
   var NewEndpoint = &EndPoint{
       Path: "/new-feature",
       Method: "POST",
       // ...
   }
   ```

2. **Changing endpoint permissions**
   ```go
   // Changed from true to false
   var GetAppointment = &EndPoint{
       Path: "/appointment/:id",
       Method: "GET",
       DenyUnauthorized: false,  // Changed!
       // ...
   }
   ```

3. **Modifying access policies**
   ```go
   // Changed policy rules
   var Policies = []*PolicyRule{
       // Modified existing policy
   }
   ```

4. **Adding/modifying system roles**
   ```go
   // Added new role or changed existing
   var Roles = []*Role{
       {Name: "NewRole", Description: "..."},
   }
   ```

❌ **Don't need to run seeding after:**
- Regular application restarts
- Bug fixes that don't affect routes/permissions
- User data changes (companies, employees, appointments)
- Frontend-only changes

### Step-by-Step Process

#### 1. Make Changes in Code

Update your endpoint configurations, roles, or policies in the codebase:

```go
// Example: core/src/config/db/model/endpoint.go
var GetAppointmentByID = &EndPoint{
    Path:             "/appointment/:id",
    Method:           "GET",
    DenyUnauthorized: false,  // Changed from true
    NeedsCompanyId:   true,
    Resource:         BranchResource,
}
```

#### 2. Test Locally First

```powershell
# Test in development
$env:APP_ENV = "dev"
go run cmd/seed/main.go

# Verify changes
# Test your application
```

#### 3. Build New Seed Binary

```powershell
# Build for production
make seed-build

# Or manually
go build -o bin/seed cmd/seed/main.go
```

#### 4. Deploy to Production

```powershell
# Deploy both application and seed binary
scp bin/mynute-go production:/app/bin/
scp bin/seed production:/app/bin/
```

#### 5. Run Seeding on Production

```bash
# SSH into production
ssh production
cd /app

# Run seeding (application can stay running)
./bin/seed
```

**The seeding will:**
- Update changed endpoints
- Add new endpoints
- Update modified policies
- Keep existing data intact

#### 6. Restart Application (If Needed)

```bash
# Restart to load new code
systemctl restart mynute-go
# Or
./mynute-go
```

#### 7. Verify Changes

```bash
# Check specific endpoint was updated
psql -d $POSTGRES_DB_PROD -c "
  SELECT path, method, deny_unauthorized, needs_company_id 
  FROM endpoints 
  WHERE path = '/appointment/:id';
"
```

**Expected output:**
```
       path        | method | deny_unauthorized | needs_company_id
-------------------+--------+-------------------+------------------
 /appointment/:id  | GET    | f                 | t
```

### Ongoing Seeding Checklist

- [ ] Made changes to endpoints/roles/policies in code
- [ ] Tested changes in development environment
- [ ] Built new seed binary
- [ ] Deployed seed binary to production
- [ ] Ran seeding on production
- [ ] Verified changes in database
- [ ] Restarted application (if needed)
- [ ] Tested functionality in production

---

## What Gets Seeded

### 1. Resources

**What:** Table configurations that can be managed via RBAC

**Location:** `core/src/config/db/model/resource.go`

**Examples:**
```go
var Resources = []*Resource{
    {Table: "appointments"},
    {Table: "branches"},
    {Table: "employees"},
    {Table: "services"},
    // ... etc
}
```

**Count:** ~15 resources

### 2. System Roles

**What:** Company-wide roles (company_id IS NULL)

**Location:** `core/src/config/db/model/role.go`

**Examples:**
```go
var Roles = []*Role{
    {Name: "Owner", Description: "Company Owner with full access"},
    {Name: "General Manager", Description: "Manages entire company"},
    {Name: "Branch Manager", Description: "Manages specific branch"},
    {Name: "Employee", Description: "Regular employee"},
    {Name: "Other", Description: "Other roles"},
}
```

**Count:** 5 system roles

### 3. API Endpoints

**What:** All API routes with permission configurations

**Location:** `core/src/config/db/model/endpoint.go`

**Examples:**
```go
var Endpoints = []*EndPoint{
    {
        Path:             "/appointment",
        Method:           "POST",
        DenyUnauthorized: false,
        NeedsCompanyId:   true,
        Resource:         BranchResource,
        CreateRecord:     true,
    },
    {
        Path:             "/appointment/:id",
        Method:           "GET",
        DenyUnauthorized: false,
        NeedsCompanyId:   true,
        Resource:         BranchResource,
    },
    // ... etc
}
```

**Count:** ~47 endpoints

**Fields seeded:**
- `path` - API route path
- `method` - HTTP method (GET, POST, PUT, DELETE)
- `deny_unauthorized` - Require authentication
- `needs_company_id` - Require company context
- `resource_id` - Linked resource for RBAC
- `create_record` - Whether this creates a new record
- `company_id` - NULL for system-wide endpoints

### 4. Access Policies

**What:** RBAC/ABAC rules determining who can access what

**Location:** `core/src/config/db/model/policy.go`

**Examples:**
```go
var Policies = []*PolicyRule{
    {
        EndPoint: CreateAppointment,
        Role:     OwnerRole,
        Allowed:  true,
        Condition: func(user *User, target *Resource) bool {
            return user.HasRole("Owner") && user.CompanyId == target.CompanyId
        },
    },
    // ... etc
}
```

**Count:** ~47 policies (one per endpoint typically)

---

## Seeding Commands

### Development

```powershell
# Quick seeding in development
make seed

# Or directly
go run cmd/seed/main.go

# Using PowerShell script
.\scripts\seed.ps1
```

### Production

#### Build the Binary

```powershell
# Using Make (recommended)
make seed-build

# Manual build for Linux (from Windows)
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -o bin/seed cmd/seed/main.go

# Manual build for Windows
go build -o bin/seed.exe cmd/seed/main.go
```

#### Run on Production Server

```bash
# After deploying binary to server
cd /app
./bin/seed
```

#### Using Go Directly (if Go installed on server)

```bash
go run cmd/seed/main.go
```

### CI/CD Integration

**GitHub Actions:**
```yaml
- name: Build Seed Binary
  run: make seed-build

- name: Deploy Seed Binary
  run: scp bin/seed production:/app/bin/

- name: Run Seeding
  run: |
    ssh production 'cd /app && ./bin/seed'
  env:
    POSTGRES_DB_PROD: ${{ secrets.DB_NAME }}
    POSTGRES_HOST: ${{ secrets.DB_HOST }}
    POSTGRES_USER: ${{ secrets.DB_USER }}
    POSTGRES_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

**Docker:**
```dockerfile
# Build stage
FROM golang:1.21 AS builder
WORKDIR /build
COPY . .
RUN go build -o seed cmd/seed/main.go

# Production stage
FROM alpine:latest
COPY --from=builder /build/seed /app/seed
CMD ["/app/seed"]
```

---

## Common Scenarios

### Scenario 1: Changed Endpoint from Authenticated to Public

**Code change:**
```go
// Before
var GetPublicSchedule = &EndPoint{
    Path:             "/schedule/:branch_id",
    Method:           "GET",
    DenyUnauthorized: true,  // Required auth
}

// After
var GetPublicSchedule = &EndPoint{
    Path:             "/schedule/:branch_id",
    Method:           "GET",
    DenyUnauthorized: false,  // Now public
}
```

**Steps:**
```powershell
# 1. Build and deploy
make seed-build
scp bin/seed production:/app/

# 2. Run seeding
ssh production
./bin/seed

# 3. Verify
psql -d maindb -c "SELECT path, deny_unauthorized FROM endpoints WHERE path = '/schedule/:branch_id';"
```

### Scenario 2: Added New API Endpoint

**Code change:**
```go
// Added new endpoint
var CreateNotification = &EndPoint{
    Path:           "/notification",
    Method:         "POST",
    DenyUnauthorized: false,
    NeedsCompanyId:   true,
    Resource:         NotificationResource,
    CreateRecord:     true,
}

// Add to Endpoints slice
var Endpoints = []*EndPoint{
    // ... existing endpoints
    CreateNotification,  // New!
}
```

**Steps:**
```powershell
# 1. Build and deploy
make seed-build
scp bin/seed production:/app/

# 2. Run seeding (will create new endpoint record)
ssh production
./bin/seed

# 3. Verify
psql -d maindb -c "SELECT * FROM endpoints WHERE path = '/notification';"
```

### Scenario 3: Updated Access Policy

**Code change:**
```go
// Modified policy to allow Branch Managers
var Policies = []*PolicyRule{
    {
        EndPoint: ViewReports,
        Role:     BranchManagerRole,  // Added this role
        Allowed:  true,
    },
}
```

**Steps:**
```powershell
# Same process - seeding updates policies
make seed-build
scp bin/seed production:/app/
ssh production './bin/seed'
```

---

## Troubleshooting

### Problem: "Foreign Key Constraint" Errors

**Symptoms:**
```
ERROR: insert or update violates foreign key constraint
```

**Cause:** Resources or roles not seeded before endpoints/policies

**Solution:**
Seeding now uses `LoadEndpointIDs()` to ensure proper order. If you still see this:
```bash
# Check that resources exist
psql -d maindb -c "SELECT COUNT(*) FROM resources;"

# Check that roles exist
psql -d maindb -c "SELECT COUNT(*) FROM roles;"

# If counts are 0, there's a code issue
# Check core/src/config/db/database.go - InitialSeed() order
```

### Problem: Endpoint Permissions Not Updating

**Symptoms:**
- Changed `DenyUnauthorized: false` in code
- Ran seeding
- Still returns 401 Unauthorized

**Solutions:**

1. **Verify seeding actually ran:**
   ```bash
   # Check last modified timestamp
   psql -d maindb -c "SELECT path, updated_at FROM endpoints WHERE path = '/your-path';"
   ```

2. **Verify the field was updated:**
   ```bash
   psql -d maindb -c "SELECT path, deny_unauthorized FROM endpoints WHERE path = '/your-path';"
   ```

3. **Restart application:**
   ```bash
   # Application may cache endpoint configurations
   systemctl restart mynute-go
   ```

4. **Check that getUpdateableFields() includes boolean fields:**
   The recent fix ensures boolean fields are always updated. Verify you have latest code.

### Problem: Boolean Fields Not Updating from `true` to `false`

**Symptoms:**
- Changed boolean from `true` to `false`
- Ran seeding
- Database still shows `true`

**Cause:** Fixed in recent update

**Solution:**
Ensure you're using latest code where `getUpdateableFields()` includes all boolean fields:

```go
// core/src/config/db/database.go
func getUpdateableFields(model interface{}) []string {
    // ... includes ALL fields including booleans
}
```

### Problem: Seeding Takes Long Time

**Symptoms:**
- Seeding process takes several minutes

**Causes:**
1. Large number of company schemas (multi-tenant)
2. Network latency to database
3. Many endpoints/policies to process

**Solutions:**
- Normal if you have many tenant schemas
- Consider running during maintenance window
- Optimize database connection (increase connection pool)

### Problem: Duplicate Records Created

**Symptoms:**
- Multiple identical roles, endpoints, or policies

**Cause:** Query keys don't match properly

**Solution:**
The seeding uses specific fields to identify existing records:
- **Endpoints:** `method + path`
- **Roles:** `name + company_id`
- **Policies:** `endpoint_id + role_id`

If you see duplicates, check that these fields are consistent in your code.

---

## Best Practices

### ✅ DO:

1. **Always run seeding after endpoint changes**
   ```powershell
   # Changed endpoint code? Run seeding!
   ./bin/seed
   ```

2. **Test seeding in development first**
   ```powershell
   $env:APP_ENV = "dev"
   go run cmd/seed/main.go
   ```

3. **Run seeding as part of deployment**
   - Include in CI/CD pipeline
   - Run before starting new application version

4. **Keep seed binary updated**
   - Rebuild when endpoint code changes
   - Version alongside application binary

5. **Verify after seeding**
   ```bash
   # Check counts
   psql -d maindb -c "SELECT COUNT(*) FROM endpoints;"
   ```

### ❌ DON'T:

1. **Don't skip seeding after endpoint changes**
   - Your authorization rules won't match your code

2. **Don't assume auto-seeding works in production**
   - It's explicitly disabled for safety

3. **Don't run seeding without testing**
   - Always test in development first

4. **Don't forget to restart app after seeding**
   - Some changes may require app restart

---

## Summary

### First-Time Production Setup
1. Verify environment configuration
2. Build seed binary
3. Deploy to production
4. Run `./bin/seed`
5. Verify data in database
6. Start application

### Ongoing Production Seeding
1. Make changes to endpoints/roles/policies in code
2. Test in development
3. Build new seed binary
4. Deploy to production
5. Run `./bin/seed`
6. Restart application (if needed)
7. Verify changes

### Remember
- ✅ Seeding is **idempotent** - safe to run multiple times
- ✅ Always run **after migrations** (needs tables to exist)
- ✅ Required for endpoint permission changes to take effect
- ✅ Updates existing records, creates missing ones
- ✅ Never deletes data

---

## Docker/Dokploy Deployment

For Docker-based deployments (including Dokploy), seeding is a **manual operation** that you run explicitly:

### Running Seeding in Docker

```bash
# If using docker-compose with profiles
docker compose -f docker-compose.prod.yml run --rm seed

# Or directly in a running container
docker exec <container-name> ./seed-tool
```

### When to Run Seeding

- **First deployment** - Populate initial resources, roles, endpoints, and policies
- **After endpoint changes** - When you add/modify API endpoints in code
- **After role/policy changes** - When authorization rules change
- **After database reset** - To repopulate system data

### Important Notes

- **No automatic seeding** - The app starts immediately without running seeding
- **Manual control** - You decide when seeding runs
- **Idempotent** - Safe to run multiple times (won't create duplicates)
- **Always after migrations** - Run seeding after migrations complete

### Complete Deployment Workflow

See `docs/DOKPLOY_DEPLOYMENT.md` for:
- First-time deployment with seeding
- Ongoing deployments with selective seeding
- Complete examples with docker-compose
- Troubleshooting common issues

---

**Need Help?** Run `go run cmd/seed/main.go -help` for more options or see [MIGRATIONS.md](./MIGRATIONS.md) for related migration commands.
