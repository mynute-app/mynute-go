//go:build ignore
// +build ignore

package main

import (
	"log"
	"mynute-go/services/core"
	_ "mynute-go/services/core/docs"
	"os"
)

// @title						Business Service API
// @version					1.0
// @description				Main Business Logic Service
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-Auth-Token
// @description				Enter the token in the format: <token>
// @contact.name				API Support
// @contact.email				support@mynute.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:4000
// @BasePath					/
func main() {
	log.Println("Starting Business Service...")

	// Override the default port for business service
	if os.Getenv("APP_PORT") == "" {
		os.Setenv("APP_PORT", "4000")
	}

	// Create and run the server using existing core.NewServer()
	// This maintains backward compatibility with your existing setup
	server := core.NewServer()
	server.Run("listen")
}
