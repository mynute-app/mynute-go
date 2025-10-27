# 🚀 Automated Migration Workflow Guide

## Overview

This guide shows you three ways to create database migrations, from fully automated to fully manual.

---

## 🧠 Method 1: Smart Migration (RECOMMENDED)

**Best for:** Adding/removing columns, changing types, modifying constraints

### How it works:
1. Compares your GORM models with the actual database schema
2. Detects differences automatically
3. Generates complete SQL with correct types

### Example: Adding a `bio` field to Employee

#### Step 1: Modify the Model

```go
// core/src/config/db/model/employee.go
type Employee struct {
    BaseModel
    Name    string `gorm:"type:varchar(100)"`
    Email   string `gorm:"type:varchar(100)"`
    Bio     string `gorm:"type:text"` // ← NEW FIELD
}
```

#### Step 2: Generate Smart Migration

```powershell
make migrate-smart NAME=add_employee_bio MODELS=Employee
```

**Output:**
```
📊 Using schema 'company_abc123' for comparison
✅ Generated smart migration files:
  migrations\20251026203000_add_employee_bio.up.sql
  migrations\20251026203000_add_employee_bio.down.sql

💡 Changes detected and SQL generated automatically!
```

#### Step 3: Review Generated SQL

**UP migration (already complete!):**
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

**DOWN migration:**
```sql
-- Removing 1 column(s) that were added
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

#### Step 4: Test

```powershell
make test-migrate
```

**Done!** ✅

---

## 📝 Method 2: Template Migration

**Best for:** When you want more control but still want a starting template

### Example: Adding a field with custom logic

#### Step 1: Modify Model (same as above)

#### Step 2: Generate Template

```powershell
make migrate-generate NAME=add_employee_bio MODELS=Employee
```

#### Step 3: Edit the Generated Template

The template comes with examples:

```sql
-- Model: Employee (Schema: company)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        -- Add your ALTER TABLE statements here for employees
        -- Example: EXECUTE format('ALTER TABLE %I.employees ADD COLUMN new_field TEXT', schema_name);
    END LOOP;
END $$;
```

**You edit it to:**

```sql
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS bio TEXT', schema_name);
        -- Add custom logic here if needed
    END LOOP;
END $$;
```

#### Step 4: Test

```powershell
make test-migrate
```

---

## ✍️ Method 3: Manual Migration

**Best for:** Complex migrations, data migrations, custom SQL

#### Step 1: Create Empty Files

```powershell
make migrate-create NAME=complex_data_migration
```

#### Step 2: Write SQL Manually

Write everything from scratch based on your needs.

#### Step 3: Test

```powershell
make test-migrate
```

---

## 📊 Comparison

| Feature | Smart | Template | Manual |
|---------|-------|----------|--------|
| Speed | ⚡⚡⚡ Fastest | ⚡⚡ Fast | ⚡ Slow |
| Auto-detection | ✅ Yes | ❌ No | ❌ No |
| Control | ⚠️ Limited | ✅ Medium | ✅ Full |
| DB Connection Required | ✅ Yes | ❌ No | ❌ No |
| Best for | Simple schema changes | Moderate changes | Complex logic |

---

## 🎯 Decision Tree

```
Need to create a migration?
│
├─ Simple schema change (add/remove column)?
│  └─> Use Smart Migration (migrate-smart)
│
├─ Need custom logic but basic structure?
│  └─> Use Template Migration (migrate-generate)
│
└─ Complex data migration or custom SQL?
   └─> Use Manual Migration (migrate-create)
```

---

## 💡 Best Practices

### 1. Always Review Generated SQL

Even smart migrations should be reviewed:
```powershell
# After generating, open the files
code migrations/20251026203000_add_employee_bio.up.sql
```

### 2. Test Before Committing

Always run automated tests:
```powershell
make test-migrate
```

This runs:
- ✅ UP migration
- ✅ Verification
- ✅ DOWN migration (rollback)
- ✅ Verification
- ✅ UP migration again

### 3. Use Descriptive Names

```powershell
# Good
make migrate-smart NAME=add_employee_bio MODELS=Employee

# Bad
make migrate-smart NAME=update MODELS=Employee
```

### 4. One Logical Change Per Migration

```powershell
# Good - focused
make migrate-smart NAME=add_employee_bio MODELS=Employee
make migrate-smart NAME=add_branch_capacity MODELS=Branch

# Bad - mixing concerns
make migrate-smart NAME=various_updates MODELS=Employee,Branch,Service
```

---

## 🔧 Advanced Usage

### Multiple Models

```powershell
# Detect changes in multiple models at once
make migrate-smart NAME=update_contact_info MODELS=Employee,Client,Branch
```

### Specific Schema

```powershell
# Use specific company schema for comparison
go run tools/smart-migration/main.go -name test -models Employee -schema company_specific_uuid
```

### Template for All Models

```powershell
# Generate template showing all models (for reference)
make migrate-generate NAME=reference MODELS=all
```

---

## 🆘 Common Issues

### Issue: "No changes detected" with Smart Migration

**Cause:** Model matches database OR table doesn't exist yet

**Solution:**
```powershell
# For new tables, use initial schema migration
make migrate-up

# Or use template instead
make migrate-generate NAME=x MODELS=Y
```

### Issue: Wrong data type detected

**Cause:** GORM tag doesn't map perfectly to PostgreSQL

**Solution:**
```powershell
# Use template and specify exact type
make migrate-generate NAME=x MODELS=Y
# Then edit the SQL manually
```

### Issue: Need to rename column

**Cause:** Smart detection sees rename as drop + add

**Solution:**
```powershell
# Use template or manual migration
make migrate-generate NAME=rename_field MODELS=Employee

# Edit to use proper RENAME:
# ALTER TABLE employees RENAME COLUMN old_name TO new_name;
```

---

## 📋 Complete Example Workflow

### Scenario: Add bio and phone to Employee, capacity to Branch

```powershell
# 1. Modify models in Go
# (edit employee.go and branch.go)

# 2. Generate migrations
make migrate-smart NAME=add_employee_fields MODELS=Employee
make migrate-smart NAME=add_branch_capacity MODELS=Branch

# 3. Review generated files
code migrations/

# 4. Test both migrations
make test-migrate

# 5. Commit
git add migrations/
git commit -m "Add employee bio/phone and branch capacity"

# 6. In production (before deploying code)
APP_ENV=prod make migrate-up
```

**Total time: ~5 minutes for 2 migrations!** 🚀

---

## 🎓 Learn More

- **Main guide:** `MIGRATIONS_AUTOMATED.md`
- **Quick reference:** `MIGRATIONS_QUICKSTART.md`
- **Full documentation:** `docs/MIGRATIONS.md`
- **Production setup:** `MIGRATION_SETUP_SUMMARY.md`
