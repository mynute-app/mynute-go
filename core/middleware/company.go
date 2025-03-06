package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type company_middleware struct {
	Gorm *handler.Gorm
}

func Company(Gorm *handler.Gorm) *company_middleware {
	return &company_middleware{Gorm: Gorm}
}

func (cm *company_middleware) GetCompany(c *fiber.Ctx) error {
	var company model.Company
	res := &lib.SendResponse{Ctx: c}
	companyID := c.Params(namespace.QueryKey.CompanyId)

	if companyID == "" {
		return errors.New("missing companyId")
	}

	if err := cm.Gorm.GetOneBy("id", companyID, &company, nil); err != nil {
		return res.Http400(err).Next()
	}

	c.Locals(namespace.CompanyKey.Model, &company)
	return c.Next()
}

func (cm *company_middleware) ValidateProps(c *fiber.Ctx) error {
	company, err := lib.GetFromCtx[*model.Company](c, namespace.GeneralKey.Model)
	if err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	err = lib.ValidateName(company.Name, "company")
	if err != nil {
		return res.Http400(err).Next()
	} else if !lib.ValidateTaxID(company.TaxID) {
		return res.Http400(errors.New("invalid tax ID")).Next()
	} else if len(company.Sectors) == 0 {
		return res.Http400(errors.New("company type is required")).Next()
	}

	for _, companyType := range company.Sectors {
		if companyType.ID == 0 {
			return res.Http400(errors.New("company type ID is missing or invalid")).Next()
		}
		var model model.Sector
		idStr := strconv.FormatUint(uint64(companyType.ID), 10)
		if err := cm.Gorm.GetOneBy("id", idStr, &model, nil); err != nil {
			errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
			return res.Http400(errors.New(errStr)).Next()
		}
		if model.Name != companyType.Name {
			return res.Http400(errors.New("company type name passed does not match the ID provided")).Next()
		}
	}
	return c.Next()
}

// type companyMiddlewareActions struct {
// 	Gorm *handler.Gorm
// }

// func Company(Gorm *handler.Gorm) *Registry {
// 	registry := NewRegistry()
// 	company := &companyMiddlewareActions{Gorm: Gorm}
// 	var CompanyMiddleActions = []MiddlewareActions{
// 		{
// 			methods: []string{"POST", "PATCH", "DELETE"},
// 			action:  WhoAreYou,
// 		},
// 		{
// 			methods: "POST",
// 			action:  company.Create,
// 		},
// 		{
// 			methods: "PATCH",
// 			action:  company.Update,
// 		},
// 	}
// 	registry.RegisterActions(namespace.CompanyKey.Name, CompanyMiddleActions)
// 	return registry
// }

// func GetCompany(Gorm *handler.Gorm) fiber.Handler {
// 	getCompany := func(c *fiber.Ctx) error {
// 		var company model.Company
// 		res := &lib.SendResponse{Ctx: c}
// 		companyID := c.Params(namespace.QueryKey.CompanyId)

// 		if companyID == "" {
// 			return res.Http400(errors.New("missing companyId")).Next()
// 		}

// 		if err := Gorm.GetOneBy("id", companyID, &company, nil); err != nil {
// 			return res.Http400(err).Next()
// 		}

// 		c.Locals(namespace.CompanyKey.Model, &company)
// 		return c.Next()
// 	}
// 	return getCompany
// }

// func (ca *companyMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
// 	company, err := lib.GetFromCtx[*model.Company](c, namespace.GeneralKey.Model)
// 	if err != nil {
// 		return 500, err
// 	}
// 	err = lib.ValidateName(company.Name, "company")
// 	if err != nil {
// 		return 400, err
// 	} else if !lib.ValidateTaxID(company.TaxID) {
// 		return 400, errors.New("invalid tax ID")
// 	} else if len(company.Sectors) == 0 {
// 		return 400, errors.New("company type is required")
// 	}

// 	for _, companyType := range company.Sectors {
// 		if companyType.ID == 0 {
// 			return 400, errors.New("company type ID is missing or invalid")
// 		}
// 		var model model.Sector
// 		idStr := strconv.FormatUint(uint64(companyType.ID), 10)
// 		if err := ca.Gorm.GetOneBy("id", idStr, &model, nil); err != nil {
// 			errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
// 			return 400, errors.New(errStr)
// 		}
// 		if model.Name != companyType.Name {
// 			return 400, errors.New("company type name passed does not match the ID provided")
// 		}
// 	}
// 	return 0, nil
// }

// func (ca *companyMiddlewareActions) Update(c *fiber.Ctx) (int, error) {
// 	changes, err := lib.GetFromCtx[map[string]any](c, namespace.GeneralKey.Changes)
// 	if err != nil {
// 		return 500, err
// 	}
// 	if changes["name"] != nil {
// 		if err = lib.ValidateName(changes["name"].(string), "company"); err != nil {
// 			return 400, err
// 		}
// 	} else if changes["tax_id"] != nil && !lib.ValidateTaxID(changes["tax_id"].(string)) {
// 		return 400, errors.New("invalid tax ID")
// 	} else if changes["sectors"] != nil {
// 		companyTypes := changes["sectors"].([]model.Sector)
// 		for _, companyType := range companyTypes {
// 			idStr := strconv.Itoa(int(companyType.ID))
// 			if err := ca.Gorm.GetOneBy("id", idStr, model.Sector{}, nil); err != nil {
// 				errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
// 				return 400, errors.New(errStr)
// 			}
// 		}
// 	}
// 	return 0, nil
// }
