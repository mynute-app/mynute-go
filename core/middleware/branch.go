package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type Branch struct {
	Gorm *handlers.Gorm
}

func (cb *Branch) Create(c fiber.Ctx) (int, error) {
	service, err := lib.GetFromCtx[*models.Branch](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(service.Name, "service"); err != nil {
		return 400, err
	}
	// Proceed to the next middleware or handler
	return 0, nil
}

