# Database Migrations Guide

## Overview

This project uses **golang-migrate** for database schema management, providing version-controlled, reversible migrations for production environments.

## 🚀 Quick Start

### Prerequisites

1. **Install golang-migrate CLI** (already installed for this project):
   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

2. **Ensure your `.env` file is configured** with the correct database credentials:
   ```env
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=your_user
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=your_prod_db
   POSTGRES_DB_DEV=your_dev_db
   POSTGRES_DB_TEST=your_test_db
   APP_ENV=dev  # or test, prod
   ```

## 📋 Common Commands

### Run Migrations

```bash
# Run all pending migrations
make migrate-up

# Or using Go directly
go run migrate/main.go -action=up
```

### Rollback Migrations

```bash
# Rollback the last migration
make migrate-down

# Rollback multiple migrations
make migrate-down-n STEPS=3
```

### Create New Migration

```bash
# Create a new migration file
make migrate-create NAME=add_user_preferences

# Or using Go directly
go run migrate/main.go -action=create add_user_preferences
```

This creates two files:
- `migrations/YYYYMMDDHHMMSS_add_user_preferences.up.sql` - Apply changes
- `migrations/YYYYMMDDHHMMSS_add_user_preferences.down.sql` - Rollback changes

### Check Migration Status

```bash
# Check current migration version
make migrate-version
```

### Force Migration Version (⚠️ Use with Caution!)

If a migration is "dirty" (failed midway), you can force the version:

```bash
make migrate-force VERSION=20251026195057
```

## 🏗️ Migration File Structure

Migration files are located in the `migrations/` directory with naming convention:
```
YYYYMMDDHHMMSS_description.up.sql    # Apply migration
YYYYMMDDHHMMSS_description.down.sql  # Rollback migration
```

### Example Migration

**`20251027120000_add_notifications.up.sql`:**
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    message TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
```

**`20251027120000_add_notifications.down.sql`:**
```sql
DROP TABLE IF EXISTS notifications;
```

## 🔄 Environment-Specific Behavior

### Development (`APP_ENV=dev`)
- ✅ Auto-migrations enabled (GORM AutoMigrate)
- ✅ Auto-seeding enabled
- Safe for rapid prototyping

### Test (`APP_ENV=test`)
- ✅ Auto-migrations enabled
- ✅ Auto-seeding enabled
- Database cleared before each test run

### Production (`APP_ENV=prod`)
- ❌ Auto-migrations **DISABLED**
- ❌ Auto-seeding **DISABLED**
- **Must run migrations manually** before deployment

## 🚢 Production Deployment Workflow

### Recommended Deployment Process

1. **Review migration files** in your PR/MR
2. **Test migrations in staging environment**:
   ```bash
   APP_ENV=staging make migrate-up
   ```

3. **Run migrations BEFORE deploying new code**:
   ```bash
   # SSH into production server or use CI/CD pipeline
   APP_ENV=prod make migrate-up
   ```

4. **Deploy application code** (migrations already applied)

5. **Verify migration status**:
   ```bash
   APP_ENV=prod make migrate-version
   ```

### CI/CD Integration Example

**GitHub Actions / GitLab CI:**
```yaml
deploy:
  steps:
    - name: Run database migrations
      run: |
        export APP_ENV=prod
        make migrate-up
    
    - name: Deploy application
      run: |
        # Your deployment commands here
```

### Rollback in Production

If you need to rollback:

```bash
# Rollback 1 migration
APP_ENV=prod make migrate-down

# Rollback to specific version
APP_ENV=prod make migrate-force VERSION=20251026195057
```

## 🔐 Best Practices

### ✅ DO:

1. **Always create both UP and DOWN migrations**
2. **Test migrations in staging before production**
3. **Make migrations backward compatible** when possible
4. **Use transactions for data migrations** (wrap in BEGIN/COMMIT)
5. **Keep migrations small and focused** (one logical change per migration)
6. **Add indexes concurrently** in PostgreSQL:
   ```sql
   CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
   ```

### ❌ DON'T:

1. **Don't edit existing migrations** after they've been deployed
2. **Don't delete migration files** from version control
3. **Don't run auto-migrate in production** (now disabled)
4. **Don't assume migrations are instant** - plan for downtime or use online migrations
5. **Don't forget to commit migration files** to git

## 🛠️ Troubleshooting

### Migration is "dirty"

If a migration fails midway:

```bash
# Check current state
make migrate-version

# Force to the previous working version
make migrate-force VERSION=<previous_version>

# Fix the migration file
# Then re-run
make migrate-up
```

### Multiple instances running migrations

Use database locks or run migrations as a separate step in deployment:

```bash
# Run migration as a one-off job before starting app instances
kubectl run migration-job --image=yourapp -- make migrate-up
```

### Schema vs Data Migrations

**Schema changes** (CREATE TABLE, ALTER TABLE):
- Safe to run during deployment
- Use standard migration files

**Data migrations** (UPDATE, INSERT):
- Consider doing these in separate steps
- May require backfilling existing data
- Test with production-sized datasets

## 📚 Additional Resources

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL migration best practices](https://www.postgresql.org/docs/current/ddl-alter.html)
- [Zero-downtime migrations guide](https://www.braintreepayments.com/blog/safe-operations-for-high-volume-postgresql/)

## 🆘 Need Help?

Run `make migrate-help` to see all available commands.
