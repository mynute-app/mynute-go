package e2e_test

import (
	"agenda-kaki-go/core"
	handler "agenda-kaki-go/core/tests/handlers"
	"testing"
)

func Test_Permissions(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Set(t)
	client := &Client{}
	client.Set(t)
	// Client tries to change something in the company
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/company/" + company.created.ID.String())
	http.ExpectStatus(401)
	http.Header("Authorization", client.auth_token)
	http.Send(map[string]any{
		"name": "New Company Name",
	})
	// Client tries to change something on himself
	http.Method("PATCH")
	http.URL("/client/" + client.created.ID.String())
	http.ExpectStatus(200)
	http.Header("Authorization", client.auth_token)
	http.Send(map[string]any{
		"name": "New Client Name",
	})

}
