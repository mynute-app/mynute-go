package handler

import (
	"agenda-kaki-go/core/config/db/model"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Policy struct {
	DB *gorm.DB
}

type AccessDecision struct {
	Allowed bool
	Reason  string // Detailed reason for denial, populated if Allowed is false and Error is nil
	Error   error  // Processing error
}

func NewPolicyEngine(db *gorm.DB) *Policy {
	return &Policy{DB: db}
}

func (p *Policy) CanAccess(subject, resource map[string]any, policy *model.PolicyRule) AccessDecision {
	decision := AccessDecision{}
	if policy == nil {
		decision.Error = fmt.Errorf("policy is nil")
		return decision
	}
	root, err := policy.GetConditionsNode()
	if err != nil {
		decision.Error = fmt.Errorf("failed to get conditions node for policy '%s': %w", policy.Name, err)
		return decision
	}
	result, reason, err := evalNode(root, subject, resource)
	if err != nil {
		decision.Error = fmt.Errorf("error evaluating conditions for policy '%s': %w", policy.Name, err)
		return decision
	}
	if policy.Effect == "Allow" {
		decision.Allowed = result
		if !decision.Allowed {
			if reason != "" {
				decision.Reason = fmt.Sprintf("Policy [Allow] '%s' denied because: %s", policy.Name, reason)
			} else {
				decision.Reason = fmt.Sprintf("Policy [Allow] '%s' conditions not met.", policy.Name)
			}
		}
		// If Allowed=true, Reason remains ""
	} else if policy.Effect == "Deny" {
		decision.Allowed = !result // Access is allowed IF the Deny condition is FALSE
		if !decision.Allowed {
			if reason != "" {
				decision.Reason = fmt.Sprintf("Policy [Deny] '%s' enforced because condition met: %s", policy.Name, reason)
			} else {
				decision.Reason = fmt.Sprintf("Policy [Deny] '%s' conditions were met.", policy.Name)
			}
		}
	} else {
		decision.Error = fmt.Errorf("unknown policy effect '%s' in policy '%s'", policy.Effect, policy.Name)
	}

	return decision
}

func evalNode(node model.ConditionNode, subject, resource map[string]any) (bool, string, error) {
	nodeDescription := ""
	if node.Description != "" {
		nodeDescription = fmt.Sprintf("'%s'", node.Description)
	}

	if node.Leaf != nil {
		// Evaluate leaf node
		res, reason, err := evalLeaf(*node.Leaf, subject, resource)
		if err != nil {
			return false, "", fmt.Errorf("leaf node `%s` evaluation failed: %w", nodeDescription, err) // Propagate processing error
		}
		if !res {
			return false, reason, nil
		}
		return true, "", nil
	}

	if len(node.Children) == 0 {
		// Should be caught by validation, but handle defensively. Branch node needs children.
		return false, "", fmt.Errorf("node `%s` has logic type '%s' but no children", nodeDescription, node.LogicType)
	}

	// Evaluate children for branch node
	childResults := make([]bool, len(node.Children))
	childReasons := make([]string, len(node.Children)) // Store reasons for potential use (especially for OR)

	for i, child := range node.Children {
		res, reason, err := evalNode(child, subject, resource)
		if err != nil {
			// Processing error in a child, propagate up
			return false, "", fmt.Errorf("child node evaluation failed within %s: %w", nodeDescription, err)
		}
		childResults[i] = res
		if !res {
			childReasons[i] = reason // Store reason only if child failed
		} else {
			childReasons[i] = "" // Clear reason if child succeeded
		}

		// Short-circuit logic
		switch node.LogicType {
		case "AND": // AND fails if any child fails. Return immediately with the reason.
			if !res {
				finalReason := fmt.Sprintf("node `%s` with logic type `AND` failed because: %s", nodeDescription, reason)
				return false, finalReason, nil
			}
		case "OR": // OR succeeds if any child succeeds. Return immediately.
			if res {
				return true, "", nil
			}
		}
	}

	switch node.LogicType {
	case "AND": // If we reached here, all children were true.
		return true, "", nil
	case "OR": // If we reached here, all children were false.
		// Combine descriptions if available.
		var reasons []string
		for i, res := range childResults {
			if !res && childReasons[i] != "" {
				reasons = append(reasons, childReasons[i])
			} else if !res {
				reasons = append(reasons, fmt.Sprintf("child condition %d evaluated to false but could not evaluate the reason why", i+1))
			}
		}
		combinedReason := strings.Join(reasons, "; AND ")
		finalReason := fmt.Sprintf("OR node `%s` failed because all conditions were false: [%s]", nodeDescription, combinedReason)
		return false, finalReason, nil
	default:
		return false, "", fmt.Errorf("node `%s` has unknown logic type: %s", nodeDescription, node.LogicType)
	}
}

func evalLeaf(leaf model.ConditionLeaf, subject, resource map[string]any) (bool, string, error) {
	leafDescription := ""
	if leaf.Description != "" {
		leafDescription = fmt.Sprintf(" ('%s')", leaf.Description) // Add parens for clarity
	}

	left, err := resolveAttr(leaf.Attribute, subject, resource)
	if err != nil {
		// Check if error is because attribute doesn't exist (common, maybe not an "error" but a false condition)
		// For now, treat it as processing error, but could be refined later.
		return false, "", fmt.Errorf("failed to resolve left attribute '%s'%s: %w", leaf.Attribute, leafDescription, err)
	}

	var right any
	var rightSource string // For reasoning string
	if leaf.ResourceAttribute != "" {
		rightSource = fmt.Sprintf("resource attribute '%s'", leaf.ResourceAttribute)
		right, err = resolveAttr(leaf.ResourceAttribute, subject, resource)
		if err != nil {
			return false, "", fmt.Errorf("failed to resolve right attribute '%s'%s: %w", leaf.ResourceAttribute, leafDescription, err)
		}
	} else if leaf.Value != nil {
		if err := json.Unmarshal(leaf.Value, &right); err != nil {
			return false, "", fmt.Errorf("invalid static value for comparison%s: %w", leafDescription, err)
		}
		rightSource = fmt.Sprintf("static value '%v'", right) // Use resolved value for clarity
		// Mask sensitive values if necessary here
		// Example: if leaf.Attribute contains "password", set rightSource = "static value [REDACTED]"
	} else if leaf.Operator != "IsNull" && leaf.Operator != "IsNotNull" {
		// If neither Value nor ResourceAttribute is set, it's only valid for IsNull/IsNotNull
		return false, "", fmt.Errorf("condition%s requires either 'value' or 'resource_attribute' for operator '%s'", leafDescription, leaf.Operator)
	}

	// Helper for creating reason string
	failReason := func(opResult bool, format string, args ...any) (bool, string, error) {
		if opResult {
			return true, "", nil // Condition passed, no reason needed
		}
		// Condition failed, generate reason
		reason := fmt.Sprintf(format, args...)
		return false, fmt.Sprintf("Condition `%s` failed: %s", leafDescription, reason), nil
	}

	switch leaf.Operator {
	case "Equals":
		res := reflect.DeepEqual(left, right)
		return failReason(res, "attribute '%s' (%v) should be equals to %s (%v)", leaf.Attribute, left, rightSource, right)
	case "IsNull":
		res := left == nil
		return failReason(res, "attribute '%s' (%v) should be null", leaf.Attribute, left)
	case "IsNotNull":
		res := left != nil
		return failReason(res, "attribute '%s' (%v) should not be null", leaf.Attribute, left) // Value isn't relevant if it was null
	case "Contains":
		// Specific checks and reasoning for Contains
		val := reflect.ValueOf(left)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			// This is a processing error, not just a false condition
			return false, "", fmt.Errorf("attribute '%s' is not a slice nor array. 'Contains' operator failed for leaf `%s`", leaf.Attribute, leafDescription)
		}
		found := false
		for i := range val.Len() {
			if reflect.DeepEqual(val.Index(i).Interface(), right) {
				found = true
				break
			}
		}
		return failReason(found, "attribute '%s' (%v) should contain %s", leaf.Attribute, left, rightSource)
	
	case "GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual":
		res, err := compareNumbers(left, right, leaf.Operator)
		if err != nil {
			return false, "", fmt.Errorf("numeric comparison failed: %s: %w", leafDescription, err) // Processing error
		}
		return failReason(res, "numeric comparison '%s %s %s' failed for values %v and %v", leaf.Attribute, leaf.Operator, rightSource, left, right)

	case "StartsWith", "EndsWith", "Includes":
		res, err := compareStrings(left, right, leaf.Operator)
		if err != nil {
			return false, "", fmt.Errorf("string comparison failed for leaf `%s`: %w", leafDescription, err) // Processing error
		}
		opDesc := map[string]string{"StartsWith": "start with", "EndsWith": "end with", "Includes": "include"}[leaf.Operator]
		return failReason(res, "string attribute '%s' (%q) did not %s %s (%q)", leaf.Attribute, left, opDesc, rightSource, right)

	case "Before", "After":
		res, err := compareTimes(left, right, leaf.Operator)
		if err != nil {
			return false, "", fmt.Errorf("time comparison failed for leaf `%s`: %w", leafDescription, err) // Processing error
		}
		opDesc := map[string]string{"Before": "before", "After": "after"}[leaf.Operator]
		// Format times for readability in reason
		lTime, _ := toTime(left) // Ignore error as already checked in compareTimes
		rTime, _ := toTime(right)
		return failReason(res, "time attribute '%s' (%s) was not %s %s (%s)", leaf.Attribute, lTime.Format(time.RFC3339), opDesc, rightSource, rTime.Format(time.RFC3339))

	default:
		// Unknown operator is a processing error
		return false, "", fmt.Errorf("unsupported operator '%s' for leaf `%s`", leaf.Operator, leafDescription)
	}
}

func resolveAttr(attr string, subject, resource map[string]any) (any, error) {
	if attr == "" {
		return nil, errors.New("empty attr")
	}
	switch {
	case len(attr) > 8 && attr[:8] == "subject.":
		return subject[attr[8:]], nil
	case len(attr) > 9 && attr[:9] == "resource.":
		return resource[attr[9:]], nil
	default:
		return nil, fmt.Errorf("invalid attr prefix: %s", attr)
	}
}

func compareNumbers(left, right any, op string) (bool, error) {
	leftFloat, ok1 := toFloat(left)
	rightFloat, ok2 := toFloat(right)

	if !ok1 || !ok2 {
		return false, fmt.Errorf("cannot compare non-numeric values: %v (%T) and %v (%T)", left, left, right, right)
	}

	switch op {
	case "GreaterThan":
		return leftFloat > rightFloat, nil
	case "GreaterThanOrEqual":
		return leftFloat >= rightFloat, nil
	case "LessThan":
		return leftFloat < rightFloat, nil
	case "LessThanOrEqual":
		return leftFloat <= rightFloat, nil
	default:
		return false, fmt.Errorf("unknown numeric comparison: %s", op)
	}
}

func toFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case json.Number:
		f, err := v.Float64()
		if err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func compareStrings(left, right any, op string) (bool, error) {
	lStr, ok1 := left.(string)
	rStr, ok2 := right.(string)
	if !ok1 || !ok2 {
		return false, fmt.Errorf("invalid string comparison: %v and %v", left, right)
	}
	switch op {
	case "StartsWith":
		return strings.HasPrefix(lStr, rStr), nil
	case "EndsWith":
		return strings.HasSuffix(lStr, rStr), nil
	case "Includes":
		return strings.Contains(lStr, rStr), nil
	default:
		return false, fmt.Errorf("unsupported string op: %s", op)
	}
}

func compareTimes(left, right any, op string) (bool, error) {
	lTime, err := toTime(left)
	if err != nil {
		return false, fmt.Errorf("left value not a valid time: %w", err)
	}
	rTime, err := toTime(right)
	if err != nil {
		return false, fmt.Errorf("right value not a valid time: %w", err)
	}

	switch op {
	case "Before":
		return lTime.Before(rTime), nil
	case "After":
		return lTime.After(rTime), nil
	default:
		return false, fmt.Errorf("unsupported time op: %s", op)
	}
}

func toTime(val any) (time.Time, error) {
	switch v := val.(type) {
	case time.Time:
		return v, nil
	case string:
		return time.Parse(time.RFC3339, v)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(i, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported time type: %T", val)
	}
}
