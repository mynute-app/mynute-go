package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"github.com/gofiber/fiber/v3"
)

var _ IMiddleware = (*Branch)(nil)

type Branch struct {
	Gorm *handlers.Gorm
}

type BranchMiddlewareActions struct {}

func (ba *BranchMiddlewareActions) CheckCompany(Gorm *handlers.Gorm) func(c fiber.Ctx) (int, error) {
	checkCompany := func(c fiber.Ctx) (int, error) {
		var company models.Company
		if s, err := GetCompany(Gorm, c, &company); err != nil {
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
	return checkCompany
}

func (ba *BranchMiddlewareActions) Create(c fiber.Ctx) (int, error) {
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

var branchActs = BranchMiddlewareActions{}

func (cb *Branch) POST() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){branchActs.CheckCompany(cb.Gorm), branchActs.Create}
}

func (cb *Branch) GET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){branchActs.CheckCompany(cb.Gorm)}
}

func (cb *Branch) DELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){branchActs.CheckCompany(cb.Gorm)}
}

func (cb *Branch) PATCH() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){branchActs.CheckCompany(cb.Gorm)}
}

func (cb *Branch) ForceGET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){branchActs.CheckCompany(cb.Gorm)}
}

func (cb *Branch) ForceDELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){branchActs.CheckCompany(cb.Gorm)}
}
