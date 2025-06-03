package e2e_test

import (
	"agenda-kaki-go/core"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"

	"testing"
)

func Test_Branch(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	company := &modelT.Company{}
	tt := handlerT.NewTestErrorHandler(t)
	tt.Test(company.Create(200), "Company creation")
	tt.Test(company.Owner.VerifyEmail(200, nil), "Company owner email verification")
	tt.Test(company.Owner.Login(200, nil), "Company owner login")
	branch := &modelT.Branch{}
	branch.Company = company
	tt.Test(branch.Create(200, company.Owner.X_Auth_Token, nil), "Branch creation")
	tt.Test(branch.Update(200, map[string]any{
		"name": branch.Created.Name,
	}, company.Owner.X_Auth_Token, nil), "Branch update")
	tt.Test(branch.GetById(200, company.Owner.X_Auth_Token, nil), "Branch get by ID")
	tt.Test(branch.GetByName(200, company.Owner.X_Auth_Token, nil), "Branch get by name")
	service := &modelT.Service{}
	service.Company = company
	tt.Test(service.Create(200, company.Owner.X_Auth_Token, nil), "Service creation")
	tt.Test(branch.AddService(200, service, company.Owner.X_Auth_Token, nil), "Adding service to branch")
	tt.Test(company.Owner.AddBranch(200, branch, nil, nil), "Adding branch to company")
	tt.Test(company.GetById(200, company.Owner.X_Auth_Token, nil), "Getting company by ID")
	tt.Test(branch.Delete(200, company.Owner.X_Auth_Token, nil), "Branch deletion")
}
