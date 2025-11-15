package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- TENANT POLICIES: Appointment-related operations (Company User Side) ---

var TenantAllowCreateAppointment = &model.TenantPolicy{
	Name:        "TP: CanCreateAppointment",
	Description: "Allows company users (managers/employees) to create appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateAppointment.ID,
	Conditions:  model.JsonRawMessage(CompanyMembershipAndManagerOrAssignedEmployeeCheck),
}

var TenantAllowGetAppointmentByID = &model.TenantPolicy{
	Name:        "TP: CanViewAppointment",
	Description: "Allows company users to view appointments based on role/relation.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(CompanyMembershipAndManagerOrAssignedEmployeeCheck),
}

var TenantAllowUpdateAppointmentByID = &model.TenantPolicy{
	Name:        "TP: CanUpdateAppointment",
	Description: "Allows company managers/assigned employees to update appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(CompanyMembershipAndManagerOrAssignedEmployeeCheck),
}

var TenantAllowCancelAppointmentByID = &model.TenantPolicy{
	Name:        "TP: CanCancelAppointment",
	Description: "Allows company managers/assigned employees to cancel appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CancelAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(CompanyMembershipAndManagerOrAssignedEmployeeCheck),
}
