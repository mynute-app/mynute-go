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
	ID                  uint            `gorm:"primaryKey"`
	CompanyID           *uint           `json:"company_id"`
	CreatedByEmployeeID *uint           `json:"created_by_employee_id"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Effect              string          `json:"effect"`      // "Allow" / "Deny"
	EndPointID          uint            `json:"endpoint_id"` // Link to EndPoint definition
	EndPoint            EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnDelete:CASCADE;" json:"endpoint"`
	Conditions          json.RawMessage `gorm:"type:jsonb" json:"conditions"`
	ResourceTable       string          `gorm:"-" json:"resource_table"`
	ResourceKey         string          `gorm:"-" json:"resource_key"`
	ResourceValueAt     string          `gorm:"-" json:"resource_value_at"`
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

// GORM Table Name (Keep as is)
func (PolicyRule) TableName() string {
	return "policy_rules"
}

// GORM Indexes (Keep as is)
func (PolicyRule) Indexes() map[string]string {
	return map[string]string{
		"idx_company_resource": "CREATE INDEX idx_company_resource ON policy_rules (company_id, resource_id)",
	}
}

// --- Helper Functions (Keep as is) ---

// MustMarshalJSON simplifies creating JSON examples
func MustMarshalJSON(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("MustMarshalJSON failed: %v", err))
	}
	return json.RawMessage(data)
}

// --- System Default Policies --- //

var client_access_check = ConditionNode{
	Description: "Client Access Check",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.company_id",
				Op:          "IsNull",
				Description: "Subject must not belong to a company",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.client_id", // Assumes the endpoint context has client_id for relevant resources (like Appointments)
				Description:     "Subject ID must match endpoint client ID",
			},
		},
	},
}

// New Check: Client Self Access (for managing own profile)
var client_self_access_check = ConditionNode{
	Description: "Client Self Access Check",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.company_id",
				Op:          "IsNull",
				Description: "Subject must not belong to a company (must be a Client)",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.id", // Check if subject ID matches the endpoint ID (the client ID being accessed)
				Description:     "Subject ID must match endpoint ID",
			},
		},
	},
}

// New Check: Employee Self Access (for managing own profile)
var employee_self_access_check = ConditionNode{
	Description: "Employee Self Access Check",
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
				ValueSourceAttr: "endpoint.id", // Check if subject ID matches the endpoint ID (the employee ID being accessed)
				Description:     "Subject ID must match endpoint ID",
			},
		},
		{ // Double check they are accessing their own company's endpoint, although endpoint.id check should suffice if IDs are unique
			Leaf: &ConditionLeaf{
				Attr:            "subject.company_id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.company_id", // Assumes employee endpoint has company_id
				Description:     "Subject company must match endpoint company",
			},
		},
	},
}

var company_membership_access_check = ConditionNode{
	Description: "Company Membership Check",
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
				ValueSourceAttr: "endpoint.company_id", // Assumes the endpoint context has company_id
				Description:     "Subject company must match endpoint company",
			},
		},
	},
}

// Simplified check just to ensure the subject belongs to *a* company
// Useful for creation where endpoint context might not be fully populated yet
var subject_is_company_member_check = ConditionNode{
	Description: "Subject Is Company Member Check",
	Leaf: &ConditionLeaf{
		Attr:        "subject.company_id",
		Op:          "IsNotNull",
		Description: "Subject must belong to a company",
	},
}

var company_owner_check = ConditionNode{
	Description: "Allow if subject is the company owner",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.role_id",
				Op:          "Equals",
				Value:       MustMarshalJSON(SystemRoleOwner.ID), // ASSUMES SystemRoleOwner is defined
				Description: "Subject role must be Owner",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.company_id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.company_id",
				Description:     "Subject company must match endpoint company",
			},
		},
	},
}

var company_general_manager_check = ConditionNode{
	Description: "Allow if subject is the company general manager",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.role_id",
				Op:          "Equals",
				Value:       MustMarshalJSON(SystemRoleGeneralManager.ID), // ASSUMES SystemRoleGeneralManager is defined
				Description: "Subject role must be General Manager",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.company_id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.company_id",
				Description:     "Subject company must match endpoint company",
			},
		},
	},
}

var company_branch_manager_check = ConditionNode{ // Simplified from assigned branch check - check if this role needs branch context
	Description: "Allow if subject is the company branch manager within the endpoint's company",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.role_id",
				Op:          "Equals",
				Value:       MustMarshalJSON(SystemRoleBranchManager.ID), // ASSUMES SystemRoleBranchManager is defined
				Description: "Subject role must be Branch Manager",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.company_id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.company_id",
				Description:     "Subject company must match endpoint company",
			},
		},
	},
}

var company_branch_manager_assigned_branch_check = ConditionNode{
	Description: "Allow if subject is the company branch manager and the endpoint relates to their assigned branch",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.role_id",
				Op:          "Equals",
				Value:       MustMarshalJSON(SystemRoleBranchManager.ID),
				Description: "Subject role must be Branch Manager",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.company_id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.company_id",
				Description:     "Subject company must match endpoint company",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.branches",   // Assumes subject context has assigned branch IDs
				Op:              "Contains",           // Check if the subject's assigned branches list contains the endpoint's branch ID
				ValueSourceAttr: "endpoint.branch_id", // Assumes endpoint context has branch_id
				Description:     "Subject branches must contain endpoint branch",
			},
		},
	},
}

var company_employee_check = ConditionNode{ // Simplified check just for Employee role within the company
	Description: "Allow if subject is an Employee within the endpoint's company",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr:        "subject.role_id",
				Op:          "Equals",
				Value:       MustMarshalJSON(SystemRoleEmployee.ID), // ASSUMES SystemRoleEmployee is defined
				Description: "Subject role must be Employee",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.company_id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.company_id",
				Description:     "Subject company must match endpoint company",
			},
		},
	},
}

var company_employee_assigned_employee_check = ConditionNode{
	Description: "Allow if subject is the company employee and the endpoint relates to them directly",
	LogicType:   "AND",
	Children: []ConditionNode{
		{
			Leaf: &ConditionLeaf{
				Attr: "subject.role_id",
				Op:   "Equals", Value: MustMarshalJSON(SystemRoleEmployee.ID),
				Description: "Subject role must be Employee",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr: "subject.company_id",
				Op:   "Equals", ValueSourceAttr: "endpoint.company_id",
				Description: "Subject company must match endpoint company",
			},
		},
		{
			Leaf: &ConditionLeaf{
				Attr:            "subject.id",
				Op:              "Equals",
				ValueSourceAttr: "endpoint.employee_id", // Assumes endpoint context has employee_id
				Description:     "Subject ID must match endpoint employee ID",
			},
		},
	},
}

// Reusable node for Admin Roles (Owner, General Manager)
var company_admin_check = ConditionNode{
	Description: "Company Admin Access Check (Owner or General Manager)",
	LogicType:   "OR",
	Children: []ConditionNode{
		company_owner_check,
		company_general_manager_check,
	},
}

// Reusable node for Manager Roles (Owner, General Manager, Branch Manager) - Context needed for BM branch check
var company_manager_check = ConditionNode{
	Description: "Company Manager Access Check (Owner, GM, or BM)",
	LogicType:   "OR",
	Children: []ConditionNode{
		company_owner_check,
		company_general_manager_check,
		company_branch_manager_check, // Using the simpler BM check here, adjust if specific branch access is always needed
	},
}

// Reusable node for all Company Internal Users (Owner, GM, BM, Employee) Check
var company_internal_user_check = ConditionNode{
	Description: "Company Internal User Check (Any Role within Company)",
	LogicType:   "AND",
	Children: []ConditionNode{
		company_membership_access_check, // Must be member of the endpoint's company
		{ // Role must be one of the company roles
			Description: "Role Check (Owner, GM, BM, Employee)",
			LogicType:   "OR",
			Children: []ConditionNode{
				{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleOwner.ID)}},
				{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleGeneralManager.ID)}},
				{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleBranchManager.ID)}},
				{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleEmployee.ID)}},
			},
		},
	},
}

// --- Policy Definitions --- //

// Policy: Allow Create appointment
var AllowCreateAppointment = &PolicyRule{ /* ... as defined in your first message ... */
	Name:            "SDP: CanCreateAppointment",
	Description:     "System Default Policy: Allows clients to create appointments, or company users based on role/relation.",
	Effect:          "Allow",
	EndPointID:      CreateAppointment.ID, // Link to the specific EndPoint definition
	ResourceTable:   "branches",
	ResourceKey:     "branch_id",
	ResourceValueAt: "body",
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Top-level OR: Allow Client Access OR Company User Access", // Top-level description
		LogicType:   "OR",
		Children: []ConditionNode{
			// Assume creation context sets `endpoint.client_id` if created by client? If not, this needs adjusting.
			client_access_check,
			{
				Description: "Company User Access Check",
				LogicType:   "AND",
				Children: []ConditionNode{
					// Requires subject.company_id = endpoint.company_id. Context must provide endpoint.company_id from input.
					company_membership_access_check,
					{
						Description: "Role/Relation Check",
						LogicType:   "OR",
						Children: []ConditionNode{
							company_owner_check,
							company_general_manager_check,
							company_branch_manager_assigned_branch_check, // Check endpoint context has branch_id
							company_employee_assigned_employee_check,     // Check endpoint context has employee_id
						},
					},
				},
			},
		},
	}),
}

// Policy: Allow GET appointment by ID.
var AllowGetAppointmentByID = &PolicyRule{ /* ... as defined in your first message ... */
	Name:            "SDP: CanViewAppointment",
	Description:     "System Default Policy: Allows clients to view own appointments, or company users based on role/relation.",
	Effect:          "Allow",
	EndPointID:      GetAppointmentByID.ID, // Link to the specific EndPoint definition
	ResourceTable:   "appointments",
	ResourceKey:     "id",
	ResourceValueAt: "path",
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Top-level OR: Allow Client Access OR Company User Access", // Top-level description
		LogicType:   "OR",
		Children: []ConditionNode{
			client_access_check, // Check endpoint context has client_id
			{
				Description: "Company User Access Check",
				LogicType:   "AND",
				Children: []ConditionNode{
					company_membership_access_check, // Check endpoint context has company_id
					{
						Description: "Role/Relation Check",
						LogicType:   "OR",
						Children: []ConditionNode{
							company_owner_check,
							company_general_manager_check,
							company_branch_manager_assigned_branch_check, // Check endpoint context has branch_id
							company_employee_assigned_employee_check,     // Check endpoint context has employee_id
						},
					},
				},
			},
		},
	}),
}

// --- NEW Policies --- //

// --- Appointment Policies (Update/Delete) ---

var AllowUpdateAppointmentByID = &PolicyRule{
	Name:            "SDP: CanUpdateAppointment",
	Description:     "Allows clients to update own appointments, or company managers/assigned employees.",
	Effect:          "Allow",
	ResourceTable:   "appointments",
	ResourceKey:     "id",
	ResourceValueAt: "path",
	EndPointID:      UpdateAppointmentByID.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Top-level OR: Allow Client Self-Update OR Company User Update",
		LogicType:   "OR",
		Children: []ConditionNode{
			client_access_check, // Check endpoint context has client_id
			{
				Description: "Company User Update Check",
				LogicType:   "AND",
				Children: []ConditionNode{
					company_membership_access_check, // Check endpoint context has company_id
					{
						Description: "Role/Relation Check (Managers or Assigned Employee)",
						LogicType:   "OR",
						Children: []ConditionNode{
							company_owner_check,
							company_general_manager_check,
							company_branch_manager_assigned_branch_check, // Check endpoint context has branch_id
							company_employee_assigned_employee_check,     // Check endpoint context has employee_id
						},
					},
				},
			},
		},
	}),
}

var AllowDeleteAppointmentByID = &PolicyRule{
	Name:            "SDP: CanDeleteAppointment",
	Description:     "Allows company managers or assigned employees to delete appointments.",
	Effect:          "Allow",
	ResourceTable:   "appointments",
	ResourceKey:     "id",
	ResourceValueAt: "path",
	EndPointID:      DeleteAppointmentByID.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		// NOTE: Clients usually cannot delete appointments directly, only company staff. Adjust if needed.
		Description: "Company User Delete Check",
		LogicType:   "AND",
		Children: []ConditionNode{
			company_membership_access_check, // Check endpoint context has company_id
			{
				Description: "Role/Relation Check (Managers or Assigned Employee)",
				LogicType:   "OR",
				Children: []ConditionNode{
					company_owner_check,
					company_general_manager_check,
					company_branch_manager_assigned_branch_check, // Check endpoint context has branch_id
					company_employee_assigned_employee_check,     // Check endpoint context has employee_id
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
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Company Admin Check for Creation",
		// Assumes the endpoint context for creation includes the target company_id based on input/subject
		LogicType: "OR", // Either Owner or GM of the target company
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
		},
	}),
}

var AllowGetBranchById = &PolicyRule{
	Name:        "SDP: CanViewBranchById",
	Description: "Allows any user belonging to the same company to view branch details by ID.",
	Effect:      "Allow",
	EndPointID:  GetBranchById.ID,
	Conditions:  MustMarshalJSON(company_internal_user_check), // Any internal user can view
}

var AllowGetBranchByName = &PolicyRule{
	Name:        "SDP: CanViewBranchByName",
	Description: "Allows any user belonging to the same company to view branch details by name.",
	Effect:      "Allow",
	EndPointID:  GetBranchByName.ID,
	Conditions:  MustMarshalJSON(company_internal_user_check), // Any internal user can view
}

var AllowUpdateBranchById = &PolicyRule{
	Name:        "SDP: CanUpdateBranch",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branches.",
	Effect:      "Allow",
	EndPointID:  UpdateBranchById.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Manager Update Access",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
			company_branch_manager_assigned_branch_check, // BM can update their own assigned branch
		},
	}),
}

var AllowDeleteBranchById = &PolicyRule{
	Name:        "SDP: CanDeleteBranch",
	Description: "Allows company Owner or General Manager to delete branches.",
	Effect:      "Allow",
	EndPointID:  DeleteBranchById.ID,
	Conditions:  MustMarshalJSON(company_admin_check), // Only Owner or GM
}

var AllowGetEmployeeServicesByBranchId = &PolicyRule{
	Name:        "SDP: CanViewEmployeeServicesInBranch",
	Description: "Allows company members to view employee services within a branch.",
	Effect:      "Allow",
	EndPointID:  GetEmployeeServicesByBranchId.ID,
	// Assumes endpoint context includes company_id based on branch_id lookip
	Conditions: MustMarshalJSON(company_internal_user_check),
}

var AllowAddServiceToBranch = &PolicyRule{
	Name:        "SDP: CanAddServiceToBranch",
	Description: "Allows company managers (Owner, GM, relevant BM) to add services to a branch.",
	Effect:      "Allow",
	EndPointID:  AddServiceToBranch.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Manager Access Check",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
			company_branch_manager_assigned_branch_check, // BM can manage services in their own branch
		},
	}),
}

var AllowRemoveServiceFromBranch = &PolicyRule{
	Name:        "SDP: CanRemoveServiceFromBranch",
	Description: "Allows company managers (Owner, GM, relevant BM) to remove services from a branch.",
	Effect:      "Allow",
	EndPointID:  RemoveServiceFromBranch.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Manager Access Check",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
			company_branch_manager_assigned_branch_check, // BM can manage services in their own branch
		},
	}),
}

// --- Client Policies ---

// GetClientByEmail - Who can search for clients by email? Assume only Admins of *any* company, or the client themselves. This is sensitive.
// Let's restrict to only client self-lookup for now unless admins need it.
var AllowGetClientByEmail = &PolicyRule{
	Name:        "SDP: CanViewClientByEmail",
	Description: "Allows a client to retrieve their own profile by email.",
	Effect:      "Allow",
	EndPointID:  GetClientByEmail.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Allow only if the subject's email matches the endpoint email being queried.",
		Leaf: &ConditionLeaf{
			Attr:            "subject.email", // Requires subject context to have email
			Op:              "Equals",
			ValueSourceAttr: "endpoint.email", // Requires endpoint context to have email from path param
		},
	}),
}

var AllowUpdateClientById = &PolicyRule{
	Name:        "SDP: CanUpdateClient",
	Description: "Allows a client to update their own profile.",
	Effect:      "Allow",
	EndPointID:  UpdateClientById.ID,
	Conditions:  MustMarshalJSON(client_self_access_check), // Client can update self
	// Add OR company_admin_check if company admins should manage clients?
}

var AllowDeleteClientById = &PolicyRule{
	Name:        "SDP: CanDeleteClient",
	Description: "Allows a client to delete their own profile.",
	Effect:      "Allow",
	EndPointID:  DeleteClientById.ID,
	Conditions:  MustMarshalJSON(client_self_access_check), // Client can delete self
	// Add OR company_admin_check if company admins should manage clients?
}

// --- Company Policies ---

var AllowGetCompanyById = &PolicyRule{
	Name:        "SDP: CanViewCompanyById",
	Description: "Allows any member (employee) of the company to view its details.",
	Effect:      "Allow",
	EndPointID:  GetCompanyById.ID,
	// Anyone belonging to the company can view it
	Conditions: MustMarshalJSON(company_membership_access_check),
}

var AllowUpdateCompanyById = &PolicyRule{
	Name:        "SDP: CanUpdateCompany",
	Description: "Allows the company Owner or General Manager to update company details.",
	Effect:      "Allow",
	EndPointID:  UpdateCompanyById.ID,
	Conditions:  MustMarshalJSON(company_admin_check), // Only Owner or GM
}

var AllowDeleteCompanyById = &PolicyRule{
	Name:        "SDP: CanDeleteCompany",
	Description: "Allows ONLY the company Owner to delete the company.",
	Effect:      "Allow",
	EndPointID:  DeleteCompanyById.ID,
	Conditions:  MustMarshalJSON(company_owner_check), // Only Owner
}

// --- Employee Policies ---

var AllowCreateEmployee = &PolicyRule{
	Name:        "SDP: CanCreateEmployee",
	Description: "Allows company Owner, GM, or BM to create employees (within their scope).",
	Effect:      "Allow",
	EndPointID:  CreateEmployee.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		// Assumes creation context provides target company_id, and potentially branch_id if created by BM
		Description: "Manager Check for Creation",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,           // Owner can create in their company
			company_general_manager_check, // GM can create in their company
			// BM creating employee likely needs to create them within an assigned branch.
			// This check requires the input payload (endpoint context) to have branch_id.
			company_branch_manager_assigned_branch_check,
		},
	}),
}

var AllowGetEmployeeById = &PolicyRule{
	Name:        "SDP: CanViewEmployeeById",
	Description: "Allows company members to view employee profiles within the same company, or employee to view their own.",
	Effect:      "Allow",
	EndPointID:  GetEmployeeById.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Allow Employee Self-View OR Any Internal Company User View",
		LogicType:   "OR",
		Children: []ConditionNode{
			employee_self_access_check,  // Can view self
			company_internal_user_check, // Any other member of the same company can view
		},
	}),
}

var AllowGetEmployeeByEmail = &PolicyRule{
	Name:        "SDP: CanViewEmployeeByEmail",
	Description: "Allows company members to find employees within the same company by email.",
	Effect:      "Allow",
	EndPointID:  GetEmployeeByEmail.ID,
	// Requires endpoint context providing company_id based on the lookup by email
	Conditions: MustMarshalJSON(company_internal_user_check),
}

var AllowUpdateEmployeeById = &PolicyRule{
	Name:        "SDP: CanUpdateEmployee",
	Description: "Allows employee to update self, or company managers (Owner, GM, BM) to update employees.",
	Effect:      "Allow",
	EndPointID:  UpdateEmployeeById.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Allow Employee Self-Update OR Manager Update",
		LogicType:   "OR",
		Children: []ConditionNode{
			employee_self_access_check, // Employee can update own profile
			company_manager_check,      // Owner, GM, BM can update employees in their company
		},
	}),
}

var AllowDeleteEmployeeById = &PolicyRule{
	Name:        "SDP: CanDeleteEmployee",
	Description: "Allows company managers (Owner, GM, BM) to delete employees.",
	Effect:      "Allow",
	EndPointID:  DeleteEmployeeById.ID,
	Conditions:  MustMarshalJSON(company_manager_check), // Owner, GM, BM can delete
}

var AllowAddServiceToEmployee = &PolicyRule{
	Name:        "SDP: CanAddServiceToEmployee",
	Description: "Allows company managers (Owner, GM, BM) to assign services to employees.",
	Effect:      "Allow",
	EndPointID:  AddServiceToEmployee.ID,
	// EndPoint context needs company_id based on employee_id lookup
	Conditions: MustMarshalJSON(company_manager_check),
}

var AllowRemoveServiceFromEmployee = &PolicyRule{
	Name:        "SDP: CanRemoveServiceFromEmployee",
	Description: "Allows company managers (Owner, GM, BM) to remove services from employees.",
	Effect:      "Allow",
	EndPointID:  RemoveServiceFromEmployee.ID,
	// EndPoint context needs company_id based on employee_id lookup
	Conditions: MustMarshalJSON(company_manager_check),
}

var AllowAddBranchToEmployee = &PolicyRule{
	Name:        "SDP: CanAddBranchToEmployee",
	Description: "Allows company managers (Owner, GM, BM) to assign employees to branches (respecting BM scope).",
	Effect:      "Allow",
	EndPointID:  AddBranchToEmployee.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		// EndPoint context needs company_id (from employee) and branch_id (from path)
		Description: "Manager Assignment Check",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,           // Owner can assign any branch in company
			company_general_manager_check, // GM can assign any branch in company
			// BM can only assign employees TO a branch THEY MANAGE
			company_branch_manager_assigned_branch_check, // Checks if subject.branches CONTAINS endpoint.branch_id
		},
	}),
}

var AllowRemoveBranchFromEmployee = &PolicyRule{
	Name:        "SDP: CanRemoveBranchFromEmployee",
	Description: "Allows company managers (Owner, GM, BM) to remove employees from branches (respecting BM scope).",
	Effect:      "Allow",
	EndPointID:  RemoveBranchFromEmployee.ID,
	Conditions: MustMarshalJSON(ConditionNode{
		// EndPoint context needs company_id (from employee) and branch_id (from path)
		Description: "Manager Assignment Check",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_owner_check,
			company_general_manager_check,
			company_branch_manager_assigned_branch_check, // BM can only unassign employees FROM a branch THEY MANAGE
		},
	}),
}

// --- Holiday Policies ---

var AllowCreateHoliday = &PolicyRule{
	Name:        "SDP: CanCreateHoliday",
	Description: "Allows company managers (Owner, GM, BM) to create holidays for the company/branch.",
	Effect:      "Allow",
	EndPointID:  CreateHoliday.ID,
	// Assumes context provides company_id. If holiday is branch-specific, context needs branch_id too.
	Conditions: MustMarshalJSON(company_manager_check),
}

var AllowGetHolidayById = &PolicyRule{
	Name:        "SDP: CanViewHolidayById",
	Description: "Allows company members to view holiday details.",
	Effect:      "Allow",
	EndPointID:  GetHolidayById.ID,
	Conditions:  MustMarshalJSON(company_internal_user_check), // Any internal user can view
}

var AllowUpdateHolidayById = &PolicyRule{
	Name:        "SDP: CanUpdateHoliday",
	Description: "Allows company managers (Owner, GM, BM) to update holidays.",
	Effect:      "Allow",
	EndPointID:  UpdateHolidayById.ID,
	Conditions:  MustMarshalJSON(company_manager_check),
}

var AllowDeleteHolidayById = &PolicyRule{
	Name:        "SDP: CanDeleteHoliday",
	Description: "Allows company managers (Owner, GM, BM) to delete holidays.",
	Effect:      "Allow",
	EndPointID:  DeleteHolidayById.ID,
	Conditions:  MustMarshalJSON(company_manager_check),
}

// --- Sector Policies --- (Assuming Sectors are company-wide, managed by Admins)

var AllowCreateSector = &PolicyRule{
	Name:        "SDP: CanCreateSector",
	Description: "Allows company Owner or General Manager to create sectors.",
	Effect:      "Allow",
	EndPointID:  CreateSector.ID,
	// Assumes context provides company_id based on subject/input
	Conditions: MustMarshalJSON(company_admin_check),
}

// GetSectorById/Name are Public, no policy needed.

var AllowUpdateSectorById = &PolicyRule{
	Name:        "SDP: CanUpdateSector",
	Description: "Allows company Owner or General Manager to update sectors.",
	Effect:      "Allow",
	EndPointID:  UpdateSectorById.ID,
	Conditions:  MustMarshalJSON(company_admin_check),
}

var AllowDeleteSectorById = &PolicyRule{
	Name:        "SDP: CanDeleteSector",
	Description: "Allows company Owner or General Manager to delete sectors.",
	Effect:      "Allow",
	EndPointID:  DeleteSectorById.ID,
	Conditions:  MustMarshalJSON(company_admin_check),
}

// --- Service Policies ---

var AllowCreateService = &PolicyRule{
	Name:        "SDP: CanCreateService",
	Description: "Allows company managers (Owner, GM, BM) to create services.",
	Effect:      "Allow",
	EndPointID:  CreateService.ID,
	// Assumes context provides company_id based on subject/input. BM might be restricted later.
	Conditions: MustMarshalJSON(company_manager_check),
}

var AllowGetServiceById = &PolicyRule{
	Name:        "SDP: CanViewServiceById",
	Description: "Allows company members to view service details.",
	Effect:      "Allow",
	EndPointID:  GetServiceById.ID,
	Conditions:  MustMarshalJSON(company_internal_user_check), // Any internal user can view
}

// GetServiceByName is Public, no policy needed.

var AllowUpdateServiceById = &PolicyRule{
	Name:        "SDP: CanUpdateService",
	Description: "Allows company managers (Owner, GM, BM) to update services.",
	Effect:      "Allow",
	EndPointID:  UpdateServiceById.ID,
	Conditions:  MustMarshalJSON(company_manager_check),
}

var AllowDeleteServiceById = &PolicyRule{
	Name:        "SDP: CanDeleteService",
	Description: "Allows company managers (Owner, GM, BM) to delete services.",
	Effect:      "Allow",
	EndPointID:  DeleteServiceById.ID,
	Conditions:  MustMarshalJSON(company_manager_check),
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
		err := tx.Where("name = ? AND company_id IS NULL AND resource_id = ?", policy.Name, policy.EndPointID).First(&existingPolicy).Error

		if err == gorm.ErrRecordNotFound {
			// Policy doesn't exist, create it
			log.Printf("Seeding new system policy: %s", policy.Name)
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

	log.Printf("System policies seeding completed. New: %d, Existing/Updated: %d", seededCount, updatedCount)
	return Policies, nil // Return the original list (or query fresh if needed)
}
