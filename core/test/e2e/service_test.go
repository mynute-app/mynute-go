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
	tt.Test(client.Set())
	company := &modelT.Company{}
	tt.Test(company.Set())
	service := &modelT.Service{}
	service.Auth_token = company.Owner.Auth_token
	service.Company = company
	tt.Test(service.Create(200))
	tt.Test(service.Update(200, map[string]any{
		"name": lib.GenerateRandomName("Updated Service"),
	}))
	tt.Test(service.GetById(200, nil))
	tt.Test(service.GetByName(200))
	branch := &modelT.Branch{}
	branch.Auth_token = company.Owner.Auth_token
	branch.Company = company
	tt.Test(branch.Create(200))
	tt.Test(branch.AddService(200, service, nil))
	tt.Test(service.Delete(200))
	tt.Test(branch.Delete(200))
}
