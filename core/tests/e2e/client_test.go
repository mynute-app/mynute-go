package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	models_test "agenda-kaki-go/core/tests/models"

	"testing"
)

type Client struct {
	Created    model.ClientFull
	Auth_token string
}

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	client := &models_test.Client{}
	client.Create(t, 200)
	client.VerifyEmail(t, 200)
	client.Login(t, 200)
	client.Update(t, 200, map[string]any{"name": "Updated Client Name"})
	client.GetByEmail(t, 200)
	client.Delete(t, 200)
}
