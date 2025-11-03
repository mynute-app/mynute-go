package e2e_test

import (
	"mynute-go/core"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/lib/file_bytes"
	"mynute-go/core/test/src/handler"
	"mynute-go/core/test/src/model"

	"testing"

	"github.com/google/uuid"
)

func Test_Service(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	client := &model.Client{}
	tt.Describe("Client setup").Test(client.Set()) // Sets up client, company, branches, and services

	company := &model.Company{}
	tt.Describe("Company setup").Test(company.Set())

	service := &model.Service{Company: company}
	tt.Describe("Service creation").Test(service.Create(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Service update").Test(service.Update(200, map[string]any{
		"name": lib.GenerateRandomName("Updated Service"),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Service get by ID").Test(service.GetById(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Service get by name").Test(service.GetByName(200, company.Owner.X_Auth_Token, nil))

	branch := &model.Branch{Company: company}
	tt.Describe("Branch creation").Test(branch.Create(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Branch add service").Test(branch.AddService(200, service, company.Owner.X_Auth_Token, nil))

	tt.Describe("Changing service company_id").Test(service.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Upload profile image").Test(service.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get profile image").Test(service.GetImage(200, service.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(service.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get overwritten profile image").Test(service.GetImage(200, service.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	img_url := service.Created.Design.Images.Profile.URL

	tt.Describe("Delete profile image").Test(service.DeleteImages(200, []string{"profile"}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get deleted profile image").Test(service.GetImage(404, img_url, nil))

	tt.Describe("Service deletion").Test(service.Delete(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Get deleted service by ID").Test(service.GetById(404, company.Owner.X_Auth_Token, nil))

	cService := company.Services[0]
	tt.Describe("Get service availability by ID").Test(cService.GetAvailability(200, nil, 0, 10))
}


