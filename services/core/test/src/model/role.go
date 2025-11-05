package model

import "mynute-go/services/core/config/db/model"

type Role struct {
	Created    *model.Role
	Company    *Company
	Auth_token string
	Employees  []*Employee
}
