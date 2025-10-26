-- Migration: initial_schema
-- Created at: 20251026195057
-- This migration creates all base tables for the application

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- PUBLIC SCHEMA TABLES
-- ============================================

-- Sectors table
CREATE TABLE IF NOT EXISTS public.sectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_sectors_deleted_at ON public.sectors(deleted_at);

-- Companies table
CREATE TABLE IF NOT EXISTS public.companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    legal_name VARCHAR(100) NOT NULL UNIQUE,
    trade_name VARCHAR(100) NOT NULL UNIQUE,
    tax_id VARCHAR(100) NOT NULL UNIQUE,
    schema_name VARCHAR(100) UNIQUE,
    design JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON public.companies(deleted_at);

-- Company-Sectors many-to-many junction table
CREATE TABLE IF NOT EXISTS public.company_sectors (
    company_id UUID NOT NULL REFERENCES public.companies(id) ON DELETE CASCADE,
    sector_id UUID NOT NULL REFERENCES public.sectors(id) ON DELETE CASCADE,
    PRIMARY KEY (company_id, sector_id)
);

-- Subdomains table
CREATE TABLE IF NOT EXISTS public.subdomains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(36) NOT NULL UNIQUE,
    company_id UUID NOT NULL REFERENCES public.companies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_subdomains_deleted_at ON public.subdomains(deleted_at);

-- Clients table
CREATE TABLE IF NOT EXISTS public.clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    meta JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON public.clients(deleted_at);

-- Client Appointments junction table (public schema)
CREATE TABLE IF NOT EXISTS public.client_appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES public.clients(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES public.companies(id) ON DELETE CASCADE,
    appointment_id UUID NOT NULL,
    schema_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_client_appointments_deleted_at ON public.client_appointments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_client_appointments_client_id ON public.client_appointments(client_id);

-- Holidays table
CREATE TABLE IF NOT EXISTS public.holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    date DATE NOT NULL,
    company_id UUID REFERENCES public.companies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_holidays_deleted_at ON public.holidays(deleted_at);
CREATE INDEX IF NOT EXISTS idx_holidays_company_id ON public.holidays(company_id);

-- Roles table
CREATE TABLE IF NOT EXISTS public.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    company_id UUID REFERENCES public.companies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_role_name_company ON public.roles(name, company_id);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON public.roles(deleted_at);

-- Resources table
CREATE TABLE IF NOT EXISTS public.resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "table" VARCHAR(100) NOT NULL UNIQUE,
    references JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_resources_deleted_at ON public.resources(deleted_at);

-- Properties table
CREATE TABLE IF NOT EXISTS public.properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    resource_id UUID NOT NULL REFERENCES public.resources(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_properties_deleted_at ON public.properties(deleted_at);
CREATE INDEX IF NOT EXISTS idx_properties_resource_id ON public.properties(resource_id);

-- Endpoints table
CREATE TABLE IF NOT EXISTS public.end_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    method VARCHAR(10) NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_endpoint_method_path ON public.end_points(method, path);
CREATE INDEX IF NOT EXISTS idx_endpoints_deleted_at ON public.end_points(deleted_at);

-- Policy Rules table
CREATE TABLE IF NOT EXISTS public.policy_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    role_id UUID NOT NULL REFERENCES public.roles(id) ON DELETE CASCADE,
    endpoint_id UUID NOT NULL REFERENCES public.end_points(id) ON DELETE CASCADE,
    resource_id UUID REFERENCES public.resources(id) ON DELETE CASCADE,
    property_id UUID REFERENCES public.properties(id) ON DELETE CASCADE,
    company_id UUID REFERENCES public.companies(id) ON DELETE CASCADE,
    created_by UUID,
    effect VARCHAR(10) NOT NULL CHECK (effect IN ('allow', 'deny')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_policy_rules_deleted_at ON public.policy_rules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_policy_rules_role_id ON public.policy_rules(role_id);
CREATE INDEX IF NOT EXISTS idx_policy_rules_endpoint_id ON public.policy_rules(endpoint_id);

-- ============================================
-- FUNCTION TO CREATE COMPANY-SPECIFIC SCHEMA
-- ============================================
-- This function will be called after a company is created to set up their tenant schema

CREATE OR REPLACE FUNCTION create_company_schema(schema_name TEXT) RETURNS VOID AS $$
BEGIN
    -- Create schema
    EXECUTE format('CREATE SCHEMA IF NOT EXISTS %I', schema_name);

    -- Branches table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.branches (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(100) NOT NULL,
            street VARCHAR(100) NOT NULL,
            number VARCHAR(100) NOT NULL,
            complement VARCHAR(100),
            neighborhood VARCHAR(100) NOT NULL,
            zip_code VARCHAR(100) NOT NULL,
            city VARCHAR(100) NOT NULL,
            state VARCHAR(100) NOT NULL,
            country VARCHAR(100) NOT NULL,
            company_id UUID NOT NULL,
            time_zone VARCHAR(100) NOT NULL,
            total_service_density INTEGER NOT NULL DEFAULT -1,
            design JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branches_deleted_at ON %I.branches(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branches_company_id ON %I.branches(company_id)', schema_name);

    -- Services table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.services (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(100) NOT NULL,
            description TEXT NOT NULL,
            price NUMERIC(10,2) NOT NULL,
            currency VARCHAR(3) DEFAULT ''BRL'',
            duration INTEGER NOT NULL,
            company_id UUID NOT NULL,
            design JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_services_deleted_at ON %I.services(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_services_company_id ON %I.services(company_id)', schema_name);

    -- Employees table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employees (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(100) NOT NULL,
            surname VARCHAR(100),
            email VARCHAR(100) NOT NULL UNIQUE,
            phone VARCHAR(20) NOT NULL UNIQUE,
            tags JSON,
            password VARCHAR(255) NOT NULL,
            slot_time_diff INTEGER DEFAULT 30,
            company_id UUID NOT NULL,
            time_zone VARCHAR(100) NOT NULL,
            total_service_density INTEGER NOT NULL DEFAULT 1,
            verified BOOLEAN DEFAULT FALSE,
            meta JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON %I.employees(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employees_company_id ON %I.employees(company_id)', schema_name);

    -- Employee-Branch junction table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employee_branches (
            employee_id UUID NOT NULL REFERENCES %I.employees(id) ON DELETE CASCADE,
            branch_id UUID NOT NULL REFERENCES %I.branches(id) ON DELETE CASCADE,
            PRIMARY KEY (employee_id, branch_id)
        )', schema_name, schema_name, schema_name);

    -- Employee-Service junction table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employee_services (
            employee_id UUID NOT NULL REFERENCES %I.employees(id) ON DELETE CASCADE,
            service_id UUID NOT NULL REFERENCES %I.services(id) ON DELETE CASCADE,
            PRIMARY KEY (employee_id, service_id)
        )', schema_name, schema_name, schema_name);

    -- Employee-Role junction table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employee_roles (
            employee_id UUID NOT NULL REFERENCES %I.employees(id) ON DELETE CASCADE,
            role_id UUID NOT NULL REFERENCES public.roles(id) ON DELETE CASCADE,
            PRIMARY KEY (employee_id, role_id)
        )', schema_name, schema_name);

    -- Branch-Service junction table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.branch_services (
            branch_id UUID NOT NULL REFERENCES %I.branches(id) ON DELETE CASCADE,
            service_id UUID NOT NULL REFERENCES %I.services(id) ON DELETE CASCADE,
            PRIMARY KEY (branch_id, service_id)
        )', schema_name, schema_name, schema_name);

    -- Branch Work Ranges
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.branch_work_ranges (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            weekday INTEGER NOT NULL,
            start_time TIMESTAMP WITH TIME ZONE NOT NULL,
            end_time TIMESTAMP WITH TIME ZONE NOT NULL,
            time_zone VARCHAR(255) NOT NULL,
            branch_id UUID NOT NULL REFERENCES %I.branches(id) ON DELETE CASCADE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name, schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branch_work_ranges_deleted_at ON %I.branch_work_ranges(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branch_work_ranges_branch_id ON %I.branch_work_ranges(branch_id)', schema_name);

    -- Employee Work Ranges
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employee_work_ranges (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            weekday INTEGER NOT NULL,
            start_time TIMESTAMP WITH TIME ZONE NOT NULL,
            end_time TIMESTAMP WITH TIME ZONE NOT NULL,
            time_zone VARCHAR(255) NOT NULL,
            employee_id UUID NOT NULL REFERENCES %I.employees(id) ON DELETE CASCADE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name, schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employee_work_ranges_deleted_at ON %I.employee_work_ranges(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employee_work_ranges_employee_id ON %I.employee_work_ranges(employee_id)', schema_name);

    -- Branch Service Density
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.branch_service_densities (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            branch_id UUID NOT NULL REFERENCES %I.branches(id) ON DELETE CASCADE,
            service_id UUID NOT NULL REFERENCES %I.services(id) ON DELETE CASCADE,
            max_schedules_overlap INTEGER NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name, schema_name, schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branch_service_densities_deleted_at ON %I.branch_service_densities(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branch_service_densities_branch_id ON %I.branch_service_densities(branch_id)', schema_name);

    -- Employee Service Density
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employee_service_densities (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            employee_id UUID NOT NULL REFERENCES %I.employees(id) ON DELETE CASCADE,
            service_id UUID NOT NULL REFERENCES %I.services(id) ON DELETE CASCADE,
            max_schedules_overlap INTEGER NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name, schema_name, schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employee_service_densities_deleted_at ON %I.employee_service_densities(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employee_service_densities_employee_id ON %I.employee_service_densities(employee_id)', schema_name);

    -- Payments table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.payments (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            amount NUMERIC(10,2) NOT NULL,
            currency VARCHAR(3) NOT NULL DEFAULT ''BRL'',
            payment_method VARCHAR(50) NOT NULL,
            status VARCHAR(20) NOT NULL,
            company_id UUID NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON %I.payments(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_payments_company_id ON %I.payments(company_id)', schema_name);

    -- Appointments table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.appointments (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            service_id UUID NOT NULL REFERENCES %I.services(id) ON DELETE CASCADE,
            employee_id UUID NOT NULL REFERENCES %I.employees(id) ON DELETE CASCADE,
            client_id UUID NOT NULL,
            branch_id UUID NOT NULL REFERENCES %I.branches(id) ON DELETE CASCADE,
            payment_id UUID UNIQUE REFERENCES %I.payments(id) ON DELETE CASCADE,
            company_id UUID NOT NULL,
            cancelled_employee_id UUID,
            start_time TIME NOT NULL,
            end_time TIME NOT NULL,
            time_zone VARCHAR(100) NOT NULL,
            actual_start_time TIME NOT NULL,
            actual_end_time TIME NOT NULL,
            cancel_time TIME NOT NULL,
            is_fulfilled BOOLEAN DEFAULT FALSE,
            is_cancelled BOOLEAN DEFAULT FALSE,
            is_cancelled_by_client BOOLEAN DEFAULT FALSE,
            is_cancelled_by_employee BOOLEAN DEFAULT FALSE,
            is_confirmed_by_client BOOLEAN DEFAULT FALSE,
            history JSONB,
            comments JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name, schema_name, schema_name, schema_name, schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON %I.appointments(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_client_id ON %I.appointments(client_id)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_company_id ON %I.appointments(company_id)', schema_name);

    -- Appointments Archive table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.appointments_archives (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            service_id UUID NOT NULL,
            employee_id UUID NOT NULL,
            client_id UUID NOT NULL,
            branch_id UUID NOT NULL,
            payment_id UUID UNIQUE,
            company_id UUID NOT NULL,
            cancelled_employee_id UUID,
            start_time TIME NOT NULL,
            end_time TIME NOT NULL,
            time_zone VARCHAR(100) NOT NULL,
            actual_start_time TIME NOT NULL,
            actual_end_time TIME NOT NULL,
            cancel_time TIME NOT NULL,
            is_fulfilled BOOLEAN DEFAULT FALSE,
            is_cancelled BOOLEAN DEFAULT FALSE,
            is_cancelled_by_client BOOLEAN DEFAULT FALSE,
            is_cancelled_by_employee BOOLEAN DEFAULT FALSE,
            is_confirmed_by_client BOOLEAN DEFAULT FALSE,
            history JSONB,
            comments JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_archives_deleted_at ON %I.appointments_archives(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_archives_client_id ON %I.appointments_archives(client_id)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_archives_company_id ON %I.appointments_archives(company_id)', schema_name);

END;
$$ LANGUAGE plpgsql;

