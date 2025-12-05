-- Smart migration - Auto-detected changes
-- Generated at: 20251205102231
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
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.employees ("email" TEXT NULL, "password" TEXT NULL, "slot_time_diff" TEXT NULL, "time_zone" TEXT NULL, "verified" BOOLEAN NULL, "surname" TEXT NULL, "company_id" UUID NOT NULL, "meta" TEXT NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "name" TEXT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "phone" TEXT NULL, "tags" TEXT NULL, "total_service_density" TEXT NOT NULL, "id" UUID NULL)', schema_name);
    END LOOP;
END $$;

-- Model: Branch (Table: branches, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.branches ("zip_code" TEXT NULL, "time_zone" TEXT NULL, "company_id" UUID NOT NULL, "id" UUID NULL, "street" TEXT NULL, "number" TEXT NULL, "complement" TEXT NULL, "neighborhood" TEXT NULL, "city" TEXT NULL, "total_service_density" INTEGER NOT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "state" TEXT NULL, "country" TEXT NULL, "design" TEXT NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "name" TEXT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: Service (Table: services, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.services ("currency" TEXT NULL, "company_id" UUID NOT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "name" TEXT NULL, "description" TEXT NULL, "duration" TEXT NOT NULL, "design" TEXT NULL, "id" UUID NULL, "price" BIGINT NOT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: Appointment (Table: appointments, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.appointments ("created_at" TIMESTAMP WITH TIME ZONE NULL, "cancel_time" TIMESTAMP WITH TIME ZONE NOT NULL, "history" TEXT NULL, "id" UUID NULL, "employee_id" UUID NOT NULL, "branch_id" UUID NOT NULL, "company_id" UUID NOT NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "time_zone" TEXT NOT NULL, "actual_start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "actual_end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "service_id" UUID NOT NULL, "payment_id" UUID NULL, "cancelled_employee_id" UUID NULL, "is_fulfilled" BOOLEAN NULL, "is_cancelled" BOOLEAN NULL, "is_cancelled_by_client" BOOLEAN NULL, "is_cancelled_by_employee" BOOLEAN NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "client_id" UUID NOT NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_confirmed_by_client" BOOLEAN NULL, "comments" TEXT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: Client (Table: clients, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.clients (
    "phone" TEXT NULL,
    "meta" TEXT NULL,
    "deleted_at" TEXT NULL,
    "surname" TEXT NULL,
    "password" TEXT NULL,
    "verified" BOOLEAN NULL,
    "id" UUID NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "name" TEXT NULL,
    "email" TEXT NULL
);

-- Model: Company (Table: companies, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.companies (
    "id" UUID NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "legal_name" TEXT NULL,
    "trade_name" TEXT NULL,
    "tax_id" TEXT NULL,
    "schema_name" TEXT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "design" TEXT NULL
);

-- Model: Sector (Table: sectors, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.sectors (
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NULL,
    "id" UUID NULL
);

-- Model: Holiday (Table: holidays, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.holidays (
    "type" TEXT NOT NULL,
    "id" UUID NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "name" TEXT NOT NULL,
    "date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "description" TEXT NOT NULL,
    "recurrent" BOOLEAN NOT NULL,
    "day_month" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL
);

-- Model: Role (Table: roles, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.roles (
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NULL,
    "company_id" UUID NULL,
    "id" UUID NULL
);

-- Model: Resource (Table: resources, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.resources (
    "references" TEXT NULL,
    "id" UUID NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NULL,
    "table" TEXT NULL
);

-- Model: Property (Table: properties, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.properties (
    "deleted_at" TEXT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NULL,
    "resource_name" TEXT NULL,
    "id" UUID NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL
);

-- Model: EndPoint (Table: endpoints, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.endpoints (
    "controller_name" TEXT NULL,
    "path" TEXT NULL,
    "deny_unauthorized" BOOLEAN NULL,
    "resource_id" UUID NULL,
    "id" UUID NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "description" TEXT NULL,
    "method" TEXT NULL,
    "needs_company_id" BOOLEAN NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL
);

-- Model: PolicyRule (Table: policy_rules, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.policy_rules (
    "name" TEXT NULL,
    "effect" TEXT NULL,
    "end_point_id" UUID NULL,
    "id" UUID NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL,
    "description" TEXT NULL,
    "conditions" TEXT NULL
);

-- Model: Subdomain (Table: subdomains, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.subdomains (
    "name" TEXT NOT NULL,
    "company_id" UUID NOT NULL,
    "id" UUID NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NULL,
    "deleted_at" TEXT NULL
);

-- Model: BranchWorkRange (Table: branch_work_ranges, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.branch_work_ranges ("created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "time_zone" TEXT NOT NULL, "branch_id" UUID NOT NULL, "id" UUID NULL, "weekday" INTEGER NOT NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: EmployeeWorkRange (Table: employee_work_ranges, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.employee_work_ranges ("deleted_at" TEXT NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "time_zone" TEXT NOT NULL, "branch_id" UUID NOT NULL, "employee_id" UUID NOT NULL, "id" UUID NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "weekday" INTEGER NOT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: BranchServiceDensity (Table: branch_service_densities, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.branch_service_densities ("branch_id" UUID NULL, "service_id" UUID NULL, "density" INTEGER NOT NULL, "id" UUID NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: EmployeeServiceDensity (Table: employee_service_densities, Schema: tenant)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.employee_service_densities ("id" UUID NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "employee_id" UUID NULL, "service_id" UUID NULL, "density" TEXT NOT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: Payment (Table: payments, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.payments ("created_at" TIMESTAMP WITH TIME ZONE NULL, "price" BIGINT NOT NULL, "currency" TEXT NOT NULL, "payment_method" TEXT NULL, "failed_at" TIMESTAMP WITH TIME ZONE NULL, "id" UUID NULL, "transaction_id" TEXT NULL, "provider" TEXT NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "status" TEXT NOT NULL, "user_id" TEXT NULL, "order_id" TEXT NULL, "metadata" TEXT NULL, "completed_at" TIMESTAMP WITH TIME ZONE NULL)', schema_name);
    END LOOP;
END $$;

-- Model: AppointmentArchive (Table: appointments_archive, Schema: company)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.appointments_archive ("created_at" TIMESTAMP WITH TIME ZONE NULL, "employee_id" UUID NOT NULL, "client_id" UUID NOT NULL, "branch_id" UUID NOT NULL, "payment_id" UUID NULL, "time_zone" TEXT NOT NULL, "is_confirmed_by_client" BOOLEAN NULL, "comments" TEXT NULL, "id" UUID NULL, "deleted_at" TEXT NULL, "company_id" UUID NOT NULL, "actual_end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "cancel_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_fulfilled" BOOLEAN NULL, "is_cancelled_by_client" BOOLEAN NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "actual_start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled_by_employee" BOOLEAN NULL, "history" TEXT NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "service_id" UUID NOT NULL, "cancelled_employee_id" UUID NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled" BOOLEAN NULL)', schema_name);
    END LOOP;
END $$;

-- Model: ClientAppointment (Table: clientappointments, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- Table does not exist, creating it
CREATE TABLE IF NOT EXISTS public.clientappointments (
    "is_cancelled" BOOLEAN NULL,
    "appointment_id" UUID NOT NULL,
    "client_id" UUID NOT NULL,
    "company_id" UUID NOT NULL,
    "start_time" TIMESTAMP WITH TIME ZONE NOT NULL,
    "end_time" TIMESTAMP WITH TIME ZONE NOT NULL,
    "time_zone" TEXT NOT NULL
);

