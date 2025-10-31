-- Migration: Rollback Admin and RoleAdmin tables
-- Purpose: Remove admin tables if needed
-- Created: 2025-10-31

-- Drop join table first (foreign key constraints)
DROP TABLE IF EXISTS public.admin_role_admins CASCADE;

-- Drop admin and role tables
DROP TABLE IF EXISTS public.admins CASCADE;
DROP TABLE IF EXISTS public.role_admins CASCADE;
