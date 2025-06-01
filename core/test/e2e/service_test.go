package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/lib"
	models_test "agenda-kaki-go/core/test/models"

	"testing"
)

func Test_Service(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	client := &models_test.Client{}
	client.Set(t)
	company := &models_test.Company{}
	company.Set(t)
	service := &models_test.Service{}
	service.Auth_token = company.Owner.Auth_token
	service.Company = company
	service.Create(t, 200)
	service.Update(t, 200, map[string]any{
		"name": lib.GenerateRandomName("Updated Service"),
	})
	service.GetById(t, 200, nil)
	service.GetByName(t, 200)
	branch := &models_test.Branch{}
	branch.Auth_token = company.Owner.Auth_token
	branch.Company = company
	branch.Create(t, 200)
	branch.AddService(t, 200, service, nil)
	service.Delete(t, 200)
	branch.Delete(t, 200)
}
