package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// --- Policy (Represents a policy for access control) ---
type Policy struct {
	BaseModel
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	EndPoint    EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}

// --- ConditionNode (Represents a logical grouping OR a leaf check) ---
type ConditionNode struct {
	// Description explains the purpose of this node
	// (e.g., "Client Access Check", "Company Membership Check")
	Description string `json:"description,omitempty"`
	// --- Fields for Branch Nodes (Logical Operators: AND, OR, NOT) ---
	// LogicType specifies how to evaluate Children ("AND", "OR").
	// Required for branch nodes.
	LogicType string `json:"logic_type,omitempty"`
	// Children contains nested nodes to be evaluated.
	// Required for branch nodes.
	Children []ConditionNode `json:"children,omitempty"`
	// --- Field for Leaf Nodes (Actual Condition Check) ---
	// Leaf points to the actual condition details.
	// Required for leaf nodes.
	Leaf *ConditionLeaf `json:"leaf,omitempty"`
}

// --- ConditionLeaf (Represents a single atomic check) ---
type ConditionLeaf struct {
	// Attribute is the primary attribute
	// e.g., "subject.role_id","resource.employee_id",
	// "path.employee_id","body.employee_id",
	// "query.employee_id","header.employee_id"
	Attribute string `json:"attribute"`
	// Operator defines the comparison operation
	// e.g., "Equals" - "==", "NotEquals" - "!=",
	// "IsNull" - "== null", "IsNotNull" - "!= null"
	Operator string `json:"operator"`
	// Description is a human-readable explanation of
	// the check being performed.
	// e.g., "Subject must be a Client"
	Description string `json:"description,omitempty"`
	// --- Use EITHER Value OR ResourceAttribute ---
	// Value is a static JSON value to compare against.
	Value json.RawMessage `json:"value,omitempty"`
	// ResourceAttribute is the name of a dynamic attribute
	// that can be used for comparison.
	ResourceAttribute string `json:"resource_attribute,omitempty"`
}

// GetConditionsNode parses and validates the stored JSON conditions
func (p *Policy) GetConditionsNode() (ConditionNode, error) {
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