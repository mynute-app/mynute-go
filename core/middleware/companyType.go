package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// var _ IMiddleware = (*CompanyType)(nil)

// type CompanyType struct {
// 	Gorm *handlers.Gorm
// }

func CompanyType(Gorm *handlers.Gorm) *Registry {
	companyType := &companyTypeMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()

	var CompanyTypeMiddleActions = []MiddlewareActions{
		{
			methods: []string{"POST", "PATCH", "DELETE"},
			action:  WhoAreYou,
		},
		{
			methods: "POST",
			action:  companyType.Create,
		},
		{
			methods: "PATCH",
			action:  companyType.Update,
		},
		{
			methods: "DELETE",
			action:  companyType.DeleteOneById,
		},
	}
	registry.RegisterActions(namespace.CompanyTypeKey.Name, CompanyTypeMiddleActions)
	return registry
}

type companyTypeMiddlewareActions struct {
	Gorm *handlers.Gorm
}

// Middleware for Create operation
func (cta *companyTypeMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
	keys := namespace.GeneralKey
	// Retrieve companyType from c.Locals
	companyType, err := lib.GetFromCtx[*models.CompanyType](c, keys.Model)
	if err != nil {
		return 500, err
	}

	err = lib.ValidateName(companyType.Name, "companyType")
	if err != nil {
		return 400, err
	}

	// Proceed to the next middleware or handler
	return 0, nil
}

// Middleware for Update operation
func (cta *companyTypeMiddlewareActions) Update(c *fiber.Ctx) (int, error) {
	keys := namespace.GeneralKey
	// Retrieve changes from c.Locals
	changes, err := lib.GetFromCtx[map[string]any](c, keys.Changes)
	if err != nil {
		return 500, err
	}

	log.Printf("Changes: %v", changes)

	// Perform validation
	if nameValue, exists := changes["name"]; exists {
		name, ok := nameValue.(string)
		if !ok {
			return 500, errors.New("invalid 'name' on 'changes' data type")
		}
		if err := lib.ValidateName(name, "companyType"); err != nil {
			return 400, err
		}
		// Check if the name already exists
		var companyType models.CompanyType
		if err := cta.Gorm.GetOneBy("name", name, &companyType, nil); err == nil {
			return 400, errors.New("companyType.Name already exists")
		}
	}
	// Proceed to the next middleware or handler
	return 0, nil
}

// Unified Middleware for Delete Validation
func (cta *companyTypeMiddlewareActions) DeleteOneById(c *fiber.Ctx) (int, error) {
	companyTypeId := c.Params("id")

	// Check if the company type is associated with any companies
	var companies []models.Company
	if err := cta.Gorm.DB.
		Model(&companies).
		Joins("JOIN company_company_types ON companies.id = company_company_types.company_id").
		Where("company_company_types.company_type_id = ?", companyTypeId).
		Find(&companies).Error; err != nil {
		return 500, err
	}

	// Return error if there are associated companies
	if len(companies) > 0 {
		return 400, errors.New("companyType is associated with companies")
	}

	// Pass the validation
	return 0, nil
}
