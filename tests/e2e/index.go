package e2e

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/tests/handlers"
	"fmt"
	"testing"
)

type IEntity interface {
	Make(n int)
	ClearDependencies()
	CreateDependencies(n int)
	GenerateTesters(n int)
}

var _ IBaseE2EActions = (*BaseE2EActions)(nil)

type IBaseE2EActions interface {
	GenerateTester(entity string, path string, postBody map[string]any, patchBody map[string]any) *BaseE2EActions
	CreateAllTesters(s int) *BaseE2EActions
	GetAllTesters(s int) *BaseE2EActions
	ForceGetAllTesters(s int) *BaseE2EActions
	UpdateAllTesters(s int) *BaseE2EActions
	DeleteAllTesters(s int) *BaseE2EActions
	ForceDeleteAllTesters(s int) *BaseE2EActions
	GetOneByIdTesters(s int) *BaseE2EActions
	ForceGetOneByIdTesters(s int) *BaseE2EActions
	RunAll() *BaseE2EActions
	SetTest(t *testing.T) *BaseE2EActions
}

type BaseE2EActions struct {
	T       *testing.T
	Testers []*handlers.Tester
}

func (b *BaseE2EActions) SetTest(t *testing.T) *BaseE2EActions {
	b.T = t
	return b
}

func (b *BaseE2EActions) GenerateTester(entity string, path string, postBody map[string]any, patchBody map[string]any) *BaseE2EActions {
	tester := &handlers.Tester{
		Entity:      path,
		BaseURL:     namespace.QueryKey.BaseURL,
		RelatedPath: path,
		PostBody:    postBody,
		PatchBody:   patchBody,
	}
	b.Testers = append(b.Testers, tester)
	return b
}

func (b *BaseE2EActions) CreateAllTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("Create%v", tester.Entity), tester.ExpectStatus(s).POST)
	}
	return b
}

func (b *BaseE2EActions) GetAllTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("Get%v", tester.Entity), tester.ExpectStatus(s).GET)
	}
	return b
}

func (b *BaseE2EActions) ForceGetAllTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("ForceGet%v", tester.Entity), tester.ExpectStatus(s).ForceGET)
	}
	return b
}

func (b *BaseE2EActions) UpdateAllTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("Update%v", tester.Entity), tester.ExpectStatus(s).PATCH)
	}
	return b
}

func (b *BaseE2EActions) DeleteAllTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("Delete%v", tester.Entity), tester.ExpectStatus(s).DELETE)
	}
	return b
}

func (b *BaseE2EActions) ForceDeleteAllTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("ForceDelete%v", tester.Entity), tester.ExpectStatus(s).ForceDELETE)
	}
	return b
}

func (b *BaseE2EActions) GetOneByIdTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("GetOne%v", tester.Entity), tester.ExpectStatus(s).GET)
	}
	return b
}

func (b *BaseE2EActions) ForceGetOneByIdTesters(s int) *BaseE2EActions {
	for _, tester := range b.Testers {
		b.T.Run(fmt.Sprintf("ForceGetOne%v", tester.Entity), tester.ExpectStatus(s).ForceGET)
	}
	return b
}

func (b *BaseE2EActions) RunAll() *BaseE2EActions {
	b.
		CreateAllTesters(201).
		GetAllTesters(200).
		UpdateAllTesters(200).
		GetAllTesters(200).
		ForceGetAllTesters(200).
		DeleteAllTesters(204).
		GetAllTesters(404).
		ForceGetAllTesters(200).
		ForceDeleteAllTesters(204).
		ForceGetAllTesters(404)
	return b
}
