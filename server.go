package main

import (
	configapi "agenda-kaki-go/core/config/api"
	configdb "agenda-kaki-go/core/config/db"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Initialize the database
	db := configdb.ConnectDB()

	// Initialize the router
	configapi.BuildRouter(db, app)

	// Close the database connection when the app closes
	defer configdb.CloseDB(db)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
