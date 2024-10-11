package e2e_test

import (
	"agenda-kaki-go/tests/handlers"
	"agenda-kaki-go/tests/lib"
	"fmt"
	"testing"
)

var company = handlers.Tester{
	Entity:  "company",
	BaseURL: baseURL,
	PostBody: map[string]interface{}{
		"name":   lib.GenerateRandomName("Company"),
		"tax_id": fmt.Sprintf("%v", lib.GenerateRandomIntOfExactly(14)),
	},
	PatchBody: map[string]interface{}{"name": lib.GenerateRandomName("Company")},
}

func TestCompanyFlow(t *testing.T) {
	t.Run("CreateCompanyType", companyType.POST)
	company.PostBody["company_types"] = []map[string]interface{}{
		{"id": companyType.EntityID, "name": companyType.PostBody["name"]},
	}
	t.Logf("company.PostBody before CreateCompany: %+v", company.PostBody)
	t.Run("CreateCompany", company.POST)
	t.Run("UpdateCompany", company.PATCH)
	t.Run("DeleteCompany", company.DELETE)
	t.Run("DeleteCompanyType", companyType.DELETE)
}
