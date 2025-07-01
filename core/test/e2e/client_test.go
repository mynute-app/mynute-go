package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
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
	tt.Describe("Client get by email").Test(client.GetByEmail(200))
	tt.Describe("Upload profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil))

	tt.Describe("Get profile image").Test(client.GetImage(200, client.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, nil))

	tt.Describe("Get overwritten profile image").Test(client.GetImage(200, client.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	tt.Describe("Client deletion").Test(client.Delete(200))
}
