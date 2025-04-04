package model

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	Name        string              `json:"name" gorm:"unique;not null"`
	Description string              `json:"description"`
	Table       string              `json:"table"`
	References  []ResourceReference `json:"references"`
}

type ResourceReference struct {
	DatabaseKey string `json:"database_key"` // The key in the database, e.g. "id", "tax_id".
	RequestKey  string `json:"request_key"`  // The key in the request body, query, path, header...
	RequestRef  string `json:"request_ref"`  // "query", "body", "header", "path".
}

func SingleQueryRef() ResourceReference {
	return ResourceReference{
		DatabaseKey: "id",
		RequestKey:  "id",
		RequestRef:  "query",
	}
}

func MultipleQueryRef(name string) ResourceReference {
	return ResourceReference{
		DatabaseKey: "id",
		RequestKey:  name,
		RequestRef:  "query",
	}
}

func MultipleBodyRef(name string) ResourceReference {
	return ResourceReference{
		DatabaseKey: "id",
		RequestKey:  name,
		RequestRef:  "body",
	}
}

var SingleIdQueryRef = ResourceReference{
	DatabaseKey: "id",
	RequestKey:  "id",
	RequestRef:  "query",
}

var CompanyResource = &Resource{
	Name:        "company",
	Description: "Company resource",
	Table:       "companies",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("company_id"),
		MultipleBodyRef("company_id"),
	},
}

var ClientResource = &Resource{
	Name:        "client",
	Description: "Client resource",
	Table:       "clients",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("client_id"),
		MultipleBodyRef("client_id"),
	},
}
