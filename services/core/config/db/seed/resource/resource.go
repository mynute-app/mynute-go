package resourceSeed

import (
	
	model "mynute-go/services/core/config/db/model"
)

var Appointment = &model.Resource{
	Name:        "appointment",
	Description: "Appointment resource",
	Table:       model.AppointmentTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("appointment_id", "id"),
		model.MultipleQueryRef("appointment_id", "id"),
		model.MultipleBodyRef("appointment_id", "id"),
		model.MultiplePathRef("name", "name"),
	},
}

var Branch = &model.Resource{
	Name:        "branch",
	Description: "Branch resource",
	Table:       model.BranchTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("branch_id", "id"),
		model.MultipleQueryRef("branch_id", "id"),
		model.MultipleBodyRef("branch_id", "id"),
		model.MultiplePathRef("name", "name"),
	},
}

var Client = &model.Resource{
	Name:        "client",
	Description: "Client resource",
	Table:       model.ClientTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("client_id", "id"),
		model.MultipleQueryRef("client_id", "id"),
		model.MultipleBodyRef("client_id", "id"),
		model.MultiplePathRef("email", "email"),
	},
}

var Company = &model.Resource{
	Name:        "company",
	Description: "Company resource",
	Table:       model.CompanyTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("company_id", "id"),
		model.MultipleQueryRef("company_id", "id"),
		model.MultipleBodyRef("company_id", "id"),
	},
}

var Employee = &model.Resource{
	Name:        "employee",
	Description: "Employee resource",
	Table:       model.EmployeeTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("employee_id", "id"),
		model.MultipleQueryRef("employee_id", "id"),
		model.MultipleBodyRef("employee_id", "id"),
		model.MultiplePathRef("email", "email"),
	},
}

var Holiday = &model.Resource{
	Name:        "holiday",
	Description: "Holiday resource",
	Table:       model.HolidayTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("holiday_id", "id"),
		model.MultipleQueryRef("holiday_id", "id"),
		model.MultipleBodyRef("holiday_id", "id"),
	},
}

var Role = &model.Resource{
	Name:        "role",
	Description: "Role resource",
	Table:       model.RoleTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("role_id", "id"),
		model.MultipleQueryRef("role_id", "id"),
		model.MultipleBodyRef("role_id", "id"),
	},
}

var Sector = &model.Resource{
	Name:        "sector",
	Description: "Sector resource",
	Table:       model.SectorTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("sector_id", "id"),
		model.MultipleQueryRef("sector_id", "id"),
		model.MultipleBodyRef("sector_id", "id"),
	},
}

var Service = &model.Resource{
	Name:        "service",
	Description: "Service resource",
	Table:       model.ServiceTableName,
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("service_id", "id"),
		model.MultipleQueryRef("service_id", "id"),
		model.MultipleBodyRef("service_id", "id"),
	},
}

var Auth = &model.Resource{
	Name:        "auth",
	Description: "Auth resource",
	Table:       "auth",
	References: model.ResourceReferences{
		model.SingleQueryRef(),
		model.SinglePathRef(),
		model.MultiplePathRef("auth_id", "id"),
		model.MultipleQueryRef("auth_id", "id"),
		model.MultipleBodyRef("auth_id", "id"),
	},
}

var Resources = []*model.Resource{
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

