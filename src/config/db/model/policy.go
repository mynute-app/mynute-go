package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var AllowNilCompanyID = false
var AllowNilCreatedBy = false
var AllowNilResourceID = false

// --- PolicyRule (Represents a policy rule for access control) ---
type PolicyRule struct {
	BaseModel
	PropertyID  *uuid.UUID      `json:"property_id"`
	Property    *Property       `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"property"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	EndPoint    EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}

func (PolicyRule) TableName() string  { return "public.policy_rules" }
func (PolicyRule) SchemaType() string { return "public" }

func (PolicyRule) Indexes() map[string]string {
	return map[string]string{
		"idx_policy_company_endpoint": "CREATE INDEX idx_policy_company_endpoint ON policy_rules (company_id, end_point_id)",
	}
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
	Attribute string `json:"attribute"` // The primary attribute (e.g., "subject.role_id"/"resource.employee_id/"path.employee_id"/"body.employee_id"/"query.employee_id"/"header.employee_id")
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
		policyIdentifier := fmt.Sprintf("ID %s", p.ID.String()) // Use p.ID which should exist on BaseModel
		if p.Name != "" {
			policyIdentifier = fmt.Sprintf("'%s' (ID %s)", p.Name, p.ID.String())
		}
		return node, fmt.Errorf("policy rule %s has missing or null conditions", policyIdentifier)
	}

	// 2. Attempt to unmarshal the JSON
	err := json.Unmarshal(p.Conditions, &node)
	if err != nil {
		policyIdentifier := fmt.Sprintf("ID %s", p.ID.String())
		if p.Name != "" {
			policyIdentifier = fmt.Sprintf("'%s' (ID %s)", p.Name, p.ID.String())
		}
		return node, fmt.Errorf("failed to unmarshal conditions JSON for policy rule %s: %w", policyIdentifier, err)
	}

	// 3. Perform recursive validation using the dedicated validator function
	if err := validateConditionNode(node); err != nil {
		policyIdentifier := fmt.Sprintf("ID %s", p.ID.String())
		if p.Name != "" {
			policyIdentifier = fmt.Sprintf("'%s' (ID %s)", p.Name, p.ID.String())
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
			// Use nodeContext which already has Description included
			return fmt.Errorf("node%s is incorrectly structured as both leaf and branch", nodeContext)
		}

		// Rule 2: Leaf must have an attribute
		if node.Leaf.Attribute == "" {
			// Use nodeContext
			return fmt.Errorf("leaf node%s is missing required 'attribute'", nodeContext)
		}
		// Rule 3: Leaf must have an operator
		if node.Leaf.Operator == "" {
			// Use nodeContext and add attribute for more detail
			return fmt.Errorf("leaf node%s (attribute '%s') is missing required 'operator'", nodeContext, node.Leaf.Attribute)
		}

		// Rule 4: Operators needing comparison values must have 'value' or 'resource_attribute'
		requiresComparisonValue := true
		switch node.Leaf.Operator {
		case "IsNull", "IsNotNull":
			requiresComparisonValue = false
			// Add any other future unary operators here (e.g., "IsEmpty", "IsTrue")
		}

		if requiresComparisonValue && node.Leaf.Value == nil && node.Leaf.ResourceAttribute == "" {
			// Use nodeContext and add attribute/operator for more detail
			return fmt.Errorf("leaf node%s (attribute '%s', operator '%s') requires either 'value' or 'resource_attribute'", nodeContext, node.Leaf.Attribute, node.Leaf.Operator)
		}

		// Rule 5: Leaf cannot have *both* 'value' and 'resource_attribute'
		if node.Leaf.Value != nil && node.Leaf.ResourceAttribute != "" {
			// Use nodeContext and add attribute for more detail
			return fmt.Errorf("leaf node%s (attribute '%s') cannot have both 'value' and 'resource_attribute' defined", nodeContext, node.Leaf.Attribute)
		}

		// Rule 6: Basic sanity check on attribute format prefixes (Updated)
		validPrefixesList := "'subject.', 'resource.', 'path.', 'body.', 'query.', or 'header.'"
		if node.Leaf.Attribute != "" {
			// Check using helper function or inline checks:
			// if !isValidAttributePrefix(node.Leaf.Attribute) {
			if !(strings.HasPrefix(node.Leaf.Attribute, "subject.") ||
				strings.HasPrefix(node.Leaf.Attribute, "resource.") ||
				strings.HasPrefix(node.Leaf.Attribute, "path.") ||
				strings.HasPrefix(node.Leaf.Attribute, "body.") ||
				strings.HasPrefix(node.Leaf.Attribute, "query.")) {
				return fmt.Errorf("leaf node%s has invalid 'attribute' ('%s'): must start with one of %s", nodeContext, node.Leaf.Attribute, validPrefixesList)
			}
		}
		if node.Leaf.ResourceAttribute != "" {
			// Check using helper function or inline checks:
			// if !isValidAttributePrefix(node.Leaf.ResourceAttribute) {
			if !(strings.HasPrefix(node.Leaf.ResourceAttribute, "subject.") ||
				strings.HasPrefix(node.Leaf.ResourceAttribute, "resource.") ||
				strings.HasPrefix(node.Leaf.ResourceAttribute, "path.") ||
				strings.HasPrefix(node.Leaf.ResourceAttribute, "body.") ||
				strings.HasPrefix(node.Leaf.ResourceAttribute, "query.")) {
				return fmt.Errorf("leaf node%s has invalid 'resource_attribute' ('%s'): must start with one of %s", nodeContext, node.Leaf.ResourceAttribute, validPrefixesList)
			}
		}

	} else { // It's intended to be a branch node (or potentially an empty root?)

		// Rule 7: A branch must have a valid LogicType (AND/OR)
		isValidBranch := false
		if node.LogicType == "AND" || node.LogicType == "OR" {
			isValidBranch = true
			// Rule 8: A branch with a logic type must have children
			if len(node.Children) == 0 {
				// Use nodeContext
				return fmt.Errorf("branch node%s has 'logic_type' %s but no 'children'", nodeContext, node.LogicType)
			}
		}

		// Rule 9: Check if it's an invalid structure (neither leaf nor valid branch)
		if node.Leaf == nil && !isValidBranch {
			// Use nodeContext in error messages
			if node.LogicType == "" && len(node.Children) == 0 && node.Description == "" {
				// An absolutely empty node {} is likely an error
				return fmt.Errorf("node is empty and invalid (must be leaf or branch)")
			} else if node.LogicType == "" && len(node.Children) == 0 && node.Description != "" {
				// Possibly allow an empty descriptive node? Treat as error for now.
				return fmt.Errorf("node%s is descriptive only - not a valid leaf (missing 'leaf') nor a valid branch (missing 'logic_type' and 'children')", nodeContext)
			} else if node.LogicType != "" && !isValidBranch { // Invalid LogicType provided
				return fmt.Errorf("branch node%s has invalid 'logic_type': '%s' (must be AND or OR)", nodeContext, node.LogicType)
			} else {
				// General catch-all for invalid structure (e.g., missing logic_type but has children)
				return fmt.Errorf("node%s is neither a valid leaf nor a valid branch (check 'logic_type' and 'children')", nodeContext)
			}
		}

		// Rule 10: Recursively validate children if it's a valid branch
		if isValidBranch {
			for i, child := range node.Children {
				if err := validateConditionNode(child); err != nil {
					// Add context about which child failed
					return fmt.Errorf("invalid child node #%d within node%s: %w", i+1, nodeContext, err)
				}
			}
		}
	}

	// All checks passed for this node and its children (if any)
	return nil
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
					Attribute: "subject.roles[*].id",
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
			{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: JsonRawMessage(SystemRoleGeneralManager.ID)}},
			company_membership_access_check, // Re-use the company match check
		},
	}

	// Checks if subject is a Branch Manager of the resource's company
	var company_branch_manager_check = ConditionNode{
		Description: "Allow if Subject is Branch Manager within the Resource's Company",
		LogicType:   "AND",
		Children: []ConditionNode{
			{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: JsonRawMessage(SystemRoleBranchManager.ID)}},
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
				Description: "Branch Manager Assigned to Branch Check",
				LogicType:   "OR",
				Children: []ConditionNode{
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.branches",   // Assumes subject context has assigned branch IDs (e.g., [10, 25])
							Operator:          "Contains",           // Checks if the list contains the value
							ResourceAttribute: "resource.branch_id", // Assumes context provides resource.branch_id from the resource/path/body
							Description:       "Subject's assigned branches must include the resource's branch",
						},
					},
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.branches",
							Operator:          "Contains",
							ResourceAttribute: "path.branch_id", // Assumes context provides branch_id from the path parameter
							Description:       "Subject ID must match the path parameter branch_id",
						},
					},
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.branches",
							Operator:          "Contains",
							ResourceAttribute: "body.branch_id", // Assumes context provides branch_id from the body
							Description:       "Subject ID must match the body branch_id",
						},
					},
					{
						Leaf: &ConditionLeaf{
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
	var company_employee_assigned_employee_check = ConditionNode{
		Description: "Allow if Subject is the Employee associated with the Resource",
		LogicType:   "AND",
		Children: []ConditionNode{
			company_membership_access_check, // Must be in the same company
			{
				Description: "Employee ID Match Check",
				LogicType:   "OR",
				Children: []ConditionNode{
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.id",
							Operator:          "Equals",
							ResourceAttribute: "resource.employee_id", // Assumes context provides resource.employee_id from the resource
							Description:       "Subject ID must match the resource's employee ID",
						},
					},
					{
						Leaf: &ConditionLeaf{
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

	var company_employee_himself = ConditionNode{
		Description: "Company Employee Check (Employee accessing their own profile/resource)",
		LogicType:   "AND",
		Children: []ConditionNode{
			company_membership_access_check, // Must be in the same company
			{
				Description: "Employee id must be on resource, path, body, or query",
				LogicType:   "OR",
				Children: []ConditionNode{
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.id",
							Operator:          "Equals",
							ResourceAttribute: "path.employee_id", // Assumes context provides employee_id from the path parameter
							Description:       "Subject ID must match the path parameter employee_id",
						},
					},
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.id",
							Operator:          "Equals",
							ResourceAttribute: "body.employee_id", // Assumes context provides employee_id from the body
							Description:       "Subject ID must match the body employee_id",
						},
					},
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.id",
							Operator:          "Equals",
							ResourceAttribute: "query.employee_id", // Assumes context provides employee_id from the query parameter
							Description:       "Subject ID must match the query parameter employee_id",
						},
					},
					{
						Leaf: &ConditionLeaf{
							Attribute:         "subject.id",
							Operator:          "Equals",
							ResourceAttribute: "resource.employee_id", // Assumes context provides employee_id from the resource
							Description:       "The employee is accessing a resource that has himself assigned.",
						},
					},
					{
						Leaf: &ConditionLeaf{
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

	var company_admin_or_employee_himself_check = ConditionNode{
		Description: "Company Admin or Employee Check (Owner, GM, or Employee)",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_admin_check,      // Owner or GM of the resource's company
			company_employee_himself, // Employee can access their own profile/resource
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
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: JsonRawMessage(SystemRoleOwner.ID)}},
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: JsonRawMessage(SystemRoleGeneralManager.ID)}},
					{Leaf: &ConditionLeaf{Attribute: "subject.roles[*].id", Operator: "Contains", Value: JsonRawMessage(SystemRoleBranchManager.ID)}},
				},
			},
		},
	}

	var company_admin_or_assigned_branch_manager_check = ConditionNode{
		Description: "Company Admin or Assigned Branch Manager Check (Owner, GM, or BM in Company)",
		LogicType:   "OR",
		Children: []ConditionNode{
			company_admin_check,                          // Owner or GM of the resource's company
			company_branch_manager_assigned_branch_check, // BM assigned to the resource's branch
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
						{Leaf: &ConditionLeaf{Attribute: "subject.id", Operator: "Equals", ResourceAttribute: "body.client_id", Description: "Client ID in body must match Subject ID"}},
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
								company_employee_assigned_employee_check,     // Employee can create *for themselves*
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
	// var AllowCancelAppointmentByID = &PolicyRule{
	// 	Name:        "SDP: CanCancelAppointment",
	// 	Description: "Allows company managers or assigned employees to delete appointments. (Clients typically cannot delete).",
	// 	Effect:      "Allow",
	// 	EndPointID:  CancelAppointmentByID.ID,
	// 	Conditions: JsonRawMessage(ConditionNode{
	// 		Description: "Company User Delete Check",
	// 		LogicType:   "AND",
	// 		Children: []ConditionNode{
	// 			company_membership_access_check, // User in same company as appointment
	// 			{
	// 				Description: "Role/Relation Check (Managers or Assigned Employee)",
	// 				LogicType:   "OR",
	// 				Children: []ConditionNode{
	// 					company_owner_check,
	// 					company_general_manager_check,
	// 					company_branch_manager_assigned_branch_check, // BM can delete appointments in their branch
	// 					company_employee_assigned_employee_check,     // Employee can delete their own appointments (?) Maybe not - adjust if needed. Let's assume managers only.
	// 				},
	// 			},
	// 		},
	// 	}),
	// }

	// Policy: Allow DELETE appointment by ID.

	var AllowCancelAppointmentByID = &PolicyRule{
		Name:        "SDP: CanCancelAppointment",
		Description: "Allows clients to cancel their own appointments, or company managers/assigned employees.",
		Effect:      "Allow",
		EndPointID:  CancelAppointmentByID.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Allow Client Self-Cancel OR Company User Cancel",
			LogicType:   "OR",
			Children: []ConditionNode{
				client_access_check, // Client can cancel if it's their appointment
				{
					Description: "Company User Cancel Check",
					LogicType:   "AND",
					Children: []ConditionNode{
						company_membership_access_check, // User in same company as appointment
						{
							Description: "Role/Relation Check (Managers or Assigned Employee)",
							LogicType:   "OR",
							Children: []ConditionNode{
								company_owner_check,
								company_general_manager_check,
								company_branch_manager_assigned_branch_check, // BM can cancel appointments in their branch
								company_employee_assigned_employee_check,     // Employee can cancel their own appointments
							},
						},
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

	// var AllowGetBranchByName = &PolicyRule{
	// 	Name:        "SDP: CanViewBranchByName",
	// 	Description: "Allows any user belonging to the same company to view branch details by name.",
	// 	Effect:      "Allow",
	// 	EndPointID:  GetBranchByName.ID,
	// 	Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company can view
	// }

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

	var AllowCreateBranchWorkSchedule = &PolicyRule{
		Name:        "SDP: CanCreateBranchWorkSchedule",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to create work schedules for a branch.",
		Effect:      "Allow",
		EndPointID:  CreateBranchWorkSchedule.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Create Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can manage work schedules in any branch
				company_branch_manager_assigned_branch_check, // BM can manage work schedules in their own branch
			},
		}),
	}
	
	var AllowGetBranchWorkRangeById = &PolicyRule{
		Name:        "SDP: CanViewBranchWorkRangeById",
		Description: "Allows company members to view branch work schedules by ID.",
		Effect:      "Allow",
		EndPointID:  GetBranchWorkRange.ID,
		Conditions:  JsonRawMessage(company_internal_user_check), // Any internal user of the branch's company can view work schedules
	}

	var AllowDeleteBranchWorkRangeById = &PolicyRule{
		Name:        "SDP: CanDeleteBranchWorkRangeById",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to delete branch work schedules.",
		Effect:      "Allow",
		EndPointID:  DeleteBranchWorkRange.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Delete Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can delete work ranges in any branch
				company_branch_manager_assigned_branch_check, // BM can delete work ranges in their own branch
			},
		}),
	}

	var AllowUpdateBranchWorkRangeById = &PolicyRule{
		Name:        "SDP: CanUpdateBranchWorkRangeById",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branch work schedules.",
		Effect:      "Allow",
		EndPointID:  UpdateBranchWorkRange.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Update Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can update work ranges in any branch
				company_branch_manager_assigned_branch_check, // BM can update work ranges in their own branch
			},
		}),
	}

	var AllowAddBranchWorkRangeService = &PolicyRule{
		Name:        "SDP: CanAddBranchWorkRangeService",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to add services to a branch work range.",
		Effect:      "Allow",
		EndPointID:  AddBranchWorkRangeServices.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Add Service Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can add services in any branch
				company_branch_manager_assigned_branch_check, // BM can add services in their own branch
			},
		}),
	}

	var AllowDeleteBranchWorkRangeService = &PolicyRule{
		Name:        "SDP: CanDeleteBranchWorkRangeService",
		Description: "Allows company Owner, General Manager or assigned Branch Manager to remove services from a branch work range.",
		Effect:      "Allow",
		EndPointID:  DeleteBranchWorkRangeService.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Remove Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can remove work ranges in any branch
				company_branch_manager_assigned_branch_check, // BM can remove work ranges in their own branch
			},
		}),
	}

	var AllowUpdateBranchImages = &PolicyRule{
		Name:        "SDP: CanUpdateBranchImages",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branch images.",
		Effect:      "Allow",
		EndPointID:  UpdateBranchImages.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Update Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can update images in any branch
				company_branch_manager_assigned_branch_check, // BM can update images in their own branch
			},
		}),
	}

	var AllowDeleteBranchImage = &PolicyRule{
		Name:        "SDP: CanDeleteBranchImage",
		Description: "Allows company Owner, General Manager, or assigned Branch Manager to delete branch images.",
		Effect:      "Allow",
		EndPointID:  DeleteBranchImage.ID,
		Conditions: JsonRawMessage(ConditionNode{
			Description: "Admin or Assigned Branch Manager Delete Access",
			LogicType:   "OR",
			Children: []ConditionNode{
				company_admin_check,                          // Owner/GM can delete images in any branch
				company_branch_manager_assigned_branch_check, // BM can delete images in their own branch
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

	var AllowGetClientById = &PolicyRule{
		Name:        "SDP: CanViewClientById",
		Description: "Allows a client to view their own profile.",
		Effect:      "Allow",
		EndPointID:  GetClientById.ID,
		Conditions:  JsonRawMessage(client_self_access_check), // Client can view self (checks subject.id == resource.id)
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

	var AllowUpdateClientImages = &PolicyRule{
		Name:        "SDP: CanUpdateClientImages",
		Description: "Allows a client to update their own profile images.",
		Effect:      "Allow",
		EndPointID:  UpdateClientImages.ID,
		Conditions:  JsonRawMessage(client_self_access_check), // Client can update self images (checks subject.id == resource.id)
	}

	var AllowDeleteClientImage = &PolicyRule{
		Name:        "SDP: CanDeleteClientImage",
		Description: "Allows a client to delete their own profile images.",
		Effect:      "Allow",
		EndPointID:  DeleteClientImage.ID,
		Conditions:  JsonRawMessage(client_self_access_check), // Client can delete self images (checks subject.id == resource.id)
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

	var AllowUpdateCompanyImages = &PolicyRule{
		Name:        "SDP: CanUpdateCompanyImages",
		Description: "Allows company Owner or General Manager to update company images.",
		Effect:      "Allow",
		EndPointID:  UpdateCompanyImages.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Only Owner or GM of this company
	}

	var AllowDeleteCompanyImage = &PolicyRule{
		Name:        "SDP: CanDeleteCompanyImage",
		Description: "Allows company Owner or General Manager to delete company images.",
		Effect:      "Allow",
		EndPointID:  DeleteCompanyImage.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Only Owner or GM of this company
	}

	var AllowUpdateCompanyColors = &PolicyRule{
		Name:        "SDP: CanUpdateCompanyColors",
		Description: "Allows company Owner or General Manager to update company colors.",
		Effect:      "Allow",
		EndPointID:  UpdateCompanyColors.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Only Owner or GM of this company
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
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check),
	}

	var AllowCreateEmployeeWorkSchedule = &PolicyRule{
		Name:        "SDP: CanCreateEmployeeWorkSchedule",
		Description: "Allows employees, or company managers (Owner, GM, BM), to create their own work schedules.",
		Effect:      "Allow",
		EndPointID:  AddEmployeeWorkSchedule.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Employee can create own schedule
	}

	var AllowGetEmployeeWorkRangeById = &PolicyRule{
		Name:        "SDP: CanViewEmployeeWorkRangeById",
		Description: "Allows employees, or company managers (Owner, GM, BM), to view their own work ranges.",
		Effect:      "Allow",
		EndPointID:  GetEmployeeWorkRange.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check),
	}

	var AllowUpdateEmployeeWorkRange = &PolicyRule{
		Name:        "SDP: CanUpdateEmployeeWorkRange",
		Description: "Allows employees, or company managers (Owner, GM, BM), to update their own work ranges.",
		Effect:      "Allow",
		EndPointID:  UpdateEmployeeWorkRange.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Employee can update own work range
	}

	var AllowDeleteEmployeeWorkRange = &PolicyRule{
		Name:        "SDP: CanDeleteEmployeeWorkRange",
		Description: "Allows employees, or company managers (Owner, GM, BM), to remove their own work ranges.",
		Effect:      "Allow",
		EndPointID:  DeleteEmployeeWorkRange.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Employee can remove own work range
	}

	var AllowAddEmployeeWorkRangeServices = &PolicyRule{
		Name:        "SDP: CanAddEmployeeWorkRangeServices",
		Description: "Allows employees, or company managers (Owner, GM, BM), to add services to their work ranges.",
		Effect:      "Allow",
		EndPointID:  AddEmployeeWorkRangeServices.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Employee can add services to own work range
	}

	var AllowDeleteEmployeeWorkRangeService = &PolicyRule{
		Name:        "SDP: CanDeleteEmployeeWorkRangeService",
		Description: "Allows employees, or company managers (Owner, GM, BM), to remove services from their work ranges.",
		Effect:      "Allow",
		EndPointID:  DeleteEmployeeWorkRangeService.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Employee can remove services from own work range
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
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Manager of the employee's company
	}

	var AllowRemoveServiceFromEmployee = &PolicyRule{
		Name:        "SDP: CanRemoveServiceFromEmployee",
		Description: "Allows company managers (Owner, GM, BM) to remove services from employees.",
		Effect:      "Allow",
		EndPointID:  RemoveServiceFromEmployee.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Manager of the employee's company
	}

	var AllowAddBranchToEmployee = &PolicyRule{
		Name:        "SDP: CanAddBranchToEmployee",
		Description: "Allows company managers (Owner, GM, BM) to assign employees to branches (respecting BM scope).",
		Effect:      "Allow",
		EndPointID:  AddBranchToEmployee.ID,
		Conditions:  JsonRawMessage(company_admin_or_assigned_branch_manager_check),
	}

	var AllowRemoveBranchFromEmployee = &PolicyRule{
		Name:        "SDP: CanRemoveBranchFromEmployee",
		Description: "Allows company managers (Owner, GM, BM) to remove employees from branches (respecting BM scope).",
		Effect:      "Allow",
		EndPointID:  RemoveBranchFromEmployee.ID,
		Conditions:  JsonRawMessage(company_admin_or_assigned_branch_manager_check),
	}

	var AllowUpdateEmployeeImages = &PolicyRule{
		Name:        "SDP: CanUpdateEmployeeImages",
		Description: "Allows company Owner, General Manager, or employee himself to update employee images.",
		Effect:      "Allow",
		EndPointID:  UpdateEmployeeImages.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check),
	}

	var AllowDeleteEmployeeImage = &PolicyRule{
		Name:        "SDP: CanDeleteEmployeeImage",
		Description: "Allows company Owner, General Manager, or employee himself to delete employee images.",
		Effect:      "Allow",
		EndPointID:  DeleteEmployeeImage.ID,
		Conditions:  JsonRawMessage(company_admin_or_employee_himself_check), // Employee can delete own images
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

	var AllowUpdateServiceImages = &PolicyRule{
		Name:        "SDP: CanUpdateServiceImages",
		Description: "Allows company managers (Owner, GM) to update service images.",
		Effect:      "Allow",
		EndPointID:  UpdateServiceImages.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Any manager of the service's company
	}

	var AllowDeleteServiceImage = &PolicyRule{
		Name:        "SDP: CanDeleteServiceImage",
		Description: "Allows company managers (Owner, GM) to delete service images.",
		Effect:      "Allow",
		EndPointID:  DeleteServiceImage.ID,
		Conditions:  JsonRawMessage(company_admin_check), // Any manager of the service's company
	}

	// --- Combined Policies List ---
	var Policies = []*PolicyRule{
		// Appointments
		AllowGetAppointmentByID,
		AllowCreateAppointment,
		AllowUpdateAppointmentByID,
		AllowCancelAppointmentByID,

		// Branches
		AllowCreateBranch,
		AllowGetBranchById,
		// AllowGetBranchByName,
		AllowUpdateBranchById,
		AllowDeleteBranchById,
		AllowGetEmployeeServicesByBranchId,
		AllowAddServiceToBranch,
		AllowRemoveServiceFromBranch,
		AllowUpdateBranchImages,
		AllowDeleteBranchImage,
		AllowCreateBranchWorkSchedule,
		AllowDeleteBranchWorkRangeById,
		AllowUpdateBranchWorkRangeById,
		AllowGetBranchWorkRangeById,
		AllowAddBranchWorkRangeService,
		AllowDeleteBranchWorkRangeService,

		// Clients (Self-Management focused)
		AllowGetClientByEmail,
		AllowUpdateClientById,
		AllowDeleteClientById,
		AllowGetClientById,
		AllowUpdateClientImages,
		AllowDeleteClientImage,

		// Company
		AllowGetCompanyById,
		AllowUpdateCompanyById,
		AllowDeleteCompanyById,
		AllowUpdateCompanyImages,
		AllowUpdateCompanyColors,
		AllowDeleteCompanyImage,

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
		AllowCreateEmployeeWorkSchedule,
		AllowGetEmployeeWorkRangeById,
		AllowUpdateEmployeeWorkRange,
		AllowDeleteEmployeeWorkRange,
		AllowAddEmployeeWorkRangeServices,
		AllowDeleteEmployeeWorkRangeService,
		AllowUpdateEmployeeImages,
		AllowDeleteEmployeeImage,

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
		AllowUpdateServiceImages,
		AllowDeleteServiceImage,
	}

	return Policies
}

type PolicyCfg struct {
	AllowNilCompanyID bool // Allow policies to be created without company_id
	AllowNilCreatedBy bool // Allow policies to be created without created_by
}

func Policies(cfg *PolicyCfg) ([]*PolicyRule, func()) {
	policies := init_policy_array()
	AllowNilCompanyID = cfg.AllowNilCompanyID
	AllowNilCreatedBy = cfg.AllowNilCreatedBy
	deferFnc := func() {
		AllowNilCompanyID = false
		AllowNilCreatedBy = false
	}
	return policies, deferFnc
}
