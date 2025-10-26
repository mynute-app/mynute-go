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

	tt.Describe("Client get by email").Test(client.GetByEmail(200))

	tt.Describe("Login with password").Test(client.Login(401, "password"))

	tt.Describe("Login with email code").Test(client.Login(200, "email_code"))

	tt.Describe("Login by password with invalid password").Test(client.LoginByPassword(401, "invalid_password"))

	tt.Describe("Login with password").Test(client.Login(200, "password"))

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

	// Re-login with new password to get a fresh token
	tt.Describe("Login with new password").Test(client.LoginByPassword(200, new_password))

	// Test password reset by email
	tt.Describe("Reset password by email").Test(client.ResetPasswordByEmail(200))

	// Test that the old password no longer works
	tt.Describe("Login with old password fails").Test(client.LoginByPassword(401, new_password))

	// Test that new password from email works
	tt.Describe("Login with password from email").Test(client.LoginByPassword(200, client.Created.Password))

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

	tt.Describe("Upload profile image again logged in with email code").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil))

	tt.Describe("Client deletion").Test(client.Delete(200))

	tt.Describe("Get deleted client by email").Test(client.GetByEmail(404))
}
