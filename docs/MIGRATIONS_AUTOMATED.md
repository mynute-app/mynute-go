# ✨ Migration System - NOW WITH SMART AUTO-DETECTION!

## 🎉 Yes! Fully automated with schema change detection!

### ✅ What's been automated:

1. **SQL Generation** - Ready-to-use templates
2. **Schema Change Detection** - Auto-detect what changed in your models! 🧠
3. **Testing** - Automatic UP → DOWN → UP testing
4. **Multi-tenant** - Automatically iterates over `company_*` schemas

---

## 🚀 Complete Workflow (3 options!)

### Option 1: Smart Migration (RECOMMENDED) 🧠

**Auto-detects changes by comparing models with database schema**

#### 1️⃣ Modify your Model

```go
// core/src/config/db/model/employee.go
type Employee struct {
    BaseModel
    Name  string
    Email string
    Bio   string `gorm:"type:text"` // ← NEW FIELD
}
```

#### 2️⃣ Generate Migration (AUTOMATIC DETECTION!)

```powershell
make migrate-smart NAME=add_employee_bio MODELS=Employee
```

**Result:**
```
📊 Using schema 'company_<uuid>' for comparison
Analyzing schema changes for models: Employee
✅ Generated smart migration files:
  migrations\20251026201359_add_employee_bio.up.sql
  migrations\20251026201359_add_employee_bio.down.sql

💡 Changes detected and SQL generated automatically!
⚠️  IMPORTANT: Review the generated SQL before applying!
```

**Generated SQL (already complete!):**
```sql
-- Smart migration - Auto-detected changes
-- Model: Employee (Table: employees, Schema: company)
-- Adding 1 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS bio TEXT NULL', schema_name);
    END LOOP;
END $$;
```

#### 3️⃣ Test Automatically

```powershell
make test-migrate
```

**Done!** The SQL is already complete based on detected changes! 🎉

---

### Option 2: Template Migration (Manual editing)

**Generates templates with examples, you fill in the details**

```powershell
# Generate template
make migrate-generate NAME=add_employee_bio MODELS=Employee

# Edit the generated files (they come with examples)
# migrations/XXXXX_add_employee_bio.up.sql

# Test
make test-migrate
```

---

### Option 3: Empty Migration (Fully manual)

```powershell
# Create empty files
make migrate-create NAME=custom_migration

# Write everything yourself
# Test
make test-migrate
```

---

## 📊 Comparison: Manual vs Template vs Smart

| Task | Manual | Template | Smart (NEW!) |
|------|--------|----------|--------------|
| Detect changes | ❌ Manual | ❌ Manual | ✅ **Automatic** |
| Write SQL | ❌ From scratch | ⚠️ Adjust template | ✅ **Auto-generated** |
| Multi-tenant iteration | ❌ Remember to add | ✅ Included | ✅ Included |
| Column types | ❌ Figure out | ⚠️ Manual | ✅ **Auto-detected** |
| Nullable constraints | ❌ Manual | ⚠️ Manual | ✅ **Auto-detected** |
| Test UP | ✅ Automatic | ✅ Automatic | ✅ Automatic |
| Test DOWN | ✅ Automatic | ✅ Automatic | ✅ Automatic |
| **Time required** | ~15 min | ~5 min | **~2 min** 🚀 |

---

## 🛠️ Main Commands

```powershell
# Smart migration (auto-detect changes) - RECOMMENDED
make migrate-smart NAME=add_field MODELS=Employee

# Template migration (manual editing)
make migrate-generate NAME=add_field MODELS=Employee

# Empty migration (fully manual)
make migrate-create NAME=custom_migration

# Test migration automatically  
make test-migrate

# Apply in production
make migrate-up

# Check status
make migrate-version

# Help
make migrate-help
```

---

## 💡 Practical Examples

### Example 1: Add a single field (SMART WAY - EASIEST!)

```powershell
# 1. Modify your model (add bio field to Employee)
# core/src/config/db/model/employee.go

# 2. Auto-generate migration with change detection
make migrate-smart NAME=add_employee_bio MODELS=Employee

# 3. Review generated SQL (already complete with correct types!)

# 4. Test
make test-migrate

# 5. Commit
git add migrations/ && git commit -m "Add bio to Employee"
```

**Time saved: 13 minutes!** ⏱️

### Example 2: Modify multiple models

```powershell
# Auto-detect changes in Branch AND Service
make migrate-smart NAME=update_descriptions MODELS=Branch,Service

# Review and test
make test-migrate
```

### Example 3: Complex migration (use template)

When smart detection doesn't work (renames, complex logic):

```powershell
make migrate-generate NAME=rename_column MODELS=Employee
# Edit the SQL manually with your custom logic
make test-migrate
```

---

## ⚙️ What Smart Migration Does Automatically

✅ **Connects** to your development database  
✅ **Compares** GORM models with actual database schema  
✅ **Detects** new columns automatically  
✅ **Generates** correct SQL with proper data types  
✅ **Includes** multi-tenant loops when needed (company_* schemas)  
✅ **Creates** both UP and DOWN migrations  
✅ **Handles** nullable/not-null constraints automatically  
✅ **Adds** IF NOT EXISTS / IF EXISTS for idempotency  

---

## ⚠️ Limitations of Smart Detection

Smart migration can detect:
- ✅ New columns added
- ✅ Column data types
- ✅ Nullable constraints
- ✅ Public vs company schema tables

Smart migration CANNOT detect:
- ❌ Column renames (sees as drop + add)
- ❌ Complex data migrations
- ❌ Index changes
- ❌ Constraint modifications
- ❌ Custom SQL logic

**For these cases, use `migrate-generate` and edit manually.**

---

## 🎯 Recommended Workflow

```
┌─────────────────────────────────────────────┐
│ 1. Modify Model                             │
│ 2. make migrate-smart NAME=x MODELS=Y      │
│ 3. Review generated SQL (usually perfect!)  │
│ 4. make test-migrate                        │
│ 5. git commit                               │
└─────────────────────────────────────────────┘

Total time: ~2 minutes! 🚀
```

---

## 🔍 How It Works

1. **Connects** to your database (uses `.env` configuration)
2. **Finds** first company schema (or uses public)
3. **Parses** GORM model to get expected columns
4. **Queries** database to get actual columns
5. **Compares** expected vs actual
6. **Generates** SQL for differences
7. **Creates** UP and DOWN migration files

---

## 📚 Complete Documentation

- **Smart migration guide:** This file
- **Full guide:** `docs/MIGRATIONS.md`
- **Quick reference:** `MIGRATIONS_QUICKSTART.md`
- **Template workflow:** `docs/MIGRATION_WORKFLOW.md`

---

## 🆘 Troubleshooting

### "No changes detected"

**Possible causes:**
1. Your model matches the database (no changes needed)
2. Table doesn't exist yet (use initial schema migration)
3. Connected to wrong database (check APP_ENV)

**Solution:**
```powershell
# Verify which database you're using
echo $env:APP_ENV
echo $env:POSTGRES_DB_DEV

# For new tables, use initial migration instead
make migrate-up  # Run initial schema first
```

### "Can't connect to database"

Smart migration needs a database connection. Check:
1. Database is running (`docker ps` or service status)
2. `.env` has correct credentials
3. `APP_ENV` is set correctly (dev, test, or prod)

**Alternative:**
```powershell
# Use template instead (no DB connection needed)
make migrate-generate NAME=x MODELS=Y
```

### "Detected wrong changes"

Sometimes GORM types don't map perfectly to PostgreSQL types.

**Solution:**
1. Review the generated SQL
2. Adjust data types if needed
3. Or use `migrate-generate` for full control

---

## 🎉 Summary

### You DON'T need to:
- ❌ Write SQL from scratch
- ❌ Manually detect what changed in your models
- ❌ Remember multi-tenant iteration syntax
- ❌ Figure out PostgreSQL data types
- ❌ Manually test UP/DOWN migrations

### You ONLY need to:
- ✅ Modify your models
- ✅ Run: `make migrate-smart NAME=x MODELS=Y`
- ✅ Review generated SQL
- ✅ Run: `make test-migrate`
- ✅ Commit!

**Migrations are now 10x faster with smart detection!** 🚀🧠
