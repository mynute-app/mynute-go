package policySeed

import (
	"mynute-go/services/core/config/db/model"
)

// GetAllClientPolicies returns all client policy definitions
func GetAllClientPolicies() []*model.ClientPolicy {
	policies := []*model.ClientPolicy{}

	// Client profile policies
	policies = append(policies,
		ClientAllowGetClientByEmail,
		ClientAllowGetClientById,
		ClientAllowUpdateClientById,
		ClientAllowDeleteClientById,
		ClientAllowUpdateClientImages,
		ClientAllowDeleteClientImage,
	)

	// Client appointment policies
	policies = append(policies,
		ClientAllowCreateAppointment,
		ClientAllowGetAppointmentByID,
		ClientAllowUpdateAppointmentByID,
		ClientAllowCancelAppointmentByID,
	)

	return policies
}
