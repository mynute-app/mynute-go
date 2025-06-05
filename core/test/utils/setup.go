package utilsT

import (
	"agenda-kaki-go/core/lib"
	modelT "agenda-kaki-go/core/test/models"
	"os"
)

func CreateCompaniesRandomly(CompaniesToCreate int) ([]*modelT.Company, error) {
	if os.Getenv("APP_ENV") == "prod" { // Make sure it never runs in production
		panic("Cannot run tests in production environment. Set APP_ENV to 'test' or 'dev'.")
	}
	var companies []*modelT.Company
	for range CompaniesToCreate {
		var company modelT.Company
		employees := lib.GenerateRandomIntFromRange(6, 12) // Random From 6 to 12 per company
		branches := lib.GenerateRandomIntFromRange(1, 4)   // Random From 1 to 4 per company
		services := lib.GenerateRandomIntFromRange(2, 32)  // Random From 2 to 32 per company
		if err := company.CreateCompanyRandomly(employees, branches, services); err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}
	return companies, nil
}
