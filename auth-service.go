//go:build ignore
// +build ignore

package main

import (
	_ "mynute-go/services/auth/docs"
	"mynute-go/services/auth"
)

// @title						Auth Service API
// @version					1.0
// @description				Authentication and Authorization Service
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-Auth-Token
// @description				Enter the token in the format: <token>
// @contact.name				API Support
// @contact.email				auth@mynute.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:4001
// @BasePath					/
func main() {
	auth.NewServer().Run("listen")
}
