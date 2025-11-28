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

	rows, err := db.Query("SELECT version, dirty FROM schema_migrations ORDER BY version")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Println("\n📋 schema_migrations table:")
	fmt.Printf("%-20s %s\n", "VERSION", "DIRTY")
	fmt.Println("────────────────────────────")

	count := 0
	for rows.Next() {
		var version int64
		var dirty bool
		if err := rows.Scan(&version, &dirty); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20d %t\n", version, dirty)
		count++
	}
	fmt.Printf("\n✅ Total: %d row(s)\n\n", count)
}
