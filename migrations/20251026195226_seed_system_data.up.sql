-- Migration: seed_system_data
-- Created at: 20251026195226
-- Seeds initial system roles, resources, and other required data

-- ============================================
-- SEED SYSTEM ROLES
-- ============================================
INSERT INTO public.roles (id, name, description, company_id, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'Owner', 'Company Owner. Can access anything within the company''s scope.', NULL, NOW(), NOW()),
    (gen_random_uuid(), 'General Manager', 'Company General Manager. Can access anything within the company''s scope besides editing the company name and taxID; and deleting the company.', NULL, NOW(), NOW()),
    (gen_random_uuid(), 'Branch Manager', 'Company Branch Manager. Can access anything within the branch''s scope besides deleting, renaming and changing its address; Can also manage appointments in the branch.', NULL, NOW(), NOW()),
    (gen_random_uuid(), 'Branch Supervisor', 'Company Branch Supervisor. Can see anything within the branch''s scope but can''t change or delete anything related to branch services, employees and properties; Can also manage appointments in the branch.', NULL, NOW(), NOW())
ON CONFLICT (name, company_id) WHERE company_id IS NULL DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW();

-- ============================================
-- SEED RESOURCES
-- ============================================
INSERT INTO public.resources (id, name, description, "table", references, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'appointment', 'Appointment resource', 'appointments', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"appointment_id","request_ref":"path"},{"database_key":"id","request_key":"appointment_id","request_ref":"query"},{"database_key":"id","request_key":"appointment_id","request_ref":"body"},{"database_key":"name","request_key":"name","request_ref":"path"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'branch', 'Branch resource', 'branches', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"branch_id","request_ref":"path"},{"database_key":"id","request_key":"branch_id","request_ref":"query"},{"database_key":"id","request_key":"branch_id","request_ref":"body"},{"database_key":"name","request_key":"name","request_ref":"path"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'client', 'Client resource', 'public.clients', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"client_id","request_ref":"path"},{"database_key":"id","request_key":"client_id","request_ref":"query"},{"database_key":"id","request_key":"client_id","request_ref":"body"},{"database_key":"email","request_key":"email","request_ref":"path"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'company', 'Company resource', 'public.companies', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"company_id","request_ref":"path"},{"database_key":"id","request_key":"company_id","request_ref":"query"},{"database_key":"id","request_key":"company_id","request_ref":"body"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'employee', 'Employee resource', 'employees', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"employee_id","request_ref":"path"},{"database_key":"id","request_key":"employee_id","request_ref":"query"},{"database_key":"id","request_key":"employee_id","request_ref":"body"},{"database_key":"email","request_key":"email","request_ref":"path"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'holiday', 'Holiday resource', 'public.holidays', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"holiday_id","request_ref":"path"},{"database_key":"id","request_key":"holiday_id","request_ref":"query"},{"database_key":"id","request_key":"holiday_id","request_ref":"body"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'role', 'Role resource', 'public.roles', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"role_id","request_ref":"path"},{"database_key":"id","request_key":"role_id","request_ref":"query"},{"database_key":"id","request_key":"role_id","request_ref":"body"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'sector', 'Sector resource', 'public.sectors', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"sector_id","request_ref":"path"},{"database_key":"id","request_key":"sector_id","request_ref":"query"},{"database_key":"id","request_key":"sector_id","request_ref":"body"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'service', 'Service resource', 'services', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"service_id","request_ref":"path"},{"database_key":"id","request_key":"service_id","request_ref":"query"},{"database_key":"id","request_key":"service_id","request_ref":"body"}]'::jsonb, 
     NOW(), NOW()),
    (gen_random_uuid(), 'auth', 'Auth resource', 'auth', 
     '[{"database_key":"id","request_key":"id","request_ref":"query"},{"database_key":"id","request_key":"id","request_ref":"path"},{"database_key":"id","request_key":"auth_id","request_ref":"path"},{"database_key":"id","request_key":"auth_id","request_ref":"query"},{"database_key":"id","request_key":"auth_id","request_ref":"body"}]'::jsonb, 
     NOW(), NOW())
ON CONFLICT ("table") DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    references = EXCLUDED.references,
    updated_at = NOW();

-- Note: Endpoints and Policies should be seeded programmatically by your application
-- as they depend on dynamically generated routes and complex business logic

