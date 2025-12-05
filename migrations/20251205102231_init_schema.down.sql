-- Smart migration rollback - Auto-detected changes
-- Generated at: 20251205102231

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
DROP TABLE IF EXISTS public.clients;

-- Rollback Model: Company
DROP TABLE IF EXISTS public.companies;

-- Rollback Model: Sector
DROP TABLE IF EXISTS public.sectors;

-- Rollback Model: Holiday
DROP TABLE IF EXISTS public.holidays;

-- Rollback Model: Role
DROP TABLE IF EXISTS public.roles;

-- Rollback Model: Resource
DROP TABLE IF EXISTS public.resources;

-- Rollback Model: Property
DROP TABLE IF EXISTS public.properties;

-- Rollback Model: EndPoint
DROP TABLE IF EXISTS public.endpoints;

-- Rollback Model: PolicyRule
DROP TABLE IF EXISTS public.policy_rules;

-- Rollback Model: Subdomain
DROP TABLE IF EXISTS public.subdomains;

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
DROP TABLE IF EXISTS public.clientappointments;

