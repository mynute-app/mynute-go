package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
	"testing"
)

var _ e2e.IEntity = (*Branch)(nil)

type Branch struct {
	*e2e.BaseE2EActions
	company *Company
}

func (b *Branch) GenerateTesters(n int) {
	for i := 0; i < n; i++ {
		b.GenerateTester(
			"branch",
			"branch",
			map[string]interface{}{
				"name": lib.GenerateRandomName("Branch"),
				"company": map[string]interface{}{
					"id": b.company.Testers[i].EntityID,
					"name": b.company.Testers[i].PostBody["name"],
				},
			},
			map[string]interface{}{
				"name": lib.GenerateRandomName("Branch"),
			},
		)
	}
}

func (b *Branch) Make(n int) {
	b.CreateDependencies(n)
	b.GenerateTesters(n)
}

func (b *Branch) CreateDependencies(n int) {
	company := &Company{}
	company.SetTest(b.T)
	company.Make(n)
	company.CreateAllTesters(201)
	b.company = company
}

func (b *Branch) ClearDependencies() {
	b.company.ForceDeleteAllTesters(204)
}

func TestBranchFlow(t *testing.T) {
	branch := &Branch{}
	branch.SetTest(t)
	branch.Make(10)
	branch.RunAll()
	branch.ClearDependencies()
}
