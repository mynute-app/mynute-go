package routes

import (
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	gormHandler := &handlers.Gorm{DB: DB}
	Company(gormHandler, App)
	CompanyType(gormHandler, App)
	Branch(gormHandler, App)
	Service(gormHandler, App)
	Employee(gormHandler, App)
	Holidays(gormHandler, App)
}
