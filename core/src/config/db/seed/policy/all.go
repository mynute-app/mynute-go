package policySeed

import (
	authModel "mynute-go/auth/model"
)

// GetAllPolicies returns all policy definitions from all modules
func GetAllPolicies() []*authModel.PolicyRule {
	policies := []*authModel.PolicyRule{}

	// Client policies
	policies = append(policies,
		AllowGetClientByEmail,
		AllowGetClientById,
		AllowUpdateClientById,
		AllowDeleteClientById,
		AllowUpdateClientImages,
		AllowDeleteClientImage,
	)

	// Add policies from other modules as they are defined
	// policies = append(policies, CompanyPolicies...)
	// policies = append(policies, EmployeePolicies...)
	// policies = append(policies, BranchPolicies...)
	// policies = append(policies, AppointmentPolicies...)
	// policies = append(policies, ServicePolicies...)
	// policies = append(policies, HolidayPolicies...)

	return policies
}
