# Migration Database Configuration

## CRITICAL CHANGE

**Migration tools now ALWAYS use `POSTGRES_DB_PROD` environment variable.**

This ensures migrations are explicit about which database they target, preventing accidental migrations to the wrong database.

## How It Works Now

### Migration Tools Behavior
All migration-related tools use **only** `POSTGRES_DB_PROD`:
- `make migrate-up`
- `make migrate-down`
- `make migrate-smart`
- `go run migrate/main.go`
- `go run tools/smart-migration/main.go`
- `go run cmd/seed/main.go`

### Application Runtime Behavior
The application (`go run main.go`) uses `APP_ENV` to determine database:
- `APP_ENV=dev` → Uses `POSTGRES_DB_DEV`
- `APP_ENV=test` → Uses `POSTGRES_DB_TEST`
- `APP_ENV=prod` → Uses `POSTGRES_DB_PROD`

## Configuration Examples

### Development Environment (.env)
```env
APP_ENV=dev
POSTGRES_DB_PROD=devdb           # Migration tools target devdb
POSTGRES_DB_DEV=devdb       # App uses this when APP_ENV=dev
POSTGRES_DB_TEST=testdb     # App uses this when APP_ENV=test
```

**Result:**
- `go run main.go` → Connects to **devdb** (via APP_ENV)
- `make migrate-up` → Migrates **devdb** (via POSTGRES_DB_PROD)
- `make migrate-smart` → Checks **devdb** (via POSTGRES_DB_PROD)

### Production Environment (.env)
```env
APP_ENV=prod
POSTGRES_DB_PROD=maindb          # Migration tools target maindb
POSTGRES_DB_DEV=devdb       # Not used in production
POSTGRES_DB_TEST=testdb     # Not used in production
```

**Result:**
- `go run main.go` → Connects to **maindb** (via APP_ENV)
- `make migrate-up` → Migrates **maindb** (via POSTGRES_DB_PROD)
- `make migrate-smart` → Checks **maindb** (via POSTGRES_DB_PROD)

## Why This Change?

### Before (DANGEROUS! ❌)
```bash
# In production .env
APP_ENV=prod
POSTGRES_DB_PROD=maindb

# Developer runs migration thinking they're on dev
make migrate-up
# ❌ Migrates PRODUCTION because APP_ENV=prod!
```

### After (SAFE! ✅)
```bash
# In production .env
APP_ENV=prod
POSTGRES_DB_PROD=maindb  # Explicitly set to maindb

# Developer sees clearly which database will be migrated
make migrate-up
# ✅ Migrates maindb because POSTGRES_DB_PROD=maindb
```

## Workflow Examples

### Development Workflow
```bash
# Set up .env
APP_ENV=dev
POSTGRES_DB_PROD=devdb

# Run migrations on dev database
make migrate-up          # Targets devdb

# Check drift on dev database
make migrate-smart NAME=check_drift MODELS=all  # Checks devdb

# Run application
go run main.go           # Connects to devdb
```

### Testing Workflow
```bash
# Set up .env for test
APP_ENV=test
POSTGRES_DB_PROD=testdb  # Point migration tools to test DB

# Run migrations on test database
make migrate-up          # Targets testdb

# Run tests
go test ./...           # App connects to testdb via APP_ENV
```

### Production Deployment
```bash
# Production .env
APP_ENV=prod
POSTGRES_DB_PROD=maindb  # ALWAYS explicit!

# Run migrations BEFORE deploying new code
make migrate-up          # Targets maindb

# Run seeding
./seed                   # Targets maindb

# Deploy application
./mynute-go             # Connects to maindb
```

## Migration Safety Checklist

Before running migrations:

1. ✅ Check `POSTGRES_DB_PROD` value
   ```bash
   echo $POSTGRES_DB_PROD  # Linux/Mac
   echo $env:POSTGRES_DB_PROD  # Windows PowerShell
   ```

2. ✅ Verify you're targeting the right database
   ```bash
   psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB_PROD -c "SELECT current_database();"
   ```

3. ✅ For production, ALWAYS double-check
   ```bash
   # Should output 'maindb' (or your production DB name)
   grep POSTGRES_DB_PROD .env
   ```

4. ✅ Run migration
   ```bash
   make migrate-up
   ```

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run Migrations
  env:
    POSTGRES_DB_PROD: ${{ secrets.PROD_DB_NAME }}  # Explicit!
    POSTGRES_HOST: ${{ secrets.PROD_DB_HOST }}
    POSTGRES_USER: ${{ secrets.PROD_DB_USER }}
    POSTGRES_PASSWORD: ${{ secrets.PROD_DB_PASSWORD }}
  run: make migrate-up
```

### Docker Compose Example
```yaml
services:
  migrations:
    image: mynute-go:latest
    environment:
      - POSTGRES_DB_PROD=maindb  # Explicit!
      - POSTGRES_HOST=postgres
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    command: ["./migrate", "-action=up", "-path=./migrations"]
```

## Breaking Change Notice

⚠️ **If you previously relied on APP_ENV to switch databases for migrations**, update your workflow:

**Old way (no longer works):**
```bash
APP_ENV=dev make migrate-up  # Used to migrate devdb
```

**New way:**
```bash
POSTGRES_DB_PROD=devdb make migrate-up  # Explicitly target devdb
```

Or update your `.env` file:
```env
POSTGRES_DB_PROD=devdb  # Set this explicitly
```

## Benefits

1. ✅ **Explicit** - Always know which database you're migrating
2. ✅ **Safe** - Can't accidentally migrate production
3. ✅ **Flexible** - Easy to target any database
4. ✅ **Consistent** - Same behavior across all migration tools
5. ✅ **Traceable** - CI/CD logs show exact database targeted

## Questions?

- **Q: How do I migrate my dev database?**
  - A: Set `POSTGRES_DB_PROD=devdb` in your `.env`

- **Q: How do I migrate production?**
  - A: Set `POSTGRES_DB_PROD=maindb` in production `.env`

- **Q: Can I override via command line?**
  - A: Yes! `POSTGRES_DB_PROD=mydb make migrate-up`

- **Q: What if POSTGRES_DB_PROD is not set?**
  - A: Migration tools will fail with clear error message
