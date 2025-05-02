package middleware

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

func GetTenant(c *fiber.Ctx) error {
	companyID := c.Get(namespace.HeadersKey.Company)
	if companyID == "" {
		return lib.Error.Auth.CompanyHeaderMissing
	}
	c.Locals(namespace.HeadersKey.Company, companyID)
	return c.Next()
}