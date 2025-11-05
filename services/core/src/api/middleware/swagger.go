package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func SwaggerAuth() fiber.Handler {
	swagger_password := os.Getenv("SWAGGER_PASSWORD")
	swagger_user := os.Getenv("SWAGGER_USER")
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			swagger_user: swagger_password,
		},
		Realm: "Restricted",
	})
}

