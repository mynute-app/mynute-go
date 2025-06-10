package utilsT

import (
	"agenda-kaki-go/core/lib"
	modelT "agenda-kaki-go/core/test/models"
	"os"
)

func CreateCompaniesRandomly(CompaniesToCreate int) ([]*modelT.Company, error) {
	app_env := os.Getenv("APP_ENV")
	if app_env == "prod" { // Make sure it never runs in production
		panic("Cannot run tests in production environment. Set APP_ENV to 'test' or 'dev'.")
	}
	var companies []*modelT.Company
	for range CompaniesToCreate {
		var company modelT.Company
		var employees int
		var branches int
		var services int
		if app_env == "test" {
			employees = 2 // Fixed number of employees for test environment
			branches = 2  // Fixed number of branches for test environment
			services = 2  // Fixed number of services for test environment
		} else {
			employees = lib.GenerateRandomIntFromRange(1, 48) // Random From 1 to 48 per company
			branches = lib.GenerateRandomIntFromRange(1, 24)   // Random From 1 to 24 per company
			services = lib.GenerateRandomIntFromRange(1, 72)   // Random From 1 to 72 per company
		}
		if err := company.CreateCompanyRandomly(employees, branches, services); err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}
	return companies, nil
}
