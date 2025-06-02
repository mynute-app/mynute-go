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
	tt.Test(client.Create(200))
	tt.Test(client.VerifyEmail(200))
	tt.Test(client.Login(200))
	tt.Test(client.Update(200, map[string]any{"name": "Updated Client Name"}))
	tt.Test(client.GetByEmail(200))
	tt.Test(client.Delete(200))
}
