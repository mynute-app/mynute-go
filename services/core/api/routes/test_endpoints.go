package routes

import (
	authClient "mynute-go/services/core/api/lib/auth_client"

	"github.com/google/uuid"
)

// GetTestEndpoints returns a minimal set of endpoints for testing
// when the auth service is not available
func GetTestEndpoints() []*authClient.EndPoint {
	return []*authClient.EndPoint{
		// Client endpoints
		{Method: "POST", Path: "/api/client", ControllerName: "CreateClient", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/client/email/:email", ControllerName: "GetClientByEmail", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/client/:id", ControllerName: "GetClientById", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "PATCH", Path: "/api/client/:id", ControllerName: "UpdateClientById", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "DELETE", Path: "/api/client/:id", ControllerName: "DeleteClientById", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/client/:id/upload-images", ControllerName: "UploadClientImages", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/client/reset-password/email/:email", ControllerName: "ResetClientPasswordByEmail", DenyUnauthorized: false, NeedsCompanyId: false},

		// Company endpoints
		{Method: "POST", Path: "/api/company", ControllerName: "CreateCompany", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/company/tax_id/:tax_id", ControllerName: "GetCompanyByTaxID", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/company/:id", ControllerName: "GetCompanyById", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "PATCH", Path: "/api/company/:id", ControllerName: "UpdateCompanyById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "DELETE", Path: "/api/company/:id", ControllerName: "DeleteCompanyById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "POST", Path: "/api/company/:id/upload-images", ControllerName: "UploadCompanyImages", DenyUnauthorized: true, NeedsCompanyId: true},

		// Employee endpoints
		{Method: "POST", Path: "/api/employee", ControllerName: "CreateEmployee", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/employee", ControllerName: "ListEmployees", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/employee/email/:email", ControllerName: "GetEmployeeByEmail", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/employee/:id", ControllerName: "GetEmployeeById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "PATCH", Path: "/api/employee/:id", ControllerName: "UpdateEmployeeById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "DELETE", Path: "/api/employee/:id", ControllerName: "DeleteEmployeeById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "POST", Path: "/api/employee/:id/upload-images", ControllerName: "UploadEmployeeImages", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "POST", Path: "/api/employee/reset-password/email/:email", ControllerName: "ResetEmployeePasswordByEmail", DenyUnauthorized: false, NeedsCompanyId: false},

		// Branch endpoints
		{Method: "POST", Path: "/api/branch", ControllerName: "CreateBranch", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/branch", ControllerName: "ListBranches", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/branch/:id", ControllerName: "GetBranchById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "PATCH", Path: "/api/branch/:id", ControllerName: "UpdateBranchById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "DELETE", Path: "/api/branch/:id", ControllerName: "DeleteBranchById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "POST", Path: "/api/branch/:id/upload-images", ControllerName: "UploadBranchImages", DenyUnauthorized: true, NeedsCompanyId: true},

		// Service endpoints
		{Method: "POST", Path: "/api/service", ControllerName: "CreateService", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/service", ControllerName: "ListServices", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/service/name/:name", ControllerName: "GetServiceByName", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/service/:id", ControllerName: "GetServiceById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "PATCH", Path: "/api/service/:id", ControllerName: "UpdateServiceById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "DELETE", Path: "/api/service/:id", ControllerName: "DeleteServiceById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "POST", Path: "/api/service/:id/upload-images", ControllerName: "UploadServiceImages", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/service/:id/availability", ControllerName: "GetServiceAvailability", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/service/availability", ControllerName: "GetServicesAvailability", DenyUnauthorized: true, NeedsCompanyId: true},

		// Appointment endpoints
		{Method: "POST", Path: "/api/appointment", ControllerName: "CreateAppointment", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/appointment", ControllerName: "ListAppointments", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/appointment/:id", ControllerName: "GetAppointmentById", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "PATCH", Path: "/api/appointment/:id", ControllerName: "UpdateAppointmentById", DenyUnauthorized: true, NeedsCompanyId: false},
		{Method: "DELETE", Path: "/api/appointment/:id", ControllerName: "DeleteAppointmentById", DenyUnauthorized: true, NeedsCompanyId: false},

		// Holiday endpoints
		{Method: "POST", Path: "/api/holiday", ControllerName: "CreateHoliday", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/holiday", ControllerName: "ListHolidays", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "GET", Path: "/api/holiday/:id", ControllerName: "GetHolidayById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "PATCH", Path: "/api/holiday/:id", ControllerName: "UpdateHolidayById", DenyUnauthorized: true, NeedsCompanyId: true},
		{Method: "DELETE", Path: "/api/holiday/:id", ControllerName: "DeleteHolidayById", DenyUnauthorized: true, NeedsCompanyId: true},

		// Sector endpoints
		{Method: "POST", Path: "/api/sector", ControllerName: "CreateSector", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/sector", ControllerName: "ListSectors", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "GET", Path: "/api/sector/:id", ControllerName: "GetSectorById", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "PATCH", Path: "/api/sector/:id", ControllerName: "UpdateSectorById", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "DELETE", Path: "/api/sector/:id", ControllerName: "DeleteSectorById", DenyUnauthorized: false, NeedsCompanyId: false},

		// Auth endpoints (handled by core service for convenience)
		{Method: "POST", Path: "/api/auth/client/login", ControllerName: "LoginClientByPassword", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/auth/client/login-with-code", ControllerName: "LoginClientByEmailCode", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/auth/client/send-login-code/email/:email", ControllerName: "SendClientLoginValidationCodeByEmail", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/auth/employee/login", ControllerName: "LoginEmployeeByPassword", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/auth/employee/login-with-code", ControllerName: "LoginEmployeeByEmailCode", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/auth/employee/send-login-code/email/:email", ControllerName: "SendEmployeeLoginValidationCodeByEmail", DenyUnauthorized: false, NeedsCompanyId: false},
		{Method: "POST", Path: "/api/auth/company/first_admin", ControllerName: "CreateCompanyFirstAdmin", DenyUnauthorized: false, NeedsCompanyId: false},
	}
}

// initTestEndpoints initializes all endpoints with IDs for testing
func initTestEndpoints() {
	endpoints := GetTestEndpoints()
	for _, ep := range endpoints {
		if ep.ID == uuid.Nil {
			ep.ID = uuid.New()
		}
	}
}
