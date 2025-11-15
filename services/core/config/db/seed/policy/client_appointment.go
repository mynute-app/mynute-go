package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- CLIENT POLICIES: Appointment operations (Client Side) ---

var ClientAllowCreateAppointment = &model.ClientPolicy{
	Name:        "CP: CanCreateAppointment",
	Description: "Allows clients to create appointments for themselves.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateAppointment.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Client Self-Creation",
		LogicType:   "AND",
		Children: []model.ConditionNode{
			{Leaf: &model.ConditionLeaf{Attribute: "subject.company_id", Operator: "IsNull", Description: "Must be a Client"}},
			{Leaf: &model.ConditionLeaf{Attribute: "subject.id", Operator: "Equals", ResourceAttribute: "body.client_id", Description: "Client ID in body must match Subject ID"}},
		},
	}),
}

var ClientAllowGetAppointmentByID = &model.ClientPolicy{
	Name:        "CP: CanViewAppointment",
	Description: "Allows clients to view their own appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(ClientAccessCheck),
}

var ClientAllowUpdateAppointmentByID = &model.ClientPolicy{
	Name:        "CP: CanUpdateAppointment",
	Description: "Allows clients to update their own appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(ClientAccessCheck),
}

var ClientAllowCancelAppointmentByID = &model.ClientPolicy{
	Name:        "CP: CanCancelAppointment",
	Description: "Allows clients to cancel their own appointments.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CancelAppointmentByID.ID,
	Conditions:  model.JsonRawMessage(ClientAccessCheck),
}
