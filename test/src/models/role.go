package modelT

import "mynute-go/src/config/db/model"

type Role struct {
	Created    *model.Role
	Company    *Company
	Auth_token string
	Employees  []*Employee
}
