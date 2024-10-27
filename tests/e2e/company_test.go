package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
	"fmt"
	"testing"
)

var _ e2e.IEntity = (*Company)(nil)

type Company struct {
	*e2e.BaseE2EActions
	companyType *CompanyType
}

func (c *Company) GenerateTesters(n int) {
	for i := 0; i < n; i++ {
		c.GenerateTester(
			"company",
			"company",
			map[string]interface{}{
				"name":   lib.GenerateRandomName("Company"),
				"tax_id": fmt.Sprintf("%v", lib.GenerateRandomIntOfExactly(14)),
				"company_types": []map[string]interface{}{
					{"id": c.companyType.Testers[i].EntityID, "name": c.companyType.Testers[i].PostBody["name"]},
				},
			},
			map[string]interface{}{
				"name": lib.GenerateRandomName("Company"),
			},
		)
	}
}

func (c *Company) Make(n int) {
	c.CreateDependencies(n)
	c.GenerateTesters(n)
}

func (c *Company) CreateDependencies(n int) {
	companyType := &CompanyType{}
	companyType.SetTest(c.T)
	companyType.Make(n)
	companyType.CreateAllTesters(201)
	c.companyType = companyType
}

func (c *Company) ClearDependencies() {
	c.companyType.ForceDeleteAllTesters(204)
}

func TestCompanyFlow(t *testing.T) {
	company := &Company{}
	company.SetTest(t)
	company.Make(10)
	company.RunAll()
	company.ClearDependencies()
}
