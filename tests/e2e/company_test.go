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

type companyActions struct {
	*Company
	Testers []handlers.Tester
}

func (c *Company) Init() *companyActions {
	return &companyActions{Company: c}
}

func (c *companyActions) GenerateTesters(n int) *companyActions {
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

func (c *companyActions) CreateAllTesters(s int) *companyActions {
	for _, tester := range c.Testers {
		c.T.Run("CreateCompany", tester.ExpectStatus(201).POST)
	}
	return c
}

func (c *companyActions) GetAllTesters(s int) *companyActions {
	for _, tester := range c.Testers {
		c.T.Run("GetCompany", tester.ExpectStatus(200).GET)
	}
	return c
}

func (c *companyActions) ForceGetAllTesters(s int) *companyActions {
	for _, tester := range c.Testers {
		c.T.Run("ForceGetCompany", tester.ExpectStatus(200).ForceGET)
	}
	return c
}

func (c *companyActions) UpdateAllTesters(s int) *companyActions {
	for _, tester := range c.Testers {
		c.T.Run("UpdateCompany", tester.ExpectStatus(200).PATCH)
	}
	return c
}

func (c *companyActions) DeleteAllTesters(s int) *companyActions {
	for _, tester := range c.Testers {
		c.T.Run("DeleteCompany", tester.ExpectStatus(204).DELETE)
	}
	return c
}

func (c *companyActions) ForceDeleteAllTesters(s int) *companyActions {
	for _, tester := range c.Testers {
		c.T.Run("ForceDeleteCompany", tester.ExpectStatus(204).ForceDELETE)
	}
	return c
}

func (c *companyActions) RunAll() {
	c.
		GenerateTesters(10).
		CreateAllTesters(201).
		UpdateAllTesters(200).
		GetAllTesters(200).
		ForceGetAllTesters(200).
		DeleteAllTesters(204).
		GetAllTesters(404).
		ForceGetAllTesters(200).
		ForceDeleteAllTesters(204).
		ForceGetAllTesters(404)
}

func (c *companyActions) GetTester(id int, s int) *companyActions {
	c.T.Run("GetCompany", c.Testers[id].ExpectStatus(s).GET)
	return c
}

func (c *companyActions) ForceGetTester(id int, s int) *companyActions {
	c.T.Run("ForceGetCompany", c.Testers[id].ExpectStatus(s).ForceGET)
	return c
}

func (c *companyActions) UpdateTester(id int, s int) *companyActions {
	c.T.Run("UpdateCompany", c.Testers[id].ExpectStatus(s).PATCH)
	return c
}

func (c *companyActions) DeleteTester(id int, s int) *companyActions {
	c.T.Run("DeleteCompany", c.Testers[id].ExpectStatus(s).DELETE)
	return c
}

func (c *companyActions) ForceDeleteTester(id int, s int) *companyActions {
	c.T.Run("ForceDeleteCompany", c.Testers[id].ExpectStatus(s).ForceDELETE)
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
