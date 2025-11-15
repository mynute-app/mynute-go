package policySeed

import (
	"mynute-go/services/core/config/db/model"
)

// --- Reusable Condition Checks for TENANT (Company) Policies --- //

// Checks if subject is an Employee AND belongs to the same company as the resource
var CompanyMembershipAccessCheck = model.ConditionNode{
	Description: "Company Membership Check (Subject & Resource in Same Company)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{
			Leaf: &model.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNotNull",
				Description: "Subject must belong to a company",
			},
		},
		{
			Description: "Either company ID is at .id or .company_id",
			LogicType:   "OR",
			Children: []model.ConditionNode{
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.company_id",
						Operator:          "Equals",
						ResourceAttribute: "resource.company_id",
						Description:       "Subject company must match the resource's company",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.company_id",
						Operator:          "Equals",
						ResourceAttribute: "resource.id",
						Description:       "Subject company must match the resource's ID",
					},
				},
			},
		},
	},
}

// Checks if subject is an Employee AND their ID matches the ID in the endpoint context
var EmployeeSelfAccessCheck = model.ConditionNode{
	Description: "Employee Self Access Check (Own Profile/Resource)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{
			Leaf: &model.ConditionLeaf{
				Attribute:   "subject.company_id",
				Operator:    "IsNotNull",
				Description: "Subject must belong to a company",
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
		{
			Leaf: &model.ConditionLeaf{
				Attribute:         "subject.company_id",
				Operator:          "Equals",
				ResourceAttribute: "resource.company_id",
				Description:       "Subject company must match the resource's company",
			},
		},
	},
}

// Checks if subject is the Owner of the resource's company
var CompanyOwnerCheck = model.ConditionNode{
	Description: "Allow if Subject is Owner of the Resource's Company",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{
			Leaf: &model.ConditionLeaf{
				Attribute: "subject.roles[*].id",
				Operator:  "Contains",
				Value:     model.JsonRawMessage(model.SystemRoleOwner.ID),
			},
		},
		CompanyMembershipAccessCheck,
	},
}

// Checks if subject is the General Manager of the resource's company
var CompanyGeneralManagerCheck = model.ConditionNode{
	Description: "Allow if Subject is General Manager of the Resource's Company",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{Leaf: &model.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: model.JsonRawMessage(model.SystemRoleGeneralManager.ID)}},
		CompanyMembershipAccessCheck,
	},
}

// Checks if subject is a Branch Manager of the resource's company
var CompanyBranchManagerCheck = model.ConditionNode{
	Description: "Allow if Subject is Branch Manager in the Resource's Company",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		{Leaf: &model.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: model.JsonRawMessage(model.SystemRoleBranchManager.ID)}},
		CompanyMembershipAccessCheck,
	},
}

// Checks if subject is a Branch Manager AND assigned to the specific branch of the resource
var CompanyBranchManagerAssignedBranchCheck = model.ConditionNode{
	Description: "Allow if Subject is Branch Manager assigned to the Resource's Branch",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		CompanyBranchManagerCheck,
		{
			Description: "Branch Manager Assigned to Branch Check",
			LogicType:   "OR",
			Children: []model.ConditionNode{
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "resource.branch_id",
						Description:       "Subject's assigned branches must include the resource's branch",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "path.branch_id",
						Description:       "Subject ID must match the path parameter branch_id",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "body.branch_id",
						Description:       "Subject ID must match the body branch_id",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.branches",
						Operator:          "Contains",
						ResourceAttribute: "query.branch_id",
						Description:       "Subject ID must match the query parameter branch_id",
					},
				},
			},
		},
	},
}

// Checks if subject is an Employee AND their ID matches the resource's employee_id
var CompanyEmployeeAssignedEmployeeCheck = model.ConditionNode{
	Description: "Allow if Subject is the Employee associated with the Resource",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		CompanyMembershipAccessCheck,
		{
			Description: "Employee ID Match Check",
			LogicType:   "OR",
			Children: []model.ConditionNode{
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "resource.employee_id",
						Description:       "Subject ID must match the resource's employee ID",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "body.employee_id",
						Description:       "Subject ID must match the body employee_id",
					},
				},
			},
		},
	},
}

// Reusable node for Admin Roles (Owner OR General Manager)
var CompanyAdminCheck = model.ConditionNode{
	Description: "Company Admin Access Check (Owner or General Manager)",
	LogicType:   "OR",
	Children: []model.ConditionNode{
		CompanyOwnerCheck,
		CompanyGeneralManagerCheck,
	},
}

// Reusable node for Manager Roles (Owner OR General Manager OR Branch Manager)
var CompanyManagerCheck = model.ConditionNode{
	Description: "Company Manager Access Check (Owner, GM, or BM in Company)",
	LogicType:   "OR",
	Children: []model.ConditionNode{
		CompanyOwnerCheck,
		CompanyGeneralManagerCheck,
		CompanyBranchManagerCheck,
	},
}

// Employee accessing their own profile/resource
var CompanyEmployeeHimself = model.ConditionNode{
	Description: "Company Employee Check (Employee accessing their own profile/resource)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		CompanyMembershipAccessCheck,
		{
			Description: "Employee id must be on resource, path, body, or query",
			LogicType:   "OR",
			Children: []model.ConditionNode{
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "path.employee_id",
						Description:       "Subject ID must match the path parameter employee_id",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "body.employee_id",
						Description:       "Subject ID must match the body employee_id",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "query.employee_id",
						Description:       "Subject ID must match the query parameter employee_id",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "resource.employee_id",
						Description:       "The employee is accessing a resource that has himself assigned.",
					},
				},
				{
					Leaf: &model.ConditionLeaf{
						Attribute:         "subject.id",
						Operator:          "Equals",
						ResourceAttribute: "resource.id",
						Description:       "The employee must be accessing himself as a resource",
					},
				},
			},
		},
	},
}

// Company Admin or Employee accessing their own resource
var CompanyAdminOrEmployeeHimselfCheck = model.ConditionNode{
	Description: "Company Admin or Employee Check (Owner, GM, or Employee)",
	LogicType:   "OR",
	Children: []model.ConditionNode{
		CompanyAdminCheck,
		CompanyEmployeeHimself,
	},
}

// All Company Internal Users (Owner, GM, BM, Employee)
var CompanyInternalUserCheck = model.ConditionNode{
	Description: "Company Internal User Check (Any Role within Resource's Company)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		CompanyMembershipAccessCheck,
		{
			Description: "Role Check (Is Owner, GM, BM, or Employee)",
			LogicType:   "OR",
			Children: []model.ConditionNode{
				{Leaf: &model.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: model.JsonRawMessage(model.SystemRoleOwner.ID)}},
				{Leaf: &model.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: model.JsonRawMessage(model.SystemRoleGeneralManager.ID)}},
				{Leaf: &model.ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: model.JsonRawMessage(model.SystemRoleBranchManager.ID)}},
			},
		},
	},
}

// Company Admin or Assigned Branch Manager
var CompanyAdminOrAssignedBranchManagerCheck = model.ConditionNode{
	Description: "Company Admin or Assigned Branch Manager Check (Owner, GM, or BM in Company)",
	LogicType:   "OR",
	Children: []model.ConditionNode{
		CompanyAdminCheck,
		CompanyBranchManagerAssignedBranchCheck,
	},
}

// Company Manager or Assigned Employee
var CompanyManagerOrAssignedEmployeeCheck = model.ConditionNode{
	Description: "Company Manager or Assigned Employee Check (Owner, GM, BM, or Employee Self)",
	LogicType:   "OR",
	Children: []model.ConditionNode{
		CompanyOwnerCheck,
		CompanyGeneralManagerCheck,
		CompanyBranchManagerAssignedBranchCheck,
		CompanyEmployeeAssignedEmployeeCheck,
	},
}

// Company Membership with Manager or Assigned Employee
var CompanyMembershipAndManagerOrAssignedEmployeeCheck = model.ConditionNode{
	Description: "Company User Check (Managers or Assigned Employee)",
	LogicType:   "AND",
	Children: []model.ConditionNode{
		CompanyMembershipAccessCheck,
		CompanyManagerOrAssignedEmployeeCheck,
	},
}

// Employee Self or Internal User
var EmployeeSelfOrInternalUserCheck = model.ConditionNode{
	Description: "Employee Self-View OR Any Internal Company User View",
	LogicType:   "OR",
	Children: []model.ConditionNode{
		EmployeeSelfAccessCheck,
		CompanyInternalUserCheck,
	},
}
