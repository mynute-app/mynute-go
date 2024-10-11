package e2e_test

import (
	"agenda-kaki-go/tests/handlers"
	"agenda-kaki-go/tests/lib"
	"testing"
)

const baseURL = "http://localhost:3000"

var companyType = handlers.Tester{
	Entity:    "companyType",
	BaseURL:   baseURL,
	PostBody:  map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")},
	PatchBody: map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")},
}

// Run the test in debug mode to avoid cache.

func TestCompanyTypeFlow(t *testing.T) {
	t.Run("CreateCompanyType", companyType.POST)
	t.Run("UpdateCompanyType", companyType.PATCH)
	t.Run("DeleteCompanyType", companyType.DELETE)
}
