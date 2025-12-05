-- Smart migration - Auto-detected changes
-- Generated at: 20251205100932
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.employees ("id" UUID NULL, "email" TEXT NULL, "slot_time_diff" TEXT NULL, "time_zone" TEXT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "phone" TEXT NULL, "tags" TEXT NULL, "password" TEXT NULL, "company_id" UUID NOT NULL, "name" TEXT NULL, "surname" TEXT NULL, "total_service_density" TEXT NOT NULL, "meta" TEXT NULL, "verified" BOOLEAN NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.branches ("deleted_at" TEXT NULL, "complement" TEXT NULL, "city" TEXT NULL, "country" TEXT NULL, "company_id" UUID NOT NULL, "total_service_density" INTEGER NOT NULL, "id" UUID NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "name" TEXT NULL, "neighborhood" TEXT NULL, "design" TEXT NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "number" TEXT NULL, "state" TEXT NULL, "street" TEXT NULL, "zip_code" TEXT NULL, "time_zone" TEXT NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.services ("design" TEXT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "description" TEXT NULL, "currency" TEXT NULL, "duration" TEXT NOT NULL, "company_id" UUID NOT NULL, "id" UUID NULL, "deleted_at" TEXT NULL, "name" TEXT NULL, "price" BIGINT NOT NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.appointments ("id" UUID NULL, "employee_id" UUID NOT NULL, "client_id" UUID NOT NULL, "cancel_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled" BOOLEAN NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "branch_id" UUID NOT NULL, "payment_id" UUID NULL, "actual_start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled_by_client" BOOLEAN NULL, "history" TEXT NULL, "service_id" UUID NOT NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_fulfilled" BOOLEAN NULL, "is_cancelled_by_employee" BOOLEAN NULL, "comments" TEXT NULL, "company_id" UUID NOT NULL, "cancelled_employee_id" UUID NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "time_zone" TEXT NOT NULL, "actual_end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_confirmed_by_client" BOOLEAN NULL)', schema_name);
    END LOOP;
END $$;

-- Model: Client (Table: clients, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

-- Model: Company (Table: companies, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

-- Model: Sector (Table: sectors, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

-- Model: Holiday (Table: holidays, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

-- Model: Role (Table: roles, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

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
-- No changes detected for this model

-- Model: EndPoint (Table: endpoints, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

-- Model: PolicyRule (Table: policy_rules, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

-- Model: Subdomain (Table: subdomains, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.branch_work_ranges ("weekday" INTEGER NOT NULL, "time_zone" TEXT NOT NULL, "branch_id" UUID NOT NULL, "id" UUID NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.employee_work_ranges ("branch_id" UUID NOT NULL, "id" UUID NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "time_zone" TEXT NOT NULL, "employee_id" UUID NOT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "weekday" INTEGER NOT NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.branch_service_densities ("density" INTEGER NOT NULL, "id" UUID NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "branch_id" UUID NULL, "service_id" UUID NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.employee_service_densities ("updated_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "employee_id" UUID NULL, "service_id" UUID NULL, "density" TEXT NOT NULL, "id" UUID NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.payments ("currency" TEXT NOT NULL, "transaction_id" TEXT NULL, "completed_at" TIMESTAMP WITH TIME ZONE NULL, "failed_at" TIMESTAMP WITH TIME ZONE NULL, "id" UUID NULL, "price" BIGINT NOT NULL, "user_id" TEXT NULL, "order_id" TEXT NULL, "metadata" TEXT NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "status" TEXT NOT NULL, "payment_method" TEXT NULL, "provider" TEXT NULL, "deleted_at" TEXT NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL)', schema_name);
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
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I.appointments_archive ("actual_end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled" BOOLEAN NULL, "created_at" TIMESTAMP WITH TIME ZONE NULL, "deleted_at" TEXT NULL, "employee_id" UUID NOT NULL, "client_id" UUID NOT NULL, "cancelled_employee_id" UUID NULL, "start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_fulfilled" BOOLEAN NULL, "comments" TEXT NULL, "service_id" UUID NOT NULL, "actual_start_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled_by_client" BOOLEAN NULL, "is_confirmed_by_client" BOOLEAN NULL, "payment_id" UUID NULL, "cancel_time" TIMESTAMP WITH TIME ZONE NOT NULL, "is_cancelled_by_employee" BOOLEAN NULL, "history" TEXT NULL, "id" UUID NULL, "updated_at" TIMESTAMP WITH TIME ZONE NULL, "branch_id" UUID NOT NULL, "company_id" UUID NOT NULL, "end_time" TIMESTAMP WITH TIME ZONE NOT NULL, "time_zone" TEXT NOT NULL)', schema_name);
    END LOOP;
END $$;

-- Model: ClientAppointment (Table: clientappointments, Schema: public)
-- Comparison: Go struct fields vs Current database columns
-- No changes detected for this model

