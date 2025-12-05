-- Smart migration rollback - Auto-detected changes
-- Generated at: 20251205100932

-- Rollback Model: Employee
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.employees', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Branch
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.branches', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Service
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.services', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Appointment
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.appointments', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Client
-- No changes to rollback

-- Rollback Model: Company
-- No changes to rollback

-- Rollback Model: Sector
-- No changes to rollback

-- Rollback Model: Holiday
-- No changes to rollback

-- Rollback Model: Role
-- No changes to rollback

-- Rollback Model: Resource
DROP TABLE IF EXISTS public.resources;

-- Rollback Model: Property
-- No changes to rollback

-- Rollback Model: EndPoint
-- No changes to rollback

-- Rollback Model: PolicyRule
-- No changes to rollback

-- Rollback Model: Subdomain
-- No changes to rollback

-- Rollback Model: BranchWorkRange
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.branch_work_ranges', schema_name);
    END LOOP;
END $$;

-- Rollback Model: EmployeeWorkRange
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.employee_work_ranges', schema_name);
    END LOOP;
END $$;

-- Rollback Model: BranchServiceDensity
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.branch_service_densities', schema_name);
    END LOOP;
END $$;

-- Rollback Model: EmployeeServiceDensity
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.employee_service_densities', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Payment
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.payments', schema_name);
    END LOOP;
END $$;

-- Rollback Model: AppointmentArchive
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I.appointments_archive', schema_name);
    END LOOP;
END $$;

-- Rollback Model: ClientAppointment
-- No changes to rollback

