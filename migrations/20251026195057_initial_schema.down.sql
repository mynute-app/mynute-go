-- Migration: initial_schema
-- Created at: 20251026195057

-- Rollback migration - drops all tables and functions

-- Drop the function for creating company schemas
DROP FUNCTION IF EXISTS create_company_schema(TEXT);

-- Drop all public schema tables in reverse order of dependencies
DROP TABLE IF EXISTS public.policy_rules CASCADE;
DROP TABLE IF EXISTS public.end_points CASCADE;
DROP TABLE IF EXISTS public.properties CASCADE;
DROP TABLE IF EXISTS public.resources CASCADE;
DROP TABLE IF EXISTS public.employee_roles CASCADE;
DROP TABLE IF EXISTS public.holidays CASCADE;
DROP TABLE IF EXISTS public.client_appointments CASCADE;
DROP TABLE IF EXISTS public.clients CASCADE;
DROP TABLE IF EXISTS public.subdomains CASCADE;
DROP TABLE IF EXISTS public.company_sectors CASCADE;
DROP TABLE IF EXISTS public.companies CASCADE;
DROP TABLE IF EXISTS public.sectors CASCADE;
DROP TABLE IF EXISTS public.roles CASCADE;

-- Note: Company-specific schemas are NOT automatically dropped
-- This is a safety measure to prevent accidental data loss
-- To drop a company schema, run: DROP SCHEMA IF EXISTS company_<uuid> CASCADE;

