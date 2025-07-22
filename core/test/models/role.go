package modelT

import "mynute-go/core/config/db/model"

type Role struct {
	Created    *model.Role
	Company    *Company
	Auth_token string
	Employees  []*Employee
}
