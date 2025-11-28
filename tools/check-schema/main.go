package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// Get first company schema
	var schemaName string
	err = db.QueryRow("SELECT nspname FROM pg_namespace WHERE nspname LIKE 'company_%' LIMIT 1").Scan(&schemaName)
	if err != nil {
		log.Fatal("No company schema found:", err)
	}

	fmt.Printf("\n📊 Columns in %s.employees table:\n\n", schemaName)

	// Query columns
	query := fmt.Sprintf(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = '%s'
		AND table_name = 'employees'
		ORDER BY ordinal_position
	`, schemaName)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Printf("%-20s %-25s %-12s %s\n", "COLUMN", "TYPE", "NULLABLE", "DEFAULT")
	fmt.Println(string(make([]byte, 80)))

	for rows.Next() {
		var colName, dataType, nullable string
		var colDefault sql.NullString
		if err := rows.Scan(&colName, &dataType, &nullable, &colDefault); err != nil {
			log.Fatal(err)
		}
		def := "NULL"
		if colDefault.Valid {
			def = colDefault.String
		}
		fmt.Printf("%-20s %-25s %-12s %s\n", colName, dataType, nullable, def)
	}
	fmt.Println()
}
