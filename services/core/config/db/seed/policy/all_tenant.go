package policySeed

import (
	"mynute-go/services/core/config/db/model"
)

// GetAllTenantPolicies returns all tenant policy definitions
func GetAllTenantPolicies() []*model.TenantPolicy {
	policies := []*model.TenantPolicy{}

	// Company policies
	policies = append(policies,
		TenantAllowGetCompanyById,
		TenantAllowUpdateCompanyById,
		TenantAllowDeleteCompanyById,
		TenantAllowUpdateCompanyImages,
		TenantAllowDeleteCompanyImage,
		TenantAllowUpdateCompanyColors,
	)

	// Employee policies
	policies = append(policies,
		TenantAllowCreateEmployee,
		TenantAllowGetEmployeeById,
		TenantAllowGetEmployeeByEmail,
		TenantAllowUpdateEmployeeById,
		TenantAllowDeleteEmployeeById,
		TenantAllowCreateEmployeeWorkSchedule,
		TenantAllowGetEmployeeWorkRangeById,
		TenantAllowUpdateEmployeeWorkRange,
		TenantAllowDeleteEmployeeWorkRange,
		TenantAllowAddEmployeeWorkRangeServices,
		TenantAllowDeleteEmployeeWorkRangeService,
		TenantAllowAddServiceToEmployee,
		TenantAllowRemoveServiceFromEmployee,
		TenantAllowAddBranchToEmployee,
		TenantAllowRemoveBranchFromEmployee,
		TenantAllowUpdateEmployeeImages,
		TenantAllowDeleteEmployeeImage,
		TenantAllowGetEmployeeAppointmentsById,
	)

	// Branch policies
	policies = append(policies,
		TenantAllowCreateBranch,
		TenantAllowGetBranchById,
		TenantAllowUpdateBranchById,
		TenantAllowDeleteBranchById,
		TenantAllowGetEmployeeServicesByBranchId,
		TenantAllowAddServiceToBranch,
		TenantAllowRemoveServiceFromBranch,
		TenantAllowCreateBranchWorkSchedule,
		TenantAllowGetBranchWorkRangeById,
		TenantAllowDeleteBranchWorkRangeById,
		TenantAllowUpdateBranchWorkRangeById,
		TenantAllowAddBranchWorkRangeService,
		TenantAllowDeleteBranchWorkRangeService,
		TenantAllowUpdateBranchImages,
		TenantAllowDeleteBranchImage,
		TenantAllowGetBranchAppointmentsById,
	)

	// Service policies
	policies = append(policies,
		TenantAllowCreateService,
		TenantAllowGetServiceById,
		TenantAllowUpdateServiceById,
		TenantAllowDeleteServiceById,
		TenantAllowUpdateServiceImages,
		TenantAllowDeleteServiceImage,
	)

	// Holiday policies
	policies = append(policies,
		TenantAllowCreateHoliday,
		TenantAllowGetHolidayById,
		TenantAllowUpdateHolidayById,
		TenantAllowDeleteHolidayById,
	)

	// Appointment policies (tenant side)
	policies = append(policies,
		TenantAllowCreateAppointment,
		TenantAllowGetAppointmentByID,
		TenantAllowUpdateAppointmentByID,
		TenantAllowCancelAppointmentByID,
	)

	return policies
}
