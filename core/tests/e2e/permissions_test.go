package e2e_test

import (
	"agenda-kaki-go/core"
	"testing"
)

func Test_Permissions(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Set(t)
	client := &Client{}
	client.Set(t)
}

