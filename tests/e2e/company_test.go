package e2e_test

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/tests/handlers"
	"agenda-kaki-go/tests/lib"
	"fmt"
	"testing"
)

type Company struct {
	T *testing.T
}

type companyActios struct {
	*Company
	Testers []handlers.Tester
}

func (c *Company) Init() *companyActios {
	return &companyActios{Company: c}
}

func (c *companyActios) LoadTester(n int) *companyActios {
	for i := 0; i < n; i++ {
		company := handlers.Tester{
			Entity:      "company",
			BaseURL:     namespace.GeneralKey.BaseURL,
			RelatedPath: "company",
			PostBody: map[string]interface{}{
				"name":   lib.GenerateRandomName("Company"),
				"tax_id": fmt.Sprintf("%v", lib.GenerateRandomIntOfExactly(14)),
			},
			PatchBody: map[string]interface{}{"name": lib.GenerateRandomName("Company")},
		}
		c.Testers = append(c.Testers, company)
	}
	return c
}

// func TestCompanyFlow(t *testing.T) {
// 	t.Run("CreateCompanyType", companyType.ExpectStatus(201).POST)
// 	company.PostBody["company_types"] = []map[string]interface{}{
// 		{"id": companyType.EntityID, "name": companyType.PostBody["name"]},
// 	}
// 	t.Run("CreateCompany", company.ExpectStatus(201).POST)
// 	t.Run("UpdateCompany", company.ExpectStatus(200).PATCH)
// 	t.Run("DeleteCompanyTypeWithError", companyType.ExpectStatus(400).DELETE)
// 	t.Run("DeleteCompany", company.ExpectStatus(204).DELETE)
// 	t.Run("DeleteCompanyType", companyType.ExpectStatus(204).DELETE)
// 	t.Run("ForceDeleteCompanyType", companyType.ExpectStatus(204).ForceDELETE)
// 	t.Run("ForceDeleteCompany", company.ExpectStatus(204).ForceDELETE)
// }
