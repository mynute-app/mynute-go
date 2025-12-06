-- Revert design columns from JSONB back to TEXT

-- Revert companies.design column
ALTER TABLE public.companies 
ALTER COLUMN "design" TYPE TEXT USING "design"::TEXT;

-- Revert branches.design column in all company schemas
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branches ALTER COLUMN "design" TYPE TEXT USING "design"::TEXT', schema_name);
    END LOOP;
END $$;

-- Revert services.design column in all company schemas
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.services ALTER COLUMN "design" TYPE TEXT USING "design"::TEXT', schema_name);
    END LOOP;
END $$;
