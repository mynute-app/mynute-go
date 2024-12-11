package main

import (
	"agenda-kaki-go/core/config/api/routes"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/handlers"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Initialize the database
	db := database.Connect()

	// Close the database connection when the app closes
	defer db.CloseDB()
	handlers.NewAuth()
	// Migrate the database
	db.Migrate()

	// Initialize the router
	routes.Build(db.Gorm, app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
