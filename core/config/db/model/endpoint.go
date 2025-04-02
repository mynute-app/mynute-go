package model

import (
	"log"

	"gorm.io/gorm"
)

var AllowResourceCreation = false

type EndPoint struct {
	gorm.Model
	Handler       string `json:"handler"`
	Description   string `json:"description"`
	Method        string `gorm:"type:varchar(10)" json:"method"`
	Path          string `json:"path"`
	IsPublic      bool   `gorm:"default:false" json:"is_public"`
}

func (r *EndPoint) BeforeCreate(tx *gorm.DB) error {
	if !AllowResourceCreation {
		panic("EndPoint creation is not allowed")
	}
	return nil
}

// Custom Composite Index
func (EndPoint) TableName() string {
	return "resources"
}

func (EndPoint) Indexes() map[string]string {
	return map[string]string{
		"idx_method_path": "CREATE UNIQUE INDEX idx_method_path ON routes (method, path)",
	}
}

// --- Appointment Resources --- //

var CreateAppointment = &EndPoint{
	Path:          "/appointment",
	Method:        "POST",
	Handler:       "CreateAppointment", // Assuming handler name matches reference ac.CreateAppointment
	Description:   "Create an appointment",
	IsPublic:      false, // Access: "private"
}
var GetAppointmentByID = &EndPoint{
	Path:          "/appointment/:id",
	Method:        "GET",
	Handler:       "GetAppointmentByID", // Assuming handler name matches reference ac.GetAppointmentByID
	Description:   "Get appointment by ID",
	IsPublic:      false, // Access: "private"
}
var UpdateAppointmentByID = &EndPoint{
	Path:          "/appointment/:id",
	Method:        "PATCH",                 // Corrected from GET based on reference
	Handler:       "UpdateAppointmentByID", // Corrected from GetAppointmentByID based on reference
	Description:   "Update appointment by ID",
	IsPublic:      false, // Access: "private"
}
var DeleteAppointmentByID = &EndPoint{
	Path:          "/appointment/:id",
	Method:        "DELETE",
	Handler:       "DeleteAppointmentByID", // Assuming handler name matches reference ac.DeleteAppointmentByID
	Description:   "Delete appointment by ID",
	IsPublic:      false, // Access: "private"
}

// --- Auth Resources --- //

var VerifyExistingAccount = &EndPoint{
	Path:        "/auth/verify-existing-account",
	Method:      "POST",
	Handler:     "VerifyExistingAccount", // From ac.VerifyExistingAccount
	Description: "Verify if an account exists",
	IsPublic:    true, // Access: "public"
}
var BeginAuthProviderCallback = &EndPoint{
	Path:        "/auth/oauth/:provider",
	Method:      "GET",
	Handler:     "BeginAuthProviderCallback", // From ac.BeginAuthProviderCallback
	Description: "Begin auth provider callback",
	IsPublic:    true, // Access: "public"
}
var GetAuthCallbackFunction = &EndPoint{
	Path:        "/auth/oauth/:provider/callback",
	Method:      "GET",
	Handler:     "GetAuthCallbackFunction", // From ac.GetAuthCallbackFunction
	Description: "Get auth callback function",
	IsPublic:    true, // Access: "public"
}
var LogoutProvider = &EndPoint{
	Path:        "/auth/oauth/logout",
	Method:      "GET",
	Handler:     "LogoutProvider", // From ac.LogoutProvider
	Description: "Logout provider",
	IsPublic:    true, // Access: "public"
}

// --- Branch Resources --- //

var CreateBranch = &EndPoint{
	Path:        "/branch",
	Method:      "POST",
	Handler:     "CreateBranch", // From bc.CreateBranch
	Description: "Create a branch",
	IsPublic:    false, // Access: "private"
}
var GetBranchById = &EndPoint{
	Path:        "/branch/:id",
	Method:      "GET",
	Handler:     "GetBranchById", // From bc.GetBranchById
	Description: "Get branch by ID",
	IsPublic:    false, // Access: "private"
}
var GetBranchByName = &EndPoint{
	Path:        "/branch/name/:name",
	Method:      "GET",
	Handler:     "GetBranchByName", // From bc.GetBranchByName
	Description: "Get branch by name",
	IsPublic:    false, // Access: "private"
}
var UpdateBranchById = &EndPoint{
	Path:        "/branch/:id",
	Method:      "PATCH",
	Handler:     "UpdateBranchById", // From bc.UpdateBranchById
	Description: "Update branch by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteBranchById = &EndPoint{
	Path:        "/branch/:id",
	Method:      "DELETE",
	Handler:     "DeleteBranchById", // From bc.DeleteBranchById
	Description: "Delete branch by ID",
	IsPublic:    false, // Access: "private"
}
var GetEmployeeServicesByBranchId = &EndPoint{
	Path:        "/branch/:branch_id/employee/:employee_id/services",
	Method:      "GET",
	Handler:     "GetEmployeeServicesByBranchId", // From bc.GetEmployeeServicesByBranchId
	Description: "Get employee offered services at the branch by branch ID",
	IsPublic:    false, // Access: "private"
}
var AddServiceToBranch = &EndPoint{
	Path:        "/branch/:branch_id/service/:service_id",
	Method:      "POST",
	Handler:     "AddServiceToBranch", // From bc.AddServiceToBranch
	Description: "Add service to branch",
	IsPublic:    false, // Access: "private"
}
var RemoveServiceFromBranch = &EndPoint{
	Path:        "/branch/:branch_id/service/:service_id",
	Method:      "DELETE",
	Handler:     "RemoveServiceFromBranch", // From bc.RemoveServiceFromBranch
	Description: "Remove service from branch",
	IsPublic:    false, // Access: "private"
}

// --- Client Resources --- //

var CreateClient = &EndPoint{
	Path:        "/client",
	Method:      "POST",
	Handler:     "CreateClient", // From cc.CreateClient
	Description: "Create client",
	IsPublic:    true, // Access: "public"
}
var LoginClient = &EndPoint{
	Path:        "/client/login",
	Method:      "POST",
	Handler:     "LoginClient", // From cc.LoginClient
	Description: "Login client",
	IsPublic:    true, // Access: "public"
}
var VerifyClientEmail = &EndPoint{
	Path:        "/client/verify-email/:email/:code",
	Method:      "POST",
	Handler:     "VerifyClientEmail", // From cc.VerifyClientEmail
	Description: "Verify client email",
	IsPublic:    true, // Access: "public"
}
var GetClientByEmail = &EndPoint{
	Path:        "/client/email/:email",
	Method:      "GET",
	Handler:     "GetClientByEmail", // From cc.GetClientByEmail
	Description: "Get client by email",
	IsPublic:    false, // Access: "private"
}
var UpdateClientById = &EndPoint{
	Path:        "/client/:id",
	Method:      "PATCH",
	Handler:     "UpdateClientById", // From cc.UpdateClientById
	Description: "Update client by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteClientById = &EndPoint{
	Path:        "/client/:id",
	Method:      "DELETE",
	Handler:     "DeleteClientById", // From cc.DeleteClientById
	Description: "Delete client by ID",
	IsPublic:    false, // Access: "private"
}

// --- Company Resources --- //

var CreateCompany = &EndPoint{
	Path:        "/company",
	Method:      "POST",
	Handler:     "CreateCompany", // From cc.CreateCompany
	Description: "Create a company",
	IsPublic:    true, // Access: "public"
}
var GetCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      "GET",
	Handler:     "GetCompanyById", // From cc.GetCompanyById
	Description: "Get company by ID",
	IsPublic:    false, // Access: "private"
}
var GetCompanyByName = &EndPoint{
	Path:        "/company/name/:name",
	Method:      "GET",
	Handler:     "GetCompanyByName", // From cc.GetCompanyByName
	Description: "Get company by name",
	IsPublic:    true, // Access: "public"
}
var GetCompanyByTaxId = &EndPoint{
	Path:        "/company/tax_id/:tax_id",
	Method:      "GET",
	Handler:     "GetCompanyByTaxId", // From cc.GetCompanyByTaxId
	Description: "Get company by tax ID",
	IsPublic:    true, // Access: "public"
}
var UpdateCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      "PATCH",
	Handler:     "UpdateCompanyById", // From cc.UpdateCompanyById
	Description: "Update company by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteCompanyById = &EndPoint{
	Path:        "/company/:id",
	Method:      "DELETE",
	Handler:     "DeleteCompanyById", // From cc.DeleteCompanyById
	Description: "Delete company by ID",
	IsPublic:    false, // Access: "private"
}

// --- Employee Resources --- //

var CreateEmployee = &EndPoint{
	Path:        "/employee",
	Method:      "POST",
	Handler:     "CreateEmployee", // From ec.CreateEmployee
	Description: "Create employee",
	IsPublic:    false, // Access: "private"
}
var LoginEmployee = &EndPoint{
	Path:        "/employee/login",
	Method:      "POST",
	Handler:     "LoginEmployee", // From ec.LoginEmployee
	Description: "Login employee",
	IsPublic:    true, // Access: "public"
}
var VerifyEmployeeEmail = &EndPoint{
	Path:        "/employee/verify-email/:email/:code",
	Method:      "POST",
	Handler:     "VerifyEmployeeEmail", // From ec.VerifyEmployeeEmail
	Description: "Verify employee email",
	IsPublic:    true, // Access: "public"
}
var GetEmployeeById = &EndPoint{
	Path:        "/employee/:id",
	Method:      "GET",
	Handler:     "GetEmployeeById", // From ec.GetEmployeeById
	Description: "Get employee by ID",
	IsPublic:    false, // Access: "private"
}
var GetEmployeeByEmail = &EndPoint{
	Path:        "/employee/email/:email",
	Method:      "GET",
	Handler:     "GetEmployeeByEmail", // From ec.GetEmployeeByEmail
	Description: "Get employee by email",
	IsPublic:    false, // Access: "private"
}
var UpdateEmployeeById = &EndPoint{
	Path:        "/employee/:id",
	Method:      "PATCH",
	Handler:     "UpdateEmployeeById", // From ec.UpdateEmployeeById
	Description: "Update employee by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteEmployeeById = &EndPoint{
	Path:        "/employee/:id",
	Method:      "DELETE",
	Handler:     "DeleteEmployeeById", // From ec.DeleteEmployeeById
	Description: "Delete employee by ID",
	IsPublic:    false, // Access: "private"
}
var AddServiceToEmployee = &EndPoint{
	Path:        "/employee/:employee_id/service/:service_id",
	Method:      "POST",
	Handler:     "AddServiceToEmployee", // From ec.AddServiceToEmployee
	Description: "Add service to employee",
	IsPublic:    false, // Access: "private"
}
var RemoveServiceFromEmployee = &EndPoint{
	Path:        "/employee/:employee_id/service/:service_id",
	Method:      "DELETE",
	Handler:     "RemoveServiceFromEmployee", // From ec.RemoveServiceFromEmployee
	Description: "Remove service from employee",
	IsPublic:    false, // Access: "private"
}
var AddBranchToEmployee = &EndPoint{
	Path:        "/employee/:employee_id/branch/:branch_id",
	Method:      "POST",
	Handler:     "AddBranchToEmployee", // From ec.AddBranchToEmployee
	Description: "Add employee to branch",
	IsPublic:    false, // Access: "private"
}
var RemoveBranchFromEmployee = &EndPoint{
	Path:        "/employee/:employee_id/branch/:branch_id",
	Method:      "DELETE",
	Handler:     "RemoveBranchFromEmployee", // From ec.RemoveBranchFromEmployee
	Description: "Remove employee from branch",
	IsPublic:    false, // Access: "private"
}

// --- Holiday Resources --- //

var CreateHoliday = &EndPoint{
	Path:        "/holiday",
	Method:      "POST",
	Handler:     "CreateHoliday", // From hc.CreateHoliday
	Description: "Create a holiday",
	IsPublic:    false, // Access: "private"
}
var GetHolidayById = &EndPoint{
	Path:        "/holiday/:id",
	Method:      "GET",
	Handler:     "GetHolidayById", // From hc.GetHolidayById
	Description: "Get holiday by ID",
	IsPublic:    false, // Access: "private"
}
var GetHolidayByName = &EndPoint{
	Path:        "/holiday/name/:name",
	Method:      "GET",
	Handler:     "GetHolidayByName", // From hc.GetHolidayByName
	Description: "Get holiday by name",
	IsPublic:    true, // Access: "public"
}
var UpdateHolidayById = &EndPoint{
	Path:        "/holiday/:id",
	Method:      "PATCH",
	Handler:     "UpdateHolidayById", // From hc.UpdateHolidayById
	Description: "Update holiday by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteHolidayById = &EndPoint{
	Path:        "/holiday/:id",
	Method:      "DELETE",
	Handler:     "DeleteHolidayById", // From hc.DeleteHolidayById
	Description: "Delete holiday by ID",
	IsPublic:    false, // Access: "private"
}

// --- Sector Resources --- //

var CreateSector = &EndPoint{
	Path:        "/sector",
	Method:      "POST",
	Handler:     "CreateSector", // From sc.CreateSector
	Description: "Creates a company sector",
	IsPublic:    false, // Access: "private"
}
var GetSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      "GET",
	Handler:     "GetSectorById", // From sc.GetSectorById
	Description: "Retrieves a company sector by ID",
	IsPublic:    true, // Access: "public"
}
var GetSectorByName = &EndPoint{
	Path:        "/sector/name/:name",
	Method:      "GET",
	Handler:     "GetSectorByName", // From sc.GetSectorByName
	Description: "Retrieves a company sector by name",
	IsPublic:    true, // Access: "public"
}
var UpdateSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      "PATCH",
	Handler:     "UpdateSectorById", // From sc.UpdateSectorById
	Description: "Updates a company sector by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteSectorById = &EndPoint{
	Path:        "/sector/:id",
	Method:      "DELETE",
	Handler:     "DeleteSectorById", // From sc.DeleteSectorById
	Description: "Deletes a company sector by ID",
	IsPublic:    false, // Access: "private"
}

// --- Service Resources --- //

var CreateService = &EndPoint{
	Path:        "/service",
	Method:      "POST",
	Handler:     "CreateService", // From sc.CreateService
	Description: "Create a service",
	IsPublic:    false, // Access: "private"
}
var GetServiceById = &EndPoint{
	Path:        "/service/:id",
	Method:      "GET",
	Handler:     "GetServiceById", // From sc.GetServiceById
	Description: "Get service by ID",
	IsPublic:    false, // Access: "private"
}
var GetServiceByName = &EndPoint{
	Path:        "/service/name/:name",
	Method:      "GET",
	Handler:     "GetServiceByName", // From sc.GetServiceByName
	Description: "Get service by name",
	IsPublic:    true, // Access: "public"
}
var UpdateServiceById = &EndPoint{
	Path:        "/service/:id",
	Method:      "PATCH",
	Handler:     "UpdateServiceById", // From sc.UpdateServiceById
	Description: "Update service by ID",
	IsPublic:    false, // Access: "private"
}
var DeleteServiceById = &EndPoint{
	Path:        "/service/:id",
	Method:      "DELETE",
	Handler:     "DeleteServiceById", // From sc.DeleteServiceById
	Description: "Delete service by ID",
	IsPublic:    false, // Access: "private"
}

// --- Combine all resources into a slice for seeding --- //
var Resources = []*EndPoint{
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

func SeedResources(db *gorm.DB) ([]*EndPoint, error) {
	AllowResourceCreation = true
	defer func() { AllowResourceCreation = false }()
	for _, rsrc := range Resources {
		err := db.Where("method = ? AND path = ?", rsrc.Method, rsrc.Path).First(rsrc).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(rsrc).Error; err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}
	log.Println("System resources seeded successfully!")
	return Resources, nil
}
