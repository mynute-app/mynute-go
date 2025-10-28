-- Auto-generated DROP TABLE migration (rollback)
-- ⚠️  REVIEW THIS SQL BEFORE APPLYING!
-- ⚠️  THIS WILL DELETE ALL DATA!

-- Drop tenant schema function
DROP FUNCTION IF EXISTS create_tenant_schema_tables(TEXT);

-- Drop multi-tenant tables (must be done per schema)
-- Note: Add DROP statements for each company_* schema

-- Drop tenant tables in public schema
DROP TABLE IF EXISTS public.payments CASCADE;
DROP TABLE IF EXISTS public.employee_service_densities CASCADE;
DROP TABLE IF EXISTS public.branch_service_densities CASCADE;
DROP TABLE IF EXISTS public.employee_work_ranges CASCADE;
DROP TABLE IF EXISTS public.branch_work_ranges CASCADE;

-- Drop public schema tables
DROP TABLE IF EXISTS public.clientappointments CASCADE;
DROP TABLE IF EXISTS public.subdomains CASCADE;
DROP TABLE IF EXISTS public.properties CASCADE;
DROP TABLE IF EXISTS public.resources CASCADE;
DROP TABLE IF EXISTS public.policy_rules CASCADE;
DROP TABLE IF EXISTS public.roles CASCADE;
DROP TABLE IF EXISTS public.endpoints CASCADE;
DROP TABLE IF EXISTS public.clients CASCADE;
DROP TABLE IF EXISTS public.holidays CASCADE;
DROP TABLE IF EXISTS public.companies CASCADE;
DROP TABLE IF EXISTS public.sectors CASCADE;
