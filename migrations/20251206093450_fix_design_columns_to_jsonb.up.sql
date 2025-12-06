-- Fix design columns from TEXT to JSONB
-- The design columns in companies, branches, and services tables should store JSON data

-- Fix companies.design column
ALTER TABLE public.companies 
ALTER COLUMN "design" TYPE JSONB USING 
  CASE 
    WHEN "design" IS NULL OR "design" = '' THEN '{}'::JSONB
    ELSE "design"::JSONB
  END;

-- Fix branches.design column in all company schemas
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branches ALTER COLUMN "design" TYPE JSONB USING 
          CASE 
            WHEN "design" IS NULL OR "design" = '''' THEN ''''{}''''::JSONB
            ELSE "design"::JSONB
          END', schema_name);
    END LOOP;
END $$;

-- Fix services.design column in all company schemas
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.services ALTER COLUMN "design" TYPE JSONB USING 
          CASE 
            WHEN "design" IS NULL OR "design" = '''' THEN ''''{}''''::JSONB
            ELSE "design"::JSONB
          END', schema_name);
    END LOOP;
END $$;
