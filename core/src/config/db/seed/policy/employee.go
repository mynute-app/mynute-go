package policySeed

import (
	"mynute-go/auth/model"
	endpointSeed "mynute-go/core/src/config/db/seed/endpoint"
)

var AllowCreateEmployee = &model.PolicyRule{
	Name:        "SDP: CanCreateEmployee",
	Description: "Allows company Owner, GM, or BM to create employees (BM restricted to their branches implicitly if data includes branch).",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateEmployee.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Admin or Branch Manager Creation Access",
		LogicType:   "OR",
		Children: []model.ConditionNode{
			CompanyAdminCheck,
			CompanyBranchManagerAssignedBranchCheck,
		},
	}),
}

var AllowGetEmployeeById = &model.PolicyRule{
	Name:        "SDP: CanViewEmployeeById",
	Description: "Allows employee to view self, or any internal user of the same company to view other employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeById.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Allow Employee Self-View OR Any Internal Company User View",
		LogicType:   "OR",
		Children: []model.ConditionNode{
			EmployeeSelfAccessCheck,  // Can view self (checks subject.id == resource.id)
			CompanyInternalUserCheck, // Any other member of the same company can view (checks subject.company_id == resource.company_id)
		},
	}),
}

var AllowGetEmployeeByEmail = &model.PolicyRule{
	Name:        "SDP: CanViewEmployeeByEmail",
	Description: "Allows company members to find employees within the same company by email.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeByEmail.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck), // Subject must be internal user of the found employee's company
}

var AllowUpdateEmployeeById = &model.PolicyRule{
	Name:        "SDP: CanUpdateEmployee",
	Description: "Allows employee to update self, or company managers (Owner, GM, BM) to update employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateEmployeeById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var AllowCreateEmployeeWorkSchedule = &model.PolicyRule{
	Name:        "SDP: CanCreateEmployeeWorkSchedule",
	Description: "Allows employees, or company managers (Owner, GM, BM), to create their own work schedules.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateEmployeeWorkSchedule.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Employee can create own schedule
}

var AllowGetEmployeeWorkRangeById = &model.PolicyRule{
	Name:        "SDP: CanViewEmployeeWorkRangeById",
	Description: "Allows employees, or company managers (Owner, GM, BM), to view their own work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var AllowUpdateEmployeeWorkRange = &model.PolicyRule{
	Name:        "SDP: CanUpdateEmployeeWorkRange",
	Description: "Allows employees, or company managers (Owner, GM, BM), to update their own work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateEmployeeWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Employee can update own work range
}

var AllowDeleteEmployeeWorkRange = &model.PolicyRule{
	Name:        "SDP: CanDeleteEmployeeWorkRange",
	Description: "Allows employees, or company managers (Owner, GM, BM), to remove their own work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Employee can remove own work range
}

var AllowAddEmployeeWorkRangeServices = &model.PolicyRule{
	Name:        "SDP: CanAddEmployeeWorkRangeServices",
	Description: "Allows employees, or company managers (Owner, GM, BM), to add services to their work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddEmployeeWorkRangeServices.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Employee can add services to own work range
}

var AllowDeleteEmployeeWorkRangeService = &model.PolicyRule{
	Name:        "SDP: CanDeleteEmployeeWorkRangeService",
	Description: "Allows employees, or company managers (Owner, GM, BM), to remove services from their work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeWorkRangeService.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Employee can remove services from own work range
}

var AllowDeleteEmployeeById = &model.PolicyRule{
	Name:        "SDP: CanDeleteEmployee",
	Description: "Allows company managers (Owner, GM, BM) to delete employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Owner, GM, BM can delete
}

var AllowAddServiceToEmployee = &model.PolicyRule{
	Name:        "SDP: CanAddServiceToEmployee",
	Description: "Allows company managers (Owner, GM, BM) to assign services to employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddServiceToEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Manager of the employee's company
}

var AllowRemoveServiceFromEmployee = &model.PolicyRule{
	Name:        "SDP: CanRemoveServiceFromEmployee",
	Description: "Allows company managers (Owner, GM, BM) to remove services from employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.RemoveServiceFromEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Manager of the employee's company
}

var AllowAddBranchToEmployee = &model.PolicyRule{
	Name:        "SDP: CanAddBranchToEmployee",
	Description: "Allows company managers (Owner, GM, BM) to assign employees to branches (respecting BM scope).",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddBranchToEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowRemoveBranchFromEmployee = &model.PolicyRule{
	Name:        "SDP: CanRemoveBranchFromEmployee",
	Description: "Allows company managers (Owner, GM, BM) to remove employees from branches (respecting BM scope).",
	Effect:      "Allow",
	EndPointID:  endpointSeed.RemoveBranchFromEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowUpdateEmployeeImages = &model.PolicyRule{
	Name:        "SDP: CanUpdateEmployeeImages",
	Description: "Allows company Owner, General Manager, or employee himself to update employee images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateEmployeeImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var AllowDeleteEmployeeImage = &model.PolicyRule{
	Name:        "SDP: CanDeleteEmployeeImage",
	Description: "Allows company Owner, General Manager, or employee himself to delete employee images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck), // Employee can delete own images
}

var AllowGetEmployeeAppointmentsById = &model.PolicyRule{
	Name:        "SDP: CanViewEmployeeAppointmentsById",
	Description: "Allows employees to view their own appointments, or company managers (Owner, GM, BM) to view employee appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeAppointmentsById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}
