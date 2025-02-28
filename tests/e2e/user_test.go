package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
)

var _ e2e.IEntity = (*User)(nil)

type User struct {
	*e2e.BaseE2EActions
}

func (c *User) CreateDependencies(n int) {}

func (c *User) ClearDependencies() {}

func (c *User) Make(n int) {
	c.GenerateTesters(n)
}

//TOOD: implement
func (c *User) GenerateTesters(n int) {
	for range n {
		c.GenerateTester(
			"user",
			"user",
			map[string]interface{}{
				"email":    lib.GenerateRandomEmail(),
				"password": lib.GenerateRandomString(10),
				"name": lib.GenerateRandomName("User Name"),
				"surname": lib.GenerateRandomName("User Surname"),
				"phone": lib.GenerateRandomIntOfExactly(10),
				"verification_code": lib.GenerateRandomString(20),
				"verified": false,
			},
			map[string]interface{}{
				"verified": true,
			},
		)
	}
}

