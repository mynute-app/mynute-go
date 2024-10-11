package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	Gorm *handlers.Gorm
}

// Middleware for Create operation
func (ctm *CompanyType) Create(c fiber.Ctx) (int, error) {
	keys := namespace.GeneralKey
	// Retrieve companyType from c.Locals
	companyType, err := lib.GetFromCtx[*models.CompanyType](c, keys.Model)
	if err != nil {
		return 500, err
	}

	// Perform validation
	if !lib.ValidateName(companyType.Name) {
		return 400, errors.New("companyType.Name must be at least 3 characters long")
	}

	// Proceed to the next middleware or handler
	return 0, nil
}

// Middleware for Update operation
func (ctm *CompanyType) Update(c fiber.Ctx) (int, error) {
	keys := namespace.GeneralKey
	// Retrieve changes from c.Locals
	changes, err := lib.GetFromCtx[map[string]interface{}](c, keys.Changes)
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
		if !lib.ValidateName(name) {
			return 400, errors.New("companyType.Name must be at least 3 characters long")
		}
		// Check if the name already exists
		var companyType models.CompanyType
		if err := ctm.Gorm.GetOneBy("name", name, &companyType, nil); err == nil {
			return 400, errors.New("companyType.Name already exists")
		}
	}

	// Proceed to the next middleware or handler
	return 0, nil
}

// Middleware for Delete operation
func (ctm *CompanyType) DeleteOneById(c fiber.Ctx) (int, error) {
	companyTypeId := c.Params("id")

	// Check if the company type is associated with any companies
	var companies []models.Company
	if err := ctm.Gorm.DB.
		Model(&companies).
		Joins("JOIN company_company_types ON companies.id = company_company_types.company_id").
		Where("company_company_types.company_type_id = ?", companyTypeId).
		Find(&companies).Error; err != nil {
		return 500, err
	}
	if len(companies) > 0 {
		return 400, errors.New("companyType is associated with companies")
	}

	// // Proceed to the next middleware or handler
	return 0, nil
}

// Middleware for Delete operation
func (ctm *CompanyType) ForceDeleteOneById(c fiber.Ctx) (int, error) {
	companyTypeId := c.Params("id")
	// Check if the company type is associated with any companies
	var companies []models.Company
	if err := ctm.Gorm.DB.
		Joins("JOIN company_company_types ON companies.id = company_company_types.company_id").
		Where("company_company_types.company_type_id = ?", companyTypeId).
		Find(&companies).Error; err != nil {
		return 500, err
	}
	if len(companies) > 0 {
		return 400, errors.New("companyType is associated with companies")
	}

	// // Proceed to the next middleware or handler
	return 0, nil
}
