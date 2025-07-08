package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
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
	servicesID := []DTO.ServiceID{{ID: service.Created.ID}}
	BranchWorkSchedule := modelT.GetExampleBranchWorkSchedule(branch.Created.ID, servicesID)
	tt.Describe("Branch work schedule fail creation").Test(branch.CreateWorkSchedule(400, BranchWorkSchedule, company.Owner.X_Auth_Token, nil))
	tt.Describe("Adding service to branch").Test(branch.AddService(200, service, company.Owner.X_Auth_Token, nil))
	tt.Describe("Branch work schedule success creation").Test(branch.CreateWorkSchedule(200, BranchWorkSchedule, company.Owner.X_Auth_Token, nil))
	wr := branch.Created.BranchWorkSchedule[0]
	tt.Describe("Updating branch work schedule").Test(branch.UpdateWorkRange(200, &wr, map[string]any{
		"start":    "07:00",
		"end":      "20:00",
		"timezone": "America/Sao_Paulo",
	}, company.Owner.X_Auth_Token, nil))
	tt.Describe("Deleting branch work range").Test(branch.DeleteWorkRange(200, &wr, company.Owner.X_Auth_Token, nil))
	tt.Describe("Adding branch to Owner").Test(company.Owner.AddBranch(200, branch, nil, nil))
	tt.Describe("Getting company by ID").Test(company.GetById(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Changing branch company_id").Test(branch.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Upload profile image").Test(branch.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get profile image").Test(branch.GetImage(200, branch.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(branch.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get overwritten profile image").Test(branch.GetImage(200, branch.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	img_url := branch.Created.Design.Images.Profile.URL

	tt.Describe("Delete profile image").Test(branch.DeleteImages(200, []string{"profile"}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get deleted profile image").Test(branch.GetImage(404, img_url, nil))

	tt.Describe("Branch deletion").Test(branch.Delete(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get deleted branch by ID").Test(branch.GetById(404, company.Owner.X_Auth_Token, nil))
}
