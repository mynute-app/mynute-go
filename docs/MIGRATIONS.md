# Database Migrations Guide (Atlas)

> **Last Updated:** December 10, 2025  
> **Migration Tool:** Atlas (https://atlasgo.io/)  
> **Status:** Production Ready ✅

This project uses **Atlas** for database migrations. Atlas automatically compares your GORM models with the database and generates migration SQL files.

This guide covers:
1. **First-time production setup** - Initial deployment
2. **Generating new migrations** - When you change GORM models
3. **Applying migrations** - In development and production

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

## Quick Start

### Development
```bash
# Generate migration after changing GORM models
go run migrate/main.go -action=diff -name=add_new_field

# Apply migrations
go run migrate/main.go -action=up

# Check status
go run migrate/main.go -action=status
```

### Production
```bash
# Apply migrations (run BEFORE starting server)
go run migrate/main.go -action=up -env=prod
```

---

## Understanding Atlas Migrations

### What Are Migrations?

Atlas generates **versioned SQL files** by comparing your GORM models with your database:

```
migrations/
├── 20251210215800_init_schema.sql          # Initial schema
├── 20251210220000_add_user_fields.sql      # Add new fields
└── atlas.sum                                # Migration checksum
```

### How Atlas Works

1. **Reads your GORM models** from `core/src/config/db/model/`
2. **Connects to your database** to check current schema
3. **Generates migration SQL** with only the differences
4. **Tracks applied migrations** in `atlas_schema_revisions` table

### Environment Behavior

| Environment | Auto-Migration | Command |
|-------------|----------------|---------|
| `dev/test`  | ❌ Manual      | `go run migrate/main.go -action=up` |
| **`prod`**  | **❌ Manual**  | **`go run migrate/main.go -action=up -env=prod`** |

**⚠️ In production, migrations MUST be run manually before starting the application.**

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
POSTGRES_DB_PROD=maindb go run migrate/main.go -action=up  # ✅ Migrates maindb
POSTGRES_DB_PROD=maindb go run cmd/seed/main.go           # ✅ Seeds maindb

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
go run migrate/main.go -action=version
```

**First-time setup will show:** `error: no migration` (this is normal)

#### 4. Run All Migrations

Apply all pending migrations to create the database schema:

```powershell
go run migrate/main.go -action=up
```

**Expected output:**
```
Migrating to version 20251128111531
Migrating to version 20251128112901
Migration complete!
```

#### 5. Verify Migration Success

```powershell
go run migrate/main.go -action=version
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
- [ ] Ran `go run migrate/main.go -action=up` successfully
- [ ] Verified with `go run migrate/main.go -action=version`
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
go run migrate/main.go -action=version
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
go run migrate/main.go -action=up

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
go run migrate/main.go -action=up
```

**Monitor the output carefully:**
```
Migrating to version 20251128112901
Migration complete!
```

#### 8. Verify Migration Success

```powershell
go run migrate/main.go -action=version
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
go run migrate/main.go -action=down -steps=1
```

#### Option 2: Restore from Backup

```powershell
# Stop application first
# Restore database from backup
pg_restore -h $env:POSTGRES_HOST -U $env:POSTGRES_USER -d $env:POSTGRES_DB_PROD -c backup_20251204_143000.dump

# Verify restoration
go run migrate/main.go -action=version
```

### Ongoing Migration Checklist

- [ ] Reviewed new migration files
- [ ] Created database backup
- [ ] Checked current migration version
- [ ] Tested in staging environment (if available)
- [ ] Scheduled maintenance window (if needed)
- [ ] Ran `go run migrate/main.go -action=up`
- [ ] Verified with `go run migrate/main.go -action=version`
- [ ] Ran `go run cmd/seed/main.go` (if needed)
- [ ] Deployed new application code
- [ ] Monitored application health
- [ ] Documented the deployment

---

## Creating New Migrations

### When You Need a New Migration

When your Go models change, you need to create a migration to update the database schema.

### Automatic Migration Generation (Recommended)

The system automatically detects changes between your Go models and the current database schema:

```powershell
# Automatically detect and generate migration SQL
go run tools/smart-migration/main.go -name=add_employee_bio -models=Employee
```

**What this does:**
- Connects to database (uses `POSTGRES_DB_PROD`)
- Compares your GORM model with actual database schema
- Detects added/removed columns automatically
- Generates both `.up.sql` and `.down.sql` files with correct SQL
- Includes multi-tenant loops for company schemas when needed

**Example output:**
```
📊 Using schema 'company_abc123' for comparison
✅ Generated smart migration files:
  migrations\20251204143000_add_employee_bio.up.sql
  migrations\20251204143000_add_employee_bio.down.sql

💡 Changes detected:
  - Added column: bio (TEXT)
```

**For multiple models:**
```powershell
# Check multiple models at once
go run tools/smart-migration/main.go -name=update_contact_info -models=Employee,Branch
```

### Testing Your Migration

Always test locally before applying to production:

```powershell
# Option 1: Automated test script (recommended)
pwsh -File scripts/test-migration.ps1 -SkipConfirmation

# Option 2: Manual testing
go run migrate/main.go -action=up      # Apply migration
go run migrate/main.go -action=version # Check it worked
go run migrate/main.go -action=down -steps=1  # Rollback
go run migrate/main.go -action=up      # Apply again
```

### Important Notes

⚠️ **Always review the generated SQL** before applying to production, even though it's auto-generated.

⚠️ **The tool compares against your current database** (specified by `POSTGRES_DB_PROD`), not migration history.

⚠️ **Test in your development environment first** with `POSTGRES_DB_PROD=devdb`.

---

## Production Migration Tutorial

This is a complete, step-by-step tutorial for applying migrations to a live production system.

### Prerequisites

- [ ] You have new migration files in `migrations/` directory
- [ ] Migration files have been tested in development
- [ ] You have access to production server
- [ ] You have database backup tools installed (`pg_dump`)
- [ ] You have scheduled a maintenance window (if needed)

### Step 1: Prepare Locally

On your development machine:

```powershell
# 1. Ensure your migration files are committed
git status
git add migrations/
git commit -m "Add migration: [description]"
git push

# 2. Test the migration locally one more time
$env:POSTGRES_DB_PROD = "devdb"
go run migrate/main.go -action=up
# Test your application thoroughly
```

### Step 2: Connect to Production Server

```powershell
# SSH into your production server
ssh user@production-server.com

# Navigate to application directory
cd /app/mynute-go
```

### Step 3: Update Code

```bash
# Pull latest code with new migrations
git pull origin main

# Verify new migration files exist
ls -la migrations/
```

### Step 4: Set Environment Variables

```bash
# Verify production environment is configured
echo $POSTGRES_DB_PROD  # Should show: maindb (or your production DB name)

# If not set, configure it
export POSTGRES_DB_PROD=maindb
export POSTGRES_HOST=your-db-host
export POSTGRES_USER=your-db-user
export POSTGRES_PASSWORD=your-db-password
```

### Step 5: Check Current State

```bash
# Check current migration version
go run migrate/main.go -action=version
```

**Example output:**
```
20251128111531
dirty: false
```

Note this version - it's your rollback point if needed.

### Step 6: Backup Database (CRITICAL!)

```bash
# Create timestamped backup
BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).dump"
pg_dump -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB_PROD -F c -f $BACKUP_FILE

# Verify backup was created
ls -lh $BACKUP_FILE

# Store backup safely (copy to backup location)
cp $BACKUP_FILE /backups/
```

### Step 7: Stop Application (Optional)

For critical schema changes, stop the application during migration:

```bash
# If using systemd
sudo systemctl stop mynute-go

# If using docker
docker stop mynute-go

# If running directly
pkill mynute-go
```

For non-critical changes (adding columns, indexes), you can keep the app running.

### Step 8: Run Migration

```bash
# Run the migration
go run migrate/main.go -action=up
```

**Watch the output carefully:**
```
Running migrations...
Migrating to version 20251204143000
Migration complete!
```

### Step 9: Verify Migration

```bash
# Check new version
go run migrate/main.go -action=version
```

**Expected output:**
```
20251204143000
dirty: false
```

✅ `dirty: false` means migration completed successfully.

### Step 10: Run Seeding (If Needed)

If you added new endpoints, roles, or policies:

```bash
# Run seeding to update system data
go run cmd/seed/main.go
```

**Output should show:**
```
Starting seeding process...
Target database: maindb
⚠️  WARNING: Seeding will modify the database specified by POSTGRES_DB_PROD

✓ Seeding completed successfully!
```

### Step 11: Start Application

```bash
# If using systemd
sudo systemctl start mynute-go
sudo systemctl status mynute-go

# If using docker
docker start mynute-go
docker logs -f mynute-go

# If running directly
./mynute-go &
```

### Step 12: Verify Application Health

```bash
# Check application is running
curl http://localhost:8080/health
# Or your health check endpoint

# Check logs for errors
tail -f /var/log/mynute-go/app.log
# Or: docker logs -f mynute-go

# Test a few API endpoints
curl http://localhost:8080/api/appointments
```

### Step 13: Monitor and Validate

Monitor for the next 15-30 minutes:

```bash
# Watch logs continuously
tail -f /var/log/mynute-go/app.log

# Check for database errors
grep -i "error" /var/log/mynute-go/app.log | tail -20

# Monitor system resources
htop
# or
docker stats mynute-go
```

### Step 14: Cleanup (After 24-48 Hours)

Once you're confident everything is working:

```bash
# Keep backup for 24-48 hours, then clean up old backups
find /backups -name "backup_*.dump" -mtime +2 -delete
```

### Rollback Procedure (If Something Goes Wrong)

If you encounter issues after migration:

#### Quick Rollback (If Migration Just Applied)

```bash
# Stop application
sudo systemctl stop mynute-go

# Rollback the migration
go run migrate/main.go -action=down -steps=1

# Verify rollback
go run migrate/main.go -action=version

# Start application
sudo systemctl start mynute-go
```

#### Full Database Restore (If Rollback Fails)

```bash
# Stop application completely
sudo systemctl stop mynute-go

# Restore from backup
pg_restore -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB_PROD -c $BACKUP_FILE

# Verify restoration
go run migrate/main.go -action=version

# Start application
sudo systemctl start mynute-go
```

### Production Migration Checklist

Before starting:
- [ ] Migration tested in development
- [ ] Code committed and pushed to repository
- [ ] Maintenance window scheduled (if needed)
- [ ] Team notified of deployment
- [ ] Backup tools verified

During migration:
- [ ] Connected to production server
- [ ] Code updated (git pull)
- [ ] Environment variables verified
- [ ] Current migration version noted
- [ ] Database backup created
- [ ] Application stopped (if needed)
- [ ] Migration executed
- [ ] Migration verified (dirty: false)
- [ ] Seeding run (if needed)
- [ ] Application started

After migration:
- [ ] Application health checked
- [ ] API endpoints tested
- [ ] Logs monitored for errors
- [ ] Team notified of completion
- [ ] Documentation updated

### Common Production Scenarios

#### Scenario 1: Zero-Downtime Migration (Adding Column)

```bash
# No need to stop application
go run migrate/main.go -action=up
# Application continues running
```

#### Scenario 2: Maintenance Window Required (Modifying Column)

```bash
# Stop application first
sudo systemctl stop mynute-go
go run migrate/main.go -action=up
sudo systemctl start mynute-go
```

#### Scenario 3: Multiple Migrations Pending

```bash
# Check how many migrations will run
ls migrations/*.up.sql | sort

# Run all pending migrations at once
go run migrate/main.go -action=up

# Verify final version
go run migrate/main.go -action=version
```

### Best Practices for Production Migrations

✅ **Always backup first** - No exceptions, ever.

✅ **Test in staging** - If you have a staging environment, test there first.

✅ **Schedule wisely** - Run during low-traffic periods.

✅ **Monitor closely** - Watch logs and metrics for at least 30 minutes after.

✅ **Communicate** - Notify your team before and after.

✅ **Document** - Keep a log of what was deployed and when.

✅ **Keep backups** - Don't delete backups for 24-48 hours minimum.

---

## Migration File Best Practices

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
go run migrate/main.go -action=version
```

**Output shows dirty:**
```
20251128112901
dirty: true  # ⚠️ Migration failed!
```

**To fix:**

```powershell
# Option 1: Force to previous working version
go run migrate/main.go -action=force -version=20251128111531

# Option 2: Force to current version (if DB is actually correct)
go run migrate/main.go -action=force -version=20251128112901

# Then fix the migration file and try again
go run migrate/main.go -action=up
```

### Application Won't Start After Migration

1. **Check migration status:**
   ```powershell
   go run migrate/main.go -action=version
   ```

2. **Check application logs** for specific error messages

3. **Rollback if necessary:**
   ```powershell
   go run migrate/main.go -action=down -steps=1
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

### Migration Commands

```powershell
# Check which database will be targeted
echo $env:POSTGRES_DB_PROD

# Check current migration version
go run migrate/main.go -action=version

# Run all pending migrations
go run migrate/main.go -action=up

# Rollback last migration
go run migrate/main.go -action=down -steps=1

# Rollback multiple migrations
go run migrate/main.go -action=down -steps=3

# Force to specific version (emergency use only)
go run migrate/main.go -action=force -version=20251128111531
```

### Creating Migrations

```powershell
# Automatically detect changes and generate migration
go run tools/smart-migration/main.go -name=add_employee_bio -models=Employee

# Multiple models
go run tools/smart-migration/main.go -name=update_fields -models=Employee,Branch
```

### Testing Migrations

```powershell
# Automated test (recommended)
pwsh -File scripts/test-migration.ps1 -SkipConfirmation

# Manual testing
go run migrate/main.go -action=up
go run migrate/main.go -action=version
go run migrate/main.go -action=down -steps=1
go run migrate/main.go -action=up
```

### Seeding

```powershell
# Run seeding (updates endpoints, roles, policies)
go run cmd/seed/main.go
```

### Building Binaries for Production

```powershell
# Build migration binary
go build -o bin/migrate migrate/main.go

# Build seed binary
go build -o bin/seed cmd/seed/main.go

# Use in production
./bin/migrate -action=up
./bin/seed
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
      run: go run migrate/main.go -action=up
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
3. Run `go run migrate/main.go -action=up`
4. Run `go run cmd/seed/main.go`
5. Start application

### Ongoing Migrations
1. Backup database
2. Check current version
3. Test in staging
4. Run `go run migrate/main.go -action=up`
5. Deploy new code
6. Monitor health

### Remember
- ⚠️ **ALWAYS backup before migrating production**
- ⚠️ **ALWAYS verify `POSTGRES_DB_PROD` value**
- ⚠️ **NEVER skip testing in staging**
- ⚠️ **NEVER edit applied migrations**

---

## Docker/Dokploy Deployment

For Docker-based deployments (including Dokploy), migrations are **manual operations** that you run explicitly:

### Running Migrations in Docker

```bash
# If using docker-compose with profiles
docker compose -f docker-compose.prod.yml run --rm migrate

# Or directly in a running container
docker exec <container-name> ./migrate-tool up
```

### Important Notes

- **No automatic migrations** - The app starts immediately without running migrations
- **Manual control** - You decide when migrations run (after backups, during maintenance windows)
- **Fast restarts** - Container restarts don't trigger unnecessary migration checks
- **Production best practice** - Migrations should be deliberate, reviewed operations

### Complete Deployment Workflow

See `docs/DOKPLOY_DEPLOYMENT.md` for:
- First-time deployment steps
- Ongoing deployment with migrations
- Troubleshooting common issues
- Complete examples with docker-compose

---

**Need Help?** Run `go run migrate/main.go -action=help` or check available actions: `up`, `down`, `version`, `force`, `create`


