package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"os"
	"testing"
)

func Test_Debug_CWD(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	fmt.Println("Current Working Directory:", dir)
}

func Test_User_Create_Success(t *testing.T) {
	server := core.NewServer().Run("test")
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
	defer server.Shutdown()
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
