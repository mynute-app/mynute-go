//go:build ignore
// +build ignore

package main

import (
	_ "mynute-go/services/email/docs"
	"mynute-go/services/email"
)

// @title						Email Service API
// @version					1.0
// @description				Email Microservice for sending emails
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-API-Key
// @description				Enter the API key
// @contact.name				API Support
// @contact.email				email@mynute.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:4002
// @BasePath					/
func main() {
	email.NewServer().Run("listen")
}
