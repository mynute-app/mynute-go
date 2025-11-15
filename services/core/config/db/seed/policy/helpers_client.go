package policySeed

import (
	"mynute-go/services/core/config/db/model"
)

// --- Reusable Condition Checks for CLIENT Policies --- //

// Checks if subject is a Client AND their ID matches the resource's client_id
var ClientAccessCheck = model.ConditionNode{
	Description: "Client Access Check (Must be Client & Match Resource's Client)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{
			Leaf: &model.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNull",
				Description: "Subject must be a Client (no company affiliation)",
			},
		},
		{
			Leaf: &model.ConditionLeaf{
				Attribute:         "subject.id",
				Operator:          "Equals",
				ResourceAttribute: "resource.client_id",
				Description:       "Subject ID must match the resource's client ID",
			},
		},
	},
}

// Checks if subject is a Client AND their ID matches the ID in the endpoint context
var ClientSelfAccessCheck = model.ConditionNode{
	Description: "Client Self Access Check (Own Profile/Resource)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{
			Leaf: &model.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNull",
				Description: "Subject must be a Client",
			},
		},
		{
			Leaf: &model.ConditionLeaf{
				Attribute:         "subject.id",
				Operator:          "Equals",
				ResourceAttribute: "resource.id",
				Description:       "Subject ID must match the resource ID being accessed",
			},
		},
	},
}
