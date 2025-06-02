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
	company := &modelT.Company{}
	tt := handlerT.NewTestErrorHandler(t)
	tt.Test(company.Create(200))
	tt.Test(company.Owner.VerifyEmail(200))
	tt.Test(company.Owner.Login(200))
	company.Auth_token = company.Owner.Auth_token
	tt.Test(company.Update(200, map[string]any{"design": mJSON.DesignConfig{
		Colors: mJSON.Colors{
			Primary:   "#FF5733",
			Secondary: "#33FF57",
			Tertiary:  "#3357FF",
		},
	}}))
	tt.Test(company.GetById(200))
	tt.Test(company.GetByName(200))
	tt.Test(company.GetBySubdomain(200))
	tt.Test(company.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	}))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_1))
	tt.Test(company.UploadImages(200, map[string][]byte{
		"banner":     FileBytes.PNG_FILE_2,
		"favicon":    FileBytes.PNG_FILE_3,
		"background": FileBytes.PNG_FILE_4,
	}))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_1))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Banner.URL, &FileBytes.PNG_FILE_2))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Favicon.URL, &FileBytes.PNG_FILE_3))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Background.URL, &FileBytes.PNG_FILE_4))
	tt.Test(company.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_3,
	}))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_3))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Banner.URL, &FileBytes.PNG_FILE_2))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Favicon.URL, &FileBytes.PNG_FILE_3))
	tt.Test(company.GetImage(200, company.Created.Design.Images.Background.URL, &FileBytes.PNG_FILE_4))
	tt.Test(company.ChangeColors(200, mJSON.Colors{
		Primary: "#123456",
	}))
	tt.Test(company.ChangeColors(200, mJSON.Colors{
		Primary:    "#654321",
		Secondary:  "#abcdef",
		Tertiary:   "#fedcba",
		Quaternary: "#123abc",
	}))
	tt.Test(company.DeleteImages(200, []string{
		"logo",
		"banner",
		"favicon",
		"background",
	}))
	tt.Test(company.ChangeColors(200, mJSON.Colors{}))
	tt.Test(company.Delete(200))
}
