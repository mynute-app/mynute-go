package main

import (
	"fmt"
	"log"
	database "mynute-go/core/src/config/db"
	"mynute-go/core/src/lib"
)

func main() {
	log.Println("Fixing employee appointments endpoint path...")

	// Load environment variables
	lib.LoadEnv()

	// Connect to database
	db := database.Connect()
	defer db.Disconnect()

	// Update the endpoint path
	result := db.Gorm.Exec(`
		UPDATE public.endpoints
		SET path = '/employee/:employee_id/appointments'
		WHERE method = 'GET' 
		  AND path = '/employee/:id/appointments'
		  AND controller_name = 'GetEmployeeAppointmentsById'
	`)

	if result.Error != nil {
		log.Fatalf("Failed to update endpoint: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Println("No rows updated. The endpoint path may already be correct or the endpoint doesn't exist.")
	} else {
		log.Printf("✓ Successfully updated endpoint path! (%d row(s) affected)\n", result.RowsAffected)
	}

	// Verify the fix
	var count int64
	db.Gorm.Raw(`
		SELECT COUNT(*) 
		FROM public.endpoints 
		WHERE method = 'GET' 
		  AND path = '/employee/:employee_id/appointments'
		  AND controller_name = 'GetEmployeeAppointmentsById'
	`).Scan(&count)

	if count > 0 {
		fmt.Println("\n✅ Endpoint is now correctly configured:")
		fmt.Println("   Path: /employee/:employee_id/appointments")
		fmt.Println("   Method: GET")
		fmt.Println("   Controller: GetEmployeeAppointmentsById")
	} else {
		fmt.Println("\n⚠️  Warning: Could not verify the endpoint configuration")
	}
}
