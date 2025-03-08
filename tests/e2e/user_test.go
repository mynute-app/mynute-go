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
		"email":    lib.GenerateRandomEmail(),
		"name":     lib.GenerateRandomName("User Name"),
		"surname":  lib.GenerateRandomName("User Surname"),
		"password": lib.GenerateRandomString(10),
		"phone":    lib.GenerateRandomStrNumber(11), // 55977747309
	})
}
