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
	tt.Test(company.Create(200))
	tt.Test(company.Owner.VerifyEmail(200))
	tt.Test(company.Owner.Login(200))
	company.Auth_token = company.Owner.Auth_token
	branch := &modelT.Branch{}
	branch.Auth_token = company.Auth_token
	branch.Company = company
	tt.Test(branch.Create(200))
	tt.Test(branch.Update(200, map[string]any{
		"name": branch.Created.Name,
	}))
	tt.Test(branch.GetById(200))
	tt.Test(branch.GetByName(200))
	service := &modelT.Service{}
	service.Auth_token = company.Auth_token
	service.Company = company
	tt.Test(service.Create(200))
	tt.Test(branch.AddService(200, service, nil))
	tt.Test(company.Owner.AddBranch(200, branch, nil))
	tt.Test(company.GetById(200))
	tt.Test(branch.Delete(200))
}
