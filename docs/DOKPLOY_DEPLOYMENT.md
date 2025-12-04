# Dokploy Deployment Guide

## Overview

Your Dockerfile includes the migration and seeding tools, but they are **run separately** as manual operations:

✅ **Migration tool** - Built and included in the image  
✅ **Seeding tool** - Built and included in the image  
✅ **Migration files** - Copied to `/mynute-go/migrations`  
✅ **Separate execution** - Run migrations/seeding explicitly when needed  
✅ **No automatic runs** - App starts immediately without migration delays  

## Deployment Workflow

### Initial Deployment (First Time)

1. **Deploy the stack** (creates database and app)
2. **Run migrations** (one-time setup)
3. **Run seeding** (one-time setup)
4. **App is ready** to serve traffic

### Ongoing Deployments (Updates)

1. **Review changes** - Check if new migrations exist
2. **Backup database** (if migrations present)
3. **Run new migrations** (only if needed)
4. **Deploy updated app** (container restarts)
5. **Run seeding** (only if new resources added)

## Deployment Commands

### Step 1: Initial Deployment

Deploy the stack (database + app):

```bash
docker compose -f docker-compose.prod.yml up -d
```

For Dokploy, use:
```bash
docker compose build --no-cache && docker compose up -d --force-recreate
```

### Step 2: Run Migrations (When Needed)

**When to run:**
- First deployment
- When new `.sql` files exist in `/migrations`
- After pulling code with database schema changes

**How to run:**

```bash
# Using docker compose with profiles
docker compose -f docker-compose.prod.yml run --rm migrate

# OR directly with docker exec (if container is running)
docker compose -f docker-compose.prod.yml exec go-backend-app ./migrate-tool up
```

**For Dokploy:**
```bash
# SSH into your Dokploy server, then:
docker exec <container-name> ./migrate-tool up
```

### Step 3: Run Seeding (When Needed)

**When to run:**
- First deployment (to populate resources, roles, endpoints, policies)
- When new resources/endpoints are added to the codebase
- After manual database cleanup

**How to run:**

```bash
# Using docker compose with profiles
docker compose -f docker-compose.prod.yml run --rm seed

# OR directly with docker exec (if container is running)
docker compose -f docker-compose.prod.yml exec go-backend-app ./seed-tool
```

**For Dokploy:**
```bash
# SSH into your Dokploy server, then:
docker exec <container-name> ./seed-tool
```

## Required Environment Variables

Ensure these are set in your Dokploy environment:

```bash
# Database Configuration (used by migration and seed tools)
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB_PROD=maindb

# Application Configuration
APP_ENV=prod
APP_PORT=4000

# JWT Configuration
JWT_SECRET=your_jwt_secret
JWT_REFRESH_SECRET=your_refresh_secret

# Other app-specific variables from your .env
```

## Important Notes

### 1. POSTGRES_DB_PROD is Critical

Both the migration and seeding tools use `POSTGRES_DB_PROD` as defined in your documentation:

- **Migration tool** (`migrate-tool`) reads `POSTGRES_DB_PROD`
- **Seeding tool** (`seed-tool`) reads `POSTGRES_DB_PROD`
- This ensures both tools target the production database

### 2. Idempotent Operations

- **Migrations**: golang-migrate tracks which migrations have been applied (safe to re-run)
- **Seeding**: Your `InitialSeed()` function checks for existing data before seeding (safe to re-run)
- **Manual control**: You decide when to run them, not automatic on every restart

### 3. Why Manual is Better

**Problems with automatic migrations/seeding on startup:**
- ❌ Runs on every container restart (unnecessary)
- ❌ Delays app startup time
- ❌ Wastes database connections
- ❌ No control over when migrations happen

**Benefits of manual migrations/seeding:**
- ✅ Run only when needed (new migrations exist)
- ✅ Fast app restarts (no migration overhead)
- ✅ Controlled timing (backup before migrate)
- ✅ Explicit operations (you know when they run)

### 4. Production Best Practice

In production, migrations should be:
1. **Reviewed** - Check SQL files before applying
2. **Backed up** - Database backup before migration
3. **Tested** - Run in staging first
4. **Monitored** - Watch migration execution
5. **Manual** - Triggered by you, not automatic

This is exactly what your current setup allows!

## Deployment Commands

### For Dokploy (Advanced > Run Command)

Use this command in Dokploy's "Advanced > Run Command" section:

```bash
docker compose build --no-cache && docker compose up -d --force-recreate
```

**After deployment, run migrations and seeding manually:**

1. SSH into your Dokploy server
2. Find your container: `docker ps | grep mynute`
3. Run migrations: `docker exec <container-name> ./migrate-tool up`
4. Run seeding: `docker exec <container-name> ./seed-tool`

### For Local Testing

Test the production setup locally:

```bash
# Deploy the stack
docker compose -f docker-compose.prod.yml up -d --build

# Run migrations (one-time or when needed)
docker compose -f docker-compose.prod.yml run --rm migrate

# Run seeding (one-time or when needed)
docker compose -f docker-compose.prod.yml run --rm seed
```

## Monitoring Deployment

### Check Application Logs

In Dokploy or locally:

```bash
docker logs -f <container-name>
```

The app should start immediately without migration/seeding output.

### Verify Migrations Applied

After running the migration command, connect to your database:

```bash
docker exec -it <postgres-container> psql -U <user> -d maindb
```

Check the schema_migrations table:

```sql
SELECT * FROM schema_migrations;
```

### Verify Seeding Completed

After running the seed command, check that resources exist:

```sql
SELECT * FROM resources;
SELECT * FROM roles;
SELECT * FROM endpoints;
```

## Troubleshooting

### Migrations Not Running

**Check if migrations exist:**
```bash
docker exec <container> ls /mynute-go/migrations
```

**Run migrations manually:**
```bash
docker exec <container> ./migrate-tool up
```

**Check migration status:**
```bash
docker exec <container> ./migrate-tool version
```

### Seeding Not Running

**Run seeding manually:**
```bash
docker exec <container> ./seed-tool
```

**Check seed tool exists:**
```bash
docker exec <container> ls /mynute-go/seed-tool
```

**Verify environment variables:**
```bash
docker exec <container> env | grep POSTGRES
```

### App Won't Start

**Check logs:**
```bash
docker logs <container-name>
```

**Common issues:**
1. Missing environment variables (check `POSTGRES_DB_PROD`, `JWT_SECRET`, etc.)
2. Database not ready (check postgres container health)
3. Port conflicts (verify `APP_PORT` is available)

### Database Connection Issues

**Test database connectivity:**
```bash
docker exec <app-container> nc -zv postgres 5432
```

**Check database exists:**
```bash
docker exec <postgres-container> psql -U <user> -l
```

## Files Modified

1. **Dockerfile**
   - Added migration tool build
   - Added seed tool build
   - Copied migration files
   - App starts immediately (no automatic migrations)

2. **docker-compose.prod.yml**
   - Added separate `migrate` service (manual execution)
   - Added separate `seed` service (manual execution)
   - Services use `profiles: [tools]` so they don't run by default
   - Main app starts without migration delays

3. **docker-entrypoint.sh** 
   - Removed (not needed - migrations/seeding are manual)

## Next Steps After Deployment

1. **Run migrations** - Execute migrations for first-time setup or after pulling schema changes
2. **Run seeding** - Populate resources, roles, endpoints, and policies
3. **Verify data** - Check that expected tables and data exist
4. **Test endpoints** - Ensure your API is functioning correctly
5. **Monitor logs** - Watch for any errors or issues
6. **Set up backups** - Follow backup procedures in MIGRATIONS.md before future updates

## Complete First Deployment Example

```bash
# 1. Deploy the stack
docker compose -f docker-compose.prod.yml up -d --build

# 2. Wait for database to be healthy (check logs)
docker compose -f docker-compose.prod.yml logs postgres

# 3. Run migrations (first time)
docker compose -f docker-compose.prod.yml run --rm migrate

# 4. Run seeding (first time)
docker compose -f docker-compose.prod.yml run --rm seed

# 5. Check app is running
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs go-backend-app

# 6. Test your API
curl http://localhost:4000/health
```

## Subsequent Deployments (Updates)

```bash
# 1. Pull latest code
git pull origin main

# 2. Check if new migrations exist
ls migrations/

# 3. Backup database (if migrations present)
docker exec <postgres-container> pg_dump -U <user> <database> > backup.sql

# 4. Rebuild and restart app
docker compose -f docker-compose.prod.yml up -d --build

# 5. Run new migrations (only if present)
docker compose -f docker-compose.prod.yml run --rm migrate

# 6. Run seeding (only if new resources/endpoints added)
docker compose -f docker-compose.prod.yml run --rm seed

# 7. Verify app is working
curl http://localhost:4000/health
```

## Summary

✅ **Migrations are manual** - Run only when needed  
✅ **Seeding is manual** - Run only when needed  
✅ **Fast app restarts** - No automatic migration overhead  
✅ **Production best practice** - Controlled, explicit operations  
✅ **Idempotent operations** - Safe to re-run migrations and seeding  

Deploy with confidence! 🚀
