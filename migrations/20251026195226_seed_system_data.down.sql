-- Migration: seed_system_data
-- Created at: 20251026195226

-- Rollback seed data

-- Delete seeded resources
DELETE FROM public.resources WHERE name IN (
    'appointment', 'branch', 'client', 'company', 'employee', 
    'holiday', 'role', 'sector', 'service', 'auth'
);

-- Delete seeded system roles (only those without company_id)
DELETE FROM public.roles WHERE company_id IS NULL AND name IN (
    'Owner', 'General Manager', 'Branch Manager', 'Branch Supervisor'
);

