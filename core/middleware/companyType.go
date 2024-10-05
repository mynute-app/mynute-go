package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"
	"errors"
)

type CompanyType struct {
	Gorm *services.Gorm
}

func (c *CompanyType) Create(companyType models.CompanyType) error {
	if !lib.ValidateName(companyType.Name) {
		return errors.New("company type name must be at least 3 characters long")
	}
	return nil
}

func (c *CompanyType) Update(changes map[string]interface{}) error {
	if changes["name"] != nil && !lib.ValidateName(changes["name"].(string)) {
		return errors.New("company type name must be at least 3 characters long")
	}
	return nil
}

func (c *CompanyType) Delete(companyType models.CompanyType) error {
	// Check if the company type is associated with any companies
	var companies []models.Company
	if err := c.Gorm.DB.
		Model(&companies).
		Joins("JOIN company_company_types ON companies.id = company_company_types.company_id").
		Where("company_company_types.company_type_id = ?", companyType.ID).
		Find(&companies).Error; err != nil {
		return err
	}
	if len(companies) > 0 {
		return errors.New("company type is associated with companies. Please remove the association before deleting")
	}
	return nil
}
