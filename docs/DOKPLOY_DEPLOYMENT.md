# Dokploy Deployment Guide

## Overview

This project uses **Docker Compose** for deployment, which includes:

✅ **PostgreSQL database** - Managed within the compose stack  
✅ **Backend application** - Go application with migration and seeding tools  
✅ **Monitoring stack** - Prometheus, Grafana, Loki  
✅ **Migration tool** - Built and included in the image  
✅ **Seeding tool** - Built and included in the image  
✅ **Migration files** - Copied to `/mynute-go/migrations`  
✅ **Separate execution** - Run migrations/seeding explicitly when needed  
✅ **No automatic runs** - App starts immediately without migration delays  

## Important: Use Docker Compose, Not Dockerfile

⚠️ **Critical:** In Dokploy, configure your project to use **Docker Compose**, not just the Dockerfile.

**Why?**
- The full stack (database + app + monitoring) is defined in `docker-compose.prod.yml`
- Using Dockerfile only will deploy the app without the database
- Your app will fail with "no such host" errors if the database isn't running
- All services need to be on the same Docker network

**In Dokploy:**
1. Set **Source Type** to "Docker Compose"
2. Point to `docker-compose.prod.yml`
3. All services will be deployed together on the same network  

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

## Dokploy Configuration

### Step 1: Create Application in Dokploy

1. **Project Settings:**
   - Name: `mynute-backend` (or your preference)
   - Repository: `github.com/mynute-app/mynute-go.git`
   - Branch: `main`

2. **Build Configuration:**
   - **Source Type:** `Docker Compose` ⚠️ (NOT "Dockerfile")
   - **Compose File:** `docker-compose.prod.yml`
   - **Compose Path:** Leave empty (file is in root)

3. **Advanced Settings > Run Command:**
   ```bash
   docker compose -f docker-compose.prod.yml build --no-cache && docker compose -f docker-compose.prod.yml up -d --force-recreate
   ```

4. **Environment Variables:**
   Add all variables from your `.env` file (see Required Environment Variables section below)

### Step 2: Initial Deployment

Deploy the stack (database + app + monitoring):

**In Dokploy:**
- Click "Deploy" button
- Wait for build to complete
- Check logs to ensure all services started

**Locally (for testing):**
```bash
docker compose -f docker-compose.prod.yml up -d --build
```

### Step 3: Run Migrations (When Needed)

**When to run:**
- First deployment
- When new `.sql` files exist in `/migrations`
- After pulling code with database schema changes

**How to run in Dokploy:**

1. **Find your project name:**
   ```bash
   # SSH into Dokploy server
   docker ps | grep mynute
   ```
   
   Look for the container with your app name (e.g., `prod-backend-fai3hk-go-backend-app-1`)

2. **Run migrations:**
   ```bash
   # Find the compose project directory
   cd /etc/dokploy/applications/<your-project-id>/
   
   # Run migrations using docker compose
   docker compose -f docker-compose.prod.yml run --rm migrate
   
   # OR directly with docker exec
   docker exec <container-name> ./migrate-tool up
   ```

**Locally (for testing):**
```bash
# Using docker compose with profiles
docker compose -f docker-compose.prod.yml run --rm migrate

# OR directly with docker exec
docker compose -f docker-compose.prod.yml exec go-backend-app ./migrate-tool up
```

### Step 4: Run Seeding (When Needed)

**When to run:**
- First deployment (to populate resources, roles, endpoints, policies)
- When new resources/endpoints are added to the codebase
- After manual database cleanup

**How to run in Dokploy:**

1. **SSH into Dokploy server and navigate to project:**
   ```bash
   cd /etc/dokploy/applications/<your-project-id>/
   ```

2. **Run seeding:**
   ```bash
   # Using docker compose with profiles
   docker compose -f docker-compose.prod.yml run --rm seed
   
   # OR directly with docker exec
   docker exec <container-name> ./seed-tool
   ```

**Locally (for testing):**
```bash
# Using docker compose with profiles
docker compose -f docker-compose.prod.yml run --rm seed

# OR directly with docker exec
docker compose -f docker-compose.prod.yml exec go-backend-app ./seed-tool
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

## Quick Reference Commands

### Dokploy Deployment Commands

**Initial setup in Dokploy UI:**
- Set Source Type: `Docker Compose`
- Compose File: `docker-compose.prod.yml`
- Advanced > Run Command:
  ```bash
  docker compose -f docker-compose.prod.yml build --no-cache && docker compose -f docker-compose.prod.yml up -d --force-recreate
  ```

**After deployment (SSH into Dokploy server):**
```bash
# Navigate to your project
cd /etc/dokploy/applications/<your-project-id>/

# Check running containers
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs go-backend-app
docker compose -f docker-compose.prod.yml logs postgres

# Run migrations (first time or when needed)
docker compose -f docker-compose.prod.yml run --rm migrate

# Run seeding (first time or when needed)
docker compose -f docker-compose.prod.yml run --rm seed

# Restart services
docker compose -f docker-compose.prod.yml restart go-backend-app
```

### Local Testing Commands

```bash
# Deploy the full stack
docker compose -f docker-compose.prod.yml up -d --build

# Check services
docker compose -f docker-compose.prod.yml ps

# Run migrations
docker compose -f docker-compose.prod.yml run --rm migrate

# Run seeding
docker compose -f docker-compose.prod.yml run --rm seed

# View logs
docker compose -f docker-compose.prod.yml logs -f go-backend-app

# Stop everything
docker compose -f docker-compose.prod.yml down
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

### ❌ Error: "no such host" when connecting to postgres

**Error message:**
```
[error] failed to initialize database, got error failed to connect to `host=postgres user=postgres database=maindb`: 
hostname resolving error (lookup postgres on 127.0.0.11:53: no such host)
```

**Cause:** Dokploy is deploying only the Dockerfile, not the full docker-compose stack.

**Solution:**

1. **In Dokploy UI, verify Source Type:**
   - Should be: `Docker Compose` ✅
   - NOT: `Dockerfile` ❌

2. **Verify Compose File setting:**
   - Compose File: `docker-compose.prod.yml`
   - Compose Path: (leave empty)

3. **Redeploy the application:**
   - This will deploy ALL services (postgres + app + monitoring)
   - All containers will be on the same network
   - The app can resolve `postgres` hostname

4. **Verify all services are running:**
   ```bash
   # SSH to Dokploy server
   cd /etc/dokploy/applications/<your-project-id>/
   docker compose -f docker-compose.prod.yml ps
   
   # You should see:
   # - postgres (healthy)
   # - go-backend-app (running)
   # - grafana (healthy)
   # - loki (running)
   # - prometheus (running)
   ```

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

### In Dokploy UI:

1. **Create New Application**
   - Set Source Type to `Docker Compose`
   - Point to `docker-compose.prod.yml`
   - Configure environment variables

2. **Deploy**
   - Click "Deploy" button
   - Wait for build and deployment

3. **Post-Deployment (SSH to server):**
   ```bash
   # Navigate to project
   cd /etc/dokploy/applications/<your-project-id>/
   
   # Wait for database to be healthy
   docker compose -f docker-compose.prod.yml logs postgres
   
   # Run migrations (first time)
   docker compose -f docker-compose.prod.yml run --rm migrate
   
   # Run seeding (first time)
   docker compose -f docker-compose.prod.yml run --rm seed
   
   # Check app is running
   docker compose -f docker-compose.prod.yml ps
   docker compose -f docker-compose.prod.yml logs go-backend-app
   
   # Test your API (from server or use your domain)
   curl http://localhost:4000/health
   ```

### Locally:

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
