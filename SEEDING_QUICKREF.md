# Production Seeding Quick Reference

## The Problem
In production, automatic seeding is disabled:
- ❌ Endpoints not seeded automatically
- ❌ Policies not seeded automatically  
- ❌ Changes to `DenyUnauthorized` etc. won't update

## The Solution
Use the dedicated seed command!

## Commands

### Development
```bash
# Quick seed
make seed

# Or directly
go run cmd/seed/main.go
```

### Production

**Option 1: Pre-built Binary**
```bash
# Build once
make seed-build

# Deploy and run on server
./bin/seed
```

**Option 2: Run Directly**
```bash
# On production server
go run cmd/seed/main.go
```

## When to Seed

✅ **Run seeding after:**
- Initial deployment
- Changing endpoint permissions (`DenyUnauthorized`, `NeedsCompanyId`, etc.)
- Adding/removing endpoints
- Updating policy rules
- Changing system roles

❌ **Don't need to seed after:**
- Regular app restarts
- User/company data changes
- Non-route code changes

## What Gets Seeded

1. **Resources** - Table configurations
2. **System Roles** - Owner, Manager, etc.
3. **Endpoints** - All API routes + permissions
4. **Policies** - RBAC/ABAC access rules

## Safety

✅ **Idempotent** - Safe to run multiple times
✅ **Updates existing** - Won't create duplicates
✅ **Creates missing** - Adds new entries
✅ **No deletions** - Won't remove data

## Example Workflow

```bash
# 1. Change endpoint in code
vim core/src/config/db/model/endpoint.go
# Change: DenyUnauthorized: false

# 2. Deploy new code
git push production

# 3. Run seeding on production
ssh production
cd /app
./seed

# 4. Restart app (if needed)
systemctl restart mynute-go
```

## Help

```bash
make seed-help          # Seeding help
make migrate-help       # Migration help
```

## Full Documentation

See [docs/SEEDING.md](../docs/SEEDING.md) for complete guide.
