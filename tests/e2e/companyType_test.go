package e2e_test

import (
	"agenda-kaki-go/tests/handlers"
	"agenda-kaki-go/tests/lib"
	"testing"
)

const baseURL = "http://localhost:3000"

// Run the test in debug mode to avoid cache.

func TestCompanyTypeFlow(t *testing.T) {
	tester := handlers.Tester{
		Entity:    "companyType",
		BaseURL:   baseURL,
		PostBody:  map[string]string{"name": lib.GenerateRandomName("CompanyType")},
		PatchBody: map[string]string{"name": lib.GenerateRandomName("CompanyType")},
	}
	t.Run("CreateCompanyType", tester.POST)

	t.Run("UpdateCompanyType", tester.PATCH)

	t.Run("DeleteCompanyType", tester.DELETE)
}
