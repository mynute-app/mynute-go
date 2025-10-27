-- Auto-generated migration
-- Generated at: 20251026201359
-- ⚠️  REVIEW THIS SQL BEFORE APPLYING!

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


-- 💡 Tips:
-- - Use 'IF NOT EXISTS' / 'IF EXISTS' for idempotency
-- - Add indexes with 'CREATE INDEX CONCURRENTLY' in production
-- - Test on a copy of production data first
