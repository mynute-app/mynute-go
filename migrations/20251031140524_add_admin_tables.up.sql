-- Migration: Add Admin and RoleAdmin tables
-- Purpose: Implement system-wide administrators with role-based access
-- Created: 2025-10-31

-- ======================
-- ADMIN TABLES (PUBLIC SCHEMA)
-- ======================

-- Admins table
CREATE TABLE IF NOT EXISTS public.admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    meta JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_admins_deleted_at ON public.admins(deleted_at);
CREATE INDEX IF NOT EXISTS idx_admins_email ON public.admins(email);
CREATE INDEX IF NOT EXISTS idx_admins_is_active ON public.admins(is_active);

COMMENT ON TABLE public.admins IS 'System-wide administrators who can access all tenants';
COMMENT ON COLUMN public.admins.is_active IS 'Whether the admin account is active and can log in';
COMMENT ON COLUMN public.admins.meta IS 'Additional metadata for the admin user (JSONB format)';

-- RoleAdmins table
CREATE TABLE IF NOT EXISTS public.role_admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_role_admins_deleted_at ON public.role_admins(deleted_at);
CREATE INDEX IF NOT EXISTS idx_role_admins_name ON public.role_admins(name);

COMMENT ON TABLE public.role_admins IS 'Roles for system administrators (e.g., superadmin, auditor)';
COMMENT ON COLUMN public.role_admins.name IS 'Unique role name (e.g., superadmin, support, auditor)';

-- Admin-Role join table (many-to-many)
CREATE TABLE IF NOT EXISTS public.admin_role_admins (
    admin_id UUID NOT NULL REFERENCES public.admins(id) ON DELETE CASCADE,
    role_admin_id UUID NOT NULL REFERENCES public.role_admins(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (admin_id, role_admin_id)
);

CREATE INDEX IF NOT EXISTS idx_admin_role_admins_admin_id ON public.admin_role_admins(admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_role_admins_role_admin_id ON public.admin_role_admins(role_admin_id);

COMMENT ON TABLE public.admin_role_admins IS 'Join table connecting admins to their roles';

-- ======================
-- SEED DEFAULT DATA
-- ======================

-- Insert default admin role: superadmin
INSERT INTO public.role_admins (name, description)
VALUES 
    ('superadmin', 'Full access to all tenants and resources'),
    ('support', 'Customer support with read-only access to tenant data'),
    ('auditor', 'Read-only access for compliance and auditing')
ON CONFLICT (name) DO NOTHING;

COMMENT ON TABLE public.role_admins IS 'System administrator roles with different permission levels';
