package e2e_test

import (
	"agenda-kaki-go/core"
	FileBytes "agenda-kaki-go/core/lib/file_bytes"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"

	"testing"

	"github.com/google/uuid"
)

func Test_Branch(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handlerT.NewTestErrorHandler(t)

	company := &modelT.Company{}
	tt.Describe("Company creation").Test(company.Create(200))

	branch := &modelT.Branch{}
	branch.Company = company
	tt.Describe("Branch creation").Test(branch.Create(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Branch update").Test(branch.Update(200, map[string]any{
		"name": branch.Created.Name,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Branch get by ID").Test(branch.GetById(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Branch get by name").Test(branch.GetByName(200, company.Owner.X_Auth_Token, nil))

	service := &modelT.Service{}
	service.Company = company

	tt.Describe("Service creation").Test(service.Create(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Adding service to branch").Test(branch.AddService(200, service, company.Owner.X_Auth_Token, nil))
	tt.Describe("Adding branch to Owner").Test(company.Owner.AddBranch(200, branch, nil, nil))
	tt.Describe("Getting company by ID").Test(company.GetById(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Changing branch company_id").Test(branch.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Upload profile image").Test(branch.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get profile image").Test(branch.GetImage(200, branch.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(branch.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_3,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get overwritten logo image").Test(branch.GetImage(200, branch.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_3))

	tt.Describe("Branch deletion").Test(branch.Delete(200, company.Owner.X_Auth_Token, nil))
}
