package model

import (
	"fmt"
	"mynute-go/core/src/config/namespace"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllowEndpointCreation = false

type EndPoint struct {
	BaseModel
	ControllerName   string     `gorm:"type:varchar(100)" json:"controller_name"`
	Description      string     `gorm:"type:text" json:"description"`
	Method           string     `gorm:"type:varchar(6)" json:"method"`
	Path             string     `gorm:"type:text" json:"path"`
	DenyUnauthorized bool       `gorm:"default:false" json:"deny_unauthorized"`
	NeedsCompanyId   bool       `gorm:"default:false" json:"needs_company_id"`
	ResourceID       *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	Resource         *Resource  `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

func (EndPoint) TableName() string  { return "public.endpoints" }
func (EndPoint) SchemaType() string { return "public" }

func (EndPoint) Indexes() map[string]string {
	return map[string]string{
		"idx_method_path": "CREATE UNIQUE INDEX idx_method_path ON routes (method, path)",
	}
}

func (r *EndPoint) BeforeCreate(tx *gorm.DB) error {
	if !AllowEndpointCreation {
		panic("EndPoint creation is not allowed")
	}
	return nil
}

// --- Appointment Endpoints --- //

var CreateAppointment = &EndPoint{
	Path:             "/appointment",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateAppointment",
	Description:      "Create an appointment",
	NeedsCompanyId:   true,
	DenyUnauthorized: false,
	Resource:         BranchResource,
}
var GetAppointmentByID = &EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAppointmentByID",
	Description:      "View appointment by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         AppointmentResource,
}
var UpdateAppointmentByID = &EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateAppointmentByID",
	Description:      "Update appointment by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         AppointmentResource,
}
var CancelAppointmentByID = &EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "CancelAppointmentByID",
	Description:      "This wil cancel appointment by ID. Deleting appointments is forbidden.",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         AppointmentResource,
}

// --- Auth Endpoints --- //

var BeginAuthProviderCallback = &EndPoint{
	Path:           "/auth/oauth/:provider",
	Method:         namespace.ViewActionMethod,
	ControllerName: "BeginAuthProviderCallback",
	Description:    "Begin auth provider callback",
}
var GetAuthCallbackFunction = &EndPoint{
	Path:           "/auth/oauth/:provider/callback",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetAuthCallbackFunction",
	Description:    "View auth callback function",
}
var LogoutProvider = &EndPoint{
	Path:           "/auth/oauth/logout",
	Method:         namespace.ViewActionMethod,
	ControllerName: "LogoutProvider",
	Description:    "Logout provider",
}

// --- Admin Auth Endpoints --- //

var AdminLoginByPassword = &EndPoint{
	Path:           "/admin/auth/login",
	Method:         namespace.CreateActionMethod,
	ControllerName: "AdminLoginByPassword",
	Description:    "Admin login",
}
var AreThereAnyAdmin = &EndPoint{
	Path:             "/admin/are_there_any_superadmin",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "AreThereAnyAdmin",
	Description:      "Check if any superadmin exists",
	DenyUnauthorized: false, // Public endpoint for initial setup check
}
var CreateFirstAdmin = &EndPoint{
	Path:             "/admin/first_superadmin",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateFirstAdmin",
	Description:      "Create the first superadmin (only works if no superadmin exists)",
	DenyUnauthorized: false, // Public endpoint for initial setup
}
var GetAdminByID = &EndPoint{
	Path:             "/admin/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAdminByID",
	Description:      "Get admin by ID",
	DenyUnauthorized: true,
}
var GetAdminByEmail = &EndPoint{
	Path:             "/admin/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAdminByEmail",
	Description:      "Get admin by email",
	DenyUnauthorized: true,
}
var ListAdmins = &EndPoint{
	Path:             "/admin",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "ListAdmins",
	Description:      "List all admins",
	DenyUnauthorized: true,
}
var CreateAdmin = &EndPoint{
	Path:             "/admin",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateAdmin",
	Description:      "Create a new admin",
	DenyUnauthorized: false, // Bootstrap case: allows first admin creation without auth
}
var UpdateAdminByID = &EndPoint{
	Path:             "/admin/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateAdminByID",
	Description:      "Update admin by ID",
	DenyUnauthorized: true,
}
var DeleteAdminByID = &EndPoint{
	Path:             "/admin/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteAdminByID",
	Description:      "Delete admin by ID",
	DenyUnauthorized: true,
}
var ListAdminRoles = &EndPoint{
	Path:             "/admin/role",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "ListAdminRoles",
	Description:      "List all admin roles",
	DenyUnauthorized: true,
}
var CreateAdminRole = &EndPoint{
	Path:             "/admin/role",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateAdminRole",
	Description:      "Create a new admin role",
	DenyUnauthorized: true,
}
var GetAdminRoleByID = &EndPoint{
	Path:             "/admin/role/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAdminRoleByID",
	Description:      "Get admin role by ID",
	DenyUnauthorized: true,
}
var UpdateAdminRoleByID = &EndPoint{
	Path:             "/admin/role/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateAdminRoleByID",
	Description:      "Update role by ID",
	DenyUnauthorized: true,
}
var DeleteAdminRoleByID = &EndPoint{
	Path:             "/admin/role/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteAdminRoleByID",
	Description:      "Delete role by ID",
	DenyUnauthorized: true,
}

// --- Branch Endpoints --- //

var CreateBranch = &EndPoint{
	Path:             "/branch",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateBranch",
	Description:      "Create a branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchById",
	Description:      "View branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}

//	var GetBranchByName = &EndPoint{
//		Path:             "/branch/name/:name",
//		Method:           namespace.ViewActionMethod,
//		ControllerName:   "GetBranchByName",
//		Description:      "View branch by name",
//		NeedsCompanyId:   true,
//		DenyUnauthorized: true,
//		Resource:         BranchResource,
//	}
var UpdateBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateBranchById",
	Description:      "Update branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var DeleteBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchById",
	Description:      "Delete branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var GetEmployeeServicesByBranchId = &EndPoint{
	Path:             "/branch/:branch_id/employee/:employee_id/services",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeServicesByBranchId",
	Description:      "View employee offered services at the branch by branch ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var AddServiceToBranch = &EndPoint{
	Path:             "/branch/:branch_id/service/:service_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddServiceToBranch",
	Description:      "Add service to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var RemoveServiceFromBranch = &EndPoint{
	Path:             "/branch/:branch_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveServiceFromBranch",
	Description:      "Remove service from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var CreateBranchWorkSchedule = &EndPoint{
	Path:             "/branch/:id/work_schedule",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateBranchWorkSchedule",
	Description:      "Add work schedule to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var UpdateBranchImages = &EndPoint{
	Path:             "/branch/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateBranchImages",
	Description:      "Update branch images",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var DeleteBranchImage = &EndPoint{
	Path:             "/branch/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchImage",
	Description:      "Delete branch image",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var GetBranchWorkSchedule = &EndPoint{
	Path:           "/branch/:id/work_schedule",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetBranchWorkSchedule",
	Description:    "View work schedule for branch",
	NeedsCompanyId: true,
	Resource:       BranchResource,
}
var GetBranchWorkRange = &EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchWorkRange",
	Description:      "View work range by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var UpdateBranchWorkRange = &EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id",
	Method:           namespace.PutActionMethod,
	ControllerName:   "UpdateBranchWorkRange",
	Description:      "Update work range in branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var DeleteBranchWorkRange = &EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchWorkRange",
	Description:      "Remove work range from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var AddBranchWorkRangeServices = &EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id/services",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddBranchWorkRangeServices",
	Description:      "Add services to work range in branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var DeleteBranchWorkRangeService = &EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchWorkRangeService",
	Description:      "Remove service from work range in branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var GetBranchAppointmentsById = &EndPoint{
	Path:             "/branch/:id/appointments",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchAppointmentsById",
	Description:      "View appointments for a branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}

// --- Client Endpoints --- //

var CreateClient = &EndPoint{
	Path:           "/client",
	Method:         namespace.CreateActionMethod,
	ControllerName: "CreateClient",
	Description:    "Create client",
}
var LoginClient = &EndPoint{
	Path:           "/client/login",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginClientByPassword",
	Description:    "Login client",
}
var LoginClientByEmailCode = &EndPoint{
	Path:           "/client/login-with-code",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginClientByEmailCode",
	Description:    "Login client by email code",
}
var SendLoginCodeToClientEmail = &EndPoint{
	Path:           "/client/send-login-code/email/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendClientLoginValidationCodeByEmail",
	Description:    "Send login code to client email",
}
var ResetClientPasswordByEmail = &EndPoint{
	Path:           "/client/reset-password/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "ResetClientPasswordByEmail",
	Description:    "Reset client password by email",
}
var SendClientVerificationCodeByEmail = &EndPoint{
	Path:           "/client/send-verification-code/email/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendClientVerificationCodeByEmail",
	Description:    "Send verification code to client email",
}
var VerifyClientEmail = &EndPoint{
	Path:           "/client/verify-email/:email/:code",
	Method:         namespace.ViewActionMethod,
	ControllerName: "VerifyClientEmail",
	Description:    "Verify client email code",
}
var GetClientByEmail = &EndPoint{
	Path:             "/client/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetClientByEmail",
	Description:      "View client by email",
	DenyUnauthorized: false,
	Resource:         ClientResource,
}
var GetClientById = &EndPoint{
	Path:             "/client/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetClientById",
	Description:      "View client by ID",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var UpdateClientById = &EndPoint{
	Path:             "/client/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateClientById",
	Description:      "Update client by ID",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var DeleteClientById = &EndPoint{
	Path:             "/client/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteClientById",
	Description:      "Delete client by ID",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var UpdateClientImages = &EndPoint{
	Path:             "/client/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateClientImages",
	Description:      "Update client design images",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var DeleteClientImage = &EndPoint{
	Path:             "/client/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteClientImage",
	Description:      "Delete client design images",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}

// --- Company Endpoints --- //

var CreateCompany = &EndPoint{
	Path:           "/company",
	Method:         namespace.CreateActionMethod,
	ControllerName: "CreateCompany",
	Description:    "Create a company",
}
var GetCompanyById = &EndPoint{
	Path:             "/company/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetCompanyById",
	Description:      "View company by ID",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         CompanyResource,
}
var GetCompanyByName = &EndPoint{
	Path:           "/company/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyByName",
	Description:    "View company by name",
}
var CheckIfCompanyExistsByTaxID = &EndPoint{
	Path:           "/company/tax_id/:tax_id/exists",
	Method:         namespace.ViewActionMethod,
	ControllerName: "CheckIfCompanyExistsByTaxID",
	Description:    "Check if company exists by tax ID",
}
var GetCompanyByTaxId = &EndPoint{
	Path:             "/company/tax_id/:tax_id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetCompanyByTaxId",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Description:      "View company by tax ID",
}
var GetCompanyBySubdomain = &EndPoint{
	Path:           "/company/subdomain/:subdomain_name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyBySubdomain",
	Description:    "View company by subdomain",
}
var UpdateCompanyById = &EndPoint{
	Path:             "/company/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateCompanyById",
	Description:      "Update company by ID",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         CompanyResource,
}
var UpdateCompanyImages = &EndPoint{
	Path:             "/company/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateCompanyImages",
	Description:      "Update company design images",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         CompanyResource,
}
var DeleteCompanyImage = &EndPoint{
	Path:             "/company/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteCompanyImage",
	Description:      "Delete company design images",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         CompanyResource,
}
var UpdateCompanyColors = &EndPoint{
	Path:             "/company/:id/design/colors",
	Method:           namespace.PutActionMethod,
	ControllerName:   "UpdateCompanyColors",
	Description:      "Update company design colors",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         CompanyResource,
}
var DeleteCompanyById = &EndPoint{
	Path:             "/company/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteCompanyById",
	Description:      "Delete company by ID",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         CompanyResource,
}

// --- Employee Endpoints --- //

var CreateEmployee = &EndPoint{
	Path:             "/employee",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateEmployee",
	Description:      "Create employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var LoginEmployee = &EndPoint{
	Path:           "/employee/login",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginEmployeeByPassword",
	Description:    "Login employee",
	NeedsCompanyId: true,
}
var LoginEmployeeByEmailCode = &EndPoint{
	Path:           "/employee/login-with-code",
	Method:         namespace.CreateActionMethod,
	ControllerName: "LoginEmployeeByEmailCode",
	Description:    "Login employee by email code",
	NeedsCompanyId: true,
}
var SendLoginCodeToEmployeeEmail = &EndPoint{
	Path:           "/employee/send-login-code/email/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendEmployeeLoginValidationCodeByEmail",
	Description:    "Send login code to employee email",
	NeedsCompanyId: true,
}
var ResetEmployeePasswordByEmail = &EndPoint{
	Path:           "/employee/reset-password/:email",
	Method:         namespace.CreateActionMethod,
	ControllerName: "ResetEmployeePasswordByEmail",
	Description:    "Reset employee password by email",
	NeedsCompanyId: true,
}
var SendEmployeeVerificationCodeByEmail = &EndPoint{
	Path:           "/employee/send-verification-code/email/:email/:company_id",
	Method:         namespace.CreateActionMethod,
	ControllerName: "SendEmployeeVerificationEmail",
	Description:    "Send verification code to employee email",
	NeedsCompanyId: true,
}
var VerifyEmployeeEmail = &EndPoint{
	Path:           "/employee/verify-email/:email/:code/:company_id",
	Method:         namespace.ViewActionMethod,
	ControllerName: "VerifyEmployeeEmail",
	Description:    "Verify employee email",
	NeedsCompanyId: true,
}
var GetEmployeeById = &EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeById",
	Description:      "View employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var GetEmployeeByEmail = &EndPoint{
	Path:             "/employee/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeByEmail",
	Description:      "View employee by email",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var UpdateEmployeeById = &EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateEmployeeById",
	Description:      "Update employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var UpdateEmployeeImages = &EndPoint{
	Path:             "/employee/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateEmployeeImages",
	Description:      "Update employee design images",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var DeleteEmployeeImage = &EndPoint{
	Path:             "/employee/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeImage",
	Description:      "Delete employee image",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var CreateEmployeeWorkSchedule = &EndPoint{
	Path:             "/employee/:id/work_schedule",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateEmployeeWorkSchedule",
	Description:      "Add work schedule to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var DeleteEmployeeWorkRange = &EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeWorkRange",
	Description:      "Remove work schedule from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var UpdateEmployeeWorkRange = &EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id",
	Method:           namespace.PutActionMethod,
	ControllerName:   "UpdateEmployeeWorkRange",
	Description:      "Update work schedule for employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var GetEmployeeWorkSchedule = &EndPoint{
	Path:           "/employee/:id/work_schedule",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetEmployeeWorkSchedule",
	Description:    "View work schedule for employee",
	NeedsCompanyId: true,
	Resource:       EmployeeResource,
}
var GetEmployeeWorkRange = &EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeWorkRangeById",
	Description:      "View work range for employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var AddEmployeeWorkRangeServices = &EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id/services",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddEmployeeWorkRangeServices",
	Description:      "Add services to work range for a employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var DeleteEmployeeWorkRangeService = &EndPoint{
	Path:             "/employee/:id/work_range/:work_range_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeWorkRangeService",
	Description:      "Remove service from work range for a employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var DeleteEmployeeById = &EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteEmployeeById",
	Description:      "Delete employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var AddServiceToEmployee = &EndPoint{
	Path:             "/employee/:employee_id/service/:service_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddServiceToEmployee",
	Description:      "Add service to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var RemoveServiceFromEmployee = &EndPoint{
	Path:             "/employee/:employee_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveServiceFromEmployee",
	Description:      "Remove service from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var AddBranchToEmployee = &EndPoint{
	Path:             "/employee/:employee_id/branch/:branch_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddBranchToEmployee",
	Description:      "Add employee to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var RemoveBranchFromEmployee = &EndPoint{
	Path:             "/employee/:employee_id/branch/:branch_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveBranchFromEmployee",
	Description:      "Remove employee from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var AddRoleToEmployee = &EndPoint{
	Path:             "/employee/:employee_id/role/:role_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddRoleToEmployee",
	Description:      "Add role to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         RoleResource,
}
var RemoveRoleFromEmployee = &EndPoint{
	Path:             "/employee/:employee_id/role/:role_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveRoleFromEmployee",
	Description:      "Remove role from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         RoleResource,
}
var GetEmployeeAppointmentsById = &EndPoint{
	Path:             "/employee/:id/appointments",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeAppointmentsById",
	Description:      "View appointments for an employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}

// --- Holiday Endpoints --- //

var CreateHoliday = &EndPoint{
	Path:             "/holiday",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateHoliday",
	Description:      "Create a holiday",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetHolidayById = &EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetHolidayById",
	Description:      "View holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         HolidayResource,
}
var GetHolidayByName = &EndPoint{
	Path:           "/holiday/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetHolidayByName",
	Description:    "View holiday by name",
	NeedsCompanyId: true,
}
var UpdateHolidayById = &EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateHolidayById",
	Description:      "Update holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         HolidayResource,
}
var DeleteHolidayById = &EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteHolidayById",
	Description:      "Delete holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         HolidayResource,
}

// --- Sector Endpoints --- //

var CreateSector = &EndPoint{
	Path:             "/sector",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateSector",
	Description:      "Creates a company sector",
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetSectorById = &EndPoint{
	Path:           "/sector/:id",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetSectorById",
	Description:    "Retrieves a company sector by ID",
}
var GetSectorByName = &EndPoint{
	Path:           "/sector/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetSectorByName",
	Description:    "Retrieves a company sector by name",
}
var UpdateSectorById = &EndPoint{
	Path:             "/sector/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateSectorById",
	Description:      "Updates a company sector by ID",
	DenyUnauthorized: true,
	Resource:         SectorResource,
}
var DeleteSectorById = &EndPoint{
	Path:             "/sector/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteSectorById",
	Description:      "Deletes a company sector by ID",
	DenyUnauthorized: true,
	Resource:         SectorResource,
}

// --- Service Endpoints --- //

var CreateService = &EndPoint{
	Path:             "/service",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateService",
	Description:      "Create a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetServiceById = &EndPoint{
	Path:             "/service/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetServiceById",
	Description:      "View service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var GetServiceByName = &EndPoint{
	Path:           "/service/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetServiceByName",
	Description:    "View service by name",
	NeedsCompanyId: true,
}
var UpdateServiceById = &EndPoint{
	Path:             "/service/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateServiceById",
	Description:      "Update service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var DeleteServiceById = &EndPoint{
	Path:             "/service/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteServiceById",
	Description:      "Delete service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var UpdateServiceImages = &EndPoint{
	Path:             "/service/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateServiceImages",
	Description:      "Update images of a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var DeleteServiceImage = &EndPoint{
	Path:             "/service/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteServiceImage",
	Description:      "Delete an image of a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var GetServiceAvailability = &EndPoint{
	Path:           "/service/:id/availability",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetServiceAvailability",
	Description:    "Get availability of a service",
	NeedsCompanyId: true,
	Resource:       ServiceResource,
}

// --- Combine all Endpoints into a slice for seeding --- //
var endpoints = []*EndPoint{
	// Appointment
	CreateAppointment,
	GetAppointmentByID,
	UpdateAppointmentByID,
	CancelAppointmentByID,
	// Auth
	BeginAuthProviderCallback,
	GetAuthCallbackFunction,
	LogoutProvider,
	// Admin Auth
	AdminLoginByPassword,
	AreThereAnyAdmin,
	CreateFirstAdmin,
	// Admin Management
	GetAdminByID,
	GetAdminByEmail,
	ListAdmins,
	CreateAdmin,
	UpdateAdminByID,
	DeleteAdminByID,
	ListAdminRoles,
	CreateAdminRole,
	GetAdminRoleByID,
	UpdateAdminRoleByID,
	DeleteAdminRoleByID,
	// Branch
	CreateBranch,
	GetBranchById,
	UpdateBranchById,
	DeleteBranchById,
	GetEmployeeServicesByBranchId,
	AddServiceToBranch,
	RemoveServiceFromBranch,
	UpdateBranchImages,
	DeleteBranchImage,
	CreateBranchWorkSchedule,
	GetBranchWorkSchedule,
	GetBranchWorkRange,
	DeleteBranchWorkRange,
	UpdateBranchWorkRange,
	AddBranchWorkRangeServices,
	DeleteBranchWorkRangeService,
	GetBranchAppointmentsById,
	// Client
	CreateClient,
	LoginClient,
	LoginClientByEmailCode,
	SendLoginCodeToClientEmail,
	SendClientVerificationCodeByEmail,
	VerifyClientEmail,
	ResetClientPasswordByEmail,
	GetClientByEmail,
	GetClientById,
	UpdateClientById,
	DeleteClientById,
	UpdateClientImages,
	DeleteClientImage,
	// Company
	CreateCompany,
	GetCompanyById,
	GetCompanyByName,
	GetCompanyByTaxId,
	CheckIfCompanyExistsByTaxID,
	GetCompanyBySubdomain,
	UpdateCompanyById,
	DeleteCompanyById,
	UpdateCompanyImages,
	DeleteCompanyImage,
	UpdateCompanyColors,
	// Employee
	CreateEmployee,
	LoginEmployee,
	LoginEmployeeByEmailCode,
	SendLoginCodeToEmployeeEmail,
	SendEmployeeVerificationCodeByEmail,
	VerifyEmployeeEmail,
	ResetEmployeePasswordByEmail,
	GetEmployeeById,
	GetEmployeeByEmail,
	UpdateEmployeeById,
	DeleteEmployeeById,
	AddServiceToEmployee,
	RemoveServiceFromEmployee,
	AddBranchToEmployee,
	RemoveBranchFromEmployee,
	UpdateEmployeeImages,
	DeleteEmployeeImage,
	CreateEmployeeWorkSchedule,
	GetEmployeeWorkSchedule,
	GetEmployeeWorkRange,
	DeleteEmployeeWorkRange,
	UpdateEmployeeWorkRange,
	AddEmployeeWorkRangeServices,
	DeleteEmployeeWorkRangeService,
	GetEmployeeAppointmentsById,
	// Holiday
	CreateHoliday,
	GetHolidayById,
	GetHolidayByName,
	UpdateHolidayById,
	DeleteHolidayById,
	// Sector
	CreateSector,
	GetSectorById,
	GetSectorByName,
	UpdateSectorById,
	DeleteSectorById,
	// Service
	CreateService,
	GetServiceById,
	GetServiceByName,
	UpdateServiceById,
	DeleteServiceById,
	UpdateServiceImages,
	DeleteServiceImage,
	GetServiceAvailability,
}

type EndpointCfg struct {
	AllowCreation bool // Allow creation of endpoints
}

func EndPoints(cfg *EndpointCfg, db *gorm.DB) ([]*EndPoint, func(), error) {
	AllowEndpointCreation = cfg.AllowCreation

	// Recuperar os recursos corretos do banco
	resourceMap := map[string]uuid.UUID{}
	var resources []Resource
	if err := db.Find(&resources).Error; err != nil {
		return nil, nil, err
	}
	for _, r := range resources {
		resourceMap[r.Table] = r.ID
	}

	for _, edp := range endpoints {
		if edp.Resource != nil {
			if id, ok := resourceMap[edp.Resource.Table]; ok {
				edp.ResourceID = &id
			} else {
				return nil, nil, fmt.Errorf("resource not found for table: %s", edp.Resource.Table)
			}
			edp.Resource = nil
		}
	}

	deferFnc := func() {
		AllowEndpointCreation = false
	}

	return endpoints, deferFnc, nil
}

// LoadEndpointIDs loads the IDs of the endpoints from the database
// and updates the endpoint variables with their corresponding IDs.
// This should be called after seeding endpoints to ensure that
// policies can reference the correct endpoint IDs.
func LoadEndpointIDs(db *gorm.DB) error {
	for _, ep := range endpoints {
		var existing EndPoint
		if err := db.
			Where("method = ? AND path = ?", ep.Method, ep.Path).
			First(&existing).Error; err != nil {
			return fmt.Errorf("failed to load endpoint ID for %s %s: %w", ep.Method, ep.Path, err)
		}
		ep.ID = existing.ID
	}
	return nil
}
