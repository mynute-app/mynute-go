package e2e_test

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/tests/handlers"
	"testing"
)

func Test_User_Create_Success(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(namespace.QueryKey.BaseURL + "/user")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"name":     lib.GenerateRandomName("User Name"),
		"surname":  lib.GenerateRandomName("User Surname"),
		"password": "1VerySecurePassword!",
		"phone":    lib.GenerateRandomStrNumber(11), // 55977747309
	})
}

func Test_Login_Success(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(namespace.QueryKey.BaseURL + "/login")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"password": "1VerySecurePassword!",
	})
}
