-- Create "client_appointments" table
CREATE TABLE "client_appointments" ("appointment_id" uuid NOT NULL, "client_id" uuid NOT NULL, "company_id" uuid NOT NULL, "start_time" timestamptz NOT NULL, "end_time" timestamptz NOT NULL, "time_zone" character varying(100) NOT NULL, "is_cancelled" boolean NULL DEFAULT false);
-- Create "appointments_archive" table
CREATE TABLE "appointments_archive" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "service_id" uuid NOT NULL, "employee_id" uuid NOT NULL, "client_id" uuid NOT NULL, "branch_id" uuid NOT NULL, "payment_id" uuid NULL, "company_id" uuid NOT NULL, "cancelled_employee_id" uuid NULL, "start_time" timestamptz NOT NULL, "end_time" timestamptz NOT NULL, "time_zone" character varying(100) NOT NULL, "actual_start_time" timestamptz NOT NULL, "actual_end_time" timestamptz NOT NULL, "cancel_time" timestamptz NOT NULL, "is_fulfilled" boolean NULL DEFAULT false, "is_cancelled" boolean NULL DEFAULT false, "is_cancelled_by_client" boolean NULL DEFAULT false, "is_cancelled_by_employee" boolean NULL DEFAULT false, "is_confirmed_by_client" boolean NULL DEFAULT false, "history" jsonb NULL, "comments" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_appointments_archive_client_id" to table: "appointments_archive"
CREATE INDEX "idx_appointments_archive_client_id" ON "appointments_archive" ("client_id");
-- Create index "idx_appointments_archive_company_id" to table: "appointments_archive"
CREATE INDEX "idx_appointments_archive_company_id" ON "appointments_archive" ("company_id");
-- Create index "idx_appointments_archive_deleted_at" to table: "appointments_archive"
CREATE INDEX "idx_appointments_archive_deleted_at" ON "appointments_archive" ("deleted_at");
-- Create index "idx_appointments_archive_payment_id" to table: "appointments_archive"
CREATE UNIQUE INDEX "idx_appointments_archive_payment_id" ON "appointments_archive" ("payment_id");
-- Create "services" table
CREATE TABLE "services" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" character varying(100) NULL, "description" text NULL, "price" bigint NOT NULL, "currency" character varying(3) NULL DEFAULT 'BRL', "duration" integer NOT NULL, "company_id" uuid NOT NULL, "design" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_services_company_id" to table: "services"
CREATE INDEX "idx_services_company_id" ON "services" ("company_id");
-- Create index "idx_services_deleted_at" to table: "services"
CREATE INDEX "idx_services_deleted_at" ON "services" ("deleted_at");
-- Create "roles" table
CREATE TABLE "roles" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" character varying(100) NOT NULL, "description" text NULL, "company_id" uuid NULL, PRIMARY KEY ("id"));
-- Create index "idx_public_roles_company_id" to table: "roles"
CREATE INDEX "idx_public_roles_company_id" ON "roles" ("company_id");
-- Create index "idx_public_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_public_roles_deleted_at" ON "roles" ("deleted_at");
-- Create index "idx_role_name_company" to table: "roles"
CREATE UNIQUE INDEX "idx_role_name_company" ON "roles" ("name", "company_id");
-- Create "holidays" table
CREATE TABLE "holidays" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" text NOT NULL, "date" timestamptz NOT NULL, "type" text NOT NULL, "description" text NOT NULL, "recurrent" boolean NOT NULL, "day_month" text NOT NULL, PRIMARY KEY ("id"));
-- Create index "idx_public_holidays_deleted_at" to table: "holidays"
CREATE INDEX "idx_public_holidays_deleted_at" ON "holidays" ("deleted_at");
-- Create index "idx_public_holidays_name" to table: "holidays"
CREATE INDEX "idx_public_holidays_name" ON "holidays" ("name");
-- Create index "idx_public_holidays_recurrent" to table: "holidays"
CREATE INDEX "idx_public_holidays_recurrent" ON "holidays" ("recurrent");
-- Create "clients" table
CREATE TABLE "clients" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" character varying(100) NULL, "surname" character varying(100) NULL, "email" character varying(100) NULL, "phone" character varying(20) NULL, "password" character varying(255) NULL, "verified" boolean NULL DEFAULT false, "meta" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_public_clients_deleted_at" to table: "clients"
CREATE INDEX "idx_public_clients_deleted_at" ON "clients" ("deleted_at");
-- Create index "idx_public_clients_email" to table: "clients"
CREATE UNIQUE INDEX "idx_public_clients_email" ON "clients" ("email");
-- Create index "idx_public_clients_phone" to table: "clients"
CREATE UNIQUE INDEX "idx_public_clients_phone" ON "clients" ("phone");
-- Create "branches" table
CREATE TABLE "branches" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" character varying(100) NULL, "street" character varying(100) NULL, "number" character varying(100) NULL, "complement" character varying(100) NULL, "neighborhood" character varying(100) NULL, "zip_code" character varying(100) NULL, "city" character varying(100) NULL, "state" character varying(100) NULL, "country" character varying(100) NULL, "company_id" text NOT NULL, "time_zone" character varying(100) NULL, "total_service_density" integer NOT NULL DEFAULT -1, "design" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_branches_company_id" to table: "branches"
CREATE INDEX "idx_branches_company_id" ON "branches" ("company_id");
-- Create index "idx_branches_deleted_at" to table: "branches"
CREATE INDEX "idx_branches_deleted_at" ON "branches" ("deleted_at");
-- Create "branch_service_densities" table
CREATE TABLE "branch_service_densities" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "branch_id" uuid NOT NULL, "service_id" uuid NOT NULL, "density" integer NOT NULL DEFAULT 1, PRIMARY KEY ("id", "branch_id", "service_id"), CONSTRAINT "fk_branch_service_densities_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_branches_service_density" FOREIGN KEY ("branch_id") REFERENCES "branches" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_branch_service_densities_deleted_at" to table: "branch_service_densities"
CREATE INDEX "idx_branch_service_densities_deleted_at" ON "branch_service_densities" ("deleted_at");
-- Create "branch_services" table
CREATE TABLE "branch_services" ("service_id" uuid NOT NULL DEFAULT gen_random_uuid(), "branch_id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("service_id", "branch_id"), CONSTRAINT "fk_branch_services_branch" FOREIGN KEY ("branch_id") REFERENCES "branches" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_branch_services_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "employees" table
CREATE TABLE "employees" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" character varying(100) NULL, "surname" character varying(100) NULL, "email" character varying(100) NULL, "phone" character varying(20) NULL, "tags" json NULL, "password" character varying(255) NULL, "slot_time_diff" bigint NULL DEFAULT 30, "company_id" text NOT NULL, "time_zone" character varying(100) NULL, "total_service_density" bigint NOT NULL DEFAULT 1, "verified" boolean NULL DEFAULT false, "meta" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_employees_company_id" to table: "employees"
CREATE INDEX "idx_employees_company_id" ON "employees" ("company_id");
-- Create index "idx_employees_deleted_at" to table: "employees"
CREATE INDEX "idx_employees_deleted_at" ON "employees" ("deleted_at");
-- Create index "idx_employees_email" to table: "employees"
CREATE UNIQUE INDEX "idx_employees_email" ON "employees" ("email");
-- Create index "idx_employees_phone" to table: "employees"
CREATE UNIQUE INDEX "idx_employees_phone" ON "employees" ("phone");
-- Create "employee_service_densities" table
CREATE TABLE "employee_service_densities" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "employee_id" uuid NOT NULL, "service_id" uuid NOT NULL, "density" bigint NOT NULL DEFAULT 1, PRIMARY KEY ("id", "employee_id", "service_id"), CONSTRAINT "fk_employee_service_densities_employee" FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_employee_service_densities_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_employee_service_densities_deleted_at" to table: "employee_service_densities"
CREATE INDEX "idx_employee_service_densities_deleted_at" ON "employee_service_densities" ("deleted_at");
-- Create "employee_work_ranges" table
CREATE TABLE "employee_work_ranges" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "weekday" bigint NOT NULL, "start_time" timestamptz NOT NULL, "end_time" timestamptz NOT NULL, "time_zone" character varying(255) NOT NULL, "branch_id" uuid NOT NULL, "employee_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_employee_work_ranges_branch" FOREIGN KEY ("branch_id") REFERENCES "branches" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "fk_employees_work_schedule" FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_employee_id" to table: "employee_work_ranges"
CREATE INDEX "idx_employee_id" ON "employee_work_ranges" ("employee_id");
-- Create index "idx_employee_work_ranges_deleted_at" to table: "employee_work_ranges"
CREATE INDEX "idx_employee_work_ranges_deleted_at" ON "employee_work_ranges" ("deleted_at");
-- Create "employee_work_range_services" table
CREATE TABLE "employee_work_range_services" ("employee_work_range_id" uuid NOT NULL DEFAULT gen_random_uuid(), "service_id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("employee_work_range_id", "service_id"), CONSTRAINT "fk_employee_work_range_services_employee_work_range" FOREIGN KEY ("employee_work_range_id") REFERENCES "employee_work_ranges" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "fk_employee_work_range_services_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE CASCADE ON DELETE CASCADE);
-- Create "employee_branches" table
CREATE TABLE "employee_branches" ("employee_id" uuid NOT NULL DEFAULT gen_random_uuid(), "branch_id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("employee_id", "branch_id"), CONSTRAINT "fk_employee_branches_branch" FOREIGN KEY ("branch_id") REFERENCES "branches" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_employee_branches_employee" FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "branch_work_ranges" table
CREATE TABLE "branch_work_ranges" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "weekday" bigint NOT NULL, "start_time" timestamptz NOT NULL, "end_time" timestamptz NOT NULL, "time_zone" character varying(255) NOT NULL, "branch_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_branches_work_schedule" FOREIGN KEY ("branch_id") REFERENCES "branches" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_branch_id" to table: "branch_work_ranges"
CREATE INDEX "idx_branch_id" ON "branch_work_ranges" ("branch_id");
-- Create index "idx_branch_work_ranges_deleted_at" to table: "branch_work_ranges"
CREATE INDEX "idx_branch_work_ranges_deleted_at" ON "branch_work_ranges" ("deleted_at");
-- Create "companies" table
CREATE TABLE "companies" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "legal_name" character varying(100) NULL, "trade_name" character varying(100) NULL, "tax_id" character varying(100) NULL, "schema_name" character varying(100) NULL, "design" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_public_companies_deleted_at" to table: "companies"
CREATE INDEX "idx_public_companies_deleted_at" ON "companies" ("deleted_at");
-- Create index "idx_public_companies_legal_name" to table: "companies"
CREATE UNIQUE INDEX "idx_public_companies_legal_name" ON "companies" ("legal_name");
-- Create index "idx_public_companies_schema_name" to table: "companies"
CREATE UNIQUE INDEX "idx_public_companies_schema_name" ON "companies" ("schema_name");
-- Create index "idx_public_companies_tax_id" to table: "companies"
CREATE UNIQUE INDEX "idx_public_companies_tax_id" ON "companies" ("tax_id");
-- Create index "idx_public_companies_trade_name" to table: "companies"
CREATE UNIQUE INDEX "idx_public_companies_trade_name" ON "companies" ("trade_name");
-- Create "sectors" table
CREATE TABLE "sectors" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" text NOT NULL, "description" text NULL, PRIMARY KEY ("id"));
-- Create index "idx_public_sectors_deleted_at" to table: "sectors"
CREATE INDEX "idx_public_sectors_deleted_at" ON "sectors" ("deleted_at");
-- Create index "uni_public_sectors_name" to table: "sectors"
CREATE UNIQUE INDEX "uni_public_sectors_name" ON "sectors" ("name");
-- Create "company_sectors" table
CREATE TABLE "company_sectors" ("company_id" uuid NOT NULL DEFAULT gen_random_uuid(), "sector_id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("company_id", "sector_id"), CONSTRAINT "fk_company_sectors_company" FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_company_sectors_sector" FOREIGN KEY ("sector_id") REFERENCES "sectors" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "resources" table
CREATE TABLE "resources" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" text NOT NULL, "description" text NULL, "table" text NULL, "references" jsonb NULL, PRIMARY KEY ("id"));
-- Create index "idx_public_resources_deleted_at" to table: "resources"
CREATE INDEX "idx_public_resources_deleted_at" ON "resources" ("deleted_at");
-- Create index "uni_public_resources_name" to table: "resources"
CREATE UNIQUE INDEX "uni_public_resources_name" ON "resources" ("name");
-- Create "properties" table
CREATE TABLE "properties" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" text NOT NULL, "description" text NULL, "resource_name" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_public_properties_resource" FOREIGN KEY ("resource_name") REFERENCES "resources" ("id") ON UPDATE CASCADE ON DELETE SET NULL);
-- Create index "idx_public_properties_deleted_at" to table: "properties"
CREATE INDEX "idx_public_properties_deleted_at" ON "properties" ("deleted_at");
-- Create index "uni_public_properties_name" to table: "properties"
CREATE UNIQUE INDEX "uni_public_properties_name" ON "properties" ("name");
-- Create "payments" table
CREATE TABLE "payments" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "price" bigint NOT NULL, "currency" character varying(3) NOT NULL DEFAULT 'BRL', "status" character varying(20) NOT NULL DEFAULT 'PENDING', "payment_method" character varying(50) NULL, "transaction_id" character varying(100) NULL, "provider" character varying(50) NULL, "user_id" bigint NULL, "order_id" bigint NULL, "metadata" jsonb NULL, "completed_at" timestamptz NULL, "failed_at" timestamptz NULL, PRIMARY KEY ("id"));
-- Create index "idx_payments_deleted_at" to table: "payments"
CREATE INDEX "idx_payments_deleted_at" ON "payments" ("deleted_at");
-- Create index "idx_payments_order_id" to table: "payments"
CREATE INDEX "idx_payments_order_id" ON "payments" ("order_id");
-- Create index "idx_payments_payment_method" to table: "payments"
CREATE INDEX "idx_payments_payment_method" ON "payments" ("payment_method");
-- Create index "idx_payments_provider" to table: "payments"
CREATE INDEX "idx_payments_provider" ON "payments" ("provider");
-- Create index "idx_payments_status" to table: "payments"
CREATE INDEX "idx_payments_status" ON "payments" ("status");
-- Create index "idx_payments_transaction_id" to table: "payments"
CREATE UNIQUE INDEX "idx_payments_transaction_id" ON "payments" ("transaction_id");
-- Create index "idx_payments_user_id" to table: "payments"
CREATE INDEX "idx_payments_user_id" ON "payments" ("user_id");
-- Create "appointments" table
CREATE TABLE "appointments" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "service_id" uuid NOT NULL, "employee_id" uuid NOT NULL, "client_id" uuid NOT NULL, "branch_id" uuid NOT NULL, "payment_id" uuid NULL, "company_id" uuid NOT NULL, "cancelled_employee_id" uuid NULL, "start_time" timestamptz NOT NULL, "end_time" timestamptz NOT NULL, "time_zone" character varying(100) NOT NULL, "actual_start_time" timestamptz NOT NULL, "actual_end_time" timestamptz NOT NULL, "cancel_time" timestamptz NOT NULL, "is_fulfilled" boolean NULL DEFAULT false, "is_cancelled" boolean NULL DEFAULT false, "is_cancelled_by_client" boolean NULL DEFAULT false, "is_cancelled_by_employee" boolean NULL DEFAULT false, "is_confirmed_by_client" boolean NULL DEFAULT false, "history" jsonb NULL, "comments" jsonb NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_appointments_payment" FOREIGN KEY ("payment_id") REFERENCES "payments" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_appointments_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_branches_appointments" FOREIGN KEY ("branch_id") REFERENCES "branches" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_employees_appointments" FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_appointments_client_id" to table: "appointments"
CREATE INDEX "idx_appointments_client_id" ON "appointments" ("client_id");
-- Create index "idx_appointments_company_id" to table: "appointments"
CREATE INDEX "idx_appointments_company_id" ON "appointments" ("company_id");
-- Create index "idx_appointments_deleted_at" to table: "appointments"
CREATE INDEX "idx_appointments_deleted_at" ON "appointments" ("deleted_at");
-- Create index "idx_appointments_payment_id" to table: "appointments"
CREATE UNIQUE INDEX "idx_appointments_payment_id" ON "appointments" ("payment_id");
-- Create "branch_work_range_services" table
CREATE TABLE "branch_work_range_services" ("branch_work_range_id" uuid NOT NULL DEFAULT gen_random_uuid(), "service_id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("branch_work_range_id", "service_id"), CONSTRAINT "fk_branch_work_range_services_branch_work_range" FOREIGN KEY ("branch_work_range_id") REFERENCES "branch_work_ranges" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "fk_branch_work_range_services_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE CASCADE ON DELETE CASCADE);
-- Create "employee_services" table
CREATE TABLE "employee_services" ("service_id" uuid NOT NULL DEFAULT gen_random_uuid(), "employee_id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("service_id", "employee_id"), CONSTRAINT "fk_employee_services_employee" FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_employee_services_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "endpoints" table
CREATE TABLE "endpoints" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "controller_name" character varying(100) NULL, "description" text NULL, "method" character varying(6) NULL, "path" text NULL, "deny_unauthorized" boolean NULL DEFAULT false, "needs_company_id" boolean NULL DEFAULT false, "resource_id" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_public_endpoints_resource" FOREIGN KEY ("resource_id") REFERENCES "resources" ("id") ON UPDATE CASCADE ON DELETE SET NULL);
-- Create index "idx_public_endpoints_deleted_at" to table: "endpoints"
CREATE INDEX "idx_public_endpoints_deleted_at" ON "endpoints" ("deleted_at");
-- Create "policy_rules" table
CREATE TABLE "policy_rules" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" text NULL, "description" text NULL, "effect" text NULL, "end_point_id" uuid NULL, "conditions" jsonb NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_public_policy_rules_end_point" FOREIGN KEY ("end_point_id") REFERENCES "endpoints" ("id") ON UPDATE CASCADE ON DELETE SET NULL);
-- Create index "idx_public_policy_rules_deleted_at" to table: "policy_rules"
CREATE INDEX "idx_public_policy_rules_deleted_at" ON "policy_rules" ("deleted_at");
-- Create "subdomains" table
CREATE TABLE "subdomains" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "name" character varying(36) NOT NULL, "company_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_public_companies_subdomains" FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_public_subdomains_company_id" to table: "subdomains"
CREATE INDEX "idx_public_subdomains_company_id" ON "subdomains" ("company_id");
-- Create index "idx_public_subdomains_deleted_at" to table: "subdomains"
CREATE INDEX "idx_public_subdomains_deleted_at" ON "subdomains" ("deleted_at");
-- Create index "idx_public_subdomains_name" to table: "subdomains"
CREATE UNIQUE INDEX "idx_public_subdomains_name" ON "subdomains" ("name");
