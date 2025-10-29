package main

import (
	"flag"
	"fmt"
	"log"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"os"
	"path/filepath"
	"strings"
)

// This tool generates SQL migration files from GORM models automatically
// Usage: go run tools/generate-migration/main.go -name add_new_field -models Employee,Branch

func main() {
	var (
		migrationName string
		modelsStr     string
		includeDown   bool
	)

	flag.StringVar(&migrationName, "name", "", "Migration name (required)")
	flag.StringVar(&modelsStr, "models", "all", "Comma-separated list of models to migrate (or 'all')")
	flag.BoolVar(&includeDown, "down", true, "Generate DOWN migration (default: true)")
	flag.Parse()

	if migrationName == "" {
		log.Fatal("Error: -name is required\nUsage: go run tools/generate-migration/main.go -name migration_name [-models ModelName1,ModelName2]")
	}

	lib.LoadEnv()

	// Determine which models to include
	var modelsToMigrate []any
	if modelsStr == "all" {
		// All models from GeneralModels and TenantModels
		modelsToMigrate = append(modelsToMigrate, model.GeneralModels...)
		modelsToMigrate = append(modelsToMigrate, model.TenantModels...)
	} else {
		// Parse specific models
		modelNames := strings.Split(modelsStr, ",")
		modelsToMigrate = getModelsByNames(modelNames)
	}

	if len(modelsToMigrate) == 0 {
		log.Fatal("No models found to migrate")
	}

	// Generate migration files
	timestamp := lib.GetTimestampVersion()
	upFile := filepath.Join("migrations", fmt.Sprintf("%s_%s.up.sql", timestamp, migrationName))
	downFile := filepath.Join("migrations", fmt.Sprintf("%s_%s.down.sql", timestamp, migrationName))

	// Generate UP SQL
	upSQL := generateUpSQL(modelsToMigrate)
	if err := os.WriteFile(upFile, []byte(upSQL), 0644); err != nil {
		log.Fatalf("Failed to write UP migration: %v", err)
	}

	// Generate DOWN SQL if requested
	if includeDown {
		downSQL := generateDownSQL(modelsToMigrate)
		if err := os.WriteFile(downFile, []byte(downSQL), 0644); err != nil {
			log.Fatalf("Failed to write DOWN migration: %v", err)
		}
		log.Printf("✅ Generated migration files:\n  %s\n  %s\n", upFile, downFile)
	} else {
		log.Printf("✅ Generated migration file:\n  %s\n", upFile)
	}

	log.Println("\n⚠️  IMPORTANT: Review and adjust the generated SQL before applying!")
	log.Println("   - Verify column types match your expectations")
	log.Println("   - Add data migrations if needed")
	log.Println("   - Consider adding indexes")
	log.Println("   - For multi-tenant tables, ensure schema iteration is correct")
	log.Println("\n💡 Next steps:")
	log.Println("   1. Edit the generated SQL files")
	log.Println("   2. Run: make test-migrate")
	log.Println("   3. If tests pass, commit your changes!")
}

func getModelsByNames(names []string) []any {
	modelMap := map[string]any{
		"Employee":               &model.Employee{},
		"Branch":                 &model.Branch{},
		"Service":                &model.Service{},
		"Appointment":            &model.Appointment{},
		"Client":                 &model.Client{},
		"Company":                &model.Company{},
		"Sector":                 &model.Sector{},
		"Holiday":                &model.Holiday{},
		"Role":                   &model.Role{},
		"Resource":               &model.Resource{},
		"Property":               &model.Property{},
		"EndPoint":               &model.EndPoint{},
		"PolicyRule":             &model.PolicyRule{},
		"Subdomain":              &model.Subdomain{},
		"BranchWorkRange":        &model.BranchWorkRange{},
		"EmployeeWorkRange":      &model.EmployeeWorkRange{},
		"BranchServiceDensity":   &model.BranchServiceDensity{},
		"EmployeeServiceDensity": &model.EmployeeServiceDensity{},
		"Payment":                &model.Payment{},
		"AppointmentArchive":     &model.AppointmentArchive{},
		"ClientAppointment":      &model.ClientAppointment{},
	}

	var result []any
	for _, name := range names {
		name = strings.TrimSpace(name)
		if m, ok := modelMap[name]; ok {
			result = append(result, m)
		} else {
			log.Printf("Warning: Model '%s' not found, skipping", name)
		}
	}
	return result
}

func generateUpSQL(models []any) string {
	var sql strings.Builder
	sql.WriteString("-- Auto-generated migration\n")
	sql.WriteString(fmt.Sprintf("-- Generated at: %s\n", lib.GetTimestampVersion()))
	sql.WriteString("-- ⚠️  REVIEW THIS SQL BEFORE APPLYING!\n\n")

	for _, m := range models {
		// Get table name and schema type
		tableName := getTableName(m)
		schemaType := getSchemaType(m)

		sql.WriteString(fmt.Sprintf("-- Model: %s (Schema: %s)\n", getModelName(m), schemaType))

		if schemaType == "company" {
			// Multi-tenant table - iterate over all company schemas
			sql.WriteString("DO $$\n")
			sql.WriteString("DECLARE\n")
			sql.WriteString("    schema_name TEXT;\n")
			sql.WriteString("BEGIN\n")
			sql.WriteString("    FOR schema_name IN \n")
			sql.WriteString("        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'\n")
			sql.WriteString("    LOOP\n")
			sql.WriteString(fmt.Sprintf("        -- Add your ALTER TABLE statements here for %s\n", tableName))
			sql.WriteString(fmt.Sprintf("        -- Example: EXECUTE format('ALTER TABLE %%I.%s ADD COLUMN new_field TEXT', schema_name);\n", tableName))
			sql.WriteString("    END LOOP;\n")
			sql.WriteString("END $$;\n\n")
		} else {
			// Public schema table
			fullTableName := fmt.Sprintf("public.%s", tableName)
			sql.WriteString(fmt.Sprintf("-- Add your ALTER TABLE statements here for %s\n", fullTableName))
			sql.WriteString(fmt.Sprintf("-- Example: ALTER TABLE %s ADD COLUMN new_field TEXT;\n\n", fullTableName))
		}
	}

	sql.WriteString("\n-- 💡 Tips:\n")
	sql.WriteString("-- - Use 'IF NOT EXISTS' / 'IF EXISTS' for idempotency\n")
	sql.WriteString("-- - Add indexes with 'CREATE INDEX CONCURRENTLY' in production\n")
	sql.WriteString("-- - Test on a copy of production data first\n")

	return sql.String()
}

func generateDownSQL(models []any) string {
	var sql strings.Builder
	sql.WriteString("-- Auto-generated rollback migration\n")
	sql.WriteString(fmt.Sprintf("-- Generated at: %s\n", lib.GetTimestampVersion()))
	sql.WriteString("-- ⚠️  REVIEW THIS SQL BEFORE APPLYING!\n\n")

	for _, m := range models {
		tableName := getTableName(m)
		schemaType := getSchemaType(m)

		sql.WriteString(fmt.Sprintf("-- Rollback Model: %s (Schema: %s)\n", getModelName(m), schemaType))

		if schemaType == "company" {
			sql.WriteString("DO $$\n")
			sql.WriteString("DECLARE\n")
			sql.WriteString("    schema_name TEXT;\n")
			sql.WriteString("BEGIN\n")
			sql.WriteString("    FOR schema_name IN \n")
			sql.WriteString("        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'\n")
			sql.WriteString("    LOOP\n")
			sql.WriteString(fmt.Sprintf("        -- Add your rollback statements here for %s\n", tableName))
			sql.WriteString(fmt.Sprintf("        -- Example: EXECUTE format('ALTER TABLE %%I.%s DROP COLUMN IF EXISTS new_field', schema_name);\n", tableName))
			sql.WriteString("    END LOOP;\n")
			sql.WriteString("END $$;\n\n")
		} else {
			fullTableName := fmt.Sprintf("public.%s", tableName)
			sql.WriteString(fmt.Sprintf("-- Add your rollback statements here for %s\n", fullTableName))
			sql.WriteString(fmt.Sprintf("-- Example: ALTER TABLE %s DROP COLUMN IF EXISTS new_field;\n\n", fullTableName))
		}
	}

	return sql.String()
}

func getModelName(m any) string {
	// Extract type name without package prefix
	typeStr := fmt.Sprintf("%T", m)
	parts := strings.Split(typeStr, ".")
	return strings.TrimPrefix(parts[len(parts)-1], "*")
}

func getTableName(m any) string {
	// Check if model has TableName method
	if tn, ok := m.(interface{ TableName() string }); ok {
		tableName := tn.TableName()
		// Remove schema prefix if exists (e.g., "public.clients" -> "clients")
		parts := strings.Split(tableName, ".")
		return parts[len(parts)-1]
	}
	// Fallback: convert model name to snake_case plural
	modelName := getModelName(m)
	return strings.ToLower(modelName) + "s"
}

func getSchemaType(m any) string {
	// Check if model has SchemaType method
	if st, ok := m.(interface{ SchemaType() string }); ok {
		return st.SchemaType()
	}
	return "public"
}
