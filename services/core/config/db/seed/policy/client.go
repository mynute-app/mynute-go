package policySeed

import (
	"mynute-go/services/auth/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

var AllowGetClientByEmail = &model.PolicyRule{
	Name:        "SDP: CanViewClientByEmail",
	Description: "Allows a client to retrieve their own profile by email.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetClientByEmail.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Allow only if the subject's email matches the email in the path.",
		LogicType:   "AND",
		Children: []model.ConditionNode{
			{Leaf: &model.ConditionLeaf{Attribute: "subject.company_id", Operator: "IsNull", Description: "Must be a Client"}},                                                         // Ensure subject is a client
			{Leaf: &model.ConditionLeaf{Attribute: "subject.email", Operator: "Equals", ResourceAttribute: "resource.email", Description: "Subject email must match email from path"}}, // Assumes context has resource.email from path
		},
	}),
}

var AllowGetClientById = &model.PolicyRule{
	Name:        "SDP: CanViewClientById",
	Description: "Allows a client to view their own profile.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetClientById.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck), // Client can view self (checks subject.id == resource.id)
}

var AllowUpdateClientById = &model.PolicyRule{
	Name:        "SDP: CanUpdateClient",
	Description: "Allows a client to update their own profile.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateClientById.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck), // Client can update self (checks subject.id == resource.id)
}

var AllowDeleteClientById = &model.PolicyRule{
	Name:        "SDP: CanDeleteClient",
	Description: "Allows a client to delete their own profile.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteClientById.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck), // Client can delete self
}

var AllowUpdateClientImages = &model.PolicyRule{
	Name:        "SDP: CanUpdateClientImages",
	Description: "Allows a client to update their own profile images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateClientImages.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck), // Client can update self images (checks subject.id == resource.id)
}

var AllowDeleteClientImage = &model.PolicyRule{
	Name:        "SDP: CanDeleteClientImage",
	Description: "Allows a client to delete their own profile images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteClientImage.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck), // Client can delete self images (checks subject.id == resource.id)
}
