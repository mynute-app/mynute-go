package routes

import (
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	postgres := &services.Postgres{DB: DB}
	Company(postgres, App)
	CompanyType(postgres, App)
}
