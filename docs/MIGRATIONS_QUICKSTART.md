# Database Migrations - Quick Reference

## 🚀 For Production

### Before Deployment
```bash
# Run this BEFORE deploying your application
APP_ENV=prod make migrate-up

# Or using the script
.\scripts\migrate.ps1 up  # Windows
./scripts/migrate.sh up   # Linux/Mac
```

### After Deployment
```bash
# Verify migrations were applied
APP_ENV=prod make migrate-version
```

## 🛠️ For Development

Development and test environments automatically run migrations on startup. No manual intervention needed!

## 📖 Full Documentation

See [docs/MIGRATIONS.md](docs/MIGRATIONS.md) for complete guide.

## 🔄 Common Commands

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-smart NAME=x MODELS=Y` | Auto-detect changes and create migration |
| `make migrate-version` | Check current version |
| `make test-migrate` | Test migration automatically |
| `make migrate-help` | Show all commands |

## ⚠️ Important Notes

- **Production**: Migrations are NOT automatic - must run manually
- **Dev/Test**: Migrations run automatically on app startup
- **Always** test migrations in staging before production
- **Never** edit migrations after they've been applied in production
