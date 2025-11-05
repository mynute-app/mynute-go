package policySeed

import (
	"mynute-go/services/auth/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

var AllowCreateAppointment = &model.PolicyRule{
	Name:        "SDP: CanCreateAppointment",
	Description: "Allows clients to create appointments for themselves, or company users based on role/relation.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateAppointment.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Allow Client Creation OR Company User Creation",
		LogicType:   "OR",
		Children: []model.ConditionNode{
			// Client creating for themselves. Assumes 'client_id' is in the body and resource.client_id gets populated.
			{
				Description: "Client Self-Creation",
				LogicType:   "AND",
				Children: []model.ConditionNode{
					{Leaf: &model.ConditionLeaf{Attribute: "subject.company_id", Operator: "IsNull", Description: "Must be a Client"}},
					{Leaf: &model.ConditionLeaf{Attribute: "subject.id", Operator: "Equals", ResourceAttribute: "body.client_id", Description: "Client ID in body must match Subject ID"}},
				},
			},
			// Company User creating. Assumes 'branch_id' and maybe 'employee_id' are in the body.
			CompanyMembershipAndManagerOrAssignedEmployeeCheck, // Company users with appropriate permissions
		},
	}),
}

var AllowGetAppointmentByID = &model.PolicyRule{
	Name:        "SDP: CanViewAppointment",
	Description: "Allows clients to view own appointments, or company users based on role/relation.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetAppointmentByID.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Allow Client Access OR Company User Access",
		LogicType:   "OR",
		Children: []model.ConditionNode{
			ClientAccessCheck, // Client can view if appointment's client_id matches subject.id
			CompanyMembershipAndManagerOrAssignedEmployeeCheck, // Company user with appropriate role/relation
		},
	}),
}

var AllowUpdateAppointmentByID = &model.PolicyRule{
	Name:        "SDP: CanUpdateAppointment",
	Description: "Allows clients to update own appointments, or company managers/assigned employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(ClientOrCompanyMembershipAndManagerOrAssignedEmployeeCheck),
}

var AllowCancelAppointmentByID = &model.PolicyRule{
	Name:        "SDP: CanCancelAppointment",
	Description: "Allows clients to cancel their own appointments, or company managers/assigned employees.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CancelAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(ClientOrCompanyMembershipAndManagerOrAssignedEmployeeCheck),
}
