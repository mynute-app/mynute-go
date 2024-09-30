package main

import (
	"agenda-kaki-company-go/api/config"
	"agenda-kaki-company-go/api/routes"
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

    // Initialize the routes
    routes.Init(db, app)

    // Start the server on port 3000
    log.Fatal(app.Listen(":3000"))
}