package policySeed

import (
	authModel "mynute-go/auth/model"
	coreModel "mynute-go/core/src/config/db/model"
)

// --- Reusable Condition Checks --- //

// Checks if subject is a Client AND their ID matches the resource's client_id
var ClientAccessCheck = authModel.ConditionNode{
	Description: "Client Access Check (Must be Client & Match Resource's Client)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNull",
				Description: "Subject must be a Client (no company affiliation)",
			},
		},
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:         "subject.id",
				Operator:          "Equals",
				ResourceAttribute: "resource.client_id", // Assumes context provides resource.client_id from the fetched resource
				Description:       "Subject ID must match the resource's client ID",
			},
		},
	},
}

// Checks if subject is a Client AND their ID matches the ID in the endpoint context (e.g., path param /clients/{id})
var ClientSelfAccessCheck = authModel.ConditionNode{
	Description: "Client Self Access Check (Own Profile/Resource)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNull",
				Description: "Subject must be a Client",
			},
		},
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:         "subject.id",
				Operator:          "Equals",
				ResourceAttribute: "resource.id", // Assumes context provides resource.id from path parameter matching subject's ID
				Description:       "Subject ID must match the resource ID being accessed",
			},
		},
	},
}

// Checks if subject is an Employee AND their ID matches the ID in the endpoint context (e.g., path param /employees/{id})
// Also double-checks company match for safety.
var EmployeeSelfAccessCheck = authModel.ConditionNode{
	Description: "Employee Self Access Check (Own Profile/Resource)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNotNull",
				Description: "Subject must belong to a company",
			},
		},
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:         "subject.id",
				Operator:          "Equals",
				ResourceAttribute: "resource.id", // Assumes context provides resource.id from path matching subject's ID
				Description:       "Subject ID must match the resource ID being accessed",
			},
		},
		{ // Belt-and-suspenders: Check company match too
			Leaf: &authModel.ConditionLeaf{
				Attribute:         "subject.company_id",
				Operator:          "Equals",
				ResourceAttribute: "resource.company_id", // Assumes context provides resource.company_id from the fetched resource
				Description:       "Subject company must match the resource's company",
			},
		},
	},
}

// Checks if subject is an Employee/Manager AND belongs to the same company as the resource
var CompanyMembershipAccessCheck = authModel.ConditionNode{
	Description: "Company Membership Check (Subject & Resource in Same Company)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNotNull",
				Description: "Subject must belong to a company",
			},
		},
		{
			Description: "Either company ID is at .id or .company_id",
			LogicType:   "OR",
			Children: []authModel.ConditionNode{
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.company_id",
						Operator:          "Equals",
						ResourceAttribute: "resource.company_id", // Assumes context provides resource.company_id from the fetched resource
						Description:       "Subject company must match the resource's company",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.company_id",
						Operator:          "Equals",
						ResourceAttribute: "resource.id", // Assumes context provides resource.id from the fetched resource
						Description:       "Subject company must match the resource's ID",
					},
				},
			},
		},
	},
}

// Checks if subject is the Owner of the resource's company
var CompanyOwnerCheck = authModel.ConditionNode{
	Description: "Allow if Subject is Owner of the Resource's Company",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{
			Leaf: &authModel.ConditionLeaf{
				Attribute: "subject.roles[*].id",
				Operator:  "Contains",
				Value:     authModel.JsonRawMessage(coreModel.SystemRoleOwner.ID),
			},
		},
		CompanyMembershipAccessCheck, // Re-use the company match check
	},
}

// Checks if subject is the General Manager of the resource's company
var CompanyGeneralManagerCheck = authModel.ConditionNode{
	Description: "Allow if Subject is General Manager of the Resource's Company",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{Leaf: &authModel.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: authModel.JsonRawMessage(coreModel.SystemRoleGeneralManager.ID)}},
		CompanyMembershipAccessCheck, // Re-use the company match check
	},
}

// Checks if subject is a Branch Manager of the resource's company
var CompanyBranchManagerCheck = authModel.ConditionNode{
	Description: "Allow if Subject is Branch Manager within the Resource's Company",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		{Leaf: &authModel.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: authModel.JsonRawMessage(coreModel.SystemRoleBranchManager.ID)}},
		CompanyMembershipAccessCheck, // Re-use the company match check
	},
}

// Checks if subject is a Branch Manager AND assigned to the specific branch of the resource
var CompanyBranchManagerAssignedBranchCheck = authModel.ConditionNode{
	Description: "Allow if Subject is Branch Manager assigned to the Resource's Branch",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		CompanyBranchManagerCheck, // Must be a BM in the correct company first
		{
			Description: "Branch Manager Assigned to Branch Check",
			LogicType:   "OR",
			Children: []authModel.ConditionNode{
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.branches",   // Assumes subject context has assigned branch IDs (e.g., [10, 25])
						Operator:          "Contains",           // Checks if the list contains the value
						ResourceAttribute: "resource.branch_id", // Assumes context provides resource.branch_id from the resource/path/body
						Description:       "Subject's assigned branches must include the resource's branch",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "path.branch_id", // Assumes context provides branch_id from the path parameter
						Description:       "Subject ID must match the path parameter branch_id",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "body.branch_id", // Assumes context provides branch_id from the body
						Description:       "Subject ID must match the body branch_id",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "query.branch_id", // Assumes context provides branch_id from the query parameter
						Description:       "Subject ID must match the query parameter branch_id",
					},
				},
			},
		},
	},
}

// Checks if subject is an Employee AND their ID matches the resource's employee_id
var CompanyEmployeeAssignedEmployeeCheck = authModel.ConditionNode{
	Description: "Allow if Subject is the Employee associated with the Resource",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		CompanyMembershipAccessCheck, // Must be in the same company
		{
			Description: "Employee ID Match Check",
			LogicType:   "OR",
			Children: []authModel.ConditionNode{
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "resource.employee_id", // Assumes context provides resource.employee_id from the resource
						Description:       "Subject ID must match the resource's employee ID",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "body.employee_id", // Assumes context provides employee_id from the body
						Description:       "Subject ID must match the body employee_id",
					},
				},
			},
		},
	},
}

// Reusable node for Admin Roles (Owner OR General Manager) within the resource's company
var CompanyAdminCheck = authModel.ConditionNode{
	Description: "Company Admin Access Check (Owner or General Manager)",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		CompanyOwnerCheck,
		CompanyGeneralManagerCheck,
	},
}

// Reusable node for Manager Roles (Owner OR General Manager OR Branch Manager) within the resource's company
// NOTE: This grants access if the user is *any* BM in the company. Use CompanyBranchManagerAssignedBranchCheck for specific branch access.
var CompanyManagerCheck = authModel.ConditionNode{
	Description: "Company Manager Access Check (Owner, GM, or BM in Company)",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		CompanyOwnerCheck,
		CompanyGeneralManagerCheck,
		CompanyBranchManagerCheck, // General check if they are a BM in that company
	},
}

var CompanyEmployeeHimself = authModel.ConditionNode{
	Description: "Company Employee Check (Employee accessing their own profile/resource)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		CompanyMembershipAccessCheck, // Must be in the same company
		{
			Description: "Employee id must be on resource, path, body, or query",
			LogicType:   "OR",
			Children: []authModel.ConditionNode{
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "path.employee_id", // Assumes context provides employee_id from the path parameter
						Description:       "Subject ID must match the path parameter employee_id",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "body.employee_id", // Assumes context provides employee_id from the body
						Description:       "Subject ID must match the body employee_id",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "query.employee_id", // Assumes context provides employee_id from the query parameter
						Description:       "Subject ID must match the query parameter employee_id",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "resource.employee_id", // Assumes context provides employee_id from the resource
						Description:       "The employee is accessing a resource that has himself assigned.",
					},
				},
				{
					Leaf: &authModel.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "resource.id", // Assumes context provides resource.id from the resource
						Description:       "The employee must be accessing himself as a resource",
					},
				},
			},
		},
	},
}

var CompanyAdminOrEmployeeHimselfCheck = authModel.ConditionNode{
	Description: "Company Admin or Employee Check (Owner, GM, or Employee)",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		CompanyAdminCheck,      // Owner or GM of the resource's company
		CompanyEmployeeHimself, // Employee can access their own profile/resource
	},
}

// Reusable node for ALL Company Internal Users (Owner, GM, BM, Employee) belonging to the resource's company
var CompanyInternalUserCheck = authModel.ConditionNode{
	Description: "Company Internal User Check (Any Role within Resource's Company)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		CompanyMembershipAccessCheck, // Base check: subject belongs to the endpoint's company
		// Additionally, ensure they have one of the internal roles (Owner, GM, BM, Employee)
		// This adds clarity but is often implicit if subject.role_id is always one of these for company members.
		{
			Description: "Role Check (Is Owner, GM, BM, or Employee)",
			LogicType:   "OR",
			Children: []authModel.ConditionNode{
				{Leaf: &authModel.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: authModel.JsonRawMessage(coreModel.SystemRoleOwner.ID)}},
				{Leaf: &authModel.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: authModel.JsonRawMessage(coreModel.SystemRoleGeneralManager.ID)}},
				{Leaf: &authModel.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: authModel.JsonRawMessage(coreModel.SystemRoleBranchManager.ID)}},
			},
		},
	},
}

var CompanyAdminOrAssignedBranchManagerCheck = authModel.ConditionNode{
	Description: "Company Admin or Assigned Branch Manager Check (Owner, GM, or BM in Company)",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		CompanyAdminCheck,                       // Owner or GM of the resource's company
		CompanyBranchManagerAssignedBranchCheck, // BM assigned to the resource's branch
	},
}

// --- Additional Composite Helpers for Common Access Patterns --- //

// CompanyManagerOrAssignedEmployeeCheck allows company admins (Owner/GM), assigned branch managers, or the employee themselves
var CompanyManagerOrAssignedEmployeeCheck = authModel.ConditionNode{
	Description: "Company Manager or Assigned Employee Check (Owner, GM, BM, or Employee Self)",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		CompanyOwnerCheck,
		CompanyGeneralManagerCheck,
		CompanyBranchManagerAssignedBranchCheck, // BM can manage in their assigned branch
		CompanyEmployeeAssignedEmployeeCheck,    // Employee can manage their own resources
	},
}

// CompanyMembershipAndManagerOrAssignedEmployeeCheck wraps CompanyManagerOrAssignedEmployeeCheck with company membership
var CompanyMembershipAndManagerOrAssignedEmployeeCheck = authModel.ConditionNode{
	Description: "Company User Check (Managers or Assigned Employee)",
	LogicType:   "AND",
	Children: []authModel.ConditionNode{
		CompanyMembershipAccessCheck,          // User must be in the same company
		CompanyManagerOrAssignedEmployeeCheck, // And have appropriate role/relation
	},
}

// ClientOrCompanyMembershipAndManagerOrAssignedEmployeeCheck allows client self-access OR company user access
var ClientOrCompanyMembershipAndManagerOrAssignedEmployeeCheck = authModel.ConditionNode{
	Description: "Client Self-Access OR Company User Access",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		ClientAccessCheck, // Client can access their own resources
		CompanyMembershipAndManagerOrAssignedEmployeeCheck, // Company users with appropriate permissions
	},
}

// EmployeeSelfOrInternalUserCheck allows employee self-access OR any internal company user
var EmployeeSelfOrInternalUserCheck = authModel.ConditionNode{
	Description: "Employee Self-View OR Any Internal Company User View",
	LogicType:   "OR",
	Children: []authModel.ConditionNode{
		EmployeeSelfAccessCheck,  // Can view self
		CompanyInternalUserCheck, // Any other member of the same company can view
	},
}
