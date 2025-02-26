package routes

import (
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handlers.Gorm{DB: DB}
	Auth(Gorm, App)
	Holidays(Gorm, App)
	CompanyType(Gorm, App)
	Company(Gorm, App)
	User(Gorm, App)
	Branch(Gorm, App)
	Service(Gorm, App)
	Swagger(App)
}
