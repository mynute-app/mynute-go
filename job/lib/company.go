package lib

import (
	"mynute-go/core/src/lib"
	"mynute-go/test/src/model"
	"os"
)

func CreateCompaniesRandomly(CompaniesToCreate int) ([]*model.Company, error) {
	app_env := os.Getenv("APP_ENV")
	if app_env == "prod" { // Make sure it never runs in production
		panic("Cannot run tests in production environment. Set APP_ENV to 'test' or 'dev'.")
	} else if app_env != "test" && app_env != "dev" {
		panic("APP_ENV must be set to 'test' or 'dev'. Current value: " + app_env)
	}
	var companies []*model.Company
	for range CompaniesToCreate {
		var company model.Company
		var employees int
		var branches int
		var services int
		if app_env == "test" {
			employees = 2 // Fixed number of employees for test environment
			branches = 2  // Fixed number of branches for test environment
			services = 2  // Fixed number of services for test environment
		} else if app_env == "dev" {
			employees = lib.GenerateRandomIntFromRange(1, 8) // Random From 1 to 8 per company
			branches = lib.GenerateRandomIntFromRange(1, 4)  // Random From 1 to 4 per company
			services = lib.GenerateRandomIntFromRange(1, 12) // Random From 1 to 12 per company
		}
		if err := company.CreateCompanyRandomly(employees, branches, services); err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}
	return companies, nil
}
