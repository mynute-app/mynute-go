# Production Database Migrations Guide

> **Last Updated:** December 4, 2025  
> **Status:** Production Ready ✅

This is your complete guide to running database migrations in **production environments**. This guide focuses on two critical scenarios:
1. **First-time production setup** - When deploying your server for the first time
2. **Ongoing migrations** - When applying new migrations to an already-running production system

---

## Table of Contents

1. [Understanding the System](#understanding-the-system)
2. [Critical Configuration](#critical-configuration)
3. [First-Time Production Setup](#first-time-production-setup)
4. [Running Migrations on Existing Production](#running-migrations-on-existing-production)
5. [Creating New Migrations](#creating-new-migrations)
6. [Emergency Procedures](#emergency-procedures)
7. [Best Practices](#best-practices)

---

## Understanding the System

### What Are Migrations?

Migrations are **version-controlled SQL files** that define your database schema. They ensure your database structure matches your application code.

```
migrations/
├── 20251128111531_change_employee_endpoint_path_parameters.up.sql
├── 20251128111531_change_employee_endpoint_path_parameters.down.sql
├── 20251128112901_fix_get_employee_work_range_path.up.sql
└── 20251128112901_fix_get_employee_work_range_path.down.sql
```

Each migration has two files:
- **`.up.sql`** - Applies changes (CREATE, ALTER, etc.)
- **`.down.sql`** - Reverts changes (for rollback)

### Environment Behavior

| Environment | Auto-Migration | Manual Migration Required |
|-------------|----------------|--------------------------|
| `dev`       | ✅ Yes         | ❌ No                    |
| `test`      | ✅ Yes         | ❌ No                    |
| **`prod`**  | **❌ No**      | **✅ YES - REQUIRED**    |

**In production, migrations MUST be run manually before starting the application.**

---

## Critical Configuration

### Database Configuration

**CRITICAL:** Both migration tools AND seeding tools ALWAYS use `POSTGRES_DB_PROD` environment variable.

This ensures you explicitly target the correct database and prevents accidental operations on the wrong environment.

**Why this matters:**
- ✅ Explicit targeting - You always know which database you're affecting
- ✅ No environment confusion - Same variable for migrations and seeding
- ✅ Production safety - Can't accidentally migrate/seed wrong database

### Production `.env` File

```env
# Application environment
APP_ENV=prod

# Database connection
POSTGRES_HOST=your-prod-db-host.com
POSTGRES_PORT=5432
POSTGRES_USER=prod_user
POSTGRES_PASSWORD=your_secure_password

# CRITICAL: Migration AND seeding tools use this variable
POSTGRES_DB_PROD=maindb

# These are used by the application runtime based on APP_ENV
POSTGRES_DB_DEV=devdb    # App uses when APP_ENV=dev
POSTGRES_DB_TEST=testdb  # App uses when APP_ENV=test
```

### Why This Matters

```bash
# Migration and seeding tools target what POSTGRES_DB_PROD points to
POSTGRES_DB_PROD=maindb make migrate-up  # ✅ Migrates maindb
POSTGRES_DB_PROD=maindb make seed        # ✅ Seeds maindb

# The application uses APP_ENV to determine the database
APP_ENV=prod ./mynute-go  # ✅ Connects to maindb (via POSTGRES_DB_PROD)
```

**Before running ANY migration command, verify your configuration:**

```powershell
# Windows PowerShell
echo $env:POSTGRES_DB_PROD

# Linux/Mac
echo $POSTGRES_DB_PROD
```

---

## First-Time Production Setup

### Scenario: Deploying Your Application for the First Time

This is when you're setting up a brand new production database that has never run the application before.

### Step-by-Step Process

#### 1. Prepare Your Environment

```powershell
# Set production environment variables
$env:APP_ENV = "prod"
$env:POSTGRES_DB_PROD = "maindb"
$env:POSTGRES_HOST = "your-prod-db-host.com"
$env:POSTGRES_USER = "prod_user"
$env:POSTGRES_PASSWORD = "your_secure_password"
```

Or create a production `.env` file:

```env
APP_ENV=prod
POSTGRES_DB_PROD=maindb
POSTGRES_HOST=your-prod-db-host.com
POSTGRES_PORT=5432
POSTGRES_USER=prod_user
POSTGRES_PASSWORD=your_secure_password
```

#### 2. Verify Database Connection

Before running migrations, ensure you can connect to the database:

```powershell
# Test connection (replace with your actual credentials)
psql -h $env:POSTGRES_HOST -U $env:POSTGRES_USER -d $env:POSTGRES_DB_PROD -c "SELECT version();"
```

Expected output: PostgreSQL version information

#### 3. Check Migration Status

See which migrations exist and their status:

```powershell
make migrate-version
```

**First-time setup will show:** `error: no migration` (this is normal)

#### 4. Run All Migrations

Apply all pending migrations to create the database schema:

```powershell
make migrate-up
```

**Expected output:**
```
Migrating to version 20251128111531
Migrating to version 20251128112901
Migration complete!
```

#### 5. Verify Migration Success

```powershell
make migrate-version
```

**Expected output:**
```
20251128112901
dirty: false
```

#### 6. Run Initial Data Seeding

After migrations, seed the database with required system data (roles, endpoints, policies):

```powershell
go run cmd/seed/main.go
```

**Expected output:**
```
Seeding Resources...
Seeding Roles...
Seeding Endpoints...
Seeding Policies...
Seeding complete!
```

#### 7. Start Your Application

Now your database is ready. Start the application:

```powershell
go run main.go
# Or if you have a compiled binary
./mynute-go
```

The application will:
- ✅ Connect to the database
- ✅ Skip auto-migrations (because APP_ENV=prod)
- ✅ Use the existing schema
- ✅ Start serving requests

### First-Time Setup Checklist

- [ ] Environment variables configured (`POSTGRES_DB_PROD`, `APP_ENV=prod`)
- [ ] Database server is running and accessible
- [ ] Database exists (create with `CREATE DATABASE maindb;` if needed)
- [ ] Ran `make migrate-up` successfully
- [ ] Verified with `make migrate-version`
- [ ] Ran `go run cmd/seed/main.go` successfully
- [ ] Started application successfully

---

## Running Migrations on Existing Production

### Scenario: Your Production Server Has Been Running, Now You Need to Apply New Migrations

This is when you've added new features that require database schema changes.

### When Do You Need This?

- Adding new tables
- Adding/removing columns
- Changing data types
- Adding indexes or constraints
- Any schema modification

### Step-by-Step Process

#### 1. Understand What's Changing

Review the new migration files in your codebase:

```powershell
# List all migration files
ls migrations/
```

Look at the newest migration files to understand what will change.

#### 2. Backup Your Database (CRITICAL!)

**Never skip this step in production:**

```powershell
# Create a backup before running migrations
pg_dump -h $env:POSTGRES_HOST -U $env:POSTGRES_USER -d $env:POSTGRES_DB_PROD -F c -f "backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').dump"
```

This creates a timestamped backup file you can restore if needed.

#### 3. Check Current Migration Status

See what version your production database is currently on:

```powershell
make migrate-version
```

**Example output:**
```
20251128111531
dirty: false
```

This tells you:
- Current version: `20251128111531`
- Status: Clean (not dirty)

#### 4. Review Pending Migrations

Check which migrations will be applied:

```powershell
# List migration files newer than current version
ls migrations/ | Where-Object { $_.Name -gt "20251128111531" }
```

#### 5. Test in Staging First (Recommended)

If you have a staging environment, test the migration there first:

```powershell
# Point to staging database
$env:POSTGRES_DB_PROD = "staging_db"
make migrate-up

# Verify staging works
# Test your application
```

#### 6. Schedule Maintenance Window (If Needed)

For critical migrations that might cause downtime:
- Notify users of planned maintenance
- Schedule during low-traffic periods
- Prepare rollback plan

#### 7. Run the Migration

```powershell
# Ensure you're pointing to production
echo $env:POSTGRES_DB_PROD  # Should show your prod database name

# Run migrations
make migrate-up
```

**Monitor the output carefully:**
```
Migrating to version 20251128112901
Migration complete!
```

#### 8. Verify Migration Success

```powershell
make migrate-version
```

Ensure:
- Version matches the latest migration
- `dirty: false` (migration completed successfully)

#### 9. Run Seeding (If New System Data Added)

If the migration added new roles, endpoints, or policies:

```powershell
go run cmd/seed/main.go
```

Seeding is idempotent - safe to run multiple times.

#### 10. Deploy New Application Code

Now deploy your updated application code:

```powershell
# Restart application with new code
./mynute-go
```

#### 11. Monitor Application Health

Watch for:
- Application starts successfully
- No database connection errors
- API endpoints respond correctly
- No error logs related to database schema

### Rollback Procedure (If Something Goes Wrong)

#### Option 1: Rollback Migration

```powershell
# Rollback the last migration
make migrate-down
```

#### Option 2: Restore from Backup

```powershell
# Stop application first
# Restore database from backup
pg_restore -h $env:POSTGRES_HOST -U $env:POSTGRES_USER -d $env:POSTGRES_DB_PROD -c backup_20251204_143000.dump

# Verify restoration
make migrate-version
```

### Ongoing Migration Checklist

- [ ] Reviewed new migration files
- [ ] Created database backup
- [ ] Checked current migration version
- [ ] Tested in staging environment (if available)
- [ ] Scheduled maintenance window (if needed)
- [ ] Ran `make migrate-up`
- [ ] Verified with `make migrate-version`
- [ ] Ran `go run cmd/seed/main.go` (if needed)
- [ ] Deployed new application code
- [ ] Monitored application health
- [ ] Documented the deployment

---

## Creating New Migrations

### When You Need a New Migration

When your Go models change, you need to create a migration to update the database schema.

### Method 1: Smart Migration (Automatic Detection)

Automatically detects changes between your Go models and database:

```powershell
# Compare model with database and generate SQL
make migrate-smart NAME=add_employee_bio MODELS=Employee
```

**This will:**
- Connect to database (uses `POSTGRES_DB_PROD`)
- Compare GORM model with actual schema
- Generate both `.up.sql` and `.down.sql` files
- Include proper multi-tenant loops if needed

**Example output:**
```
📊 Using schema 'company_abc123' for comparison
✅ Generated smart migration files:
  migrations\20251204143000_add_employee_bio.up.sql
  migrations\20251204143000_add_employee_bio.down.sql

💡 Changes detected:
  - Added column: bio (TEXT)
```

### Method 2: Generate Template

Creates migration template with examples:

```powershell
make migrate-generate NAME=add_employee_bio MODELS=Employee
```

Then edit the generated files to add your specific SQL.

### Method 3: Manual Migration

For complex migrations or data transformations:

```powershell
# Create empty migration files
make migrate-create NAME=complex_data_migration
```

Then write the SQL manually in both files.

### Testing Your Migration

Always test before applying to production:

```powershell
# Automated test: UP → verify → DOWN → verify → UP
make test-migrate
```

### Migration File Best Practices

#### Always Use Idempotent SQL

```sql
-- ✅ Good - Can run multiple times safely
CREATE TABLE IF NOT EXISTS notifications (...);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS bio TEXT;

-- ❌ Bad - Will fail if run twice
CREATE TABLE notifications (...);
ALTER TABLE employees ADD COLUMN bio TEXT;
```

#### Always Handle Multi-Tenant Schemas

For tenant-specific tables:

```sql
-- Iterate over all company schemas
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS bio TEXT', schema_name);
    END LOOP;
END $$;
```

#### Always Provide Rollback (DOWN migration)

```sql
-- down.sql must undo what up.sql did
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS bio', schema_name);
    END LOOP;
END $$;
```

---

## Emergency Procedures

### Migration Failed Midway ("Dirty" State)

If a migration fails partway through:

```powershell
# Check status
make migrate-version
```

**Output shows dirty:**
```
20251128112901
dirty: true  # ⚠️ Migration failed!
```

**To fix:**

```powershell
# Option 1: Force to previous working version
make migrate-force VERSION=20251128111531

# Option 2: Force to current version (if DB is actually correct)
make migrate-force VERSION=20251128112901

# Then fix the migration file and try again
make migrate-up
```

### Application Won't Start After Migration

1. **Check migration status:**
   ```powershell
   make migrate-version
   ```

2. **Check application logs** for specific error messages

3. **Rollback if necessary:**
   ```powershell
   make migrate-down
   ```

4. **Restore from backup if rollback doesn't work:**
   ```powershell
   pg_restore -h $env:POSTGRES_HOST -U $env:POSTGRES_USER -d $env:POSTGRES_DB_PROD -c backup_file.dump
   ```

### Wrong Database Migrated

If you accidentally migrated the wrong database:

1. **Immediately stop the application**

2. **Check which database was affected:**
   ```powershell
   echo $env:POSTGRES_DB_PROD
   ```

3. **Restore that database from backup**

4. **Fix your environment configuration**

5. **Re-run migration on correct database**

---

## Best Practices

### Before Running Migrations

✅ **ALWAYS check which database you're targeting:**
```powershell
echo $env:POSTGRES_DB_PROD
```

✅ **ALWAYS backup production database:**
```powershell
pg_dump -F c -f backup.dump ...
```

✅ **ALWAYS test in staging first** (if you have staging environment)

✅ **ALWAYS review migration SQL files** before applying

### During Migrations

✅ **Monitor the output** - watch for errors

✅ **Don't interrupt** - let migrations complete

✅ **Keep terminal session alive** - use tmux/screen for SSH sessions

### After Migrations

✅ **Verify migration version** matches expectations

✅ **Test application functionality** immediately

✅ **Monitor error logs** for issues

✅ **Keep backup for 24-48 hours** before deleting

### General Best Practices

✅ **One logical change per migration** - don't combine unrelated changes

✅ **Use descriptive migration names** - `add_employee_bio` not `update_stuff`

✅ **Never edit applied migrations** - create new ones instead

✅ **Commit migration files to git** - version control is critical

✅ **Use transactions where possible** - wrap in BEGIN/COMMIT

✅ **Add indexes concurrently** - `CREATE INDEX CONCURRENTLY` for zero downtime

---

## Command Reference

### Essential Commands

```powershell
# Check which database will be targeted
echo $env:POSTGRES_DB_PROD

# Check current migration version
make migrate-version

# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Rollback multiple migrations
make migrate-down-n STEPS=3

# Force to specific version (emergency use)
make migrate-force VERSION=20251128111531

# Create new migration (automatic detection)
make migrate-smart NAME=description MODELS=ModelName

# Create new migration (template)
make migrate-generate NAME=description MODELS=ModelName

# Create new migration (manual)
make migrate-create NAME=description

# Test migration automatically
make test-migrate

# Run seeding
go run cmd/seed/main.go

# Show all available commands
make migrate-help
```

### CI/CD Integration

**GitHub Actions:**
```yaml
deploy-production:
  steps:
    - name: Backup Database
      run: |
        pg_dump -h ${{ secrets.DB_HOST }} -U ${{ secrets.DB_USER }} \
          -d ${{ secrets.DB_NAME }} -F c -f backup.dump
      env:
        PGPASSWORD: ${{ secrets.DB_PASSWORD }}
    
    - name: Run Migrations
      run: make migrate-up
      env:
        POSTGRES_DB_PROD: ${{ secrets.DB_NAME }}
        POSTGRES_HOST: ${{ secrets.DB_HOST }}
        POSTGRES_USER: ${{ secrets.DB_USER }}
        POSTGRES_PASSWORD: ${{ secrets.DB_PASSWORD }}
    
    - name: Run Seeding
      run: go run cmd/seed/main.go
      env:
        POSTGRES_DB_PROD: ${{ secrets.DB_NAME }}
        POSTGRES_HOST: ${{ secrets.DB_HOST }}
        POSTGRES_USER: ${{ secrets.DB_USER }}
        POSTGRES_PASSWORD: ${{ secrets.DB_PASSWORD }}
    
    - name: Deploy Application
      run: |
        # Your deployment commands here
```

---

## Multi-Tenant Architecture Notes

Your application uses a multi-tenant architecture with separate schemas per company:

```
maindb (database)
├── public schema (system-wide tables)
│   ├── companies
│   ├── roles
│   ├── endpoints
│   └── policy_rules
│
└── company_<uuid> schemas (per-tenant tables)
    ├── employees
    ├── branches
    ├── appointments
    └── services
```

**Important:** When creating migrations for tenant-specific tables, always iterate over all `company_*` schemas (see examples above).

---

## Troubleshooting

### "No change" error

**Meaning:** All migrations are already applied. This is normal if your database is up to date.

### "Connection refused"

**Fix:**
1. Verify database server is running
2. Check `POSTGRES_HOST` and `POSTGRES_PORT`
3. Verify firewall rules allow connection
4. Test with `psql` command directly

### "Authentication failed"

**Fix:**
1. Verify `POSTGRES_USER` and `POSTGRES_PASSWORD`
2. Check database user has proper permissions
3. Ensure user has access to `POSTGRES_DB_PROD` database

### "Database does not exist"

**Fix:**
```powershell
# Create the database
psql -h $env:POSTGRES_HOST -U $env:POSTGRES_USER -c "CREATE DATABASE $env:POSTGRES_DB_PROD;"
```

---

## Summary

### First-Time Setup
1. Configure environment variables
2. Verify database connection
3. Run `make migrate-up`
4. Run `go run cmd/seed/main.go`
5. Start application

### Ongoing Migrations
1. Backup database
2. Check current version
3. Test in staging
4. Run `make migrate-up`
5. Deploy new code
6. Monitor health

### Remember
- ⚠️ **ALWAYS backup before migrating production**
- ⚠️ **ALWAYS verify `POSTGRES_DB_PROD` value**
- ⚠️ **NEVER skip testing in staging**
- ⚠️ **NEVER edit applied migrations**

---

**Need Help?** Run `make migrate-help` for command reference


