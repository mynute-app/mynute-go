package main

import (
	"agenda-kaki-go/api"
	"agenda-kaki-go/api/config"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Initialize the database
	db := config.ConnectDB()

	// Close the database connection when the app closes
	defer config.CloseDB(db)

	// Initialize the API content
	api.Build(db, app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
