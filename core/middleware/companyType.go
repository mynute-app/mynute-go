package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	Gorm *handlers.Gorm
}

// Define context keys using a non-exported type to avoid collisions
type contextKey string

const (
	companyTypeKey contextKey = "companyType"
	changesKey     contextKey = "changes"
)

// Middleware for Create operation
func (ctm *CompanyType) Create(c fiber.Ctx) (int, error) {
	// Retrieve companyType from c.Locals
	companyTypeInterface := c.Locals(companyTypeKey)
	if companyTypeInterface == nil {
		return 500, interfaceDataNotFound("companyType")
	}
	companyType, ok := companyTypeInterface.(models.CompanyType)
	if !ok {
		return 500, invalidDataType("companyType")
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
	// Retrieve changes from c.Locals
	changesInterface := c.Locals(changesKey)
	if changesInterface == nil {
		return 500, interfaceDataNotFound("'changes'")
	}
	changes, ok := changesInterface.(map[string]interface{})
	if !ok {
		return 500, invalidDataType("'changes'")
	}

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
func (ctm *CompanyType) Delete(c fiber.Ctx) (int, error) {
	// Retrieve companyType from c.Locals
	companyTypeInterface := c.Locals(companyTypeKey)
	if companyTypeInterface == nil {
		return 500, interfaceDataNotFound("companyType")
	}
	companyType, ok := companyTypeInterface.(models.CompanyType)
	if !ok {
		return 500, invalidDataType("companyType")
	}

	// Check if the company type is associated with any companies
	var companies []models.Company
	if err := ctm.Gorm.DB.
		Model(&companies).
		Joins("JOIN company_company_types ON companies.id = company_company_types.company_id").
		Where("company_company_types.company_type_id = ?", companyType.ID).
		Find(&companies).Error; err != nil {
		return 500, err
	}
	if len(companies) > 0 {
		return 400, errors.New("companyType is associated with companies")
	}

	// Proceed to the next middleware or handler
	return 0, nil
}

func interfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in context", interfaceName)
	return errors.New(errStr)
}

func invalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}
