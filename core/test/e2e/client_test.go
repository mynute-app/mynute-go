package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	FileBytes "agenda-kaki-go/core/lib/file_bytes"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"

	"testing"
)

type Client struct {
	Created    model.Client
	Auth_token string
}

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handlerT.NewTestErrorHandler(t)
	client := &modelT.Client{}

	tt.Describe("Client creation").Test(client.Create(200))
	tt.Describe("Client email verification").Test(client.VerifyEmail(200))
	tt.Describe("Client login").Test(client.Login(200))
	tt.Describe("Client update").Test(client.Update(200, map[string]any{
		"name": "Updated Client Name",
	}))
	tt.Describe("Client update").Test(client.Update(400, map[string]any{
		"name":     "Should Fail Update on Client Name",
		"password": "newpswrd123",
	}))
	new_password := lib.GenerateValidPassword()
	tt.Describe("Client update").Test(client.Update(200, map[string]any{
		"name":     "Should Succeed Update on Client Name",
		"password": new_password,
	}))
	tt.Describe("Client update").Test(client.Update(401, map[string]any{
		"password": "NewPswrd1@!",
	}))
	client.Created.Password = new_password // Update the password in the client model
	tt.Describe("Client get by email").Test(client.GetByEmail(401))
	tt.Describe("Client login").Test(client.Login(200))
	tt.Describe("Client get by email").Test(client.GetByEmail(200))

	tt.Describe("Upload profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil))

	tt.Describe("Get profile image").Test(client.GetImage(200, client.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, nil))

	tt.Describe("Get overwritten profile image").Test(client.GetImage(200, client.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	img_url := client.Created.Design.Images.Profile.URL

	tt.Describe("Delete profile image").Test(client.DeleteImages(200, []string{"profile"}, nil))

	tt.Describe("Get deleted profile image").Test(client.GetImage(404, img_url, nil))

	tt.Describe("Client deletion").Test(client.Delete(200))

	tt.Describe("Get deleted client by email").Test(client.GetByEmail(404))
}
