package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"log"

	"github.com/gofiber/fiber/v2"
)

type HolidaysMiddlewareActions struct {
	Gorm *handler.Gorm
}

func Holidays(Gorm *handler.Gorm) *Registry {
	holidays := &HolidaysMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()
	registry.RegisterAction(namespace.HolidaysKey.Name, "POST", holidays.Create)

	return registry
}

func (hm *HolidaysMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
	Holidays, err := lib.GetFromCtx[*model.Holidays](c, namespace.GeneralKey.Model)
	if err != nil {
		log.Println(err)
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(Holidays.Name, "holidays"); err != nil {
		return 400, err
	}

	return 0, nil

}
