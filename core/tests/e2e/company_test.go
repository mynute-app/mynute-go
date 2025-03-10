package e2e_test

import "testing"

type Company struct{}

func Test_Company(t *testing.T) {
	user := &User{}
	user.Create(t)
	user.VerifyEmail(t)
}

func (c *Company) Create(t *testing.T, user *User) {

}
