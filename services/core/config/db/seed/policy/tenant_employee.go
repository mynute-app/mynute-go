package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- TENANT POLICIES: Employee-related operations ---

var TenantAllowCreateEmployee = &model.TenantPolicy{
	Name:        "TP: CanCreateEmployee",
	Description: "Allows company Owner, GM, or BM to create employees.",
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

var TenantAllowGetEmployeeById = &model.TenantPolicy{
	Name:        "TP: CanViewEmployeeById",
	Description: "Allows employee to view self, or any internal user of the same company to view other employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeById.ID,
	Conditions:  model.JsonRawMessage(EmployeeSelfOrInternalUserCheck),
}

var TenantAllowGetEmployeeByEmail = &model.TenantPolicy{
	Name:        "TP: CanViewEmployeeByEmail",
	Description: "Allows company members to find employees within the same company by email.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeByEmail.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck),
}

var TenantAllowUpdateEmployeeById = &model.TenantPolicy{
	Name:        "TP: CanUpdateEmployee",
	Description: "Allows employee to update self, or company managers to update employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateEmployeeById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowDeleteEmployeeById = &model.TenantPolicy{
	Name:        "TP: CanDeleteEmployee",
	Description: "Allows company managers (Owner, GM, BM) to delete employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}

var TenantAllowCreateEmployeeWorkSchedule = &model.TenantPolicy{
	Name:        "TP: CanCreateEmployeeWorkSchedule",
	Description: "Allows employees or company managers to create work schedules.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateEmployeeWorkSchedule.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowGetEmployeeWorkRangeById = &model.TenantPolicy{
	Name:        "TP: CanViewEmployeeWorkRangeById",
	Description: "Allows employees or company managers to view work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowUpdateEmployeeWorkRange = &model.TenantPolicy{
	Name:        "TP: CanUpdateEmployeeWorkRange",
	Description: "Allows employees or company managers to update work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateEmployeeWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowDeleteEmployeeWorkRange = &model.TenantPolicy{
	Name:        "TP: CanDeleteEmployeeWorkRange",
	Description: "Allows employees or company managers to remove work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowAddEmployeeWorkRangeServices = &model.TenantPolicy{
	Name:        "TP: CanAddEmployeeWorkRangeServices",
	Description: "Allows employees or company managers to add services to work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddEmployeeWorkRangeServices.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowDeleteEmployeeWorkRangeService = &model.TenantPolicy{
	Name:        "TP: CanDeleteEmployeeWorkRangeService",
	Description: "Allows employees or company managers to remove services from work ranges.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeWorkRangeService.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowAddServiceToEmployee = &model.TenantPolicy{
	Name:        "TP: CanAddServiceToEmployee",
	Description: "Allows company managers to assign services to employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddServiceToEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowRemoveServiceFromEmployee = &model.TenantPolicy{
	Name:        "TP: CanRemoveServiceFromEmployee",
	Description: "Allows company managers to remove services from employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.RemoveServiceFromEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowAddBranchToEmployee = &model.TenantPolicy{
	Name:        "TP: CanAddBranchToEmployee",
	Description: "Allows company managers to assign employees to branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddBranchToEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowRemoveBranchFromEmployee = &model.TenantPolicy{
	Name:        "TP: CanRemoveBranchFromEmployee",
	Description: "Allows company managers to remove employees from branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.RemoveBranchFromEmployee.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowUpdateEmployeeImages = &model.TenantPolicy{
	Name:        "TP: CanUpdateEmployeeImages",
	Description: "Allows company managers or employee himself to update employee images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateEmployeeImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowDeleteEmployeeImage = &model.TenantPolicy{
	Name:        "TP: CanDeleteEmployeeImage",
	Description: "Allows company managers or employee himself to delete employee images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteEmployeeImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}

var TenantAllowGetEmployeeAppointmentsById = &model.TenantPolicy{
	Name:        "TP: CanViewEmployeeAppointmentsById",
	Description: "Allows employees or company managers to view employee appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeAppointmentsById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrEmployeeHimselfCheck),
}
