-- Rollback: Re-add cross-schema foreign key constraints
-- Generated at: 20251206120000

-- Re-add FK constraint from roles.company_id to companies.id
ALTER TABLE public.roles 
    ADD CONSTRAINT fk_roles_company 
    FOREIGN KEY (company_id) REFERENCES public.companies(id);

-- Re-add FK constraint from services.company_id to companies.id (in all company schemas)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('
            ALTER TABLE %I.services 
            ADD CONSTRAINT fk_services_company 
            FOREIGN KEY (company_id) REFERENCES public.companies(id)
        ', schema_name);
    END LOOP;
END $$;
