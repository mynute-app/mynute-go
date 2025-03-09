package main

import (
	"agenda-kaki-go/core"
	_ "agenda-kaki-go/docs"
)

//	@title						Fiber Example API
//	@version					1.0
//	@description				Swagger API for testing and debugging
//	@termsOfService				http://swagger.io/terms/
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Enter the token in the format: <token>
//	@contact.name				API Support
//	@contact.email				fiber@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:4000
//	@BasePath					/
func main() {
	core.NewServer().Run("listen")
}
