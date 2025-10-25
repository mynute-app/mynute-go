package e2e_test

import (
	"mynute-go/core"
	"mynute-go/core/src/lib"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"

	"testing"
)

type Client struct {
	Created    model.Client
	Auth_token string
}

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)
	client := &model.Client{}

	tt.Describe("Client creation").Test(client.Create(200))

	tt.Describe("Send Login code by email").Test(client.SendLoginCode(200))

	loginCode, err := client.GetLoginCodeFromEmail()
	if err != nil {
		tt.Describe("Get login code from email").Test(err)
	}

	tt.Describe("Login by email code").Test(client.LoginByEmailCode(200, loginCode))

	tt.Describe("Client login by password").Test(client.Login(200))

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

	tt.Describe("Get profile image").Test(client.GetImage(200, client.Created.Meta.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, nil))

	tt.Describe("Get overwritten profile image").Test(client.GetImage(200, client.Created.Meta.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	img_url := client.Created.Meta.Design.Images.Profile.URL

	tt.Describe("Delete profile image").Test(client.DeleteImages(200, []string{"profile"}, nil))

	tt.Describe("Get deleted profile image").Test(client.GetImage(404, img_url, nil))

	tt.Describe("Login by email code with invalid code").Test(client.LoginByEmailCode(400, "000000"))

	tt.Describe("Send Login code by email").Test(client.SendLoginCode(200))

	loginCode, err = client.GetLoginCodeFromEmail()
	if err != nil {
		tt.Describe("Get login code from email").Test(err)
	}

	tt.Describe("Login by email code").Test(client.LoginByEmailCode(200, loginCode))

	tt.Describe("Upload profile image again logged in with email code").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil))

	tt.Describe("Client deletion").Test(client.Delete(200))

	tt.Describe("Get deleted client by email").Test(client.GetByEmail(404))
}
