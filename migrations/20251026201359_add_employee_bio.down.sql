-- Auto-generated rollback migration
-- Generated at: 20251026201359
-- ⚠️  REVIEW THIS SQL BEFORE APPLYING!

-- Rollback Model: Employee (Schema: company)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        -- Add your rollback statements here for employees
        -- Example: EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS new_field', schema_name);
    END LOOP;
END $$;

