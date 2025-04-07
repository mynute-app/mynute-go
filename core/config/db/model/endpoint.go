package model

import (
	"agenda-kaki-go/core/config/namespace"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllowEndpointCreation = false

type EndPoint struct {
	BaseModel
	Handler     string     `json:"handler"`
	Description string     `json:"description"`
	Method      string     `gorm:"type:varchar(10)" json:"method"`
	Path        string     `json:"path"`
	IsPublic    bool       `gorm:"default:false" json:"is_public"`
	ResourceID  *uuid.UUID `json:"resource_id"`
	Resource    *Resource  `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

func (r *EndPoint) BeforeCreate(tx *gorm.DB) error {
	if !AllowEndpointCreation {
		panic("EndPoint creation is not allowed")
	}
	return nil
}

// Custom Composite Index
func (EndPoint) TableName() string {
	return "endpoints"
}

func (EndPoint) Indexes() map[string]string {
	return map[string]string{
		"idx_method_path": "CREATE UNIQUE INDEX idx_method_path ON routes (method, path)",
	}
}

// --- Appointment Endpoints --- //

var CreateAppointment = &EndPoint{
	Path:        "/appointment",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateAppointment", // Assuming handler name matches reference ac.CreateAppointment
	Description: "Create an appointment",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var GetAppointmentByID = &EndPoint{
	Path:        "/appointment/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetAppointmentByID", // Assuming handler name matches reference ac.GetAppointmentByID
	Description: "View appointment by ID",
	IsPublic:    false, // Access: "private"
	Resource:    AppointmentResource,
}
var UpdateAppointmentByID = &EndPoint{
	Path:        "/appointment/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateAppointmentByID", // Corrected from GetAppointmentByID based on reference
	Description: "Update appointment by ID",
	IsPublic:    false, // Access: "private"
	Resource:    AppointmentResource,
}
var DeleteAppointmentByID = &EndPoint{
	Path:        "/appointment/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteAppointmentByID", // Assuming handler name matches reference ac.DeleteAppointmentByID
	Description: "Delete appointment by ID",
	IsPublic:    false, // Access: "private"
	Resource:    AppointmentResource,
}

// --- Auth Endpoints --- //

var VerifyExistingAccount = &EndPoint{
	Path:        "/auth/verify-existing-account",
	Method:      namespace.CreateActionMethod,
	Handler:     "VerifyExistingAccount", // From ac.VerifyExistingAccount
	Description: "Verify if an account exists",
	IsPublic:    true, // Access: "public"
}
var BeginAuthProviderCallback = &EndPoint{
	Path:        "/auth/oauth/:provider",
	Method:      namespace.ViewActionMethod,
	Handler:     "BeginAuthProviderCallback", // From ac.BeginAuthProviderCallback
	Description: "Begin auth provider callback",
	IsPublic:    true, // Access: "public"
}
var GetAuthCallbackFunction = &EndPoint{
	Path:        "/auth/oauth/:provider/callback",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetAuthCallbackFunction", // From ac.GetAuthCallbackFunction
	Description: "View auth callback function",
	IsPublic:    true, // Access: "public"
}
var LogoutProvider = &EndPoint{
	Path:        "/auth/oauth/logout",
	Method:      namespace.ViewActionMethod,
	Handler:     "LogoutProvider", // From ac.LogoutProvider
	Description: "Logout provider",
	IsPublic:    true, // Access: "public"
}

// --- Branch Endpoints --- //

var CreateBranch = &EndPoint{
	Path:        "/branch",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateBranch", // From bc.CreateBranch
	Description: "Create a branch",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var GetBranchById = &EndPoint{
	Path:        "/branch/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetBranchById", // From bc.GetBranchById
	Description: "View branch by ID",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var GetBranchByName = &EndPoint{
	Path:        "/branch/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetBranchByName", // From bc.GetBranchByName
	Description: "View branch by name",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var UpdateBranchById = &EndPoint{
	Path:        "/branch/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateBranchById", // From bc.UpdateBranchById
	Description: "Update branch by ID",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var DeleteBranchById = &EndPoint{
	Path:        "/branch/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteBranchById", // From bc.DeleteBranchById
	Description: "Delete branch by ID",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var GetEmployeeServicesByBranchId = &EndPoint{
	Path:        "/branch/:branch_id/employee/:employee_id/services",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetEmployeeServicesByBranchId", // From bc.GetEmployeeServicesByBranchId
	Description: "View employee offered services at the branch by branch ID",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var AddServiceToBranch = &EndPoint{
	Path:        "/branch/:branch_id/service/:service_id",
	Method:      namespace.CreateActionMethod,
	Handler:     "AddServiceToBranch", // From bc.AddServiceToBranch
	Description: "Add service to branch",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var RemoveServiceFromBranch = &EndPoint{
	Path:        "/branch/:branch_id/service/:service_id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "RemoveServiceFromBranch", // From bc.RemoveServiceFromBranch
	Description: "Remove service from branch",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}

// --- Client Endpoints --- //

var CreateClient = &EndPoint{
	Path:        "/client",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateClient", // From cc.CreateClient
	Description: "Create client",
	IsPublic:    true, // Access: "public"
}
var LoginClient = &EndPoint{
	Path:        "/client/login",
	Method:      namespace.CreateActionMethod,
	Handler:     "LoginClient", // From cc.LoginClient
	Description: "Login client",
	IsPublic:    true, // Access: "public"
}
var VerifyClientEmail = &EndPoint{
	Path:        "/client/verify-email/:email/:code",
	Method:      namespace.CreateActionMethod,
	Handler:     "VerifyClientEmail", // From cc.VerifyClientEmail
	Description: "Verify client email",
	IsPublic:    true, // Access: "public"
}
var GetClientByEmail = &EndPoint{
	Path:        "/client/email/:email",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetClientByEmail", // From cc.GetClientByEmail
	Description: "View client by email",
	IsPublic:    false, // Access: "private"
	Resource:    ClientResource,
}
var UpdateClientById = &EndPoint{
	Path:        "/client/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateClientById", // From cc.UpdateClientById
	Description: "Update client by ID",
	IsPublic:    false, // Access: "private"
	Resource:    ClientResource,
}
var DeleteClientById = &EndPoint{
	Path:        "/client/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteClientById", // From cc.DeleteClientById
	Description: "Delete client by ID",
	IsPublic:    false, // Access: "private"
	Resource:    ClientResource,
}

// --- Company Endpoints --- //

var CreateCompany = &EndPoint{
	Path:        "/company",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateCompany", // From cc.CreateCompany
	Description: "Create a company",
	IsPublic:    true, // Access: "public"
}
var GetCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetCompanyById", // From cc.GetCompanyById
	Description: "View company by ID",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var GetCompanyByName = &EndPoint{
	Path:        "/company/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetCompanyByName", // From cc.GetCompanyByName
	Description: "View company by name",
	IsPublic:    true, // Access: "public"
}
var GetCompanyByTaxId = &EndPoint{
	Path:        "/company/tax_id/:tax_id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetCompanyByTaxId", // From cc.GetCompanyByTaxId
	Description: "View company by tax ID",
	IsPublic:    true, // Access: "public"
}
var UpdateCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateCompanyById", // From cc.UpdateCompanyById
	Description: "Update company by ID",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var DeleteCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteCompanyById", // From cc.DeleteCompanyById
	Description: "Delete company by ID",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}

// --- Employee Endpoints --- //

var CreateEmployee = &EndPoint{
	Path:        "/employee",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateEmployee", // From ec.CreateEmployee
	Description: "Create employee",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var LoginEmployee = &EndPoint{
	Path:        "/employee/login",
	Method:      namespace.CreateActionMethod,
	Handler:     "LoginEmployee", // From ec.LoginEmployee
	Description: "Login employee",
	IsPublic:    true, // Access: "public"
}
var VerifyEmployeeEmail = &EndPoint{
	Path:        "/employee/verify-email/:email/:code",
	Method:      namespace.CreateActionMethod,
	Handler:     "VerifyEmployeeEmail", // From ec.VerifyEmployeeEmail
	Description: "Verify employee email",
	IsPublic:    true, // Access: "public"
}
var GetEmployeeById = &EndPoint{
	Path:        "/employee/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetEmployeeById", // From ec.GetEmployeeById
	Description: "View employee by ID",
	IsPublic:    false, // Access: "private"
	Resource:    EmployeeResource,
}
var GetEmployeeByEmail = &EndPoint{
	Path:        "/employee/email/:email",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetEmployeeByEmail", // From ec.GetEmployeeByEmail
	Description: "View employee by email",
	IsPublic:    false, // Access: "private"
	Resource:    EmployeeResource,
}
var UpdateEmployeeById = &EndPoint{
	Path:        "/employee/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateEmployeeById", // From ec.UpdateEmployeeById
	Description: "Update employee by ID",
	IsPublic:    false, // Access: "private"
	Resource:    EmployeeResource,
}
var DeleteEmployeeById = &EndPoint{
	Path:        "/employee/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteEmployeeById", // From ec.DeleteEmployeeById
	Description: "Delete employee by ID",
	IsPublic:    false, // Access: "private"
	Resource:    EmployeeResource,
}
var AddServiceToEmployee = &EndPoint{
	Path:        "/employee/:employee_id/service/:service_id",
	Method:      namespace.CreateActionMethod,
	Handler:     "AddServiceToEmployee", // From ec.AddServiceToEmployee
	Description: "Add service to employee",
	IsPublic:    false, // Access: "private"
	Resource:    ServiceResource,
}
var RemoveServiceFromEmployee = &EndPoint{
	Path:        "/employee/:employee_id/service/:service_id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "RemoveServiceFromEmployee", // From ec.RemoveServiceFromEmployee
	Description: "Remove service from employee",
	IsPublic:    false, // Access: "private"
	Resource:    ServiceResource,
}
var AddBranchToEmployee = &EndPoint{
	Path:        "/employee/:employee_id/branch/:branch_id",
	Method:      namespace.CreateActionMethod,
	Handler:     "AddBranchToEmployee", // From ec.AddBranchToEmployee
	Description: "Add employee to branch",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}
var RemoveBranchFromEmployee = &EndPoint{
	Path:        "/employee/:employee_id/branch/:branch_id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "RemoveBranchFromEmployee", // From ec.RemoveBranchFromEmployee
	Description: "Remove employee from branch",
	IsPublic:    false, // Access: "private"
	Resource:    BranchResource,
}

// --- Holiday Endpoints --- //

var CreateHoliday = &EndPoint{
	Path:        "/holiday",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateHoliday", // From hc.CreateHoliday
	Description: "Create a holiday",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var GetHolidayById = &EndPoint{
	Path:        "/holiday/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetHolidayById", // From hc.GetHolidayById
	Description: "View holiday by ID",
	IsPublic:    false, // Access: "private"
	Resource:    HolidayResource,
}
var GetHolidayByName = &EndPoint{
	Path:        "/holiday/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetHolidayByName", // From hc.GetHolidayByName
	Description: "View holiday by name",
	IsPublic:    true, // Access: "public"
}
var UpdateHolidayById = &EndPoint{
	Path:        "/holiday/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateHolidayById", // From hc.UpdateHolidayById
	Description: "Update holiday by ID",
	IsPublic:    false, // Access: "private"
	Resource:    HolidayResource,
}
var DeleteHolidayById = &EndPoint{
	Path:        "/holiday/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteHolidayById", // From hc.DeleteHolidayById
	Description: "Delete holiday by ID",
	IsPublic:    false, // Access: "private"
	Resource:    HolidayResource,
}

// --- Sector Endpoints --- //

var CreateSector = &EndPoint{
	Path:        "/sector",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateSector", // From sc.CreateSector
	Description: "Creates a company sector",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var GetSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetSectorById", // From sc.GetSectorById
	Description: "Retrieves a company sector by ID",
	IsPublic:    true, // Access: "public"
}
var GetSectorByName = &EndPoint{
	Path:        "/sector/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetSectorByName", // From sc.GetSectorByName
	Description: "Retrieves a company sector by name",
	IsPublic:    true, // Access: "public"
}
var UpdateSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateSectorById", // From sc.UpdateSectorById
	Description: "Updates a company sector by ID",
	IsPublic:    false, // Access: "private"
	Resource:    SectorResource,
}
var DeleteSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteSectorById", // From sc.DeleteSectorById
	Description: "Deletes a company sector by ID",
	IsPublic:    false, // Access: "private"
	Resource:    SectorResource,
}

// --- Service Endpoints --- //

var CreateService = &EndPoint{
	Path:        "/service",
	Method:      namespace.CreateActionMethod,
	Handler:     "CreateService", // From sc.CreateService
	Description: "Create a service",
	IsPublic:    false, // Access: "private"
	Resource:    CompanyResource,
}
var GetServiceById = &EndPoint{
	Path:        "/service/:id",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetServiceById", // From sc.GetServiceById
	Description: "View service by ID",
	IsPublic:    false, // Access: "private"
	Resource:    ServiceResource,
}
var GetServiceByName = &EndPoint{
	Path:        "/service/name/:name",
	Method:      namespace.ViewActionMethod,
	Handler:     "GetServiceByName", // From sc.GetServiceByName
	Description: "View service by name",
	IsPublic:    true, // Access: "public"
}
var UpdateServiceById = &EndPoint{
	Path:        "/service/:id",
	Method:      namespace.UpdateActionMethod,
	Handler:     "UpdateServiceById", // From sc.UpdateServiceById
	Description: "Update service by ID",
	IsPublic:    false, // Access: "private"
	Resource:    ServiceResource,
}
var DeleteServiceById = &EndPoint{
	Path:        "/service/:id",
	Method:      namespace.DeleteActionMethod,
	Handler:     "DeleteServiceById", // From sc.DeleteServiceById
	Description: "Delete service by ID",
	IsPublic:    false, // Access: "private"
	Resource:    ServiceResource,
}

// --- Combine all Endpoints into a slice for seeding --- //
var Endpoints = []*EndPoint{
	// Appointment
	CreateAppointment,
	GetAppointmentByID,
	UpdateAppointmentByID,
	DeleteAppointmentByID,
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

func LoadEndpoints() {
	for _, edp := range Endpoints {
		if edp.Resource != nil {
			edp.ResourceID = &edp.Resource.ID
			edp.Resource = nil // Avoid circular reference
		}
	}
}

func SeedEndpoints(db *gorm.DB) ([]*EndPoint, error) {
	AllowEndpointCreation = true
	tx := db.Begin()
	defer func() {
		AllowEndpointCreation = false
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic occurred during policy seeding: %v", r)
		}
		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit transaction: %v", err)
		}
		log.Print("Resources seeded successfully")
	}()
	LoadEndpoints()
	for _, edp := range Endpoints {
		err := tx.Where("method = ? AND path = ?", edp.Method, edp.Path).First(edp).Error
		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(edp).Error; err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}	
	log.Println("System endpoints seeded successfully!")
	return Endpoints, nil
}
