package routes

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handlers.Gorm{DB: DB}

	// independent routes
	rootRouter := App.Group("/", middleware.WhoAreYou)
	Company(Gorm, rootRouter)
	CompanyType(Gorm, rootRouter)

	// CompanyId dependent routes
	companyPrefix := fmt.Sprintf("/company/:%s", namespace.QueryKey.CompanyId)
	companyCheck := middleware.GetCompany(Gorm)
	companyRouter := rootRouter.Group(companyPrefix, companyCheck)
	Branch(Gorm, companyRouter)
	Service(Gorm, companyRouter)
	Employee(Gorm, companyRouter)
}
