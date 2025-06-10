package e2e_test

import (
	"agenda-kaki-go/core"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib/FileBytes"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	"testing"
)

func Test_Company(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handlerT.NewTestErrorHandler(t)
	company := &modelT.Company{}

	tt.Describe("Company creation").Test(company.Create(200))

	tt.Describe("Company update design config").Test(company.Update(200, map[string]any{
		"design": mJSON.DesignConfig{
			Colors: mJSON.Colors{
				Primary:   "#FF5733",
				Secondary: "#33FF57",
				Tertiary:  "#3357FF",
			},
		},
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Company get by ID").Test(company.GetById(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Company get by name").Test(company.GetByName(200, company.Owner.X_Auth_Token, nil))
	tt.Describe("Company get by subdomain").Test(company.GetBySubdomain(200))

	tt.Describe("Upload logo image").Test(company.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get logo image").Test(company.GetImage(200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Upload banner, favicon, and background images").Test(company.UploadImages(200, map[string][]byte{
		"banner":     FileBytes.PNG_FILE_2,
		"favicon":    FileBytes.PNG_FILE_3,
		"background": FileBytes.PNG_FILE_4,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get logo image again").Test(company.GetImage(200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_1))
	tt.Describe("Get banner image").Test(company.GetImage(200, company.Created.Design.Images.Banner.URL, &FileBytes.PNG_FILE_2))
	tt.Describe("Get favicon image").Test(company.GetImage(200, company.Created.Design.Images.Favicon.URL, &FileBytes.PNG_FILE_3))
	tt.Describe("Get background image").Test(company.GetImage(200, company.Created.Design.Images.Background.URL, &FileBytes.PNG_FILE_4))

	tt.Describe("Overwrite logo image").Test(company.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_3,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get overwritten logo image").Test(company.GetImage(200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_3))
	tt.Describe("Recheck banner image").Test(company.GetImage(200, company.Created.Design.Images.Banner.URL, &FileBytes.PNG_FILE_2))
	tt.Describe("Recheck favicon image").Test(company.GetImage(200, company.Created.Design.Images.Favicon.URL, &FileBytes.PNG_FILE_3))
	tt.Describe("Recheck background image").Test(company.GetImage(200, company.Created.Design.Images.Background.URL, &FileBytes.PNG_FILE_4))

	tt.Describe("Change primary color only").Test(company.ChangeColors(200, mJSON.Colors{
		Primary: "#123456",
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Change all colors").Test(company.ChangeColors(200, mJSON.Colors{
		Primary:    "#654321",
		Secondary:  "#abcdef",
		Tertiary:   "#fedcba",
		Quaternary: "#123abc",
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Delete all images").Test(company.DeleteImages(200, []string{
		"logo",
		"banner",
		"favicon",
		"background",
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Reset colors").Test(company.ChangeColors(200, mJSON.Colors{}, company.Owner.X_Auth_Token, nil))
	tt.Describe("Company deletion").Test(company.Delete(200, company.Owner.X_Auth_Token, nil))
}
