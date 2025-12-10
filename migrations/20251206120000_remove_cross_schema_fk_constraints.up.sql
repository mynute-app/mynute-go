-- Remove cross-schema foreign key constraints
-- Generated at: 20251206120000
-- Reason: Cross-schema FK constraints cause issues with schema-per-tenant architecture

-- Drop FK constraint from roles.company_id to companies.id
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_roles_company' 
        AND table_schema = 'public' 
        AND table_name = 'roles'
    ) THEN
        ALTER TABLE public.roles DROP CONSTRAINT fk_roles_company;
    END IF;
END $$;

-- Drop FK constraint from services.company_id to companies.id (in all company schemas)
DO $$
DECLARE
    schema_name TEXT;
    constraint_exists BOOLEAN;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('
            SELECT EXISTS (
                SELECT 1 
                FROM information_schema.table_constraints 
                WHERE constraint_name = ''fk_services_company'' 
                AND table_schema = %L
                AND table_name = ''services''
            )', schema_name) INTO constraint_exists;
        
        IF constraint_exists THEN
            EXECUTE format('ALTER TABLE %I.services DROP CONSTRAINT fk_services_company', schema_name);
        END IF;
    END LOOP;
END $$;

-- Note: employee_roles join table should only have FK to employees, not to roles
-- This is handled in the application code (company.go MigrateSchema method)
