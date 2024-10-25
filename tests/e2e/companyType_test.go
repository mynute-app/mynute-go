package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
	"testing"
)

func TestCompanyTypeFlow(t *testing.T) {
	companyType := &e2e.BaseE2EActions{T: t}
	postBody := map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")}
	patchBody := map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")}
	companyType.GenerateTesters(5, "companyType", "companyType", postBody, patchBody).RunAll()
}

// type CompanyType struct {
// 	T *testing.T
// }

// type companyTypeActios struct {
// 	*CompanyType
// 	Testers []handlers.Tester
// }

// func (c *CompanyType) Init() *companyTypeActios {
// 	return &companyTypeActios{CompanyType: c}
// }

// func (c *companyTypeActios) LoadTester(n int) *companyTypeActios {
// 	for i := 0; i < n; i++ {
// 		companyType := handlers.Tester{
// 			Entity:    "companyType",
// 			RelatedPath: "companyType",
// 			BaseURL:   namespace.GeneralKey.BaseURL,
// 			PostBody:  map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")},
// 			PatchBody: map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")},
// 		}
// 		c.Testers = append(c.Testers, companyType)
// 	}
// 	return c
// }

// Run the test in debug mode to avoid cache.

// func TestCompanyTypeFlow(t *testing.T) {
// 	t.Run("CreateCompanyType", companyType.ExpectStatus(201).POST)
// 	t.Run("UpdateCompanyType", companyType.ExpectStatus(200).PATCH)
// 	t.Run("DeleteCompanyType", companyType.ExpectStatus(204).DELETE)
// 	t.Run("ForceDeleteCompanyType", companyType.ExpectStatus(204).ForceDELETE)
// }