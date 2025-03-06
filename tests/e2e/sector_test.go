package e2e_test

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/tests/e2e"
	"testing"
)

var _ e2e.IEntity = (*Sector)(nil)

type Sector struct {
	*e2e.BaseE2EActions
}

func (c *Sector) GenerateTesters(n int) {
	for i := 0; i < n; i++ {
		c.GenerateTester(
			"companyType",
			"companyType",
			map[string]any{"name": lib.GenerateRandomName("Sector")},
			map[string]any{"name": lib.GenerateRandomName("Sector")},
		)
	}
}

func (c *Sector) Make(n int) {
	c.GenerateTesters(n)
}

func (c *Sector) CreateDependencies(n int) {}

func (c *Sector) ClearDependencies() {}

func TestCompanyTypeFlow(t *testing.T) {
	companyType := &Sector{
		BaseE2EActions: &e2e.BaseE2EActions{},
	}
	companyType.SetTest(t)
	companyType.Make(10)
	companyType.RunAll()
}

// type Sector struct {
// 	T *testing.T
// }

// type companyTypeActios struct {
// 	*Sector
// 	Testers []handler.Tester
// }

// func (c *Sector) Init() *companyTypeActios {
// 	return &companyTypeActios{Sector: c}
// }

// func (c *companyTypeActios) LoadTester(n int) *companyTypeActios {
// 	for i := 0; i < n; i++ {
// 		companyType := handler.Tester{
// 			Entity:    "companyType",
// 			RelatedPath: "companyType",
// 			BaseURL:   namespace.GeneralKey.BaseURL,
// 			PostBody:  map[string]any{"name": lib.GenerateRandomName("Sector")},
// 			PatchBody: map[string]any{"name": lib.GenerateRandomName("Sector")},
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
