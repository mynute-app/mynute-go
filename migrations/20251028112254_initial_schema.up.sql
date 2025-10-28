-- Auto-generated CREATE TABLE migration
-- Generated from GORM models
-- ⚠️  REVIEW THIS SQL BEFORE APPLYING!

-- ======================
-- PUBLIC SCHEMA TABLES
-- ======================

-- Sectors table
CREATE TABLE IF NOT EXISTS public.sectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_sectors_deleted_at ON public.sectors(deleted_at);

-- Companies table
CREATE TABLE IF NOT EXISTS public.companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    sector_id UUID NOT NULL REFERENCES public.sectors(id),
    tier VARCHAR(50),
    subscription_status VARCHAR(50),
    trial_ends_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON public.companies(deleted_at);
CREATE INDEX IF NOT EXISTS idx_companies_sector_id ON public.companies(sector_id);

-- Holidays table
CREATE TABLE IF NOT EXISTS public.holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT false,
    company_id UUID REFERENCES public.companies(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_holidays_deleted_at ON public.holidays(deleted_at);
CREATE INDEX IF NOT EXISTS idx_holidays_company_id ON public.holidays(company_id);
CREATE INDEX IF NOT EXISTS idx_holidays_date ON public.holidays(date);

-- Clients table
CREATE TABLE IF NOT EXISTS public.clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES public.companies(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON public.clients(deleted_at);
CREATE INDEX IF NOT EXISTS idx_clients_company_id ON public.clients(company_id);
CREATE INDEX IF NOT EXISTS idx_clients_email ON public.clients(email);

-- Endpoints table
CREATE TABLE IF NOT EXISTS public.endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    group_name VARCHAR(100),
    deny_unauthorized BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(path, method)
);

CREATE INDEX IF NOT EXISTS idx_endpoints_deleted_at ON public.endpoints(deleted_at);
CREATE INDEX IF NOT EXISTS idx_endpoints_path_method ON public.endpoints(path, method);

-- Roles table
CREATE TABLE IF NOT EXISTS public.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON public.roles(deleted_at);

-- Policy Rules table
CREATE TABLE IF NOT EXISTS public.policy_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES public.roles(id) ON DELETE CASCADE,
    endpoint_id UUID NOT NULL REFERENCES public.endpoints(id) ON DELETE CASCADE,
    effect VARCHAR(10) NOT NULL CHECK (effect IN ('allow', 'deny')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(role_id, endpoint_id)
);

CREATE INDEX IF NOT EXISTS idx_policy_rules_deleted_at ON public.policy_rules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_policy_rules_role_id ON public.policy_rules(role_id);
CREATE INDEX IF NOT EXISTS idx_policy_rules_endpoint_id ON public.policy_rules(endpoint_id);

-- Resources table
CREATE TABLE IF NOT EXISTS public.resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_resources_deleted_at ON public.resources(deleted_at);

-- Properties table
CREATE TABLE IF NOT EXISTS public.properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL REFERENCES public.resources(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(resource_id, key)
);

CREATE INDEX IF NOT EXISTS idx_properties_deleted_at ON public.properties(deleted_at);
CREATE INDEX IF NOT EXISTS idx_properties_resource_id ON public.properties(resource_id);

-- Subdomains table
CREATE TABLE IF NOT EXISTS public.subdomains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL UNIQUE REFERENCES public.companies(id) ON DELETE CASCADE,
    subdomain VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_subdomains_deleted_at ON public.subdomains(deleted_at);
CREATE INDEX IF NOT EXISTS idx_subdomains_company_id ON public.subdomains(company_id);
CREATE INDEX IF NOT EXISTS idx_subdomains_subdomain ON public.subdomains(subdomain);

-- Client Appointments table
CREATE TABLE IF NOT EXISTS public.clientappointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES public.clients(id) ON DELETE CASCADE,
    appointment_id UUID NOT NULL,
    company_id UUID NOT NULL REFERENCES public.companies(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_clientappointments_deleted_at ON public.clientappointments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_clientappointments_client_id ON public.clientappointments(client_id);
CREATE INDEX IF NOT EXISTS idx_clientappointments_appointment_id ON public.clientappointments(appointment_id);
CREATE INDEX IF NOT EXISTS idx_clientappointments_company_id ON public.clientappointments(company_id);

-- ===============================================
-- TENANT SCHEMA TABLES (in public for now)
-- Note: These should be moved to company_* schemas in production
-- ===============================================

-- Branch Work Ranges
CREATE TABLE IF NOT EXISTS public.branch_work_ranges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_branch_work_ranges_deleted_at ON public.branch_work_ranges(deleted_at);
CREATE INDEX IF NOT EXISTS idx_branch_work_ranges_branch_id ON public.branch_work_ranges(branch_id);

-- Employee Work Ranges
CREATE TABLE IF NOT EXISTS public.employee_work_ranges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_employee_work_ranges_deleted_at ON public.employee_work_ranges(deleted_at);
CREATE INDEX IF NOT EXISTS idx_employee_work_ranges_employee_id ON public.employee_work_ranges(employee_id);

-- Branch Service Densities
CREATE TABLE IF NOT EXISTS public.branch_service_densities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID NOT NULL,
    service_id UUID NOT NULL,
    max_concurrent INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_branch_service_densities_deleted_at ON public.branch_service_densities(deleted_at);
CREATE INDEX IF NOT EXISTS idx_branch_service_densities_branch_id ON public.branch_service_densities(branch_id);

-- Employee Service Densities
CREATE TABLE IF NOT EXISTS public.employee_service_densities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    service_id UUID NOT NULL,
    max_concurrent INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_employee_service_densities_deleted_at ON public.employee_service_densities(deleted_at);
CREATE INDEX IF NOT EXISTS idx_employee_service_densities_employee_id ON public.employee_service_densities(employee_id);

-- Payments
CREATE TABLE IF NOT EXISTS public.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(50),
    payment_method VARCHAR(50),
    transaction_id VARCHAR(255),
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON public.payments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_appointment_id ON public.payments(appointment_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON public.payments(status);

-- ===============================================
-- MULTI-TENANT TABLES (company_* schemas)
-- These tables are created per company schema
-- ===============================================

-- Function to create tenant schema tables
CREATE OR REPLACE FUNCTION create_tenant_schema_tables(schema_name TEXT) RETURNS VOID AS $$
BEGIN
    -- Employees table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.employees (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            first_name VARCHAR(100) NOT NULL,
            last_name VARCHAR(100) NOT NULL,
            email VARCHAR(255) NOT NULL,
            phone VARCHAR(50),
            role_id UUID,
            branch_id UUID,
            is_active BOOLEAN DEFAULT true,
            bio TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);

    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON %I.employees(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_employees_email ON %I.employees(email)', schema_name);

    -- Branches table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.branches (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(255) NOT NULL,
            address TEXT,
            city VARCHAR(100),
            state VARCHAR(100),
            postal_code VARCHAR(20),
            country VARCHAR(100),
            phone VARCHAR(50),
            email VARCHAR(255),
            is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);

    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_branches_deleted_at ON %I.branches(deleted_at)', schema_name);

    -- Services table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.services (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(255) NOT NULL,
            description TEXT,
            duration INTEGER NOT NULL,
            price DECIMAL(10, 2),
            is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);

    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_services_deleted_at ON %I.services(deleted_at)', schema_name);

    -- Appointments table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.appointments (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            employee_id UUID,
            service_id UUID,
            branch_id UUID,
            start_time TIMESTAMP WITH TIME ZONE NOT NULL,
            end_time TIMESTAMP WITH TIME ZONE NOT NULL,
            status VARCHAR(50) DEFAULT ''scheduled'',
            notes TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);

    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON %I.appointments(deleted_at)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_employee_id ON %I.appointments(employee_id)', schema_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_start_time ON %I.appointments(start_time)', schema_name);

    -- Appointments Archive table
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I.appointments_archive (
            id UUID PRIMARY KEY,
            employee_id UUID,
            service_id UUID,
            branch_id UUID,
            start_time TIMESTAMP WITH TIME ZONE NOT NULL,
            end_time TIMESTAMP WITH TIME ZONE NOT NULL,
            status VARCHAR(50),
            notes TEXT,
            archived_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            original_created_at TIMESTAMP WITH TIME ZONE,
            original_updated_at TIMESTAMP WITH TIME ZONE,
            original_deleted_at TIMESTAMP WITH TIME ZONE
        )', schema_name);

    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_appointments_archive_archived_at ON %I.appointments_archive(archived_at)', schema_name);
END;
$$ LANGUAGE plpgsql;

-- Note: Run this function for each company schema:
-- SELECT create_tenant_schema_tables('company_<uuid>');
