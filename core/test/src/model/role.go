package model

import "mynute-go/core/src/config/db/model"

type Role struct {
	Created    *model.Role
	Company    *Company
	Auth_token string
	Employees  []*Employee
}

