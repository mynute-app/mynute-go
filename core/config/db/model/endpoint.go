package model

import (
	"agenda-kaki-go/core/config/namespace"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllowEndpointCreation = false

type EndPoint struct {
	BaseModel
	Handler          string     `json:"handler"`
	Description      string     `json:"description"`
	Method           string     `gorm:"type:varchar(10)" json:"method"`
	Path             string     `json:"path"`
	DenyUnauthorized bool       `gorm:"default:false" json:"deny_unauthorized"`
	NeedsCompanyId   bool       `gorm:"default:false" json:"needs_company_id"`
	ResourceID       *uuid.UUID `json:"resource_id"`
	Resource         *Resource  `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

func (EndPoint) TableName() string {
	return "public.endpoints"
}

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
	Handler:          "CreateAppointment",
	Description:      "Create an appointment",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var GetAppointmentByID = &EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetAppointmentByID",
	Description:      "View appointment by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         AppointmentResource,
}
var UpdateAppointmentByID = &EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateAppointmentByID",
	Description:      "Update appointment by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         AppointmentResource,
}
var CancelAppointmentByID = &EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "CancelAppointmentByID",
	Description:      "This wil cancel appointment by ID. Deleting appointments is forbidden.",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         AppointmentResource,
}

// --- Auth Endpoints --- //

var VerifyExistingAccount = &EndPoint{
	Path:        "/auth/verify-existing-account",
	Method:      namespace.CreateActionMethod,
	Handler:     "VerifyExistingAccount",
	Description: "Verify if an account exists",
}
var BeginAuthProviderCallback = &EndPoint{
	Path:        "/auth/oauth/:provider",
	Method:      namespace.ViewActionMethod,
	Handler:     "BeginAuthProviderCallback",
	Description: "Begin auth provider callback",
}
var GetAuthCallbackFunction = &EndPoint{
	Path:        "/auth/oauth/:provider/callback",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetAuthCallbackFunction",
	Description: "View auth callback function",
}
var LogoutProvider = &EndPoint{
	Path:        "/auth/oauth/logout",
	Method:      namespace.ViewActionMethod,
	Handler:     "LogoutProvider",
	Description: "Logout provider",
}

// --- Branch Endpoints --- //

var CreateBranch = &EndPoint{
	Path:             "/branch",
	Method:           namespace.CreateActionMethod,
	Handler:          "CreateBranch",
	Description:      "Create a branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetBranchById",
	Description:      "View branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var GetBranchByName = &EndPoint{
	Path:             "/branch/name/:name",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetBranchByName",
	Description:      "View branch by name",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var UpdateBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateBranchById",
	Description:      "Update branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var DeleteBranchById = &EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteBranchById",
	Description:      "Delete branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var GetEmployeeServicesByBranchId = &EndPoint{
	Path:             "/branch/:branch_id/employee/:employee_id/services",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetEmployeeServicesByBranchId",
	Description:      "View employee offered services at the branch by branch ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var AddServiceToBranch = &EndPoint{
	Path:             "/branch/:branch_id/service/:service_id",
	Method:           namespace.CreateActionMethod,
	Handler:          "AddServiceToBranch",
	Description:      "Add service to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var RemoveServiceFromBranch = &EndPoint{
	Path:             "/branch/:branch_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "RemoveServiceFromBranch",
	Description:      "Remove service from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}

// --- Client Endpoints --- //

var CreateClient = &EndPoint{
	Path:        "/client",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateClient",
	Description: "Create client",
}
var LoginClient = &EndPoint{
	Path:        "/client/login",
	Method:      namespace.CreateActionMethod,
	Handler:     "LoginClient",
	Description: "Login client",
}
var VerifyClientEmail = &EndPoint{
	Path:        "/client/verify-email/:email/:code",
	Method:      namespace.CreateActionMethod,
	Handler:     "VerifyClientEmail",
	Description: "Verify client email",
}
var GetClientByEmail = &EndPoint{
	Path:             "/client/email/:email",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetClientByEmail",
	Description:      "View client by email",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var UpdateClientById = &EndPoint{
	Path:             "/client/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateClientById",
	Description:      "Update client by ID",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}
var DeleteClientById = &EndPoint{
	Path:             "/client/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteClientById",
	Description:      "Delete client by ID",
	DenyUnauthorized: true,
	Resource:         ClientResource,
}

// --- Company Endpoints --- //

var CreateCompany = &EndPoint{
	Path:        "/company",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateCompany",
	Description: "Create a company",
}
var GetCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetCompanyById",
	Description: "View company by ID",
	Resource:    CompanyResource,
}
var GetCompanyByName = &EndPoint{
	Path:        "/company/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetCompanyByName",
	Description: "View company by name",
}
var GetCompanyByTaxId = &EndPoint{
	Path:        "/company/tax_id/:tax_id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetCompanyByTaxId",
	Description: "View company by tax ID",
}
var UpdateCompanyById = &EndPoint{
	Path:             "/company/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateCompanyById",
	Description:      "Update company by ID",
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var DeleteCompanyById = &EndPoint{
	Path:             "/company/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteCompanyById",
	Description:      "Delete company by ID",
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}

// --- Employee Endpoints --- //

var CreateEmployee = &EndPoint{
	Path:             "/employee",
	Method:           namespace.CreateActionMethod,
	Handler:          "CreateEmployee",
	Description:      "Create employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var LoginEmployee = &EndPoint{
	Path:           "/employee/login",
	Method:         namespace.CreateActionMethod,
	Handler:        "LoginEmployee",
	Description:    "Login employee",
	NeedsCompanyId: true,
}
var VerifyEmployeeEmail = &EndPoint{
	Path:        "/employee/verify-email/:email/:code",
	Method:      namespace.CreateActionMethod,
	Handler:     "VerifyEmployeeEmail",
	Description: "Verify employee email",
}
var GetEmployeeById = &EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetEmployeeById",
	Description:      "View employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var GetEmployeeByEmail = &EndPoint{
	Path:             "/employee/email/:email",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetEmployeeByEmail",
	Description:      "View employee by email",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var UpdateEmployeeById = &EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateEmployeeById",
	Description:      "Update employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var DeleteEmployeeById = &EndPoint{
	Path:             "/employee/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteEmployeeById",
	Description:      "Delete employee by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         EmployeeResource,
}
var AddServiceToEmployee = &EndPoint{
	Path:             "/employee/:employee_id/service/:service_id",
	Method:           namespace.CreateActionMethod,
	Handler:          "AddServiceToEmployee",
	Description:      "Add service to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var RemoveServiceFromEmployee = &EndPoint{
	Path:             "/employee/:employee_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "RemoveServiceFromEmployee",
	Description:      "Remove service from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var AddBranchToEmployee = &EndPoint{
	Path:             "/employee/:employee_id/branch/:branch_id",
	Method:           namespace.CreateActionMethod,
	Handler:          "AddBranchToEmployee",
	Description:      "Add employee to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var RemoveBranchFromEmployee = &EndPoint{
	Path:             "/employee/:employee_id/branch/:branch_id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "RemoveBranchFromEmployee",
	Description:      "Remove employee from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         BranchResource,
}
var AddRoleToEmployee = &EndPoint{
	Path:             "/employee/:employee_id/role/:role_id",
	Method:           namespace.CreateActionMethod,
	Handler:          "AddRoleToEmployee",
	Description:      "Add role to employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         RoleResource,
}
var RemoveRoleFromEmployee = &EndPoint{
	Path:             "/employee/:employee_id/role/:role_id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "RemoveRoleFromEmployee",
	Description:      "Remove role from employee",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         RoleResource,
}

// --- Holiday Endpoints --- //

var CreateHoliday = &EndPoint{
	Path:             "/holiday",
	Method:           namespace.CreateActionMethod,
	Handler:          "CreateHoliday",
	Description:      "Create a holiday",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetHolidayById = &EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetHolidayById",
	Description:      "View holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         HolidayResource,
}
var GetHolidayByName = &EndPoint{
	Path:           "/holiday/name/:name",
	Method:         namespace.ViewActionMethod,
	Handler:        "GetHolidayByName",
	Description:    "View holiday by name",
	NeedsCompanyId: true,
}
var UpdateHolidayById = &EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateHolidayById",
	Description:      "Update holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         HolidayResource,
}
var DeleteHolidayById = &EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteHolidayById",
	Description:      "Delete holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         HolidayResource,
}

// --- Sector Endpoints --- //

var CreateSector = &EndPoint{
	Path:             "/sector",
	Method:           namespace.CreateActionMethod,
	Handler:          "CreateSector",
	Description:      "Creates a company sector",
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetSectorById",
	Description: "Retrieves a company sector by ID",
}
var GetSectorByName = &EndPoint{
	Path:        "/sector/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetSectorByName",
	Description: "Retrieves a company sector by name",
}
var UpdateSectorById = &EndPoint{
	Path:             "/sector/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateSectorById",
	Description:      "Updates a company sector by ID",
	DenyUnauthorized: true,
	Resource:         SectorResource,
}
var DeleteSectorById = &EndPoint{
	Path:             "/sector/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteSectorById",
	Description:      "Deletes a company sector by ID",
	DenyUnauthorized: true,
	Resource:         SectorResource,
}

// --- Service Endpoints --- //

var CreateService = &EndPoint{
	Path:             "/service",
	Method:           namespace.CreateActionMethod,
	Handler:          "CreateService",
	Description:      "Create a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         CompanyResource,
}
var GetServiceById = &EndPoint{
	Path:             "/service/:id",
	Method:           namespace.ViewActionMethod,
	Handler:          "GetServiceById",
	Description:      "View service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var GetServiceByName = &EndPoint{
	Path:           "/service/name/:name",
	Method:         namespace.ViewActionMethod,
	Handler:        "GetServiceByName",
	Description:    "View service by name",
	NeedsCompanyId: true,
}
var UpdateServiceById = &EndPoint{
	Path:             "/service/:id",
	Method:           namespace.UpdateActionMethod,
	Handler:          "UpdateServiceById",
	Description:      "Update service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         ServiceResource,
}
var DeleteServiceById = &EndPoint{
	Path:             "/service/:id",
	Method:           namespace.DeleteActionMethod,
	Handler:          "DeleteServiceById",
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
	VerifyExistingAccount,
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
