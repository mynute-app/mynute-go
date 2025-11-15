package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- CLIENT POLICIES: Client profile operations ---

var ClientAllowGetClientByEmail = &model.ClientPolicy{
	Name:        "CP: CanViewClientByEmail",
	Description: "Allows a client to retrieve their own profile by email.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetClientByEmail.ID,
	Conditions: model.JsonRawMessage(model.ConditionNode{
		Description: "Allow only if the subject's email matches the email in the path.",
		LogicType:   "AND",
		Children: []model.ConditionNode{
			{Leaf: &model.ConditionLeaf{Attribute: "subject.company_id", Operator: "IsNull", Description: "Must be a Client"}},
			{Leaf: &model.ConditionLeaf{Attribute: "subject.email", Operator: "Equals", ResourceAttribute: "resource.email", Description: "Subject email must match email from path"}},
		},
	}),
}

var ClientAllowGetClientById = &model.ClientPolicy{
	Name:        "CP: CanViewClientById",
	Description: "Allows a client to view their own profile.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetClientById.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck),
}

var ClientAllowUpdateClientById = &model.ClientPolicy{
	Name:        "CP: CanUpdateClient",
	Description: "Allows a client to update their own profile.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateClientById.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck),
}

var ClientAllowDeleteClientById = &model.ClientPolicy{
	Name:        "CP: CanDeleteClient",
	Description: "Allows a client to delete their own profile.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteClientById.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck),
}

var ClientAllowUpdateClientImages = &model.ClientPolicy{
	Name:        "CP: CanUpdateClientImages",
	Description: "Allows a client to update their own profile images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateClientImages.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck),
}

var ClientAllowDeleteClientImage = &model.ClientPolicy{
	Name:        "CP: CanDeleteClientImage",
	Description: "Allows a client to delete their own profile images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteClientImage.ID,
	Conditions:  model.JsonRawMessage(ClientSelfAccessCheck),
}
