package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
	"testing"
)

var _ e2e.IEntity = (*CompanyType)(nil)

type CompanyType struct {
	*e2e.BaseE2EActions
}

func (c *CompanyType) GenerateTesters(n int) {
	for i := 0; i < n; i++ {
		c.GenerateTester(
			"companyType",
			"companyType",
			map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")},
			map[string]interface{}{"name": lib.GenerateRandomName("CompanyType")},
		)
	}
}

func (c *CompanyType) Make(n int) {
	c.GenerateTesters(n)
}

func (c *CompanyType) CreateDependencies(n int) {}

func (c *CompanyType) ClearDependencies() {}


func TestCompanyTypeFlow(t *testing.T) {
	companyType := &CompanyType{
		BaseE2EActions: &e2e.BaseE2EActions{},
	}
	companyType.SetTest(t)
	companyType.Make(1)
	companyType.RunAll()
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