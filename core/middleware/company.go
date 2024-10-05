package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"strconv"
)

type Company struct {
	Gorm *handlers.Gorm
}

func (c *Company) Create(company models.Company) error {
	if !lib.ValidateName(company.Name) {
		return errors.New("company name must be at least 3 characters long")
	} else if !lib.ValidateTaxID(company.TaxID) {
		return errors.New("invalid tax ID")
	} else if len(company.CompanyTypes) == 0 {
		return errors.New("company type is required")
	}
	for _, companyType := range company.CompanyTypes {
		if companyType.ID == 0 {
			return errors.New("company type ID is missing or invalid")
		}
		var model models.CompanyType
		idStr := strconv.FormatUint(uint64(companyType.ID), 10)
		if err := c.Gorm.GetOneBy("id", idStr, &model, nil); err != nil {
			errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
			return errors.New(errStr)
		}
		if model.Name != companyType.Name {
			return errors.New("company type name passed does not match the ID provided")
		}
	}
	return nil
}

func (c *Company) Update(changes map[string]interface{}) error {
	if changes["name"] != nil && !lib.ValidateName(changes["name"].(string)) {
		return errors.New("company name must be at least 3 characters long")
	} else if changes["tax_id"] != nil && !lib.ValidateTaxID(changes["tax_id"].(string)) {
		return errors.New("invalid tax ID")
	} else if changes["company_types"] != nil {
		companyTypes := changes["company_types"].([]models.CompanyType)
		for _, companyType := range companyTypes {
			idStr := strconv.Itoa(int(companyType.ID))
			if err := c.Gorm.GetOneBy("id", idStr, models.CompanyType{}, nil); err != nil {
				errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
				return errors.New(errStr)
			}
		}
	}
	return nil
}