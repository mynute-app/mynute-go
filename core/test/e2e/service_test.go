package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/lib"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"

	"testing"

	"github.com/google/uuid"
)

func Test_Service(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handlerT.NewTestErrorHandler(t)

	client := &modelT.Client{}
	tt.Describe("Client setup").Test(client.Set()) // Sets up client, company, branches, and services

	company := &modelT.Company{}
	tt.Describe("Company setup").Test(company.Set())

	service := &modelT.Service{Company: company}
	tt.Describe("Service creation").Test(service.Create(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Service update").Test(service.Update(200, map[string]any{
		"name": lib.GenerateRandomName("Updated Service"),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Service get by ID").Test(service.GetById(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Service get by name").Test(service.GetByName(200, company.Owner.X_Auth_Token, nil))

	branch := &modelT.Branch{Company: company}
	tt.Describe("Branch creation").Test(branch.Create(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Branch add service").Test(branch.AddService(200, service, company.Owner.X_Auth_Token, nil))

	tt.Describe("Changing service company_id").Test(service.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Service deletion").Test(service.Delete(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Branch deletion").Test(branch.Delete(200, company.Owner.X_Auth_Token, nil))
}
