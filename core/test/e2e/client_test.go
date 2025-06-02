package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"

	"testing"
)

type Client struct {
	Created    model.ClientFull
	Auth_token string
}

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	client := &modelT.Client{}
	tt := handlerT.NewTestErrorHandler(t)
	tt.Test(client.Create(200), "Client creation")
	tt.Test(client.VerifyEmail(200), "Client email verification")
	tt.Test(client.Login(200), "Client login")
	tt.Test(client.Update(200, map[string]any{"name": "Updated Client Name"}), "Client update")
	tt.Test(client.GetByEmail(200), "Client get by email")
	tt.Test(client.Delete(200), "Client deletion")
}
