package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

var _ IMiddleware = (*Company)(nil)

type Company struct {
	Gorm *handlers.Gorm
}

type CompanyMiddlewareActions struct {}

func GetCompany(Gorm *handlers.Gorm, c fiber.Ctx, company *models.Company) (int, error) {
	companyID := c.Params(string(namespace.GeneralKey.CompanyId))
	if companyID == "" {
		return 400, errors.New("missing companyId")
	}

	if err := Gorm.GetOneBy("id", companyID, &company, nil); err != nil {
		return 400, err
	}

	return 0, nil
}

func (ca *CompanyMiddlewareActions) Create(Gorm *handlers.Gorm) func(c fiber.Ctx) (int, error) {
	validateCreation := func(c fiber.Ctx) (int, error) {
		company, err := lib.GetFromCtx[*models.Company](c, namespace.GeneralKey.Model)
		if err != nil {
			return 500, err
		}
		err = lib.ValidateName(company.Name, "company")
		if err != nil {
			return 400, err
		} else if !lib.ValidateTaxID(company.TaxID) {
			return 400, errors.New("invalid tax ID")
		} else if len(company.CompanyTypes) == 0 {
			return 400, errors.New("company type is required")
		}

		for _, companyType := range company.CompanyTypes {
			if companyType.ID == 0 {
				return 400, errors.New("company type ID is missing or invalid")
			}
			var model models.CompanyType
			idStr := strconv.FormatUint(uint64(companyType.ID), 10)
			if err := Gorm.GetOneBy("id", idStr, &model, nil); err != nil {
				errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
				return 400, errors.New(errStr)
			}
			if model.Name != companyType.Name {
				return 400, errors.New("company type name passed does not match the ID provided")
			}
		}
		return 0, nil
	}
	return validateCreation
}

func (ca *CompanyMiddlewareActions) Update(Gorm *handlers.Gorm) func(c fiber.Ctx) (int, error) {
	validateUpdates := func(c fiber.Ctx) (int, error) {
		changes, err := lib.GetFromCtx[map[string]interface{}](c, namespace.GeneralKey.Changes)
		if err != nil {
			return 500, err
		}
		if changes["name"] != nil {
			if err = lib.ValidateName(changes["name"].(string), "company"); err != nil {
				return 400, err
			}
		} else if changes["tax_id"] != nil && !lib.ValidateTaxID(changes["tax_id"].(string)) {
			return 400, errors.New("invalid tax ID")
		} else if changes["company_types"] != nil {
			companyTypes := changes["company_types"].([]models.CompanyType)
			for _, companyType := range companyTypes {
				idStr := strconv.Itoa(int(companyType.ID))
				if err := Gorm.GetOneBy("id", idStr, models.CompanyType{}, nil); err != nil {
					errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
					return 400, errors.New(errStr)
				}
			}
		}
		return 0, nil
	}
	return validateUpdates
}

var companyActs = CompanyMiddlewareActions{}

func (cc *Company) POST() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){companyActs.Create(cc.Gorm)}
}

func (cc *Company) PATCH() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){companyActs.Update(cc.Gorm)}
}

func (cc *Company) GET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (cc *Company) DELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (cc *Company) ForceGET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (cc *Company) ForceDELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}
