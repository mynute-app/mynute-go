package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/lib"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"

	"testing"
)

func Test_Service(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	client := &modelT.Client{}
	tt := handlerT.NewTestErrorHandler(t)
	tt.Test(client.Set(), "Client setup") // This sets up client, company, branches, and services.
	company := &modelT.Company{}
	tt.Test(company.Set(), "Company setup")
	service := &modelT.Service{}
	service.X_Auth_Token = company.Owner.X_Auth_Token
	service.Company = company
	tt.Test(service.Create(200), "Service creation")
	tt.Test(service.Update(200, map[string]any{
		"name": lib.GenerateRandomName("Updated Service"),
	}), "Service update")
	tt.Test(service.GetById(200, nil), "Service get by ID")
	tt.Test(service.GetByName(200), "Service get by name")
	branch := &modelT.Branch{}
	branch.X_Auth_Token = company.Owner.X_Auth_Token
	branch.Company = company
	tt.Test(branch.Create(200), "Branch creation")
	tt.Test(branch.AddService(200, service, nil), "Branch add service")
	tt.Test(service.Delete(200), "Service deletion")
	tt.Test(branch.Delete(200), "Branch deletion")
}
