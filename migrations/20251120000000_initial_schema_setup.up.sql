-- Migration: initial_schema_setup
-- Created at: 20251120000000
-- Description: Initial database schema setup for Mynute Go application
--              Creates all tables in public schema for multi-tenant architecture

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ======================
-- PUBLIC SCHEMA TABLES
-- ======================

-- Companies table (multi-tenant root)
CREATE TABLE IF NOT EXISTS public.companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(255) UNIQUE NOT NULL,
    schema_name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT companies_subdomain_check CHECK (subdomain ~ '^[a-z0-9-]+$')
);

CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON public.companies(deleted_at);
CREATE INDEX IF NOT EXISTS idx_companies_subdomain ON public.companies(subdomain) WHERE deleted_at IS NULL;

-- Subdomains table
CREATE TABLE IF NOT EXISTS public.subdomains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subdomain VARCHAR(255) UNIQUE NOT NULL,
    company_id UUID NOT NULL REFERENCES public.companies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_subdomains_deleted_at ON public.subdomains(deleted_at);
CREATE INDEX IF NOT EXISTS idx_subdomains_company_id ON public.subdomains(company_id);

-- Sectors table
CREATE TABLE IF NOT EXISTS public.sectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_sectors_deleted_at ON public.sectors(deleted_at);

-- Holidays table
CREATE TABLE IF NOT EXISTS public.holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_holidays_deleted_at ON public.holidays(deleted_at);
CREATE INDEX IF NOT EXISTS idx_holidays_date ON public.holidays(date);

-- Clients table (shared across companies)
CREATE TABLE IF NOT EXISTS public.clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    meta JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON public.clients(deleted_at);
CREATE INDEX IF NOT EXISTS idx_clients_email ON public.clients(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clients_phone ON public.clients(phone) WHERE deleted_at IS NULL;

-- Client Appointments (bridge table)
CREATE TABLE IF NOT EXISTS public.client_appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES public.clients(id) ON DELETE CASCADE,
    appointment_id UUID NOT NULL,
    company_id UUID NOT NULL REFERENCES public.companies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_client_appointments_deleted_at ON public.client_appointments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_client_appointments_client_id ON public.client_appointments(client_id);
CREATE INDEX IF NOT EXISTS idx_client_appointments_appointment_id ON public.client_appointments(appointment_id);
CREATE INDEX IF NOT EXISTS idx_client_appointments_company_id ON public.client_appointments(company_id);

-- Resources table (RBAC)
CREATE TABLE IF NOT EXISTS public.resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    "table" VARCHAR(255),
    "references" JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_resources_deleted_at ON public.resources(deleted_at);
CREATE INDEX IF NOT EXISTS idx_resources_name ON public.resources(name) WHERE deleted_at IS NULL;

-- Roles table (RBAC)
CREATE TABLE IF NOT EXISTS public.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON public.roles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_roles_name ON public.roles(name) WHERE deleted_at IS NULL;

-- Endpoints table (RBAC - API routes)
CREATE TABLE IF NOT EXISTS public.endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    controller_name VARCHAR(255) NOT NULL,
    description TEXT,
    resource_id UUID REFERENCES public.resources(id) ON DELETE SET NULL,
    needs_company_id BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT endpoints_method_check CHECK (method IN ('GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD'))
);

CREATE INDEX IF NOT EXISTS idx_endpoints_deleted_at ON public.endpoints(deleted_at);
CREATE INDEX IF NOT EXISTS idx_endpoints_path_method ON public.endpoints(path, method) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_endpoints_resource_id ON public.endpoints(resource_id);
CREATE INDEX IF NOT EXISTS idx_endpoints_controller_name ON public.endpoints(controller_name);

-- Policy Rules table (RBAC - authorization rules)
CREATE TABLE IF NOT EXISTS public.policy_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES public.roles(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES public.resources(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    effect VARCHAR(10) NOT NULL DEFAULT 'allow',
    conditions JSONB,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT policy_rules_effect_check CHECK (effect IN ('allow', 'deny')),
    CONSTRAINT policy_rules_action_check CHECK (action IN ('create', 'read', 'update', 'delete', '*'))
);

CREATE INDEX IF NOT EXISTS idx_policy_rules_deleted_at ON public.policy_rules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_policy_rules_role_id ON public.policy_rules(role_id);
CREATE INDEX IF NOT EXISTS idx_policy_rules_resource_id ON public.policy_rules(resource_id);
CREATE INDEX IF NOT EXISTS idx_policy_rules_action ON public.policy_rules(action);

-- Properties table (key-value configuration)
CREATE TABLE IF NOT EXISTS public.properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_properties_deleted_at ON public.properties(deleted_at);
CREATE INDEX IF NOT EXISTS idx_properties_key ON public.properties(key) WHERE deleted_at IS NULL;

-- ======================
-- TENANT SCHEMA TABLES
-- ======================
-- Note: These tables will be created dynamically in each company_{uuid} schema
-- This is documented here for reference. Actual creation happens via application code.

-- Services table (per company schema)
-- CREATE TABLE IF NOT EXISTS services (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(255) NOT NULL,
--     description TEXT,
--     duration INTEGER NOT NULL,
--     price DECIMAL(10, 2),
--     color VARCHAR(7),
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP WITH TIME ZONE
-- );

-- Branches table (per company schema)
-- CREATE TABLE IF NOT EXISTS branches (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(255) NOT NULL,
--     address TEXT,
--     phone VARCHAR(50),
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP WITH TIME ZONE
-- );

-- Employees table (per company schema)
-- CREATE TABLE IF NOT EXISTS employees (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(255) NOT NULL,
--     email VARCHAR(255),
--     phone VARCHAR(50),
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP WITH TIME ZONE
-- );

-- Appointments table (per company schema)
-- CREATE TABLE IF NOT EXISTS appointments (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     client_id UUID NOT NULL,
--     employee_id UUID NOT NULL,
--     service_id UUID NOT NULL,
--     branch_id UUID NOT NULL,
--     start_time TIMESTAMP WITH TIME ZONE NOT NULL,
--     end_time TIMESTAMP WITH TIME ZONE NOT NULL,
--     status VARCHAR(50) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP WITH TIME ZONE
-- );

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
