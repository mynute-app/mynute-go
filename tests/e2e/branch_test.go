package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
	"fmt"
	"testing"
)

var _ e2e.IEntity = (*Branch)(nil)

type Branch struct {
	*e2e.BaseE2EActions
	company *Company
}

func (b *Branch) GenerateTesters(n int) {
	for i := 0; i < n; i++ {
		path := fmt.Sprintf("company/%d/branch", b.company.Testers[i].EntityID)
		b.GenerateTester(
			"branch",
			path,
			map[string]any{
				"name": lib.GenerateRandomName("Branch"),
			},
			map[string]any{
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
	company := &Company{
		BaseE2EActions: &e2e.BaseE2EActions{},
	}
	company.SetTest(b.T)
	company.Make(n)
	company.CreateAllTesters(201)
	b.company = company
}

func (b *Branch) ClearDependencies() {
	b.company.ForceDeleteAllTesters(204)
}

func TestBranchFlow(t *testing.T) {
	branch := &Branch{
		BaseE2EActions: &e2e.BaseE2EActions{},
	}
	branch.SetTest(t)
	branch.Make(2)
	branch.RunAll()
	branch.ClearDependencies()
}
