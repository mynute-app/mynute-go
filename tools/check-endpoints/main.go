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
		os.Getenv("POSTGRES_DB_PROD"),
	)

	// Connect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// Query
	rows, err := db.Query(`
		SELECT path, method, controller_name 
		FROM public.endpoints 
		WHERE (controller_name LIKE '%Employee%' AND (path LIKE '%work%' OR path LIKE '%appointment%'))
		OR path = '/employee/:employee_id/appointments'
		ORDER BY controller_name, method, path
	`)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                          Employee Endpoint Paths                             ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ %-40s %-8s %-30s ║\n", "PATH", "METHOD", "CONTROLLER")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")

	for rows.Next() {
		var path, method, controller string
		if err := rows.Scan(&path, &method, &controller); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("║ %-40s %-8s %-30s ║\n", path, method, controller)
	}

	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}
