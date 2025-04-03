package model

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var AllowNilCompanyID = false
var AllowNilCreatedBy = false
var AllowNilResourceID = false

type PolicyRule struct {
	gorm.Model
	CompanyID           *uint           `json:"company_id"`
	CreatedByEmployeeID *uint           `json:"created_by_employee_id"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Effect              string          `json:"effect"`       // "Allow" / "Deny"
	EndPointID          uint            `json:"end_point_id"` // Link to EndPoint definition
	EndPoint            EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnDelete:CASCADE;" json:"endpoint"`
	Conditions          json.RawMessage `gorm:"type:jsonb" json:"conditions"`
	ResourceTable       string          `json:"resource_table"`
	ResourceKey         string          `json:"resource_key"`
	ResourceValueAt     string          `json:"resource_value_at"`
}

// --- ConditionLeaf (Represents a single atomic check) ---
type ConditionLeaf struct {
	Attr string `json:"attr"` // The primary attribute (e.g., "subject.role_id")
	Op   string `json:"op"`   // The comparison operator (e.g., "Equals", "IsNull", "Contains")

	// Use EITHER Value OR ValueSourceAttr for comparison. Omitempty ensures only one (or neither for ops like IsNull) appears in JSON.
	Value           json.RawMessage `json:"value,omitempty"`             // Static value to compare against
	ValueSourceAttr string          `json:"value_source_attr,omitempty"` // Other attribute's name to compare against

	Description string `json:"description,omitempty"` // Optional human-readable description of the check
}

// --- ConditionNode (Represents a logical grouping OR a leaf check) ---
type ConditionNode struct {
	// Description explains the purpose of this node (e.g., "Client Access Check", "Company Membership Check")
	Description string `json:"description,omitempty"`

	// --- Fields for Branch Nodes (Logical Operators: AND, OR, NOT) ---
	// LogicType specifies how to evaluate Children ("AND", "OR"). Required for branch nodes.
	LogicType string `json:"logic_type,omitempty"` // omitempty: Only present for branch nodes

	// Children contains nested nodes to be evaluated. Required for branch nodes.
	Children []ConditionNode `json:"children,omitempty"` // omitempty: Only present for branch nodes

	// --- Field for Leaf Nodes (Actual Condition Check) ---
	// Leaf points to the actual condition details. Required for leaf nodes.
	Leaf *ConditionLeaf `json:"leaf,omitempty"` // omitempty: Only present for leaf nodes
}

// GetConditionsNode parses the stored JSON conditions into the ConditionNode struct
func (p *PolicyRule) GetConditionsNode() (ConditionNode, error) {
	var node ConditionNode
	if len(p.Conditions) == 0 || string(p.Conditions) == "null" {
		// If conditions are empty or null, return an error or a default node based on your logic
		return node, fmt.Errorf("policy rule %d has missing or null conditions", p.ID) // Or return (node, nil) if allowed
	}
	err := json.Unmarshal(p.Conditions, &node)
	if err != nil {
		return node, fmt.Errorf("failed to unmarshal conditions JSON for policy rule %d: %w", p.ID, err)
	}
	// Basic validation: Ensure a node is either a valid branch or a valid leaf
	if node.Leaf == nil && (node.LogicType == "" || len(node.Children) == 0) {
		// It's a branch, but has no children or logic type - indicates malformed JSON perhaps
		return node, fmt.Errorf("node `%s` is neither a valid branch nor a valid leaf", node.Description)
	}
	if node.Leaf != nil && (node.LogicType != "" || len(node.Children) > 0) {
		// It's a leaf, but also has branch properties - indicates malformed JSON perhaps
		return node, fmt.Errorf("policy rule %d condition node %s is incorrectly structured as both leaf and branch", p.ID, node.Description)
	}
	return node, nil
}

func (PolicyRule) TableName() string {
	return "policy_rules"
}

func (PolicyRule) Indexes() map[string]string {
	return map[string]string{
		"idx_policy_company_endpoint": "CREATE INDEX idx_policy_company_endpoint ON policy_rules (company_id, end_point_id)",
	}
}

// --- Helper Functions ---

// JsonRawMessage simplifies creating JSON.RawMessage from any value.
func JsonRawMessage(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("JsonRawMessage failed: %v", err))
	}
	return json.RawMessage(data)
}

func init_policy_array() []*PolicyRule { // --- Reusable Condition Checks --- //

	// Checks if subject is a Client AND their ID matches the resource's client_id
	var client_access_check = ConditionNode{
		Description: "Client Access Check (Must be Client & Match Resource's Client)",
		LogicType:   "AND",
		Children: []ConditionNode{
			{
				Leaf: &ConditionLeaf{
					Attr:        "subject.company_id",
					Op:          "IsNull",
					Description: "Subject must be a Client (no company affiliation)",
				},
			},
			{
				Leaf: &ConditionLeaf{
					Attr:            "subject.id",
					Op:              "Equals",
					ValueSourceAttr: "resource.client_id", // Assumes context provides resource.client_id from the fetched resource
					Description:     "Subject ID must match the resource's client ID",
				},
			},
		},
	}

	// Checks if subject is a Client AND their ID matches the ID in the endpoint context (e.g., path param /clients/{id})
	var client_self_access_check = ConditionNode{
		Description: "Client Self Access Check (Own Profile/Resource)",
		LogicType:   "AND",
		Children: []ConditionNode{
			{
				Leaf: &ConditionLeaf{
					Attr:        "subject.company_id",
					Op:          "IsNull",
					Description: "Subject must be a Client",
				},
			},
			{
				Leaf: &ConditionLeaf{
					Attr:            "subject.id",
					Op:              "Equals",
					ValueSourceAttr: "resource.id", // Assumes context provides resource.id from path parameter matching subject's ID
					Description:     "Subject ID must match the resource ID being accessed",
				},
			},
		},
	}

	// Checks if subject is an Employee AND their ID matches the ID in the endpoint context (e.g., path param /employees/{id})
	// Also double-checks company match for safety.
	var employee_self_access_check = ConditionNode{
		Description: "Employee Self Access Check (Own Profile/Resource)",
		LogicType:   "AND",
		Children: []ConditionNode{
			{
				Leaf: &ConditionLeaf{
					Attr:        "subject.company_id",
					Op:          "IsNotNull",
					Description: "Subject must belong to a company",
				},
			},
			{
				Leaf: &ConditionLeaf{
					Attr:            "subject.id",
					Op:              "Equals",
					ValueSourceAttr: "resource.id", // Assumes context provides resource.id from path matching subject's ID
					Description:     "Subject ID must match the resource ID being accessed",
				},
			},
			{ // Belt-and-suspenders: Check company match too
				Leaf: &ConditionLeaf{
					Attr:            "subject.company_id",
					Op:              "Equals",
					ValueSourceAttr: "resource.company_id", // Assumes context provides resource.company_id from the fetched resource
					Description:     "Subject company must match the resource's company",
				},
			},
		},
	}

	// Checks if subject is an Employee/Manager AND belongs to the same company as the resource
	var company_membership_access_check = ConditionNode{
		Description: "Company Membership Check (Subject & Resource in Same Company)",
		LogicType:   "AND",
		Children: []ConditionNode{
			{
				Leaf: &ConditionLeaf{
					Attr:        "subject.company_id",
					Op:          "IsNotNull",
					Description: "Subject must belong to a company",
				},
			},
			{
				Leaf: &ConditionLeaf{
					Attr:            "subject.company_id",
					Op:              "Equals",
					ValueSourceAttr: "resource.company_id", // Assumes context provides resource.company_id from resource/path/body
					Description:     "Subject company must match the resource's company",
				},
			},
		},
	}

	// Checks if subject is the Owner of the resource's company
	var company_owner_check = ConditionNode{
		Description: "Allow if Subject is Owner of the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleOwner.ID)}},
			company_membership_access_check, // Re-use the company match check
		},
	}

	// Checks if subject is the General Manager of the resource's company
	var company_general_manager_check = ConditionNode{
		Description: "Allow if Subject is General Manager of the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleGeneralManager.ID)}},
			company_membership_access_check, // Re-use the company match check
		},
	}

	// Checks if subject is a Branch Manager of the resource's company
	var company_branch_manager_check = ConditionNode{
		Description: "Allow if Subject is Branch Manager within the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleBranchManager.ID)}},
			company_membership_access_check, // Re-use the company match check
		},
	}

	// Checks if subject is a Branch Manager AND assigned to the specific branch of the resource
	var company_branch_manager_assigned_branch_check = ConditionNode{
		Description: "Allow if Subject is Branch Manager assigned to the Resource's Branch",
		LogicType:   "AND",
		Children: []ConditionNode{
			company_branch_manager_check, // Must be a BM in the correct company first
			{
				Leaf: &ConditionLeaf{
					Attr:            "subject.branches",   // Assumes subject context has assigned branch IDs (e.g., [10, 25])
					Op:              "Contains",           // Checks if the list contains the value
					ValueSourceAttr: "resource.branch_id", // Assumes context provides resource.branch_id from the resource/path/body
					Description:     "Subject's assigned branches must include the resource's branch",
				},
			},
		},
	}

	// Checks if subject is an Employee AND their ID matches the resource's employee_id
	var company_employee_assigned_employee_check = ConditionNode{
		Description: "Allow if Subject is the Employee associated with the Resource",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleEmployee.ID)}},
			company_membership_access_check, // Must be in the same company
			{
				Leaf: &ConditionLeaf{
					Attr:            "subject.id",
					Op:              "Equals",
					ValueSourceAttr: "resource.employee_id", // Assumes context provides resource.employee_id from the resource
					Description:     "Subject ID must match the resource's employee ID",
				},
			},
		},
	}

	// Reusable node for Admin Roles (Owner OR General Manager) within the resource's company
	var company_admin_check = ConditionNode{
		Description: "Company Admin Access Check (Owner or General Manager)",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
		},
	}

	// Reusable node for Manager Roles (Owner OR General Manager OR Branch Manager) within the resource's company
	// NOTE: This grants access if the user is *any* BM in the company. Use company_branch_manager_assigned_branch_check for specific branch access.
	var company_manager_check = ConditionNode{
		Description: "Company Manager Access Check (Owner, GM, or BM in Company)",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
			company_branch_manager_check, // General check if they are a BM in that company
		},
	}

	// Reusable node for ALL Company Internal Users (Owner, GM, BM, Employee) belonging to the resource's company
	var company_internal_user_check = ConditionNode{
		Description: "Company Internal User Check (Any Role within Resource's Company)",
		LogicType:   "AND",
		Children: []ConditionNode{
			company_membership_access_check, // Base check: subject belongs to the endpoint's company
			// Additionally, ensure they have one of the internal roles (Owner, GM, BM, Employee)
			// This adds clarity but is often implicit if subject.role_id is always one of these for company members.
			{
				Description: "Role Check (Is Owner, GM, BM, or Employee)",
				LogicType:   "OR",
				Children: []ConditionNode{
					{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleOwner.ID)}},
					{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleGeneralManager.ID)}},
					{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleBranchManager.ID)}},
					{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: JsonRawMessage(SystemRoleEmployee.ID)}},
				},
			},
		},
	}

	// --- Policy Definitions with Resource Context --- //

	// Policy: Allow Create appointment
	var AllowCreateAppointment = &PolicyRule{
		Name:        "SDP: CanCreateAppointment",
		Description: "Allows clients to create appointments for themselves, or company users based on role/relation.",
		Effect:      "Allow",
		EndPointID:  CreateAppointment.ID,
		// Resource context often refers to the container/parent for creation. Here, likely the branch.
		ResourceTable:   "branches",  // Resource being referenced for context/permissions
		ResourceKey:     "branch_id", // Key in the request body identifying the branch
		ResourceValueAt: "body",      // Where to find the value for ResourceKey
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Client Creation OR Company User Creation",
			LogicType:   "OR",
			Children: []ConditionNode{
				// Client creating for themselves. Assumes 'client_id' is in the body and resource.client_id gets populated.
				{
					Description: "Client Self-Creation",
					LogicType:   "AND",
					Children: []ConditionNode{
						{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "IsNull", Description: "Must be a Client"}},
						{Leaf: &ConditionLeaf{Attr: "subject.id", Op: "Equals", ValueSourceAttr: "resource.client_id", Description: "Client ID in body must match Subject ID"}},
					},
				},
				// Company User creating. Assumes 'branch_id' and maybe 'employee_id' are in the body.
				// 'company_id' will be derived from the branch_id for context checks.
				{
					Description: "Company User Creation Check",
					LogicType:   "AND",
					Children: []ConditionNode{
						company_membership_access_check, // Checks subject.company_id == resource.company_id (derived from branch_id)
						{
							Description: "Role/Relation Check (Managers for Branch, Assigned Employee)",
							LogicType:   "OR",
							Children: []ConditionNode{
								company_owner_check,                          // Owner can create in any branch of their company
								company_general_manager_check,                // GM can create in any branch of their company
								company_branch_manager_assigned_branch_check, // BM can create *for their assigned branch* (checks resource.branch_id)
								company_employee_assigned_employee_check,     // Employee can create *for themselves* (checks resource.employee_id from body)
							},
						},
					},
				},
			},
		}),
	}

	// Policy: Allow GET appointment by ID.
	var AllowGetAppointmentByID = &PolicyRule{
		Name:            "SDP: CanViewAppointment",
		Description:     "Allows clients to view own appointments, or company users based on role/relation.",
		Effect:          "Allow",
		EndPointID:      GetAppointmentByID.ID,
		ResourceTable:   "appointments", // The resource being acted upon
		ResourceKey:     "id",           // The identifier in the path
		ResourceValueAt: "path",         // Found in the URL path
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Client Access OR Company User Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				// Client can view if appointment's client_id matches subject.id
				client_access_check, // Needs resource.client_id populated from the fetched appointment
				// Company user can view if they are in the same company and meet role/relation criteria
				{
					Description: "Company User View Check",
					LogicType:   "AND",
					Children: []ConditionNode{
						company_membership_access_check, // Needs resource.company_id from appointment
						{
							Description: "Role/Relation Check",
							LogicType:   "OR",
							Children: []ConditionNode{
								company_owner_check,                          // Needs resource.company_id
								company_general_manager_check,                // Needs resource.company_id
								company_branch_manager_assigned_branch_check, // Needs resource.company_id and resource.branch_id
								company_employee_assigned_employee_check,     // Needs resource.company_id and resource.employee_id
							},
						},
					},
				},
			},
		}),
	}

	// Policy: Allow UPDATE appointment by ID.
	var AllowUpdateAppointmentByID = &PolicyRule{
		Name:            "SDP: CanUpdateAppointment",
		Description:     "Allows clients to update own appointments, or company managers/assigned employees.",
		Effect:          "Allow",
		EndPointID:      UpdateAppointmentByID.ID,
		ResourceTable:   "appointments",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Client Self-Update OR Company User Update",
			LogicType:   "OR",
			Children: []ConditionNode{
				client_access_check, // Client can update if it's their appointment
				{
					Description: "Company User Update Check",
					LogicType:   "AND",
					Children: []ConditionNode{
						company_membership_access_check, // User in same company as appointment
						{
							Description: "Role/Relation Check (Managers or Assigned Employee)",
							LogicType:   "OR",
							Children: []ConditionNode{
								company_owner_check,
								company_general_manager_check,
								company_branch_manager_assigned_branch_check, // BM can update appointments in their branch
								company_employee_assigned_employee_check,     // Employee can update their own appointments
							},
						},
					},
				},
			},
		}),
	}

	// Policy: Allow DELETE appointment by ID.
	var AllowDeleteAppointmentByID = &PolicyRule{
		Name:            "SDP: CanDeleteAppointment",
		Description:     "Allows company managers or assigned employees to delete appointments. (Clients typically cannot delete).",
		Effect:          "Allow",
		EndPointID:      DeleteAppointmentByID.ID,
		ResourceTable:   "appointments",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Company User Delete Check",
			LogicType:   "AND",
			Children: []ConditionNode{
				company_membership_access_check, // User in same company as appointment
				{
					Description: "Role/Relation Check (Managers or Assigned Employee)",
					LogicType:   "OR",
					Children: []ConditionNode{
						company_owner_check,
						company_general_manager_check,
						company_branch_manager_assigned_branch_check, // BM can delete appointments in their branch
						company_employee_assigned_employee_check,     // Employee can delete their own appointments (?) Maybe not - adjust if needed. Let's assume managers only.
					},
				},
			},
		}),
	}

	// --- Branch Policies ---

	var AllowCreateBranch = &PolicyRule{
		Name:        "SDP: CanCreateBranch",
		Description: "Allows company Owner or General Manager to create branches.",
		Effect:      "Allow",
		EndPointID:  CreateBranch.ID,
		// Context is the company under which the branch is created
		ResourceTable:   "companies",
		ResourceKey:     "company_id", // Assuming company_id is provided in the request body
		ResourceValueAt: "body",
		// Condition checks against resource.company_id derived from the body
		Conditions: JsonRawMessage(company_admin_check), // Owner or GM of the target company
	}

	var AllowGetBranchById = &PolicyRule{
		Name:            "SDP: CanViewBranchById",
		Description:     "Allows any user belonging to the same company to view branch details by ID.",
		Effect:          "Allow",
		EndPointID:      GetBranchById.ID,
		ResourceTable:   "branches",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company can view
	}

	var AllowGetBranchByName = &PolicyRule{
		Name:            "SDP: CanViewBranchByName",
		Description:     "Allows any user belonging to the same company to view branch details by name.",
		Effect:          "Allow",
		EndPointID:      GetBranchByName.ID,
		ResourceTable:   "branches", // The lookup happens on branches table
		ResourceKey:     "name",     // The key used for lookup in the path
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company can view
	}

	var AllowUpdateBranchById = &PolicyRule{
		Name:            "SDP: CanUpdateBranch",
		Description:     "Allows company Owner, General Manager, or assigned Branch Manager to update branches.",
		Effect:          "Allow",
		EndPointID:      UpdateBranchById.ID,
		ResourceTable:   "branches",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Update Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner or GM can update any branch in their company
				company_branch_manager_assigned_branch_check, // BM can update their own assigned branch (needs resource.branch_id from path->resource)
			},
		}),
	}

	var AllowDeleteBranchById = &PolicyRule{
		Name:            "SDP: CanDeleteBranch",
		Description:     "Allows company Owner or General Manager to delete branches.",
		Effect:          "Allow",
		EndPointID:      DeleteBranchById.ID,
		ResourceTable:   "branches",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(company_admin_check), // Only Owner or GM
	}

	var AllowGetEmployeeServicesByBranchId = &PolicyRule{
		Name:            "SDP: CanViewEmployeeServicesInBranch",
		Description:     "Allows company members to view employee services within a branch.",
		Effect:          "Allow",
		EndPointID:      GetEmployeeServicesByBranchId.ID,
		ResourceTable:   "branches", // Primary resource context is the branch
		ResourceKey:     "id",       // Branch ID from the path
		ResourceValueAt: "path",
		// Context needs resource.company_id derived from the branch ID
		Conditions: JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company
	}

	var AllowAddServiceToBranch = &PolicyRule{
		Name:            "SDP: CanAddServiceToBranch",
		Description:     "Allows company managers (Owner, GM, relevant BM) to add services to a branch.",
		Effect:          "Allow",
		EndPointID:      AddServiceToBranch.ID,
		ResourceTable:   "branches", // Action is on the branch
		ResourceKey:     "id",       // Branch ID from the path
		ResourceValueAt: "path",
		// Service ID likely comes from body/query
		// Context needs resource.company_id and resource.branch_id from the branch resource
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can manage services in any branch
				company_branch_manager_assigned_branch_check, // BM can manage services in their own branch
			},
		}),
	}

	var AllowRemoveServiceFromBranch = &PolicyRule{
		Name:            "SDP: CanRemoveServiceFromBranch",
		Description:     "Allows company managers (Owner, GM, relevant BM) to remove services from a branch.",
		Effect:          "Allow",
		EndPointID:      RemoveServiceFromBranch.ID,
		ResourceTable:   "branches", // Action is on the branch
		ResourceKey:     "id",       // Branch ID from the path
		ResourceValueAt: "path",
		// Service ID likely comes from body/query/path
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can manage services in any branch
				company_branch_manager_assigned_branch_check, // BM can manage services in their own branch
			},
		}),
	}

	// --- Client Policies ---

	var AllowGetClientByEmail = &PolicyRule{
		Name:            "SDP: CanViewClientByEmail",
		Description:     "Allows a client to retrieve their own profile by email.",
		Effect:          "Allow",
		EndPointID:      GetClientByEmail.ID,
		ResourceTable:   "clients", // Lookup is on clients table
		ResourceKey:     "email",   // Key used for lookup from path
		ResourceValueAt: "path",
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow only if the subject's email matches the email in the path.",
			LogicType:   "AND",
			Children: []ConditionNode{
				{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "IsNull", Description: "Must be a Client"}},                                                       // Ensure subject is a client
				{Leaf: &ConditionLeaf{Attr: "subject.email", Op: "Equals", ValueSourceAttr: "resource.email", Description: "Subject email must match email from path"}}, // Assumes context has resource.email from path
			},
		}),
	}

	var AllowUpdateClientById = &PolicyRule{
		Name:            "SDP: CanUpdateClient",
		Description:     "Allows a client to update their own profile.",
		Effect:          "Allow",
		EndPointID:      UpdateClientById.ID,
		ResourceTable:   "clients",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(client_self_access_check), // Client can update self (checks subject.id == resource.id)
	}

	var AllowDeleteClientById = &PolicyRule{
		Name:            "SDP: CanDeleteClient",
		Description:     "Allows a client to delete their own profile.",
		Effect:          "Allow",
		EndPointID:      DeleteClientById.ID,
		ResourceTable:   "clients",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(client_self_access_check), // Client can delete self
	}

	// --- Company Policies ---

	var AllowGetCompanyById = &PolicyRule{
		Name:            "SDP: CanViewCompanyById",
		Description:     "Allows any member (employee/manager) of the company to view its details.",
		Effect:          "Allow",
		EndPointID:      GetCompanyById.ID,
		ResourceTable:   "companies",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Checks subject.company_id equals resource.company_id (derived from the path ID)
		Conditions: JsonRawMessage(company_membership_access_check),
	}

	var AllowUpdateCompanyById = &PolicyRule{
		Name:            "SDP: CanUpdateCompany",
		Description:     "Allows the company Owner or General Manager to update company details.",
		Effect:          "Allow",
		EndPointID:      UpdateCompanyById.ID,
		ResourceTable:   "companies",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(company_admin_check), // Only Owner or GM of this company
	}

	var AllowDeleteCompanyById = &PolicyRule{
		Name:            "SDP: CanDeleteCompany",
		Description:     "Allows ONLY the company Owner to delete the company.",
		Effect:          "Allow",
		EndPointID:      DeleteCompanyById.ID,
		ResourceTable:   "companies",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions:      JsonRawMessage(company_owner_check), // Only Owner of this company
	}

	// --- Employee Policies ---

	var AllowCreateEmployee = &PolicyRule{
		Name:        "SDP: CanCreateEmployee",
		Description: "Allows company Owner, GM, or BM to create employees (BM restricted to their branches implicitly if data includes branch).",
		Effect:      "Allow",
		EndPointID:  CreateEmployee.ID,
		// Context is the company they are being created under
		ResourceTable:   "companies",
		ResourceKey:     "company_id", // Assume company_id is in the request body or derived
		ResourceValueAt: "body",       // Or "context" if derived from subject/environment
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Branch Manager Creation Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check, // Owner/GM can create employees in their company
				// BM Check: If resource.branch_id is provided in body and needs checking:
				company_branch_manager_assigned_branch_check, // Ensures BM is assigned to the branch the employee might be added to (if resource.branch_id is provided/checked)
				// If a BM can create *any* employee in the company (even unassigned):
				// company_branch_manager_check,
			},
		}),
	}

	var AllowGetEmployeeById = &PolicyRule{
		Name:            "SDP: CanViewEmployeeById",
		Description:     "Allows employee to view self, or any internal user of the same company to view other employees.",
		Effect:          "Allow",
		EndPointID:      GetEmployeeById.ID,
		ResourceTable:   "employees",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Employee Self-View OR Any Internal Company User View",
			LogicType:   "OR",
			Children: []ConditionNode{
				employee_self_access_check,  // Can view self (checks subject.id == resource.id)
				company_internal_user_check, // Any other member of the same company can view (checks subject.company_id == resource.company_id)
			},
		}),
	}

	var AllowGetEmployeeByEmail = &PolicyRule{
		Name:            "SDP: CanViewEmployeeByEmail",
		Description:     "Allows company members to find employees within the same company by email.",
		Effect:          "Allow",
		EndPointID:      GetEmployeeByEmail.ID,
		ResourceTable:   "employees", // Lookup on employees table
		ResourceKey:     "email",     // Key from path
		ResourceValueAt: "path",
		// Requires resource.company_id populated from the employee found by email
		Conditions: JsonRawMessage(company_internal_user_check), // Subject must be internal user of the found employee's company
	}

	var AllowUpdateEmployeeById = &PolicyRule{
		Name:            "SDP: CanUpdateEmployee",
		Description:     "Allows employee to update self, or company managers (Owner, GM, BM) to update employees.",
		Effect:          "Allow",
		EndPointID:      UpdateEmployeeById.ID,
		ResourceTable:   "employees",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Employee Self-Update OR Manager Update",
			LogicType:   "OR",
			Children: []ConditionNode{
				employee_self_access_check, // Employee can update own profile
				// Any manager (Owner, GM, BM) can update employees in the same company.
				// Note: BM is not restricted to only employees in their assigned branch here.
				// If stricter control is needed, replace company_manager_check with a more complex OR:
				// OR (company_admin_check, specific_bm_check_for_employee_branch )
				company_manager_check,
			},
		}),
	}

	var AllowDeleteEmployeeById = &PolicyRule{
		Name:            "SDP: CanDeleteEmployee",
		Description:     "Allows company managers (Owner, GM, BM) to delete employees.",
		Effect:          "Allow",
		EndPointID:      DeleteEmployeeById.ID,
		ResourceTable:   "employees",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Same scoping comment as Update: BM can delete any employee in the company currently.
		Conditions: JsonRawMessage(company_manager_check), // Owner, GM, BM can delete
	}

	var AllowAddServiceToEmployee = &PolicyRule{
		Name:            "SDP: CanAddServiceToEmployee",
		Description:     "Allows company managers (Owner, GM, BM) to assign services to employees.",
		Effect:          "Allow",
		EndPointID:      AddServiceToEmployee.ID,
		ResourceTable:   "employees", // Primary action context is the employee
		ResourceKey:     "id",        // Employee ID from path
		ResourceValueAt: "path",
		// Service ID likely from body/query
		// Requires resource.company_id from the employee resource
		Conditions: JsonRawMessage(company_manager_check), // Manager of the employee's company
	}

	var AllowRemoveServiceFromEmployee = &PolicyRule{
		Name:            "SDP: CanRemoveServiceFromEmployee",
		Description:     "Allows company managers (Owner, GM, BM) to remove services from employees.",
		Effect:          "Allow",
		EndPointID:      RemoveServiceFromEmployee.ID,
		ResourceTable:   "employees", // Primary action context is the employee
		ResourceKey:     "id",        // Employee ID from path
		ResourceValueAt: "path",
		// Service ID likely from body/query/path
		// Requires resource.company_id from the employee resource
		Conditions: JsonRawMessage(company_manager_check), // Manager of the employee's company
	}

	var AllowAddBranchToEmployee = &PolicyRule{
		Name:            "SDP: CanAddBranchToEmployee",
		Description:     "Allows company managers (Owner, GM, BM) to assign employees to branches (respecting BM scope).",
		Effect:          "Allow",
		EndPointID:      AddBranchToEmployee.ID,
		ResourceTable:   "employees", // Primary action context is the employee
		ResourceKey:     "id",        // Employee ID from path
		ResourceValueAt: "path",
		// Branch ID comes from path/body, used for resource.branch_id in BM check
		// Requires resource.company_id from the employee resource
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Assignment Check",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check, // Owner/GM can assign any branch in company
				// BM can only assign employees TO a branch THEY MANAGE
				company_branch_manager_assigned_branch_check, // Checks subject.branches CONTAINS resource.branch_id (branch being added)
			},
		}),
	}

	var AllowRemoveBranchFromEmployee = &PolicyRule{
		Name:            "SDP: CanRemoveBranchFromEmployee",
		Description:     "Allows company managers (Owner, GM, BM) to remove employees from branches (respecting BM scope).",
		Effect:          "Allow",
		EndPointID:      RemoveBranchFromEmployee.ID,
		ResourceTable:   "employees", // Primary action context is the employee
		ResourceKey:     "id",        // Employee ID from path
		ResourceValueAt: "path",
		// Branch ID comes from path/body, used for resource.branch_id in BM check
		// Requires resource.company_id from the employee resource
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Assignment Check",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check, // Owner/GM can unassign from any branch
				// BM can only unassign employees FROM a branch THEY MANAGE
				company_branch_manager_assigned_branch_check, // Checks subject.branches CONTAINS resource.branch_id (branch being removed)
			},
		}),
	}

	// --- Holiday Policies ---

	var AllowCreateHoliday = &PolicyRule{
		Name:        "SDP: CanCreateHoliday",
		Description: "Allows company managers (Owner, GM, BM) to create holidays for the company/branch.",
		Effect:      "Allow",
		EndPointID:  CreateHoliday.ID,
		// Context depends on whether holidays are company-wide or branch-specific
		ResourceTable:   "companies",  // Assuming company-wide context needed
		ResourceKey:     "company_id", // Assume company_id provided or derived
		ResourceValueAt: "body",       // Or "context"
		// If branch-specific is possible via BM:
		// ResourceTable: "branches", ResourceKey: "branch_id", ResourceValueAt: "body"
		Conditions: JsonRawMessage(company_manager_check), // Any manager in the relevant company context. Add branch check if needed.
	}

	var AllowGetHolidayById = &PolicyRule{
		Name:            "SDP: CanViewHolidayById",
		Description:     "Allows company members to view holiday details.",
		Effect:          "Allow",
		EndPointID:      GetHolidayById.ID,
		ResourceTable:   "holidays",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the holiday resource
		Conditions: JsonRawMessage(company_internal_user_check), // Any internal user of the holiday's company
	}

	var AllowUpdateHolidayById = &PolicyRule{
		Name:            "SDP: CanUpdateHoliday",
		Description:     "Allows company managers (Owner, GM, BM) to update holidays.",
		Effect:          "Allow",
		EndPointID:      UpdateHolidayById.ID,
		ResourceTable:   "holidays",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the holiday resource
		Conditions: JsonRawMessage(company_manager_check), // Any manager of the holiday's company
	}

	var AllowDeleteHolidayById = &PolicyRule{
		Name:            "SDP: CanDeleteHoliday",
		Description:     "Allows company managers (Owner, GM, BM) to delete holidays.",
		Effect:          "Allow",
		EndPointID:      DeleteHolidayById.ID,
		ResourceTable:   "holidays",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the holiday resource
		Conditions: JsonRawMessage(company_manager_check), // Any manager of the holiday's company
	}

	// --- Sector Policies --- (Assuming Company Specific)

	var AllowCreateSector = &PolicyRule{
		Name:            "SDP: CanCreateSector",
		Description:     "Allows company Owner or General Manager to create sectors.",
		Effect:          "Allow",
		EndPointID:      CreateSector.ID,
		ResourceTable:   "companies",  // Context is the company
		ResourceKey:     "company_id", // Assume provided or derived
		ResourceValueAt: "body",       // Or "context"
		// Requires resource.company_id derived from request/context
		Conditions: JsonRawMessage(company_admin_check), // Only Owner/GM of that company
	}

	// GetSectorById/Name assumed Public - No Policy Needed

	var AllowUpdateSectorById = &PolicyRule{
		Name:            "SDP: CanUpdateSector",
		Description:     "Allows company Owner or General Manager to update sectors.",
		Effect:          "Allow",
		EndPointID:      UpdateSectorById.ID,
		ResourceTable:   "sectors",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the sector resource
		Conditions: JsonRawMessage(company_admin_check), // Only Owner/GM of the sector's company
	}

	var AllowDeleteSectorById = &PolicyRule{
		Name:            "SDP: CanDeleteSector",
		Description:     "Allows company Owner or General Manager to delete sectors.",
		Effect:          "Allow",
		EndPointID:      DeleteSectorById.ID,
		ResourceTable:   "sectors",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the sector resource
		Conditions: JsonRawMessage(company_admin_check), // Only Owner/GM of the sector's company
	}

	// --- Service Policies --- (Assuming Company Specific)

	var AllowCreateService = &PolicyRule{
		Name:            "SDP: CanCreateService",
		Description:     "Allows company managers (Owner, GM, BM) to create services.",
		Effect:          "Allow",
		EndPointID:      CreateService.ID,
		ResourceTable:   "companies",  // Context is the company
		ResourceKey:     "company_id", // Assume provided or derived
		ResourceValueAt: "body",       // Or "context"
		// Requires resource.company_id derived from request/context
		Conditions: JsonRawMessage(company_manager_check), // Any manager of the company context
	}

	var AllowGetServiceById = &PolicyRule{
		Name:            "SDP: CanViewServiceById",
		Description:     "Allows company members to view service details.",
		Effect:          "Allow",
		EndPointID:      GetServiceById.ID,
		ResourceTable:   "services",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the service resource
		Conditions: JsonRawMessage(company_internal_user_check), // Any internal user of the service's company
	}

	// GetServiceByName assumed Public - No Policy Needed

	var AllowUpdateServiceById = &PolicyRule{
		Name:            "SDP: CanUpdateService",
		Description:     "Allows company managers (Owner, GM, BM) to update services.",
		Effect:          "Allow",
		EndPointID:      UpdateServiceById.ID,
		ResourceTable:   "services",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the service resource
		Conditions: JsonRawMessage(company_manager_check), // Any manager of the service's company
	}

	var AllowDeleteServiceById = &PolicyRule{
		Name:            "SDP: CanDeleteService",
		Description:     "Allows company managers (Owner, GM, BM) to delete services.",
		Effect:          "Allow",
		EndPointID:      DeleteServiceById.ID,
		ResourceTable:   "services",
		ResourceKey:     "id",
		ResourceValueAt: "path",
		// Requires resource.company_id from the service resource
		Conditions: JsonRawMessage(company_manager_check), // Any manager of the service's company
	}

	// --- Combined Policies List ---
	var Policies = []*PolicyRule{
		// Appointments
		AllowGetAppointmentByID,
		AllowCreateAppointment,
		AllowUpdateAppointmentByID,
		AllowDeleteAppointmentByID,

		// Branches
		AllowCreateBranch,
		AllowGetBranchById,
		AllowGetBranchByName,
		AllowUpdateBranchById,
		AllowDeleteBranchById,
		AllowGetEmployeeServicesByBranchId,
		AllowAddServiceToBranch,
		AllowRemoveServiceFromBranch,

		// Clients (Self-Management focused)
		AllowGetClientByEmail,
		AllowUpdateClientById,
		AllowDeleteClientById,

		// Company
		AllowGetCompanyById,
		AllowUpdateCompanyById,
		AllowDeleteCompanyById,

		// Employees
		AllowCreateEmployee,
		AllowGetEmployeeById,
		AllowGetEmployeeByEmail,
		AllowUpdateEmployeeById,
		AllowDeleteEmployeeById,
		AllowAddServiceToEmployee,
		AllowRemoveServiceFromEmployee,
		AllowAddBranchToEmployee,
		AllowRemoveBranchFromEmployee,

		// Holidays
		AllowCreateHoliday,
		AllowGetHolidayById,
		AllowUpdateHolidayById,
		AllowDeleteHolidayById,

		// Sectors (Managed by Admins)
		AllowCreateSector,
		AllowUpdateSectorById,
		AllowDeleteSectorById,

		// Services
		AllowCreateService,
		AllowGetServiceById,
		AllowUpdateServiceById,
		AllowDeleteServiceById,
	}

	return Policies
}

var Policies []*PolicyRule

// SeedPolicies function needs to iterate through the *full* Policies list
func SeedPolicies(db *gorm.DB) ([]*PolicyRule, error) {
	AllowNilCompanyID = true // Allow seeding system policies without company_id
	AllowNilCreatedBy = true // Allow seeding system policies without created_by
	defer func() {
		AllowNilCompanyID = false
		AllowNilCreatedBy = false
	}()

	seededCount := 0
	updatedCount := 0 // Optionally track updates if you modify seeding logic

	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	Policies = init_policy_array()

	for _, policy := range Policies {
		if policy == nil {
			log.Println("Warning: Encountered nil policy in Policies list.")
			continue
		}
		if policy.EndPointID == 0 {
			// This assumes EndPoint IDs are assigned elsewhere before seeding.
			// If EndPoint structs need creating/finding first, do that here.
			log.Printf("Warning: Policy '%s' has EndPointID 0. Skipping seeding.", policy.Name)
			continue
		}

		// Create a placeholder to find existing policy
		var existingPolicy PolicyRule
		// Find system policies specifically (company_id IS NULL)
		err := tx.Where("name = ? AND company_id IS NULL AND end_point_id = ?", policy.Name, policy.EndPointID).First(&existingPolicy).Error

		if err == gorm.ErrRecordNotFound {
			// Policy doesn't exist, create it
			if err := tx.Create(policy).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create policy '%s': %w", policy.Name, err)
			}
			seededCount++
		} else if err != nil {
			// Other database error
			tx.Rollback()
			return nil, fmt.Errorf("failed to query policy '%s': %w", policy.Name, err)
		} else {
			// Policy exists - Optional: Update if definition changed (more complex)
			// For simplicity, we'll just log that it exists.
			// To update: Compare fields (e.g., Description, Conditions JSON) and tx.Save(&existingPolicy) if needed.
			// log.Printf("System policy '%s' already exists.", policy.Name)
			updatedCount++ // Increment if you implement update logic
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("System policies seeded successfully. New: %d, Existing/Updated: %d", seededCount, updatedCount)
	return Policies, nil // Return the original list (or query fresh if needed)
}
