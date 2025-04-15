package e2e_test

import "agenda-kaki-go/core/config/db/model"

type Role struct {
	created model.Role
	company *Company
	auth_token string
	employees []*Employee
}