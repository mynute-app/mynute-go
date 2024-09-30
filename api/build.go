package api

import (
	"agenda-kaki-go/api/routes"
	"agenda-kaki-go/api/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	postgres := &services.Postgres{DB: DB}
	routes.Company(postgres, App)
	routes.CompanyType(postgres, App)
}
