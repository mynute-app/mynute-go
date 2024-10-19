package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"errors"

	"github.com/gofiber/fiber/v3"
)

type Branch struct {
	Gorm *handlers.Gorm
}

func (cb *Branch) CheckCompany(c fiber.Ctx) (int, error) {
	companyId := c.Params("companyId")
	if companyId == "" {
		return 400, errors.New("missing companyId")
	}

	// Check if the company with the following id exists on the database.
	var company models.Company
	if err := cb.Gorm.GetOneBy("id", companyId, &company, []string{}); err != nil {
		return 400, err
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

func (cb *Branch) Create(c fiber.Ctx) (int, error) {
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

