package e2e_test

import (
	"agenda-kaki-go/core"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib/FileBytes"
	models_test "agenda-kaki-go/core/test/models"

	"testing"
)

func Test_Company(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &models_test.Company{}
	company.Create(t, 200)
	company.Owner.VerifyEmail(t, 200)
	company.Owner.Login(t, 200)
	company.Auth_token = company.Owner.Auth_token
	company.Update(t, 200, map[string]any{"design": mJSON.DesignConfig{
		Colors: mJSON.Colors{
			Primary:   "#FF5733",
			Secondary: "#33FF57",
			Tertiary:  "#3357FF",
		},
	}})
	company.GetById(t, 200)
	company.GetByName(t, 200)
	company.GetBySubdomain(t, 200)
	company.UploadImages(t, 200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	})
	company.GetImage(t, 200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_1)
	company.UploadImages(t, 200, map[string][]byte{
		"banner":     FileBytes.PNG_FILE_2,
		"favicon":    FileBytes.PNG_FILE_3,
		"background": FileBytes.PNG_FILE_4,
	})
	company.GetImage(t, 200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_1)
	company.GetImage(t, 200, company.Created.Design.Images.Banner.URL, &FileBytes.PNG_FILE_2)
	company.GetImage(t, 200, company.Created.Design.Images.Favicon.URL, &FileBytes.PNG_FILE_3)
	company.GetImage(t, 200, company.Created.Design.Images.Background.URL, &FileBytes.PNG_FILE_4)
	company.UploadImages(t, 200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_3,
	})
	company.GetImage(t, 200, company.Created.Design.Images.Logo.URL, &FileBytes.PNG_FILE_3)
	company.GetImage(t, 200, company.Created.Design.Images.Banner.URL, &FileBytes.PNG_FILE_2)
	company.GetImage(t, 200, company.Created.Design.Images.Favicon.URL, &FileBytes.PNG_FILE_3)
	company.GetImage(t, 200, company.Created.Design.Images.Background.URL, &FileBytes.PNG_FILE_4)
	company.ChangeColors(t, 200, mJSON.Colors{
		Primary: "#123456",
	})
	company.ChangeColors(t, 200, mJSON.Colors{
		Primary:    "#654321",
		Secondary:  "#abcdef",
		Tertiary:   "#fedcba",
		Quaternary: "#123abc",
	})
	company.DeleteImages(t, 200, []string{
		"logo",
		"banner",
		"favicon",
		"background",
	})
	company.ChangeColors(t, 200, mJSON.Colors{})
	company.Delete(t, 200)
}
