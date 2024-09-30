package config

import (
	"agenda-kaki-go/core/config/api/routes"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func BuildRouter(DB *gorm.DB, App *fiber.App) {
	postgres := &services.Postgres{DB: DB}
	routes.Company(postgres, App)
	routes.CompanyType(postgres, App)
}
