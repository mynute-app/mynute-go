package routes

import (
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handler.Gorm{DB: DB}
	Auth(Gorm, App)
	Holidays(Gorm, App)
	Sector(Gorm, App)
	Company(Gorm, App)
	User(Gorm, App)
	Branch(Gorm, App)
	Service(Gorm, App)
	Swagger(App)
}
