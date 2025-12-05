-- Smart migration rollback - Auto-detected changes
-- Generated at: 20251205095228

-- Rollback Model: Employee
-- Removing 16 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS surname', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS tags', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS slot_time_diff', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS time_zone', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS email', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS verified', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS password', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS company_id', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS total_service_density', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS name', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS phone', schema_name);
        EXECUTE format('ALTER TABLE %I.employees DROP COLUMN IF EXISTS meta', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Branch
-- Removing 17 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS time_zone', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS total_service_density', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS number', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS country', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS neighborhood', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS zip_code', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS company_id', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS design', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS street', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS complement', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS state', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS name', schema_name);
        EXECUTE format('ALTER TABLE %I.branches DROP COLUMN IF EXISTS city', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Service
-- Removing 11 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS name', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS duration', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS company_id', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS design', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS description', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS price', schema_name);
        EXECUTE format('ALTER TABLE %I.services DROP COLUMN IF EXISTS currency', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Appointment
-- Removing 24 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS end_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS employee_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS company_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS cancelled_employee_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS is_cancelled', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS is_cancelled_by_client', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS is_confirmed_by_client', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS service_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS payment_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS actual_start_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS history', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS comments', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS client_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS branch_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS is_cancelled_by_employee', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS start_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS time_zone', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS actual_end_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS cancel_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments DROP COLUMN IF EXISTS is_fulfilled', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Client
-- Removing 11 column(s) that were added
ALTER TABLE public.clients DROP COLUMN IF EXISTS id;
ALTER TABLE public.clients DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.clients DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.clients DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.clients DROP COLUMN IF EXISTS name;
ALTER TABLE public.clients DROP COLUMN IF EXISTS phone;
ALTER TABLE public.clients DROP COLUMN IF EXISTS meta;
ALTER TABLE public.clients DROP COLUMN IF EXISTS surname;
ALTER TABLE public.clients DROP COLUMN IF EXISTS email;
ALTER TABLE public.clients DROP COLUMN IF EXISTS password;
ALTER TABLE public.clients DROP COLUMN IF EXISTS verified;

-- Rollback Model: Company
-- Removing 9 column(s) that were added
ALTER TABLE public.companies DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.companies DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.companies DROP COLUMN IF EXISTS legal_name;
ALTER TABLE public.companies DROP COLUMN IF EXISTS id;
ALTER TABLE public.companies DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.companies DROP COLUMN IF EXISTS trade_name;
ALTER TABLE public.companies DROP COLUMN IF EXISTS tax_id;
ALTER TABLE public.companies DROP COLUMN IF EXISTS schema_name;
ALTER TABLE public.companies DROP COLUMN IF EXISTS design;

-- Rollback Model: Sector
-- Removing 6 column(s) that were added
ALTER TABLE public.sectors DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.sectors DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.sectors DROP COLUMN IF EXISTS name;
ALTER TABLE public.sectors DROP COLUMN IF EXISTS description;
ALTER TABLE public.sectors DROP COLUMN IF EXISTS id;
ALTER TABLE public.sectors DROP COLUMN IF EXISTS created_at;

-- Rollback Model: Holiday
-- Removing 10 column(s) that were added
ALTER TABLE public.holidays DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS name;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS type;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS description;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS recurrent;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS date;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS day_month;
ALTER TABLE public.holidays DROP COLUMN IF EXISTS id;

-- Rollback Model: Role
-- Removing 7 column(s) that were added
ALTER TABLE public.roles DROP COLUMN IF EXISTS name;
ALTER TABLE public.roles DROP COLUMN IF EXISTS description;
ALTER TABLE public.roles DROP COLUMN IF EXISTS company_id;
ALTER TABLE public.roles DROP COLUMN IF EXISTS id;
ALTER TABLE public.roles DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.roles DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.roles DROP COLUMN IF EXISTS deleted_at;

-- Rollback Model: Resource
-- Removing 8 column(s) that were added
ALTER TABLE public.resources DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.resources DROP COLUMN IF EXISTS name;
ALTER TABLE public.resources DROP COLUMN IF EXISTS description;
ALTER TABLE public.resources DROP COLUMN IF EXISTS table;
ALTER TABLE public.resources DROP COLUMN IF EXISTS references;
ALTER TABLE public.resources DROP COLUMN IF EXISTS id;
ALTER TABLE public.resources DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.resources DROP COLUMN IF EXISTS updated_at;

-- Rollback Model: Property
-- Removing 7 column(s) that were added
ALTER TABLE public.properties DROP COLUMN IF EXISTS description;
ALTER TABLE public.properties DROP COLUMN IF EXISTS resource_name;
ALTER TABLE public.properties DROP COLUMN IF EXISTS id;
ALTER TABLE public.properties DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.properties DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.properties DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.properties DROP COLUMN IF EXISTS name;

-- Rollback Model: EndPoint
-- Removing 11 column(s) that were added
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS controller_name;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS description;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS method;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS needs_company_id;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS id;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS path;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS deny_unauthorized;
ALTER TABLE public.endpoints DROP COLUMN IF EXISTS resource_id;

-- Rollback Model: PolicyRule
-- Removing 9 column(s) that were added
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS name;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS description;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS id;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS created_at;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS effect;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS end_point_id;
ALTER TABLE public.policy_rules DROP COLUMN IF EXISTS conditions;

-- Rollback Model: Subdomain
-- Removing 6 column(s) that were added
ALTER TABLE public.subdomains DROP COLUMN IF EXISTS updated_at;
ALTER TABLE public.subdomains DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.subdomains DROP COLUMN IF EXISTS name;
ALTER TABLE public.subdomains DROP COLUMN IF EXISTS company_id;
ALTER TABLE public.subdomains DROP COLUMN IF EXISTS id;
ALTER TABLE public.subdomains DROP COLUMN IF EXISTS created_at;

-- Rollback Model: BranchWorkRange
-- Removing 9 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS weekday', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS start_time', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS end_time', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS time_zone', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges DROP COLUMN IF EXISTS branch_id', schema_name);
    END LOOP;
END $$;

-- Rollback Model: EmployeeWorkRange
-- Removing 10 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS weekday', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS start_time', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS time_zone', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS branch_id', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS end_time', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS employee_id', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges DROP COLUMN IF EXISTS updated_at', schema_name);
    END LOOP;
END $$;

-- Rollback Model: BranchServiceDensity
-- Removing 7 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS branch_id', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS service_id', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS density', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities DROP COLUMN IF EXISTS updated_at', schema_name);
    END LOOP;
END $$;

-- Rollback Model: EmployeeServiceDensity
-- Removing 7 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS density', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS employee_id', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities DROP COLUMN IF EXISTS service_id', schema_name);
    END LOOP;
END $$;

-- Rollback Model: Payment
-- Removing 15 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS order_id', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS failed_at', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS price', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS provider', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS metadata', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS status', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS currency', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS completed_at', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS payment_method', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS transaction_id', schema_name);
        EXECUTE format('ALTER TABLE %I.payments DROP COLUMN IF EXISTS user_id', schema_name);
    END LOOP;
END $$;

-- Rollback Model: AppointmentArchive
-- Removing 24 column(s) that were added
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS cancel_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS payment_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS is_cancelled_by_client', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS history', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS comments', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS branch_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS company_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS is_fulfilled', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS updated_at', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS end_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS time_zone', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS actual_start_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS actual_end_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS is_cancelled', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS is_cancelled_by_employee', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS is_confirmed_by_client', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS start_time', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS created_at', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS deleted_at', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS service_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS employee_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS client_id', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive DROP COLUMN IF EXISTS cancelled_employee_id', schema_name);
    END LOOP;
END $$;

-- Rollback Model: ClientAppointment
-- Removing 7 column(s) that were added
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS end_time;
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS time_zone;
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS is_cancelled;
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS appointment_id;
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS client_id;
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS company_id;
ALTER TABLE public.clientappointments DROP COLUMN IF EXISTS start_time;

