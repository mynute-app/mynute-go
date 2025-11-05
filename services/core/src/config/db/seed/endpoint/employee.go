package endpointSeed

import (
	"mynute-go/services/auth/config/db/model"
	resourceSeed "mynute-go/services/core/src/config/db/seed/resource"
	"mynute-go/services/core/src/config/namespace"
)

var CreateEmployee = &model.EndPoint{
	Path:             "/employee",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateEmployee",
	Description:      "Create employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Company,
}

var LoginEmployee = &model.EndPoint{
	Path:           "/employee/login",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginEmployeeByPassword",
	Description:    "Login employee",
	NeedsCompanyId: true,
}

var LoginEmployeeByEmailCode = &model.EndPoint{
	Path:           "/employee/login-with-code",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginEmployeeByEmailCode",
	Description:    "Login employee by email code",
	NeedsCompanyId: true,
}

var SendLoginCodeToEmployeeEmail = &model.EndPoint{
	Path:           "/employee/send-login-code/email/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendEmployeeLoginValidationCodeByEmail",
	Description:    "Send login code to employee email",
	NeedsCompanyId: true,
}

var ResetEmployeePasswordByEmail = &model.EndPoint{
	Path:           "/employee/reset-password/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "ResetEmployeePasswordByEmail",
	Description:    "Reset employee password by email",
	NeedsCompanyId: true,
}

var SendEmployeeVerificationCodeByEmail = &model.EndPoint{
	Path:           "/employee/send-verification-code/email/:email/:company_id",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendEmployeeVerificationEmail",
	Description:    "Send verification code to employee email",
	NeedsCompanyId: true,
}

var VerifyEmployeeEmail = &model.EndPoint{
	Path:           "/employee/verify-email/:email/:code/:company_id",
	Method:         namespace.ViewActionMethod,
	ControllerName: "VerifyEmployeeEmail",
	Description:    "Verify employee email",
	NeedsCompanyId: true,
}

var GetEmployeeById = &model.EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeById",
	Description:      "View employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var GetEmployeeByEmail = &model.EndPoint{
	Path:             "/employee/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeByEmail",
	Description:      "View employee by email",
	NeedsCompanyId:   true,
	DenyUnauthorized: false,
	Resource:         resourceSeed.Employee,
}

var UpdateEmployeeById = &model.EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateEmployeeById",
	Description:      "Update employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var UpdateEmployeeImages = &model.EndPoint{
	Path:             "/employee/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateEmployeeImages",
	Description:      "Update employee design images",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var DeleteEmployeeImage = &model.EndPoint{
	Path:             "/employee/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeImage",
	Description:      "Delete employee image",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var CreateEmployeeWorkSchedule = &model.EndPoint{
	Path:             "/employee/:id/work_schedule",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateEmployeeWorkSchedule",
	Description:      "Add work schedule to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var DeleteEmployeeWorkRange = &model.EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeWorkRange",
	Description:      "Remove work schedule from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var UpdateEmployeeWorkRange = &model.EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id",
	Method:           namespace.PutActionMethod,
	ControllerName:   "UpdateEmployeeWorkRange",
	Description:      "Update work schedule for employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var GetEmployeeWorkSchedule = &model.EndPoint{
	Path:           "/employee/:id/work_schedule",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetEmployeeWorkSchedule",
	Description:    "View work schedule for employee",
	NeedsCompanyId: true,
	Resource:       resourceSeed.Employee,
}

var GetEmployeeWorkRange = &model.EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeWorkRangeById",
	Description:      "View work range for employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var AddEmployeeWorkRangeServices = &model.EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id/services",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddEmployeeWorkRangeServices",
	Description:      "Add services to work range for a employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var DeleteEmployeeWorkRangeService = &model.EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeWorkRangeService",
	Description:      "Remove service from work range for a employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var DeleteEmployeeById = &model.EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeById",
	Description:      "Delete employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

var AddServiceToEmployee = &model.EndPoint{
	Path:             "/employee/:employee_id/service/:service_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddServiceToEmployee",
	Description:      "Add service to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var RemoveServiceFromEmployee = &model.EndPoint{
	Path:             "/employee/:employee_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveServiceFromEmployee",
	Description:      "Remove service from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var AddBranchToEmployee = &model.EndPoint{
	Path:             "/employee/:employee_id/branch/:branch_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddBranchToEmployee",
	Description:      "Add employee to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var RemoveBranchFromEmployee = &model.EndPoint{
	Path:             "/employee/:employee_id/branch/:branch_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveBranchFromEmployee",
	Description:      "Remove employee from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var AddRoleToEmployee = &model.EndPoint{
	Path:             "/employee/:employee_id/role/:role_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddRoleToEmployee",
	Description:      "Add role to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Role,
}

var RemoveRoleFromEmployee = &model.EndPoint{
	Path:             "/employee/:employee_id/role/:role_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveRoleFromEmployee",
	Description:      "Remove role from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Role,
}

var GetEmployeeAppointmentsById = &model.EndPoint{
	Path:             "/employee/:id/appointments",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeAppointmentsById",
	Description:      "View appointments for an employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Employee,
}

