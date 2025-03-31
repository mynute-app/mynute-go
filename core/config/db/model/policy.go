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
	ResourceID          uint            `json:"resource_id"` // Link to Resource definition
	Resource            Resource        `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE;" json:"resource"`
	Conditions          json.RawMessage `gorm:"type:jsonb" json:"conditions"`
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
	if len(p.Conditions) == 0 || string(p.Conditions) == "null" { // Check for empty or explicit null
		// Decide behavior: Is an empty condition set allowed? Does it mean always true/false?
		// Returning an empty node might be acceptable if evaluator handles it.
		// Or return an error if conditions are mandatory.
		return node, fmt.Errorf("policy rule %d has missing or null conditions", p.ID) // Or return (node, nil) if allowed
	}
	err := json.Unmarshal(p.Conditions, &node)
	if err != nil {
		return node, fmt.Errorf("failed to unmarshal conditions JSON for policy rule %d: %w", p.ID, err)
	}
	// Basic validation: Ensure a node is either a valid branch or a valid leaf
	if node.Leaf == nil && (node.LogicType == "" || len(node.Children) == 0) {
		// It's not a leaf, but lacks LogicType or Children - likely invalid structure
		// However, an empty top-level object "{}" might parse without error but be invalid policy.
		// Consider adding more validation if needed, based on how robust parsing/validation should be.
		// For now, relying on evaluator logic to handle potentially underspecified nodes.
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
func MustMarshalJSON(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("MustMarshalJSON failed: %v", err))
	}
	return json.RawMessage(data)
}

var client_access_check = ConditionNode{
	Description: "Client Access Check",
	LogicType:   "AND",
	Children: []ConditionNode{
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "IsNull", Description: "Subject must not belong to a company"}},
		{Leaf: &ConditionLeaf{Attr: "subject.id", Op: "Equals", ValueSourceAttr: "resource.client_id", Description: "Subject ID must match appointment client ID"}},
	},
}

var company_membership_access_check = ConditionNode{
	Description: "Company Membership Check",
	LogicType:   "AND",
	Children: []ConditionNode{
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "IsNotNull", Description: "Subject must belong to a company"}},
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "Equals", ValueSourceAttr: "resource.company_id", Description: "Subject company must match resource company"}},
	},
}

var company_owner_check = ConditionNode{
	Description: "Allow if subject is the company owner",
	LogicType:   "AND",
	Children: []ConditionNode{
		{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleOwner.ID), Description: "Subject role must be Owner"}},
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "Equals", ValueSourceAttr: "resource.company_id", Description: "Subject company must match resource company"}},
	},
}

var company_general_manager_check = ConditionNode{
	Description: "Allow if subject is the company general manager",
	LogicType:   "AND",
	Children: []ConditionNode{
		{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleGeneralManager.ID), Description: "Subject role must be General Manager"}},
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "Equals", ValueSourceAttr: "resource.company_id", Description: "Subject company must match resource company"}},
	},
}

var company_branch_manager_assigned_branch_check = ConditionNode{
	Description: "Allow if subject is the company branch manager and the resource has the branch ID assigned to it",
	LogicType:   "AND",
	Children: []ConditionNode{
		{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleBranchManager.ID), Description: "Subject role must be Branch Manager"}},
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "Equals", ValueSourceAttr: "resource.company_id", Description: "Subject company must match resource company"}},
		{Leaf: &ConditionLeaf{Attr: "subject.branches", Op: "Contains", ValueSourceAttr: "resource.branch_id", Description: "Subject branches must contain resource branch"}},
	},
}

var company_employee_assigned_employee_check = ConditionNode{
	Description: "Allow if subject is the company employee and the resource has the employee ID assigned to it",
	LogicType:   "AND",
	Children: []ConditionNode{
		{Leaf: &ConditionLeaf{Attr: "subject.role_id", Op: "Equals", Value: MustMarshalJSON(SystemRoleEmployee.ID), Description: "Subject role must be Employee"}},
		{Leaf: &ConditionLeaf{Attr: "subject.company_id", Op: "Equals", ValueSourceAttr: "resource.company_id", Description: "Subject company must match resource company"}},
		{Leaf: &ConditionLeaf{Attr: "subject.id", Op: "Equals", ValueSourceAttr: "resource.employee_id", Description: "Subject ID must match resource employee ID"}},
	},
}
// Policy: Allow Create appointment by ID.
var AllowCreateAppointment = &PolicyRule{
	Name:        "SDP: CanCreateAppointment",
	Description: "System Default Policy: Allows clients to create appointments, or company users based on role/relation.",
	Effect:      "Allow",
	ResourceID:  CreateAppointment.ID, // Link to the specific Resource definition
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Top-level OR: Allow Client Access OR Company User Access", // Top-level description
		LogicType:   "OR",
		Children: []ConditionNode{
			client_access_check,
			{
				Description: "Company User Access Check",
				LogicType:   "AND",
				Children: []ConditionNode{
					// Part 2.A: Company Check
					company_membership_access_check,
					// Part 2.B: Role/Relation Check
					{
						Description: "Employee Role/Relation Check",
						LogicType:   "OR",
						Children: []ConditionNode{
							company_owner_check,
							company_general_manager_check,
							company_branch_manager_assigned_branch_check,
							company_employee_assigned_employee_check,
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
	Description: "System Default Policy: Allows clients to view own appointments, or company users based on role/relation.",
	Effect:      "Allow",
	ResourceID:  GetAppointmentByID.ID, // Link to the specific Resource definition
	Conditions: MustMarshalJSON(ConditionNode{
		Description: "Top-level OR: Allow Client Access OR Company User Access", // Top-level description
		LogicType:   "OR",
		Children: []ConditionNode{
			client_access_check,
			{
				Description: "Company User Access Check",
				LogicType:   "AND",
				Children: []ConditionNode{
					// Part 2.A: Company Check
					company_membership_access_check,
					// Part 2.B: Role/Relation Check
					{
						Description: "Employee Role/Relation Check",
						LogicType:   "OR",
						Children: []ConditionNode{
							company_owner_check,
							company_general_manager_check,
							company_branch_manager_assigned_branch_check,
							company_employee_assigned_employee_check,
						},
					},
				},
			},
		},
	}),
}

var Policies = []*PolicyRule{
	AllowGetAppointmentByID,
}

func SeedPolicies(db *gorm.DB) ([]*PolicyRule, error) {
	AllowNilCompanyID = true
	AllowNilCreatedBy = true
	defer func() {
		AllowNilCompanyID = false
		AllowNilCreatedBy = false
	}()
	for _, policy := range Policies {
		err := db.Where("name = ? AND company_id IS NULL", policy.Name).First(policy).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(policy).Error; err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}
	log.Println("System policies seeded successfully!")
	return Policies, nil
}
