package e2e_test

import (
	"agenda-kaki-go/core"
	models_test "agenda-kaki-go/core/tests/models"

	"testing"
)

func Test_Branch(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &models_test.Company{}
	company.Create(t, 200)
	company.Owner.VerifyEmail(t, 200)
	company.Owner.Login(t, 200)
	company.Auth_token = company.Owner.Auth_token
	branch := &models_test.Branch{}
	branch.Auth_token = company.Auth_token
	branch.Company = company
	branch.Create(t, 200)
	branch.Update(t, 200, map[string]any{
		"name": branch.Created.Name,
	})
	branch.GetById(t, 200)
	branch.GetByName(t, 200)
	service := &models_test.Service{}
	service.Auth_token = company.Auth_token
	service.Company = company
	service.Create(t, 200)
	branch.AddService(t, 200, service, nil)
	company.Owner.AddBranch(t, 200, branch, nil)
	company.GetById(t, 200)
	branch.Delete(t, 200)
}

