-- Migration: initial_schema_setup (DOWN)
-- Created at: 20251120000000
-- Description: Rollback initial database schema setup

-- Drop tables in reverse order (respecting foreign key constraints)

-- Drop RBAC tables
DROP TABLE IF EXISTS public.policy_rules CASCADE;
DROP TABLE IF EXISTS public.endpoints CASCADE;
DROP TABLE IF EXISTS public.roles CASCADE;
DROP TABLE IF EXISTS public.resources CASCADE;

-- Drop configuration
DROP TABLE IF EXISTS public.properties CASCADE;

-- Drop client tables
DROP TABLE IF EXISTS public.client_appointments CASCADE;
DROP TABLE IF EXISTS public.clients CASCADE;

-- Drop reference tables
DROP TABLE IF EXISTS public.holidays CASCADE;
DROP TABLE IF EXISTS public.sectors CASCADE;

-- Drop company tables
DROP TABLE IF EXISTS public.subdomains CASCADE;
DROP TABLE IF EXISTS public.companies CASCADE;

-- Note: Tenant schemas (company_*) should be dropped separately if needed
-- DROP SCHEMA IF EXISTS company_{uuid} CASCADE;
