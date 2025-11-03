package resourceSeed

import (
	authModel "mynute-go/auth/config/db/model"
	coreModel "mynute-go/core/src/config/db/model"
)

var Appointment = &authModel.Resource{
	Name:        "appointment",
	Description: "Appointment resource",
	Table:       coreModel.AppointmentTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("appointment_id", "id"),
		authModel.MultipleQueryRef("appointment_id", "id"),
		authModel.MultipleBodyRef("appointment_id", "id"),
		authModel.MultiplePathRef("name", "name"),
	},
}

var Branch = &authModel.Resource{
	Name:        "branch",
	Description: "Branch resource",
	Table:       coreModel.BranchTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("branch_id", "id"),
		authModel.MultipleQueryRef("branch_id", "id"),
		authModel.MultipleBodyRef("branch_id", "id"),
		authModel.MultiplePathRef("name", "name"),
	},
}

var Client = &authModel.Resource{
	Name:        "client",
	Description: "Client resource",
	Table:       coreModel.ClientTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("client_id", "id"),
		authModel.MultipleQueryRef("client_id", "id"),
		authModel.MultipleBodyRef("client_id", "id"),
		authModel.MultiplePathRef("email", "email"),
	},
}

var Company = &authModel.Resource{
	Name:        "company",
	Description: "Company resource",
	Table:       coreModel.CompanyTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("company_id", "id"),
		authModel.MultipleQueryRef("company_id", "id"),
		authModel.MultipleBodyRef("company_id", "id"),
	},
}

var Employee = &authModel.Resource{
	Name:        "employee",
	Description: "Employee resource",
	Table:       coreModel.EmployeeTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("employee_id", "id"),
		authModel.MultipleQueryRef("employee_id", "id"),
		authModel.MultipleBodyRef("employee_id", "id"),
		authModel.MultiplePathRef("email", "email"),
	},
}

var Holiday = &authModel.Resource{
	Name:        "holiday",
	Description: "Holiday resource",
	Table:       coreModel.HolidayTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("holiday_id", "id"),
		authModel.MultipleQueryRef("holiday_id", "id"),
		authModel.MultipleBodyRef("holiday_id", "id"),
	},
}

var Role = &authModel.Resource{
	Name:        "role",
	Description: "Role resource",
	Table:       coreModel.RoleTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("role_id", "id"),
		authModel.MultipleQueryRef("role_id", "id"),
		authModel.MultipleBodyRef("role_id", "id"),
	},
}

var Sector = &authModel.Resource{
	Name:        "sector",
	Description: "Sector resource",
	Table:       coreModel.SectorTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("sector_id", "id"),
		authModel.MultipleQueryRef("sector_id", "id"),
		authModel.MultipleBodyRef("sector_id", "id"),
	},
}

var Service = &authModel.Resource{
	Name:        "service",
	Description: "Service resource",
	Table:       coreModel.ServiceTableName,
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("service_id", "id"),
		authModel.MultipleQueryRef("service_id", "id"),
		authModel.MultipleBodyRef("service_id", "id"),
	},
}

var Auth = &authModel.Resource{
	Name:        "auth",
	Description: "Auth resource",
	Table:       "auth",
	References: authModel.ResourceReferences{
		authModel.SingleQueryRef(),
		authModel.SinglePathRef(),
		authModel.MultiplePathRef("auth_id", "id"),
		authModel.MultipleQueryRef("auth_id", "id"),
		authModel.MultipleBodyRef("auth_id", "id"),
	},
}

var Resources = []*authModel.Resource{
	Appointment,
	Branch,
	Client,
	Company,
	Employee,
	Holiday,
	Role,
	Sector,
	Service,
	Auth,
}

