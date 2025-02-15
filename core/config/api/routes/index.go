package routes

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handlers.Gorm{DB: DB}
	//will pass in the middleware and if not authenticated will return 401
	Auth(Gorm, App)
	Holidays(Gorm, App)
	CompanyType(Gorm, App)
	Company(Gorm, App)

	companyPrefix := fmt.Sprintf("/company/:%s", namespace.QueryKey.CompanyId)
	companyCheck := middleware.GetCompany(Gorm)
	authRouter := App.Group("/", middleware.WhoAreYou)
	companyRouter := authRouter.Group(companyPrefix, companyCheck)
	Branch(Gorm, companyRouter)
	Service(Gorm, companyRouter)
	User(Gorm, authRouter)
}
