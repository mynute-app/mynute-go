-- Smart migration - Auto-detected changes
-- Generated at: 20251205095228
-- Compared against schema: public
--
-- ⚠️  IMPORTANT: This tool compares Go models against the CURRENT database schema.
--     It does NOT track migration history. A 'new' column means:
--     - The Go model has this field, AND
--     - The database table does NOT have this column
--
--     Review carefully before applying!

-- Model: Employee (Table: employees, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Adding 16 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS surname TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS tags TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS slot_time_diff TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS time_zone TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS email TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS verified BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS password TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS total_service_density TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS name TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS phone TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employees ADD COLUMN IF NOT EXISTS meta TEXT NULL', schema_name);
    END LOOP;
END $$;

-- Model: Branch (Table: branches, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Adding 17 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS time_zone TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS total_service_density INTEGER NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS number TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS country TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS neighborhood TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS zip_code TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS design TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS street TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS complement TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS state TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS name TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branches ADD COLUMN IF NOT EXISTS city TEXT NULL', schema_name);
    END LOOP;
END $$;

-- Model: Service (Table: services, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Adding 11 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS name TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS duration TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS design TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS description TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS price BIGINT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.services ADD COLUMN IF NOT EXISTS currency TEXT NULL', schema_name);
    END LOOP;
END $$;

-- Model: Appointment (Table: appointments, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Adding 24 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS employee_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS cancelled_employee_id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS is_cancelled BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS is_cancelled_by_client BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS is_confirmed_by_client BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS service_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS payment_id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS actual_start_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS history TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS comments TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS client_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS branch_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS is_cancelled_by_employee BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS time_zone TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS actual_end_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS cancel_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments ADD COLUMN IF NOT EXISTS is_fulfilled BOOLEAN NULL', schema_name);
    END LOOP;
END $$;

-- Model: Client (Table: clients, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 11 new column(s)
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS name TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS phone TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS meta TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS surname TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS email TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS password TEXT NULL;
ALTER TABLE public.clients ADD COLUMN IF NOT EXISTS verified BOOLEAN NULL;

-- Model: Company (Table: companies, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 9 new column(s)
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS legal_name TEXT NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS trade_name TEXT NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS tax_id TEXT NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS schema_name TEXT NULL;
ALTER TABLE public.companies ADD COLUMN IF NOT EXISTS design TEXT NULL;

-- Model: Sector (Table: sectors, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 6 new column(s)
ALTER TABLE public.sectors ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.sectors ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.sectors ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
ALTER TABLE public.sectors ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE public.sectors ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.sectors ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;

-- Model: Holiday (Table: holidays, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 10 new column(s)
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS type TEXT NOT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS description TEXT NOT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS recurrent BOOLEAN NOT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS date TIMESTAMP WITH TIME ZONE NOT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS day_month TEXT NOT NULL;
ALTER TABLE public.holidays ADD COLUMN IF NOT EXISTS id UUID NULL;

-- Model: Role (Table: roles, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 7 new column(s)
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS company_id UUID NULL;
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;

-- Model: Resource (Table: resources, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 8 new column(s)
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS table TEXT NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS references TEXT NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.resources ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;

-- Model: Property (Table: properties, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 7 new column(s)
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS resource_name TEXT NULL;
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.properties ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;

-- Model: EndPoint (Table: endpoints, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 11 new column(s)
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS controller_name TEXT NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS method TEXT NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS needs_company_id BOOLEAN NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS path TEXT NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS deny_unauthorized BOOLEAN NULL;
ALTER TABLE public.endpoints ADD COLUMN IF NOT EXISTS resource_id UUID NULL;

-- Model: PolicyRule (Table: policy_rules, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 9 new column(s)
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS name TEXT NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS effect TEXT NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS end_point_id UUID NULL;
ALTER TABLE public.policy_rules ADD COLUMN IF NOT EXISTS conditions TEXT NULL;

-- Model: Subdomain (Table: subdomains, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 6 new column(s)
ALTER TABLE public.subdomains ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE public.subdomains ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL;
ALTER TABLE public.subdomains ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
ALTER TABLE public.subdomains ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL;
ALTER TABLE public.subdomains ADD COLUMN IF NOT EXISTS id UUID NULL;
ALTER TABLE public.subdomains ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL;

-- Model: BranchWorkRange (Table: branch_work_ranges, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Adding 9 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS weekday INTEGER NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS time_zone TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_work_ranges ADD COLUMN IF NOT EXISTS branch_id UUID NOT NULL', schema_name);
    END LOOP;
END $$;

-- Model: EmployeeWorkRange (Table: employee_work_ranges, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Adding 10 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS weekday INTEGER NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS time_zone TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS branch_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS employee_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_work_ranges ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
    END LOOP;
END $$;

-- Model: BranchServiceDensity (Table: branch_service_densities, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Adding 7 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS branch_id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS service_id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS density INTEGER NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.branch_service_densities ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
    END LOOP;
END $$;

-- Model: EmployeeServiceDensity (Table: employee_service_densities, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Adding 7 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS density TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS employee_id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.employee_service_densities ADD COLUMN IF NOT EXISTS service_id UUID NULL', schema_name);
    END LOOP;
END $$;

-- Model: Payment (Table: payments, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Adding 15 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS order_id TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS failed_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS price BIGINT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS provider TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS metadata TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS status TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS currency TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS payment_method TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS transaction_id TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.payments ADD COLUMN IF NOT EXISTS user_id TEXT NULL', schema_name);
    END LOOP;
END $$;

-- Model: AppointmentArchive (Table: appointments_archive, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Adding 24 new column(s)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS cancel_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS payment_id UUID NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS is_cancelled_by_client BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS history TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS comments TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS branch_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS is_fulfilled BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS time_zone TEXT NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS actual_start_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS actual_end_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS is_cancelled BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS is_cancelled_by_employee BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS is_confirmed_by_client BOOLEAN NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS deleted_at TEXT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS service_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS employee_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS client_id UUID NOT NULL', schema_name);
        EXECUTE format('ALTER TABLE %I.appointments_archive ADD COLUMN IF NOT EXISTS cancelled_employee_id UUID NULL', schema_name);
    END LOOP;
END $$;

-- Model: ClientAppointment (Table: clientappointments, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Adding 7 new column(s)
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE NOT NULL;
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS time_zone TEXT NOT NULL;
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS is_cancelled BOOLEAN NULL;
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS appointment_id UUID NOT NULL;
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS client_id UUID NOT NULL;
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL;
ALTER TABLE public.clientappointments ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE NOT NULL;

