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

var AppointmentResource = &Resource{
	Name:        "appointment",
	Description: "Appointment resource",
	Table:       "appointments",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("appointment_id"),
		MultipleBodyRef("appointment_id"),
		MultipleQueryRef("name"),
	},
}

var BranchResource = &Resource{
	Name:        "branch",
	Description: "Branch resource",
	Table:       "branches",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("branch_id"),
		MultipleBodyRef("branch_id"),
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

var EmployeeResource = &Resource{
	Name:        "employee",
	Description: "Employee resource",
	Table:       "employees",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("employee_id"),
		MultipleBodyRef("employee_id"),
	},
}

var HolidayResource = &Resource{
	Name:        "holiday",
	Description: "Holiday resource",
	Table:       "holidays",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("holiday_id"),
		MultipleBodyRef("holiday_id"),
	},
}

var RoleResource = &Resource{
	Name:        "role",
	Description: "Role resource",
	Table:       "roles",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("role_id"),
		MultipleBodyRef("role_id"),
	},
}

var SectorResource = &Resource{
	Name:        "sector",
	Description: "Sector resource",
	Table:       "sectors",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("sector_id"),
		MultipleBodyRef("sector_id"),
	},
}

var ServiceResource = &Resource{
	Name:        "service",
	Description: "Service resource",
	Table:       "services",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("service_id"),
		MultipleBodyRef("service_id"),
	},
}

var AuthResource = &Resource{
	Name:        "auth",
	Description: "Auth resource",
	Table:       "auth",
	References: []ResourceReference{
		SingleQueryRef(),
		MultipleQueryRef("auth_id"),
		MultipleBodyRef("auth_id"),
	},
}
