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
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

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

	// List all schemas
	rows, err := db.Query("SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' AND nspname != 'information_schema' ORDER BY nspname")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Println("\n📂 Schemas in database:")
	companyCount := 0
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  - %s\n", schema)
		if len(schema) > 8 && schema[:8] == "company_" {
			companyCount++
		}
	}
	fmt.Printf("\n✅ Found %d company schema(s)\n\n", companyCount)
}
