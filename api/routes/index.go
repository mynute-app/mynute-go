package routes

import (
	"agenda-kaki-company-go/api/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Init(DB *gorm.DB, App *fiber.App) {
	postgres := &services.Postgres{DB: DB}
	Company(postgres, App)
}