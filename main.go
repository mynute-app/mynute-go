package main

import (
	"agenda-kaki-go/core/config/api/routes"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/handlers"
	"log"
	_ "github.com/gofiber/swagger/example/docs"
	"github.com/gofiber/fiber/v2"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Initialize the database
	db := database.Connect()

	// Close the database connection when the app closes
	defer db.CloseDB()

	//Initialize Auth handlers
	session := handlers.NewCookieStore(handlers.SessionOpts())
	handlers.NewAuth(session)

	// Migrate the database
	db.Migrate()

	// Initialize the router
	routes.Build(db.Gorm, app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":4000"))
}
