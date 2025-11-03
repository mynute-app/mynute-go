package main

import (
	"flag"
	"fmt"
	"log"
	authModel "mynute-go/auth/model"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Smart migration generator - detects schema changes automatically
// Usage: go run tools/smart-migration/main.go -name add_new_field -models Employee

type ColumnInfo struct {
	Name     string
	Type     string
	Nullable bool
	Default  *string
}

func main() {
	var (
		migrationName string
		modelsStr     string
		schemaName    string
	)

	flag.StringVar(&migrationName, "name", "", "Migration name (required)")
	flag.StringVar(&modelsStr, "models", "", "Comma-separated list of models to analyze, or 'all' for all models (required)")
	flag.StringVar(&schemaName, "schema", "", "Schema to compare (leave empty for first company schema)")
	flag.Parse()

	if migrationName == "" || modelsStr == "" {
		log.Fatal("Error: -name and -models are required\nUsage: go run tools/smart-migration/main.go -name migration_name -models ModelName1,ModelName2\n       or: -models all")
	}

	lib.LoadEnv()

	// Connect to database
	db := connectToDatabase()

	// Parse models
	var modelNames []string
	if strings.ToLower(strings.TrimSpace(modelsStr)) == "all" {
		modelNames = getAllModelNames()
		log.Printf("🔍 Analyzing ALL models (%d total)\n", len(modelNames))
	} else {
		modelNames = strings.Split(modelsStr, ",")
	}
	modelsToAnalyze := getModelsByNames(modelNames)

	if len(modelsToAnalyze) == 0 {
		log.Fatal("No models found to analyze")
	}

	// Detect target schema
	if schemaName == "" {
		schemaName = findFirstCompanySchema(db)
		if schemaName == "" {
			log.Println("⚠️  No company schemas found. Using public schema for comparison.")
			schemaName = "public"
		} else {
			log.Printf("📊 Using schema '%s' for comparison\n", schemaName)
		}
	}

	// Generate migrations
	timestamp := lib.GetTimestampVersion()
	upFile := filepath.Join("migrations", fmt.Sprintf("%s_%s.up.sql", timestamp, migrationName))
	downFile := filepath.Join("migrations", fmt.Sprintf("%s_%s.down.sql", timestamp, migrationName))

	upSQL, downSQL := generateSmartMigrations(db, modelsToAnalyze, schemaName)

	// Write files
	if err := os.WriteFile(upFile, []byte(upSQL), 0644); err != nil {
		log.Fatalf("Failed to write UP migration: %v", err)
	}

	if err := os.WriteFile(downFile, []byte(downSQL), 0644); err != nil {
		log.Fatalf("Failed to write DOWN migration: %v", err)
	}

	log.Printf("✅ Generated smart migration files:\n  %s\n  %s\n", upFile, downFile)
	log.Println("\n💡 Changes detected and SQL generated automatically!")
	log.Println("⚠️  IMPORTANT: Review the generated SQL before applying!")
	log.Println("\n� How this works:")
	log.Println("   - Reads Go model struct definitions")
	log.Println("   - Queries current database schema")
	log.Println("   - Compares: Model fields vs Database columns")
	log.Println("   - Does NOT track migration file history")
	log.Println("\n�� Next steps:")
	log.Println("   1. Review the generated SQL files")
	log.Println("   2. Verify these are truly new changes")
	log.Println("   3. Run: make test-migrate")
	log.Println("   4. If tests pass, commit your changes!")
}

func connectToDatabase() *gorm.DB {
	// Get DB config from environment
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if dbName == "" {
		log.Fatal("POSTGRES_DB environment variable is required")
	}

	log.Printf("Connecting to database: %s\n", dbName)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbName, port)

	db, err := gorm.Open(lib.PostgresDialector(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func findFirstCompanySchema(db *gorm.DB) string {
	var schemaName string
	err := db.Raw("SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%' LIMIT 1").Scan(&schemaName).Error
	if err != nil || schemaName == "" {
		return ""
	}
	return schemaName
}

func generateSmartMigrations(db *gorm.DB, models []any, schemaName string) (string, string) {
	var upSQL, downSQL strings.Builder

	upSQL.WriteString("-- Smart migration - Auto-detected changes\n")
	upSQL.WriteString(fmt.Sprintf("-- Generated at: %s\n", lib.GetTimestampVersion()))
	upSQL.WriteString(fmt.Sprintf("-- Compared against schema: %s\n", schemaName))
	upSQL.WriteString("--\n")
	upSQL.WriteString("-- ⚠️  IMPORTANT: This tool compares Go models against the CURRENT database schema.\n")
	upSQL.WriteString("--     It does NOT track migration history. A 'new' column means:\n")
	upSQL.WriteString("--     - The Go model has this field, AND\n")
	upSQL.WriteString("--     - The database table does NOT have this column\n")
	upSQL.WriteString("--\n")
	upSQL.WriteString("--     Review carefully before applying!\n")
	upSQL.WriteString("\n")

	downSQL.WriteString("-- Smart migration rollback - Auto-detected changes\n")
	downSQL.WriteString(fmt.Sprintf("-- Generated at: %s\n\n", lib.GetTimestampVersion()))

	hasChanges := false

	for _, m := range models {
		modelName := getModelName(m)
		tableName := getTableName(m)
		schemaType := getSchemaType(m)

		upSQL.WriteString(fmt.Sprintf("-- Model: %s (Table: %s, Schema: %s)\n", modelName, tableName, schemaType))
		upSQL.WriteString("-- Comparison: Go struct fields vs Current database columns\n")
		downSQL.WriteString(fmt.Sprintf("-- Rollback Model: %s\n", modelName))

		// Get expected columns from GORM model
		expectedCols := getGormColumns(db, m)

		// Get actual columns from database
		fullTableName := tableName
		// Both "company" and "tenant" schema types should use company_* schemas
		if schemaType == "company" || schemaType == "tenant" {
			fullTableName = schemaName + "." + tableName
		} else {
			fullTableName = "public." + tableName
		}

		actualCols := getDatabaseColumns(db, fullTableName)

		// Find differences
		addedCols, removedCols := compareColumns(expectedCols, actualCols)

		if len(addedCols) > 0 || len(removedCols) > 0 {
			hasChanges = true
		}

		// Generate ALTER TABLE statements for added columns
		if len(addedCols) > 0 {
			upSQL.WriteString(fmt.Sprintf("-- Adding %d new column(s)\n", len(addedCols)))

			// Both "company" and "tenant" schema types should use company_* schema pattern
			if schemaType == "company" || schemaType == "tenant" {
				upSQL.WriteString("DO $$\n")
				upSQL.WriteString("DECLARE\n")
				upSQL.WriteString("    schema_name TEXT;\n")
				upSQL.WriteString("BEGIN\n")
				upSQL.WriteString("    FOR schema_name IN \n")
				upSQL.WriteString("        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'\n")
				upSQL.WriteString("    LOOP\n")

				for _, col := range addedCols {
					nullable := "NOT NULL"
					if col.Nullable {
						nullable = "NULL"
					}
					upSQL.WriteString(fmt.Sprintf("        EXECUTE format('ALTER TABLE %%I.%s ADD COLUMN IF NOT EXISTS %s %s %s', schema_name);\n",
						tableName, col.Name, col.Type, nullable))
				}

				upSQL.WriteString("    END LOOP;\n")
				upSQL.WriteString("END $$;\n\n")
			} else {
				for _, col := range addedCols {
					nullable := "NOT NULL"
					if col.Nullable {
						nullable = "NULL"
					}
					upSQL.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s %s;\n",
						fullTableName, col.Name, col.Type, nullable))
				}
				upSQL.WriteString("\n")
			}
		}

		// Generate DROP COLUMN for removed columns (for DOWN migration)
		if len(removedCols) > 0 {
			downSQL.WriteString(fmt.Sprintf("-- Removing %d column(s) that were added\n", len(addedCols)))

			// Both "company" and "tenant" schema types should use company_* schema pattern
			if schemaType == "company" || schemaType == "tenant" {
				downSQL.WriteString("DO $$\n")
				downSQL.WriteString("DECLARE\n")
				downSQL.WriteString("    schema_name TEXT;\n")
				downSQL.WriteString("BEGIN\n")
				downSQL.WriteString("    FOR schema_name IN \n")
				downSQL.WriteString("        SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%'\n")
				downSQL.WriteString("    LOOP\n")

				for _, col := range addedCols {
					downSQL.WriteString(fmt.Sprintf("        EXECUTE format('ALTER TABLE %%I.%s DROP COLUMN IF EXISTS %s', schema_name);\n",
						tableName, col.Name))
				}

				downSQL.WriteString("    END LOOP;\n")
				downSQL.WriteString("END $$;\n\n")
			} else {
				for _, col := range addedCols {
					downSQL.WriteString(fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s;\n", fullTableName, col.Name))
				}
				downSQL.WriteString("\n")
			}
		}

		if len(addedCols) == 0 && len(removedCols) == 0 {
			upSQL.WriteString("-- No changes detected for this model\n\n")
			downSQL.WriteString("-- No changes to rollback\n\n")
		}
	}

	if !hasChanges {
		upSQL.WriteString("\n⚠️  WARNING: No schema changes detected!\n")
		upSQL.WriteString("-- Either:\n")
		upSQL.WriteString("--   1. Your models match the database schema\n")
		upSQL.WriteString("--   2. The database doesn't have these tables yet\n")
		upSQL.WriteString("--   3. You need to create the initial schema first\n")
	}

	return upSQL.String(), downSQL.String()
}

func getGormColumns(db *gorm.DB, model any) map[string]ColumnInfo {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(model)

	cols := make(map[string]ColumnInfo)
	for _, field := range stmt.Schema.Fields {
		if field.DBName == "" {
			continue
		}

		col := ColumnInfo{
			Name:     field.DBName,
			Type:     getPostgresType(field),
			Nullable: !field.NotNull,
		}
		cols[field.DBName] = col
	}

	return cols
}

func getDatabaseColumns(db *gorm.DB, tableName string) map[string]ColumnInfo {
	parts := strings.Split(tableName, ".")
	var schema, table string
	if len(parts) == 2 {
		schema = parts[0]
		table = parts[1]
	} else {
		schema = "public"
		table = tableName
	}

	type DBColumn struct {
		ColumnName string
		DataType   string
		IsNullable string
	}

	var dbCols []DBColumn
	err := db.Raw(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
	`, schema, table).Scan(&dbCols).Error

	if err != nil {
		log.Printf("Warning: Could not query columns for %s: %v", tableName, err)
		return make(map[string]ColumnInfo)
	}

	cols := make(map[string]ColumnInfo)
	for _, dbCol := range dbCols {
		cols[dbCol.ColumnName] = ColumnInfo{
			Name:     dbCol.ColumnName,
			Type:     dbCol.DataType,
			Nullable: dbCol.IsNullable == "YES",
		}
	}

	return cols
}

func compareColumns(expected, actual map[string]ColumnInfo) (added, removed []ColumnInfo) {
	// Find added columns (in expected but not in actual)
	for name, col := range expected {
		if _, exists := actual[name]; !exists {
			added = append(added, col)
		}
	}

	// Find removed columns (in actual but not in expected)
	for name, col := range actual {
		if _, exists := expected[name]; !exists {
			removed = append(removed, col)
		}
	}

	return added, removed
}

func getPostgresType(field *schema.Field) string {
	// Check if there's a type tag specified in gorm
	if field.Tag.Get("type") != "" {
		return strings.ToUpper(field.Tag.Get("type"))
	}

	// Map GORM types to PostgreSQL types based on Go type
	switch field.FieldType.Kind() {
	case reflect.String:
		// Check for varchar tag
		if size := field.Tag.Get("size"); size != "" {
			return fmt.Sprintf("VARCHAR(%s)", size)
		}
		return "TEXT"
	case reflect.Int, reflect.Int32:
		return "INTEGER"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		return "DOUBLE PRECISION"
	default:
		// Check if it's a UUID or Time type
		typeName := field.FieldType.Name()
		if typeName == "UUID" || strings.Contains(field.FieldType.String(), "uuid.UUID") {
			return "UUID"
		}
		if typeName == "Time" || strings.Contains(field.FieldType.String(), "time.Time") {
			return "TIMESTAMP WITH TIME ZONE"
		}
		return "TEXT"
	}
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
		"Resource":               &authModel.Resource{},
		"Property":               &authModel.Property{},
		"EndPoint":               &authModel.EndPoint{},
		"PolicyRule":             &authModel.PolicyRule{},
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

func getAllModelNames() []string {
	return []string{
		"Employee",
		"Branch",
		"Service",
		"Appointment",
		"Client",
		"Company",
		"Sector",
		"Holiday",
		"Role",
		"Resource",
		"Property",
		"EndPoint",
		"PolicyRule",
		"Subdomain",
		"BranchWorkRange",
		"EmployeeWorkRange",
		"BranchServiceDensity",
		"EmployeeServiceDensity",
		"Payment",
		"AppointmentArchive",
		"ClientAppointment",
	}
}

func getModelName(m any) string {
	typeStr := fmt.Sprintf("%T", m)
	parts := strings.Split(typeStr, ".")
	return strings.TrimPrefix(parts[len(parts)-1], "*")
}

func getTableName(m any) string {
	if tn, ok := m.(interface{ TableName() string }); ok {
		tableName := tn.TableName()
		parts := strings.Split(tableName, ".")
		return parts[len(parts)-1]
	}
	modelName := getModelName(m)
	return strings.ToLower(modelName) + "s"
}

func getSchemaType(m any) string {
	if st, ok := m.(interface{ SchemaType() string }); ok {
		return st.SchemaType()
	}
	return "public"
}
