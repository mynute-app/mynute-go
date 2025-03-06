package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

func Branch(Gorm *handler.Gorm) *Registry {
	branch := &branchMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()
	var BranchMiddleActions = []MiddlewareActions{
		{
			methods: []string{"POST", "PUT", "DELETE"},
			action:  WhoAreYou,
		},
		{
			methods: "POST",
			action:  branch.Create,
		},
	}
	registry.RegisterActions(namespace.BranchKey.Name, BranchMiddleActions)
	return registry
}

type branchMiddlewareActions struct {
	Gorm *handler.Gorm
}

// Check if the company exists and attach the company ID to the branch.
func (ba *branchMiddlewareActions) CheckCompany(c *fiber.Ctx) (int, error) {
	company, err := lib.GetFromCtx[*model.Company](c, namespace.CompanyKey.Model)
	if err != nil {
		return 500, err
	}

	if c.Method() == "GET" {
		return 0, nil
	}

	branch, err := lib.GetFromCtx[*model.Branch](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}

	branch.CompanyID = company.ID

	return 0, nil
}

func (ba *branchMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
	branch, err := lib.GetFromCtx[*model.Branch](c, namespace.GeneralKey.Model)
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
