package modelT

import "agenda-kaki-go/core/config/db/model"

type Role struct {
	Created    model.Role
	Company    *Company
	Auth_token string
	Employees  []*Employee
}
