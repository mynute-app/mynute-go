package main

import (
	"mynute-go/core"
	_ "mynute-go/docs"
)

// @title						GO-Fiber API
// @version					1.0
// @description				Swagger API for testing and debugging
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-Auth-Token
// @description				Enter the token in the format: <token>
// @in							header
// @name						X-Company-ID
// @description				Enter the company ID in the format: <company_id>
// @contact.name				API Support
// @contact.email				fiber@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:4000
// @BasePath					/
func main() {
	core.NewServer().Run("listen")
}

