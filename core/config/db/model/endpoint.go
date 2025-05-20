package model

import (
	"agenda-kaki-go/core/config/namespace"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllowEndpointCreation = false

type EndPoint struct {
	BaseModel
	ControllerName   string     `json:"controller_name"`
	Description      string     `json:"description"`
	Method           string     `gorm:"type:varchar(10)" json:"method"`
	Path             string     `json:"path"`
	DenyUnauthorized bool       `gorm:"default:false" json:"deny_unauthorized"`
	NeedsCompanyId   bool       `gorm:"default:false" json:"needs_company_id"`
	ResourceID       *uuid.UUID `json:"resource_id"`
	Resource         *Resource  `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

func (EndPoint) TableName() string { return "public.endpoints" }
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
	DenyUnauthorized: true,
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
	Method:           namespace.UpdateActionMethod,
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
var GetBranchByName = &EndPoint{
	Path:             "/branch/name/:name",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchByName",
	Description:      "View branch by name",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var UpdateBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.UpdateActionMethod,
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
	ControllerName: "LoginClient",
	Description:    "Login client",
}
var VerifyClientEmail = &EndPoint{
	Path:           "/client/verify-email/:email/:code",
	Method:         namespace.CreateActionMethod,
	ControllerName: "VerifyClientEmail",
	Description:    "Verify client email",
}
var GetClientByEmail = &EndPoint{
	Path:             "/client/email/:email",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetClientByEmail",
	Description:      "View client by email",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var UpdateClientById = &EndPoint{
	Path:             "/client/:id",
	Method:           namespace.UpdateActionMethod,
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

// --- Company Endpoints --- //

var CreateCompany = &EndPoint{
	Path:           "/company",
	Method:         namespace.CreateActionMethod,
	ControllerName: "CreateCompany",
	Description:    "Create a company",
}
var GetCompanyById = &EndPoint{
	Path:           "/company/:id",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyById",
	Description:    "View company by ID",
	Resource:       CompanyResource,
}
var GetCompanyByName = &EndPoint{
	Path:           "/company/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyByName",
	Description:    "View company by name",
}
var GetCompanyByTaxId = &EndPoint{
	Path:           "/company/tax_id/:tax_id",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyByTaxId",
	Description:    "View company by tax ID",
}
var GetCompanyIdBySubdomain = &EndPoint{
	Path:           "/company/subdomain/:subdomain_name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyIdBySubdomain",
	Description:    "View company by subdomain",
}
var UpdateCompanyById = &EndPoint{
	Path:             "/company/:id",
	Method:           namespace.UpdateActionMethod,
	ControllerName:   "UpdateCompanyById",
	Description:      "Update company by ID",
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
	ControllerName: "LoginEmployee",
	Description:    "Login employee",
	NeedsCompanyId: true,
}
var VerifyEmployeeEmail = &EndPoint{
	Path:           "/employee/verify-email/:email/:code",
	Method:         namespace.CreateActionMethod,
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
	Method:           namespace.UpdateActionMethod,
	ControllerName:   "UpdateEmployeeById",
	Description:      "Update employee by ID",
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
	Method:           namespace.UpdateActionMethod,
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
	Method:           namespace.UpdateActionMethod,
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
	Method:           namespace.UpdateActionMethod,
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
	// Branch
	CreateBranch,
	GetBranchById,
	GetBranchByName,
	UpdateBranchById,
	DeleteBranchById,
	GetEmployeeServicesByBranchId,
	AddServiceToBranch,
	RemoveServiceFromBranch,
	// Client
	CreateClient,
	LoginClient,
	VerifyClientEmail,
	GetClientByEmail,
	UpdateClientById,
	DeleteClientById,
	// Company
	CreateCompany,
	GetCompanyById,
	GetCompanyByName,
	GetCompanyByTaxId,
	GetCompanyIdBySubdomain,
	UpdateCompanyById,
	DeleteCompanyById,
	// Employee
	CreateEmployee,
	LoginEmployee,
	VerifyEmployeeEmail,
	GetEmployeeById,
	GetEmployeeByEmail,
	UpdateEmployeeById,
	DeleteEmployeeById,
	AddServiceToEmployee,
	RemoveServiceFromEmployee,
	AddBranchToEmployee,
	RemoveBranchFromEmployee,
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
}

type EndpointCfg struct {
	AllowCreation bool // Allow creation of endpoints
}

func EndPoints(cfg *EndpointCfg) ([]*EndPoint, func()) {
	AllowEndpointCreation = cfg.AllowCreation
	for _, edp := range endpoints {
		if edp.Resource != nil {
			edp.ResourceID = &edp.Resource.ID
			edp.Resource = nil // Avoid circular reference
		}
	}
	deferFnc := func() {
		AllowEndpointCreation = false
	}
	return endpoints, deferFnc
}

// func SeedEndpoints(db *gorm.DB) ([]*EndPoint, error) {
// 	AllowEndpointCreation = true
// 	tx := db.Begin()
// 	defer func() {
// 		AllowEndpointCreation = false
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			log.Printf("Panic occurred during policy seeding: %v", r)
// 		}
// 		if err := tx.Commit().Error; err != nil {
// 			log.Printf("Failed to commit transaction: %v", err)
// 		}
// 		log.Print("System Endpoints seeded successfully")
// 	}()
// 	LoadEndpoints()
// 	for _, edp := range Endpoints {
// 		err := tx.Where("method = ? AND path = ?", edp.Method, edp.Path).First(edp).Error
// 		if err == gorm.ErrRecordNotFound {
// 			if err := tx.Create(edp).Error; err != nil {
// 				return nil, err
// 			}
// 		} else if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return Endpoints, nil
// }
