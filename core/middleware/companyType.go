package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"
	"errors"
)

type CompanyType struct {
	DB *services.Postgres
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