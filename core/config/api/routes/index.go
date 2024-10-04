package routes

import (
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	gorm := &services.Gorm{DB: DB}
	Company(gorm, App)
	CompanyType(gorm, App)
}
