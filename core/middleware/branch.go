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

// Check if the company exists and attach the company ID to the branch.
func (ba *branchMiddlewareActions) CheckCompany(c fiber.Ctx) (int, error) {
	company, err := lib.GetFromCtx[*models.Company](c, namespace.CompanyKey.Model)
	if err != nil {
		return 500, err
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
	branch, err := lib.GetFromCtx[*models.Branch](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(branch.Name, "branch"); err != nil {
		return 400, err
	}
	// Check if the company exists
	if s, err := ba.CheckCompany(c); err != nil {
		return s, err
	}
	// Proceed to the next middleware or handler
	return 0, nil
}
