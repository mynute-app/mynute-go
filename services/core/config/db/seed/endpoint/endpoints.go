package endpointSeed

import (
	authModel "mynute-go/services/core/config/db/model"
)

// GetAllEndpoints returns all endpoint definitions for seeding
func GetAllEndpoints() []*authModel.EndPoint {
	return []*authModel.EndPoint{
		// Admin endpoints
		AdminLoginByPassword,
		AreThereAnyAdmin,
		CreateFirstAdmin,
		GetAdminByID,
		GetAdminByEmail,
		ListAdmins,
		CreateAdmin,
		UpdateAdminByID,
		DeleteAdminByID,
		ResetAdminPasswordByEmail,
		SendAdminVerificationCodeByEmail,
		VerifyAdminEmail,
		ListAdminRoles,
		CreateAdminRole,
		GetAdminRoleByID,
		UpdateAdminRoleByID,
		DeleteAdminRoleByID,

		// Appointment endpoints
		CreateAppointment,
		GetAppointmentByID,
		UpdateAppointmentByID,
		CancelAppointmentByID,

		// Auth endpoints
		BeginAuthProviderCallback,
		GetAuthCallbackFunction,
		LogoutProvider,

		// Branch endpoints
		CreateBranch,
		GetBranchById,
		UpdateBranchById,
		DeleteBranchById,
		GetEmployeeServicesByBranchId,
		AddServiceToBranch,
		RemoveServiceFromBranch,
		CreateBranchWorkSchedule,
		UpdateBranchImages,
		DeleteBranchImage,
		GetBranchWorkSchedule,
		GetBranchWorkRange,
		UpdateBranchWorkRange,
		DeleteBranchWorkRange,
		AddBranchWorkRangeServices,
		DeleteBranchWorkRangeService,
		GetBranchAppointmentsById,

		// Client endpoints
		CreateClient,
		LoginClient,
		LoginClientByEmailCode,
		SendLoginCodeToClientEmail,
		ResetClientPasswordByEmail,
		SendClientVerificationCodeByEmail,
		VerifyClientEmail,
		GetClientByEmail,
		GetClientById,
		UpdateClientById,
		DeleteClientById,
		UpdateClientImages,
		DeleteClientImage,

		// Company endpoints
		CreateCompany,
		GetCompanyById,
		GetCompanyByName,
		CheckIfCompanyExistsByTaxID,
		GetCompanyByTaxId,
		GetCompanyBySubdomain,
		UpdateCompanyById,
		UpdateCompanyImages,
		DeleteCompanyImage,
		UpdateCompanyColors,
		DeleteCompanyById,

		// Employee endpoints
		CreateEmployee,
		LoginEmployee,
		LoginEmployeeByEmailCode,
		SendLoginCodeToEmployeeEmail,
		ResetEmployeePasswordByEmail,
		SendEmployeeVerificationCodeByEmail,
		VerifyEmployeeEmail,
		GetEmployeeById,
		GetEmployeeByEmail,
		UpdateEmployeeById,
		UpdateEmployeeImages,
		DeleteEmployeeImage,
		CreateEmployeeWorkSchedule,
		DeleteEmployeeWorkRange,
		UpdateEmployeeWorkRange,
		GetEmployeeWorkSchedule,
		GetEmployeeWorkRange,
		AddEmployeeWorkRangeServices,
		DeleteEmployeeWorkRangeService,
		DeleteEmployeeById,
		AddServiceToEmployee,
		RemoveServiceFromEmployee,
		AddBranchToEmployee,
		RemoveBranchFromEmployee,
		AddRoleToEmployee,
		RemoveRoleFromEmployee,
		GetEmployeeAppointmentsById,

		// Holiday endpoints
		CreateHoliday,
		GetHolidayById,
		GetHolidayByName,
		UpdateHolidayById,
		DeleteHolidayById,

		// Sector endpoints
		CreateSector,
		GetSectorById,
		GetSectorByName,
		UpdateSectorById,
		DeleteSectorById,

		// Service endpoints
		CreateService,
		GetServiceById,
		GetServiceByName,
		UpdateServiceById,
		DeleteServiceById,
		UpdateServiceImages,
		DeleteServiceImage,
		GetServiceAvailability,
	}
}

