package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"
	"errors"
	"fmt"
	"strconv"
)

type Company struct {
	DB *services.Postgres
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
		idStr := strconv.Itoa(int(companyType.ID))
		if err := c.DB.GetOneBy("id", idStr, models.CompanyType{}, nil); err != nil {
			errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
			return errors.New(errStr)
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
			if err := c.DB.GetOneBy("id", idStr, models.CompanyType{}, nil); err != nil {
				errStr := fmt.Sprintf("company type with ID %s does not exist", idStr)
				return errors.New(errStr)
			}
		}
	}
	return nil
}