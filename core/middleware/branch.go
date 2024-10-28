package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"github.com/gofiber/fiber/v3"
)

func Branch(Gorm *handlers.Gorm) *Registry {
	branch := &branchMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()

	registry.RegisterAction(namespace.BranchKey.Name, "POST", branch.Create)

	return registry
}

type branchMiddlewareActions struct {
	Gorm *handlers.Gorm
}

func (ba *branchMiddlewareActions) CheckCompany(c fiber.Ctx) (int, error) {
	var company models.Company
	if s, err := GetCompany(ba.Gorm, c, &company); err != nil {
		return s, err
	}

	if c.Method() == "GET" {
		return 0, nil
	}

	branch, err := lib.GetFromCtx[*models.Branch](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}

	branch.CompanyID = company.ID

	return 0, nil
}

func (ba *branchMiddlewareActions) Create(c fiber.Ctx) (int, error) {
	service, err := lib.GetFromCtx[*models.Branch](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(service.Name, "service"); err != nil {
		return 400, err
	}
	// Proceed to the next middleware or handler
	return 0, nil
}
