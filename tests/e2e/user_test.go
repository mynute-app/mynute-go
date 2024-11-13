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
	for i := 0; i < n; i++ {
		c.GenerateTester(
			"companyType",
			"companyType",
			map[string]interface{}{"name": lib.GenerateRandomName("EmployeeInfo")},
			map[string]interface{}{"name": lib.GenerateRandomName("EmployeeInfo")},
		)
	}
}

