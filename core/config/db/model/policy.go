package model

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllowNilCompanyID = false
var AllowNilCreatedBy = false
var AllowNilResourceID = false

// --- PolicyRule (Represents a policy rule for access control) ---
// [Conditions] : Effect : Method : Resource : Property
// --- [if is company owner] : Allow : PATCH : /company/{id} : nil
// --- [if is not company owner] : Deny : PATCH : /company/{id} : tax_id
type PolicyRule struct {
	BaseModel
	CompanyID           *uuid.UUID      `json:"company_id"`
	Company             Company         `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company"`
	CreatedByEmployeeID *uuid.UUID      `json:"created_by_employee_id"`
	CreatedByEmployee   Employee        `gorm:"foreignKey:CreatedByEmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"created_by_employee"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Effect              string          `json:"effect"` // "Allow" / "Deny"
	EndPointID          uuid.UUID       `json:"end_point_id"`
	EndPoint            EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	PropertyID          *uuid.UUID      `json:"property_id"`
	Property            *Property       `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"property"`
	Conditions          json.RawMessage `gorm:"type:jsonb" json:"conditions"`
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

// --- ConditionLeaf (Represents a single atomic check) ---
type ConditionLeaf struct {
	Attribute string `json:"attribute"` // The primary attribute (e.g., "subject.role_id")
	// The comparison operator
	// Equals - "==", NotEquals - "!=", IsNull - "== null", IsNotNull - "!= null"
	Operator    string `json:"operator"`
	Description string `json:"description,omitempty"` // Optional human-readable description of the check

	// Use EITHER Value OR ResourceAttribute for comparison. Omitempty ensures only one (or neither for ops like IsNull) appears in JSON.
	Value             json.RawMessage `json:"value,omitempty"`              // Static value to compare against
	ResourceAttribute string          `json:"resource_attribute,omitempty"` // Other attribute's name to compare against
}

// GetConditionsNode parses and validates the stored JSON conditions
func (p *PolicyRule) GetConditionsNode() (ConditionNode, error) {
	var node ConditionNode // Initialize empty node

	// 1. Check for empty or null JSON content
	if len(p.Conditions) == 0 || string(p.Conditions) == "null" {
		// Decide if this is an error or implies a default. Forcing explicit
		// conditions is usually safer. Use p.Name for better context if available.
		policyIdentifier := fmt.Sprintf("ID %d", p.ID)
		if p.Name != "" {
			policyIdentifier = fmt.Sprintf("'%s' (ID %d)", p.Name, p.ID)
		}
		return node, fmt.Errorf("policy rule %s has missing or null conditions", policyIdentifier)
	}

	// 2. Attempt to unmarshal the JSON
	err := json.Unmarshal(p.Conditions, &node)
	if err != nil {
		policyIdentifier := fmt.Sprintf("ID %d", p.ID)
		if p.Name != "" {
			policyIdentifier = fmt.Sprintf("'%s' (ID %d)", p.Name, p.ID)
		}
		return node, fmt.Errorf("failed to unmarshal conditions JSON for policy rule %s: %w", policyIdentifier, err)
	}

	// 3. Perform recursive validation using the dedicated validator function
	// This replaces your original two specific validation checks.
	if err := validateConditionNode(node); err != nil {
		policyIdentifier := fmt.Sprintf("ID %d", p.ID)
		if p.Name != "" {
			policyIdentifier = fmt.Sprintf("'%s' (ID %d)", p.Name, p.ID)
		}
		// Wrap the validation error with policy context
		return node, fmt.Errorf("invalid conditions structure for policy rule %s: %w", policyIdentifier, err)
	}

	// 4. Return the successfully parsed and validated node
	return node, nil
}

// validateConditionNode performs recursive validation on a condition node structure.
func validateConditionNode(node ConditionNode) error {
	nodeContext := ""
	if node.Description != "" {
		nodeContext = fmt.Sprintf(" ('%s')", node.Description)
	}

	if node.Leaf != nil { // It's intended to be a leaf node
		// Rule 1: A leaf node cannot have branch properties
		if node.LogicType != "" || len(node.Children) > 0 {
			return fmt.Errorf("node%s is incorrectly structured as both leaf and branch", nodeContext)
		}

		// Rule 2: Leaf must have an attribute
		if node.Leaf.Attribute == "" {
			return fmt.Errorf("leaf node `%s` is missing required 'attribute'", nodeContext)
		}
		// Rule 3: Leaf must have an operator
		if node.Leaf.Operator == "" {
			return fmt.Errorf("leaf node `%s` (attribute '%s') is missing required 'operator'", nodeContext, node.Leaf.Attribute)
		}

		// Rule 4: Operators needing comparison values must have 'value' or 'resource_attribute'
		requiresComparisonValue := true
		switch node.Leaf.Operator {
		case "IsNull", "IsNotNull":
			requiresComparisonValue = false
			// Add any other future unary operators here (e.g., "IsEmpty", "IsTrue")
		}

		if requiresComparisonValue && node.Leaf.Value == nil && node.Leaf.ResourceAttribute == "" {
			return fmt.Errorf("leaf node `%s` (attribute '%s', operator '%s') requires either 'value' or 'resource_attribute'", nodeContext, node.Leaf.Attribute, node.Leaf.Operator)
		}

		// Rule 5: Leaf cannot have *both* 'value' and 'resource_attribute'
		if node.Leaf.Value != nil && node.Leaf.ResourceAttribute != "" {
			return fmt.Errorf("leaf node `%s` (attribute '%s') cannot have both 'value' and 'resource_attribute' defined", nodeContext, node.Leaf.Attribute)
		}

		// Rule 6: Basic sanity check on attribute format
		if node.Leaf.Attribute != "" && !strings.HasPrefix(node.Leaf.Attribute, "subject.") && !strings.HasPrefix(node.Leaf.Attribute, "resource.") {
			return fmt.Errorf("leaf node `%s` has invalid 'attribute' ('%s'): must start with 'subject.' or 'resource.'", nodeContext, node.Leaf.Attribute)
		}
		if node.Leaf.ResourceAttribute != "" && !strings.HasPrefix(node.Leaf.ResourceAttribute, "subject.") && !strings.HasPrefix(node.Leaf.ResourceAttribute, "resource.") {
			return fmt.Errorf("leaf node `%s` has invalid 'resource_attribute' ('%s'): must start with 'subject.' or 'resource.'", nodeContext, node.Leaf.ResourceAttribute)
		}

	} else { // It's intended to be a branch node (or potentially an empty root? - see notes)

		// Rule 7: A branch must have a valid LogicType (unless it's truly empty - decide if that's allowed)
		isValidBranch := false
		if node.LogicType == "AND" || node.LogicType == "OR" {
			isValidBranch = true
			// Rule 8: A branch with a logic type must have children
			if len(node.Children) == 0 {
				return fmt.Errorf("branch node `%s` has 'logic_type' %s but no 'children'", nodeContext, node.LogicType)
			}
		}

		// Rule 9: Check if it's an invalid structure (neither leaf nor valid branch)
		// This replaces your first original check `node.Leaf == nil && (node.LogicType == "" || len(node.Children) == 0)`
		// It's invalid if it's not a leaf AND it's not a valid branch structure defined above.
		// Exception: Is an *empty* root node (`{}`, Description only) acceptable?
		// If so, add logic here. Assuming for now it's invalid if it's not leaf/branch.
		if node.Leaf == nil && !isValidBranch {
			// This condition means: It's not a leaf. AND (LogicType is missing/invalid OR Children are missing even if logic type is present)
			// If Description is the *only* thing present, maybe allow? Needs careful consideration.
			if node.LogicType == "" && len(node.Children) == 0 && node.Description != "" {
				// Possibly allow an empty descriptive node? Or treat as error? Treat as error for stricter policy defs.
				return fmt.Errorf("node%s is not a valid leaf (missing 'leaf') nor a valid branch (missing/invalid 'logic_type' or missing 'children')", nodeContext)
			} else if node.LogicType != "" && !isValidBranch {
				return fmt.Errorf("branch node `%s` has invalid 'logic_type': '%s' (must be AND or OR)", nodeContext, node.LogicType)
			} else {
				// Generic catch-all for invalid structure
				return fmt.Errorf("node%s is neither a valid leaf nor a valid branch", nodeContext)
			}
		}

		// Rule 10: Recursively validate children if it's a valid branch
		if isValidBranch {
			for i, child := range node.Children {
				if err := validateConditionNode(child); err != nil {
					// Add context about which child failed
					return fmt.Errorf("invalid child node %d under node%s: %w", i+1, nodeContext, err)
				}
			}
		}
	}

	// All checks passed for this node and its children (if any)
	return nil
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
					Attribute:   "subject.company_id",
					Operator:    "IsNull",
					Description: "Subject must be a Client (no company affiliation)",
				},
			},
			{
				Leaf: &ConditionLeaf{
					Attribute:         "subject.id",
					Operator:          "Equals",
					ResourceAttribute: "resource.client_id", // Assumes context provides resource.client_id from the fetched resource
					Description:       "Subject ID must match the resource's client ID",
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
					Attribute:   "subject.company_id",
					Operator:    "IsNull",
					Description: "Subject must be a Client",
				},
			},
			{
				Leaf: &ConditionLeaf{
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
	var employee_self_access_check = ConditionNode{
		Description: "Employee Self Access Check (Own Profile/Resource)",
		LogicType:   "AND",
		Children: []ConditionNode{
			{
				Leaf: &ConditionLeaf{
					Attribute:   "subject.company_id",
					Operator:    "IsNotNull",
					Description: "Subject must belong to a company",
				},
			},
			{
				Leaf: &ConditionLeaf{
					Attribute:         "subject.id",
					Operator:          "Equals",
					ResourceAttribute: "resource.id", // Assumes context provides resource.id from path matching subject's ID
					Description:       "Subject ID must match the resource ID being accessed",
				},
			},
			{ // Belt-and-suspenders: Check company match too
				Leaf: &ConditionLeaf{
					Attribute:         "subject.company_id",
					Operator:          "Equals",
					ResourceAttribute: "resource.company_id", // Assumes context provides resource.company_id from the fetched resource
					Description:       "Subject company must match the resource's company",
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
					Attribute:   "subject.company_id",
					Operator:    "IsNotNull",
					Description: "Subject must belong to a company",
				},
			},
			{
				Description: "Either company ID is at .id or .company_id",
				LogicType:   "OR",
				Children: []ConditionNode{
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.company_id",
							Operator:          "Equals",
							ResourceAttribute: "resource.company_id", // Assumes context provides resource.company_id from the fetched resource
							Description:       "Subject company must match the resource's company",
						},
					},
					{
						Leaf: &ConditionLeaf{
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
	var company_owner_check = ConditionNode{
		Description: "Allow if Subject is Owner of the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{
				Leaf: &ConditionLeaf{
					Attribute: "subject.roles[*].ID",
					Operator:  "Contains",
					Value:     JsonRawMessage(SystemRoleOwner.ID),
				},
			},
			company_membership_access_check, // Re-use the company match check
		},
	}

	// Checks if subject is the General Manager of the resource's company
	var company_general_manager_check = ConditionNode{
		Description: "Allow if Subject is General Manager of the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleGeneralManager.ID)}},
			company_membership_access_check, // Re-use the company match check
		},
	}

	// Checks if subject is a Branch Manager of the resource's company
	var company_branch_manager_check = ConditionNode{
		Description: "Allow if Subject is Branch Manager within the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleBranchManager.ID)}},
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
					Attribute:         "subject.branches",   // Assumes subject context has assigned branch IDs (e.g., [10, 25])
					Operator:          "Contains",           // Checks if the list contains the value
					ResourceAttribute: "resource.branch_id", // Assumes context provides resource.branch_id from the resource/path/body
					Description:       "Subject's assigned branches must include the resource's branch",
				},
			},
		},
	}

	// Checks if subject is an Employee AND their ID matches the resource's employee_id
	var company_employee_assigned_employee_check = ConditionNode{
		Description: "Allow if Subject is the Employee associated with the Resource",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleEmployee.ID)}},
			company_membership_access_check, // Must be in the same company
			{
				Leaf: &ConditionLeaf{
					Attribute:         "subject.id",
					Operator:          "Equals",
					ResourceAttribute: "resource.employee_id", // Assumes context provides resource.employee_id from the resource
					Description:       "Subject ID must match the resource's employee ID",
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
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleOwner.ID)}},
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleGeneralManager.ID)}},
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleBranchManager.ID)}},
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].ID", Operator: "Contains", Value: JsonRawMessage(SystemRoleEmployee.ID)}},
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
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Client Creation OR Company User Creation",
			LogicType:   "OR",
			Children: []ConditionNode{
				// Client creating for themselves. Assumes 'client_id' is in the body and resource.client_id gets populated.
				{
					Description: "Client Self-Creation",
					LogicType:   "AND",
					Children: []ConditionNode{
						{Leaf: &ConditionLeaf{Attribute: "subject.company_id", Operator: "IsNull", Description: "Must be a Client"}},
						{Leaf: &ConditionLeaf{Attribute: "subject.id", Operator: "Equals", ResourceAttribute: "resource.client_id", Description: "Client ID in body must match Subject ID"}},
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
		Name:        "SDP: CanViewAppointment",
		Description: "Allows clients to view own appointments, or company users based on role/relation.",
		Effect:      "Allow",
		EndPointID:  GetAppointmentByID.ID,
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
		Name:        "SDP: CanUpdateAppointment",
		Description: "Allows clients to update own appointments, or company managers/assigned employees.",
		Effect:      "Allow",
		EndPointID:  UpdateAppointmentByID.ID,
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
		Name:        "SDP: CanDeleteAppointment",
		Description: "Allows company managers or assigned employees to delete appointments. (Clients typically cannot delete).",
		Effect:      "Allow",
		EndPointID:  DeleteAppointmentByID.ID,
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
		Conditions:  JsonRawMessage(company_admin_check), // Owner or GM of the target company
	}

	var AllowGetBranchById = &PolicyRule{
		Name:        "SDP: CanViewBranchById",
		Description: "Allows any user belonging to the same company to view branch details by ID.",
		Effect:      "Allow",
		EndPointID:  GetBranchById.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company can view
	}

	var AllowGetBranchByName = &PolicyRule{
		Name:        "SDP: CanViewBranchByName",
		Description: "Allows any user belonging to the same company to view branch details by name.",
		Effect:      "Allow",
		EndPointID:  GetBranchByName.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company can view
	}

	var AllowUpdateBranchById = &PolicyRule{
		Name:        "SDP: CanUpdateBranch",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branches.",
		Effect:      "Allow",
		EndPointID:  UpdateBranchById.ID,
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
		Name:        "SDP: CanDeleteBranch",
		Description: "Allows company Owner or General Manager to delete branches.",
		Effect:      "Allow",
		EndPointID:  DeleteBranchById.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Only Owner or GM
	}

	var AllowGetEmployeeServicesByBranchId = &PolicyRule{
		Name:        "SDP: CanViewEmployeeServicesInBranch",
		Description: "Allows company members to view employee services within a branch.",
		Effect:      "Allow",
		EndPointID:  GetEmployeeServicesByBranchId.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company
	}

	var AllowAddServiceToBranch = &PolicyRule{
		Name:        "SDP: CanAddServiceToBranch",
		Description: "Allows company managers (Owner, GM, relevant BM) to add services to a branch.",
		Effect:      "Allow",
		EndPointID:  AddServiceToBranch.ID,
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
		Name:        "SDP: CanRemoveServiceFromBranch",
		Description: "Allows company managers (Owner, GM, relevant BM) to remove services from a branch.",
		Effect:      "Allow",
		EndPointID:  RemoveServiceFromBranch.ID,
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
		Name:        "SDP: CanViewClientByEmail",
		Description: "Allows a client to retrieve their own profile by email.",
		Effect:      "Allow",
		EndPointID:  GetClientByEmail.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow only if the subject's email matches the email in the path.",
			LogicType:   "AND",
			Children: []ConditionNode{
				{Leaf: &ConditionLeaf{Attribute: "subject.company_id", Operator: "IsNull", Description: "Must be a Client"}},                                                         // Ensure subject is a client
				{Leaf: &ConditionLeaf{Attribute: "subject.email", Operator: "Equals", ResourceAttribute: "resource.email", Description: "Subject email must match email from path"}}, // Assumes context has resource.email from path
			},
		}),
	}

	var AllowUpdateClientById = &PolicyRule{
		Name:        "SDP: CanUpdateClient",
		Description: "Allows a client to update their own profile.",
		Effect:      "Allow",
		EndPointID:  UpdateClientById.ID,
		Conditions:  JsonRawMessage(client_self_access_check), // Client can update self (checks subject.id == resource.id)
	}

	var AllowDeleteClientById = &PolicyRule{
		Name:        "SDP: CanDeleteClient",
		Description: "Allows a client to delete their own profile.",
		Effect:      "Allow",
		EndPointID:  DeleteClientById.ID,
		Conditions:  JsonRawMessage(client_self_access_check), // Client can delete self
	}

	// --- Company Policies ---

	var AllowGetCompanyById = &PolicyRule{
		Name:        "SDP: CanViewCompanyById",
		Description: "Allows any member (employee/manager) of the company to view its details.",
		Effect:      "Allow",
		EndPointID:  GetCompanyById.ID,
		Conditions:  JsonRawMessage(company_membership_access_check),
	}

	var AllowUpdateCompanyById = &PolicyRule{
		Name:        "SDP: CanUpdateCompany",
		Description: "Allows the company Owner or General Manager to update company details.",
		Effect:      "Allow",
		EndPointID:  UpdateCompanyById.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Only Owner or GM of this company
	}

	var AllowDeleteCompanyById = &PolicyRule{
		Name:        "SDP: CanDeleteCompany",
		Description: "Allows ONLY the company Owner to delete the company.",
		Effect:      "Allow",
		EndPointID:  DeleteCompanyById.ID,
		Conditions:  JsonRawMessage(company_owner_check), // Only Owner of this company
	}

	// --- Employee Policies ---

	var AllowCreateEmployee = &PolicyRule{
		Name:        "SDP: CanCreateEmployee",
		Description: "Allows company Owner, GM, or BM to create employees (BM restricted to their branches implicitly if data includes branch).",
		Effect:      "Allow",
		EndPointID:  CreateEmployee.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Branch Manager Creation Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,
				company_branch_manager_assigned_branch_check,
			},
		}),
	}

	var AllowGetEmployeeById = &PolicyRule{
		Name:        "SDP: CanViewEmployeeById",
		Description: "Allows employee to view self, or any internal user of the same company to view other employees.",
		Effect:      "Allow",
		EndPointID:  GetEmployeeById.ID,
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
		Name:        "SDP: CanViewEmployeeByEmail",
		Description: "Allows company members to find employees within the same company by email.",
		Effect:      "Allow",
		EndPointID:  GetEmployeeByEmail.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Subject must be internal user of the found employee's company
	}

	var AllowUpdateEmployeeById = &PolicyRule{
		Name:        "SDP: CanUpdateEmployee",
		Description: "Allows employee to update self, or company managers (Owner, GM, BM) to update employees.",
		Effect:      "Allow",
		EndPointID:  UpdateEmployeeById.ID,
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
		Name:        "SDP: CanDeleteEmployee",
		Description: "Allows company managers (Owner, GM, BM) to delete employees.",
		Effect:      "Allow",
		EndPointID:  DeleteEmployeeById.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Owner, GM, BM can delete
	}

	var AllowAddServiceToEmployee = &PolicyRule{
		Name:        "SDP: CanAddServiceToEmployee",
		Description: "Allows company managers (Owner, GM, BM) to assign services to employees.",
		Effect:      "Allow",
		EndPointID:  AddServiceToEmployee.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Manager of the employee's company
	}

	var AllowRemoveServiceFromEmployee = &PolicyRule{
		Name:        "SDP: CanRemoveServiceFromEmployee",
		Description: "Allows company managers (Owner, GM, BM) to remove services from employees.",
		Effect:      "Allow",
		EndPointID:  RemoveServiceFromEmployee.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Manager of the employee's company
	}

	var AllowAddBranchToEmployee = &PolicyRule{
		Name:        "SDP: CanAddBranchToEmployee",
		Description: "Allows company managers (Owner, GM, BM) to assign employees to branches (respecting BM scope).",
		Effect:      "Allow",
		EndPointID:  AddBranchToEmployee.ID,
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
		Name:        "SDP: CanRemoveBranchFromEmployee",
		Description: "Allows company managers (Owner, GM, BM) to remove employees from branches (respecting BM scope).",
		Effect:      "Allow",
		EndPointID:  RemoveBranchFromEmployee.ID,
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
		Conditions:  JsonRawMessage(company_manager_check), // Any manager in the relevant company context. Add branch check if needed.
	}

	var AllowGetHolidayById = &PolicyRule{
		Name:        "SDP: CanViewHolidayById",
		Description: "Allows company members to view holiday details.",
		Effect:      "Allow",
		EndPointID:  GetHolidayById.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the holiday's company
	}

	var AllowUpdateHolidayById = &PolicyRule{
		Name:        "SDP: CanUpdateHoliday",
		Description: "Allows company managers (Owner, GM, BM) to update holidays.",
		Effect:      "Allow",
		EndPointID:  UpdateHolidayById.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Any manager of the holiday's company
	}

	var AllowDeleteHolidayById = &PolicyRule{
		Name:        "SDP: CanDeleteHoliday",
		Description: "Allows company managers (Owner, GM, BM) to delete holidays.",
		Effect:      "Allow",
		EndPointID:  DeleteHolidayById.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Any manager of the holiday's company
	}

	// --- Sector Policies --- (Assuming Company Specific)

	// var AllowCreateSector = &PolicyRule{
	// 	Name:        "SDP: CanCreateSector",
	// 	Description: "Allows company Owner or General Manager to create sectors.",
	// 	Effect:      "Allow",
	// 	Conditions:  JsonRawMessage(company_admin_check), // Only Owner/GM of that company
	// }

	// // GetSectorById/Name assumed Public - No Policy Needed

	// var AllowUpdateSectorById = &PolicyRule{
	// 	Name:        "SDP: CanUpdateSector",
	// 	Description: "Allows company Owner or General Manager to update sectors.",
	// 	Effect:      "Allow",
	// 	Conditions:  JsonRawMessage(company_admin_check), // Only Owner/GM of the sector's company
	// }

	// var AllowDeleteSectorById = &PolicyRule{
	// 	Name:        "SDP: CanDeleteSector",
	// 	Description: "Allows company Owner or General Manager to delete sectors.",
	// 	Effect:      "Allow",
	// 	Conditions:  JsonRawMessage(company_admin_check), // Only Owner/GM of the sector's company
	// }

	// --- Service Policies --- (Assuming Company Specific)

	var AllowCreateService = &PolicyRule{
		Name:        "SDP: CanCreateService",
		Description: "Allows company managers (Owner, GM, BM) to create services.",
		Effect:      "Allow",
		EndPointID:  CreateService.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Any manager of the company context
	}

	var AllowGetServiceById = &PolicyRule{
		Name:        "SDP: CanViewServiceById",
		Description: "Allows company members to view service details.",
		Effect:      "Allow",
		EndPointID:  GetServiceById.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the service's company
	}

	// GetServiceByName assumed Public - No Policy Needed

	var AllowUpdateServiceById = &PolicyRule{
		Name:        "SDP: CanUpdateService",
		Description: "Allows company managers (Owner, GM, BM) to update services.",
		Effect:      "Allow",
		EndPointID:  UpdateServiceById.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Any manager of the service's company
	}

	var AllowDeleteServiceById = &PolicyRule{
		Name:        "SDP: CanDeleteService",
		Description: "Allows company managers (Owner, GM, BM) to delete services.",
		Effect:      "Allow",
		EndPointID:  DeleteServiceById.ID,
		Conditions:  JsonRawMessage(company_manager_check), // Any manager of the service's company
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
		// AllowCreateSector,
		// AllowUpdateSectorById,
		// AllowDeleteSectorById,

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

	seededCount := 0
	updatedCount := 0 // Optionally track updates if you modify seeding logic

	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		AllowNilCompanyID = false
		AllowNilCreatedBy = false
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic occurred during policy seeding: %v", r)
		}
		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit transaction: %v", err)
		}
		log.Printf("System policies seeded successfully. New: %d, Existing/Updated: %d", seededCount, updatedCount)
	}()

	Policies = init_policy_array()

	for _, policy := range Policies {
		if policy == nil {
			log.Println("Warning: Encountered nil policy in Policies list.")
			continue
		}

		// Create a placeholder to find existing policy
		var existingPolicy PolicyRule
		// Find system policies specifically (company_id IS NULL)
		err := tx.Where("name = ? AND company_id IS NULL", policy.Name).First(&existingPolicy).Error

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
			updatedCount++ // Increment if you implement update logic
		}
	}

	return Policies, nil // Return the original list (or query fresh if needed)
}
