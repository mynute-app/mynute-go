package handler

import (
	"agenda-kaki-go/core/config/db/model"
	"encoding/json"
	"errors"
	"fmt"
    "math" // Added for toTime float handling
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

// Indentation helper string
const indentStep = "  "

func (p *Policy) CanAccess(subject, resource map[string]any, policy *model.PolicyRule) AccessDecision {
	decision := AccessDecision{}
	if policy == nil {
		decision.Error = errors.New("policy rule cannot be nil") // Use errors.New for simple static errors
		return decision
	}
	// Ensure policy Name is available for errors/reasons
	policyIdentifier := fmt.Sprintf("ID %s", policy.ID.String()) // Assuming p.ID exists
	if policy.Name != "" {
		policyIdentifier = fmt.Sprintf("'%s' (ID %s)", policy.Name, policy.ID.String())
	}

	root, err := policy.GetConditionsNode()
	if err != nil {
		// ***** CORRECTED: Use fmt.Errorf *****
		decision.Error = fmt.Errorf("failed to get/validate conditions node for policy %s: %w", policyIdentifier, err)
		return decision
	}

	// Start evaluation with depth 0
	result, reason, err := evalNode(root, subject, resource, 0) // Pass initial depth
	if err != nil {
        // ***** CORRECTED: Use fmt.Errorf *****
		decision.Error = fmt.Errorf("error evaluating conditions for policy %s: %w", policyIdentifier, err)
		return decision
	}

	effectDesc := fmt.Sprintf("[%s]", policy.Effect)

	if policy.Effect == "Allow" {
		decision.Allowed = result
		if !decision.Allowed {
			if reason != "" {
				// Add a newline before the detailed reason for better separation
				decision.Reason = fmt.Sprintf("Policy %s %s denied:\n%s", effectDesc, policyIdentifier, reason)
			} else {
				// Fallback should be rare now but keep it safe
				decision.Reason = fmt.Sprintf("Policy %s %s conditions not met (no specific reason provided).", effectDesc, policyIdentifier)
			}
		}
	} else if policy.Effect == "Deny" {
		decision.Allowed = !result // Access allowed if Deny conditions are FALSE
		if !decision.Allowed { // Access denied if Deny conditions are TRUE
			if reason != "" {
				decision.Reason = fmt.Sprintf("Policy %s %s enforced because condition evaluated to true:\n%s", effectDesc, policyIdentifier, reason)
			} else {
				decision.Reason = fmt.Sprintf("Policy %s %s conditions were met (no specific reason provided).", effectDesc, policyIdentifier)
			}
		}
	} else {
        // ***** CORRECTED: Use fmt.Errorf *****
		decision.Error = fmt.Errorf("unknown policy effect '%s' in policy %s", policy.Effect, policyIdentifier)
	}

	return decision
}

// evalNode now accepts depth for indentation
func evalNode(node model.ConditionNode, subject, resource map[string]any, depth int) (bool, string, error) {
	indent := strings.Repeat(indentStep, depth)
	nodeDescription := ""
	if node.Description != "" {
		nodeDescription = fmt.Sprintf("'%s'", node.Description)
	} else {
        nodeDescription = "(unnamed node)" // Handle unnamed nodes
    }

	if node.Leaf != nil {
        // Evaluate leaf node, passing depth
		res, reason, err := evalLeaf(*node.Leaf, subject, resource, depth) // Pass depth
		if err != nil {
            // Add indentation to error context if needed, though errors usually bubble up raw
			return false, "", fmt.Errorf("%sLeaf Node %s evaluation failed: %w", indent, nodeDescription, err)
		}
		// Leaf reason already includes indentation if it failed (!res)
		return res, reason, nil
	}

    // Validate branch node structure (should be caught by GetConditionsNode but belt-and-suspenders)
	if len(node.Children) == 0 {
		return false, "", fmt.Errorf("%sNode %s has logic type '%s' but no children", indent, nodeDescription, node.LogicType)
	}
	if node.LogicType != "AND" && node.LogicType != "OR" {
		return false, "", fmt.Errorf("%sNode %s has unknown logic type: %s", indent, nodeDescription, node.LogicType)
	}


	// --- Evaluate Branch Node Children ---
	childResults := make([]bool, len(node.Children))
	childReasons := make([]string, len(node.Children))

	for i, child := range node.Children {
        // Evaluate child with increased depth
		res, reason, err := evalNode(child, subject, resource, depth+1)
		if err != nil {
            // Error detail will come from lower levels
			return false, "", fmt.Errorf("%sChild evaluation failed within %s node: %w", indent, nodeDescription, err)
		}
		childResults[i] = res
		if !res {
			childReasons[i] = reason // Store failure reason (already includes its own indent)
		} else {
            childReasons[i] = "" // Not needed if child passed
        }


		// --- Short-circuit logic ---
		switch node.LogicType {
		case "AND":
			if !res {
				// AND fails: Use the reason from the first failing child.
				// The child's reason already has the deeper indent.
				finalReason := fmt.Sprintf("%s- [%s Node: %s] failed because child condition failed:\n%s", indent, node.LogicType, nodeDescription, reason)
				return false, finalReason, nil
			}
		case "OR":
			if res {
				// OR succeeds on first success. No reason needed.
				return true, "", nil
			}
		}
	}


	// --- Process final results (if no short-circuit occurred) ---
	switch node.LogicType {
	case "AND":
		// If we reached here, all children were true. Node succeeds.
		return true, "", nil
	case "OR":
		// If we reached here, all children were false. Node fails.
		var failingReasons []string
		for _, r := range childReasons {
			if r != "" { // Only collect non-empty reasons (failures)
				failingReasons = append(failingReasons, r)
			}
		}
        // Combine failing child reasons, separated by newline. They already have correct indentation.
		combinedReason := strings.Join(failingReasons, "\n")
		finalReason := fmt.Sprintf("%s- [%s Node: %s] failed because all child conditions were false:\n%s", indent, node.LogicType, nodeDescription, combinedReason)
		return false, finalReason, nil
	default:
        // Should be unreachable due to earlier check, but satisfy compiler
         return false, "", fmt.Errorf("%sInternal Error: Unhandled logic type %s in node %s after child evaluation", indent, node.LogicType, nodeDescription)
	}
}

// evalLeaf now accepts depth for indentation
func evalLeaf(leaf model.ConditionLeaf, subject, resource map[string]any, depth int) (bool, string, error) {
	indent := strings.Repeat(indentStep, depth)
	leafDescription := ""
	if leaf.Description != "" {
		leafDescription = fmt.Sprintf("'%s' ", leaf.Description) // Note the trailing space
	}

	// --- Resolve Attributes ---
	left, err := resolveAttr(leaf.Attribute, subject, resource)
	if err != nil {
		return false, "", fmt.Errorf("%sFailed to resolve attribute '%s' for condition %s: %w", indent, leaf.Attribute, leafDescription, err)
	}

	var right any
	var rightSourceDesc string // Describes the source of the right-hand value
	if leaf.ResourceAttribute != "" {
		rightSourceDesc = fmt.Sprintf("resource attribute '%s'", leaf.ResourceAttribute)
		right, err = resolveAttr(leaf.ResourceAttribute, subject, resource)
		if err != nil {
			return false, "", fmt.Errorf("%sFailed to resolve %s for condition %s: %w", indent, rightSourceDesc, leafDescription, err)
		}
	} else if leaf.Value != nil {
		if err := json.Unmarshal(leaf.Value, &right); err != nil {
			return false, "", fmt.Errorf("%sInvalid static JSON value for condition %s: %w", indent, leafDescription, err)
		}
		rightSourceDesc = fmt.Sprintf("static value '%v'", right) // TODO: Mask sensitive values
	} else if leaf.Operator != "IsNull" && leaf.Operator != "IsNotNull" {
		return false, "", fmt.Errorf("%sCondition %srequires 'value' or 'resource_attribute' for operator '%s'", indent, leafDescription, leaf.Operator)
	} else {
		rightSourceDesc = "(not applicable)" // For IsNull/IsNotNull
	}

	// Format left value for logging (handle nil)
	leftStr := fmt.Sprintf("%v", left)
	if left == nil {
		leftStr = "<nil>"
	}
    // ***** REMOVED unused rightStr variable *****
	// var rightStr string // Removed
	// if leaf.ResourceAttribute != "" && right == nil {
	// 	 rightStr = "<nil>" // Not used anymore
	// } else if right != nil {
    //     rightStr = fmt.Sprintf("%v", right) // Not used anymore
    // }

	// --- Helper for creating indented failure reason ---
	failReason := func(opResult bool, comparisonDesc string) (bool, string) {
		if opResult {
			return true, "" // Condition passed
		}
		// Condition failed: Create indented reason string
		reason := fmt.Sprintf("%s- Condition %s(%s: '%s' %s)", indent, leafDescription, leaf.Attribute, leftStr, comparisonDesc)
		return false, reason
	}

	// --- Perform Comparison ---
	var res bool // Comparison result
	var compareErr error // Error during comparison itself (e.g., type mismatch)
	var comparisonDescription string // Describes the failed comparison

	switch leaf.Operator {
	case "Equals":
		res = reflect.DeepEqual(left, right)
		comparisonDescription = fmt.Sprintf("!= %s)", rightSourceDesc) // Note: uses rightSourceDesc which already has value info
	case "IsNull":
		res = left == nil
		comparisonDescription = "is not null)"
	case "IsNotNull":
		res = left != nil
		comparisonDescription = "is null)"
	case "Contains":
		val := reflect.ValueOf(left)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			compareErr = fmt.Errorf("attribute '%s' is %T, not a slice/array", leaf.Attribute, left)
		} else {
			found := false
			for i := 0; i < val.Len(); i++ {
				if reflect.DeepEqual(val.Index(i).Interface(), right) {
					found = true
					break
				}
			}
			res = found
			comparisonDescription = fmt.Sprintf("does not contain %s)", rightSourceDesc)
		}

	case "GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual":
		res, compareErr = compareNumbers(left, right, leaf.Operator)
		opSymbols := map[string]string{
			"GreaterThan": ">", "GreaterThanOrEqual": ">=", "LessThan": "<", "LessThanOrEqual": "<=",
		}
		comparisonDescription = fmt.Sprintf("is not %s %s)", opSymbols[leaf.Operator], rightSourceDesc)

	case "StartsWith", "EndsWith", "Includes":
		res, compareErr = compareStrings(left, right, leaf.Operator)
		opDesc := map[string]string{
			"StartsWith": "starts with", "EndsWith": "ends with", "Includes": "includes",
		}
		comparisonDescription = fmt.Sprintf("does not %s %s)", opDesc[leaf.Operator], rightSourceDesc)

	case "Before", "After":
		res, compareErr = compareTimes(left, right, leaf.Operator)
		opDesc := map[string]string{"Before": "before", "After": "after"}[leaf.Operator]
		// Format times only if needed and conversion works
		lTime, lErr := toTime(left)
		rTime, rErr := toTime(right)
		rValDesc := rightSourceDesc // Use base description
		if rErr == nil && right != nil { // Only add formatted time if conversion worked and it wasn't IsNull check
			rValDesc = fmt.Sprintf("%s (%s)", rightSourceDesc, rTime.Format(time.RFC3339)) // Add formatted time
		}
		// Describe failed check
		if lErr == nil {
			comparisonDescription = fmt.Sprintf("(%s) is not %s %s)", lTime.Format(time.RFC3339), opDesc, rValDesc)
		} else {
			comparisonDescription = fmt.Sprintf("is not %s %s)", opDesc, rValDesc)
		}

	default:
		compareErr = fmt.Errorf("unsupported operator '%s'", leaf.Operator)
	}

	// Check for comparison errors first
	if compareErr != nil {
		return false, "", fmt.Errorf("%sComparison Error for condition %s: %w", indent, leafDescription, compareErr)
	}

	// Generate reason string if comparison failed
	finalRes, finalReason := failReason(res, comparisonDescription)
	return finalRes, finalReason, nil
}


// --- Helper functions (resolveAttr, compare*, to*) ---
// Minor improvements & error messages

func resolveAttr(attr string, subject, resource map[string]any) (any, error) {
	if attr == "" {
		return nil, errors.New("attribute name cannot be empty")
	}
	var sourceMap map[string]any
	var key string
    var contextName string

	switch {
	case strings.HasPrefix(attr, "subject."):
		sourceMap = subject
		key = strings.TrimPrefix(attr, "subject.")
        contextName = "subject"
		if key == "" { return nil, errors.New("invalid subject attribute key (e.g., 'subject.id')")}
	case strings.HasPrefix(attr, "resource."):
		sourceMap = resource
		key = strings.TrimPrefix(attr, "resource.")
        contextName = "resource"
        if key == "" { return nil, errors.New("invalid resource attribute key (e.g., 'resource.owner_id')")}
	default:
		return nil, fmt.Errorf("invalid attribute format: '%s' (must start with 'subject.' or 'resource.')", attr)
	}

    if sourceMap == nil {
         // Handle case where subject or resource map itself is nil
         return nil, fmt.Errorf("context map '%s' is nil, cannot resolve attribute '%s'", contextName, attr)
    }

    // Handle potential nested keys (split by '.' - basic example)
    // For "subject.profile.email": key = "profile.email", sourceMap = subject
    parts := strings.Split(key, ".")
    currentVal, ok := sourceMap[parts[0]]
    if !ok {
        return nil, nil // Treat top-level key not found as nil value, not error
    }

    // Traverse nested parts if any
    for i := 1; i < len(parts); i++ {
        // Check if currentVal is a map[string]any
        if nestedMap, isMap := currentVal.(map[string]any); isMap {
            currentVal, ok = nestedMap[parts[i]]
             if !ok {
                 return nil, nil // Nested key not found, treat as nil
             }
        } else {
             // Trying to access nested key on non-map value
            return nil, fmt.Errorf("attribute '%s' found non-map value at '%s' while resolving '%s'", attr, strings.Join(parts[:i],"."), parts[i])
        }
    }

    return currentVal, nil // Return the final value found
}


func compareNumbers(left, right any, op string) (bool, error) {
	leftFloat, ok1 := toFloat(left)
	rightFloat, ok2 := toFloat(right)

    // Be more explicit about *which* value failed conversion
	if !ok1 {
        return false, fmt.Errorf("cannot compare: left value %v (%T) is not a recognized number", left, left)
    }
    if !ok2 {
        return false, fmt.Errorf("cannot compare: right value %v (%T) is not a recognized number", right, right)
    }

	// ... rest of switch statement is fine ...
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
		// This should be caught by validation earlier, but defensively:
		return false, fmt.Errorf("internal error: unknown numeric comparison operator '%s'", op)
	}
}

func toFloat(val any) (float64, bool) {
     if val == nil { return 0, false } // Explicitly handle nil
	switch v := val.(type) {
	case int: return float64(v), true
	case int8: return float64(v), true
	case int16: return float64(v), true
	case int32: return float64(v), true
	case int64: return float64(v), true
	case uint: return float64(v), true
	case uint8: return float64(v), true
	case uint16: return float64(v), true
	case uint32: return float64(v), true
	case uint64: return float64(v), true // Potential precision loss but generally okay
	case float32: return float64(v), true
	case float64: return v, true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	default:
		// Maybe add string conversion carefully? e.g., using strconv.ParseFloat
		return 0, false
	}
}


func compareStrings(left, right any, op string) (bool, error) {
	lStr, ok1 := left.(string)
	if !ok1 { return false, fmt.Errorf("cannot compare: left value %v (%T) is not a string", left, left)}
    rStr, ok2 := right.(string)
    if !ok2 { return false, fmt.Errorf("cannot compare: right value %v (%T) is not a string", right, right)}

	// ... rest of switch statement is fine ...
	switch op {
	case "StartsWith":
		return strings.HasPrefix(lStr, rStr), nil
	case "EndsWith":
		return strings.HasSuffix(lStr, rStr), nil
	case "Includes":
		return strings.Contains(lStr, rStr), nil
	default:
		return false, fmt.Errorf("internal error: unsupported string comparison operator '%s'", op)
	}
}


func compareTimes(left, right any, op string) (bool, error) {
	lTime, err := toTime(left)
	if err != nil {
		return false, fmt.Errorf("cannot compare times: left value %v (%T) failed conversion: %w", left, left, err)
	}
	rTime, err := toTime(right)
	if err != nil {
		return false, fmt.Errorf("cannot compare times: right value %v (%T) failed conversion: %w", right, right, err)
	}

	// ... rest of switch statement is fine ...
	switch op {
	case "Before":
		return lTime.Before(rTime), nil
	case "After":
		return lTime.After(rTime), nil
	default:
		return false, fmt.Errorf("internal error: unsupported time comparison operator '%s'", op)
	}
}

// Enhanced toTime with more formats and robustness
func toTime(val any) (time.Time, error) {
    if val == nil { return time.Time{}, errors.New("cannot convert nil to time.Time")}

	switch v := val.(type) {
	case time.Time:
		// Ensure time has location info if possible, default to UTC if ambiguous
        if v.Location() == nil {
            return v.UTC(), nil // Or maybe keep local? Depends on source. UTC is safer standard.
        }
		return v, nil
	case string:
        // Try common formats - add more as needed by your data sources
        formats := []string{
            time.RFC3339,
            time.RFC3339Nano,
            "2006-01-02T15:04:05Z07:00", // RFC3339 slightly simplified
            "2006-01-02 15:04:05 Z07:00", // Space separation
            "2006-01-02 15:04:05", // Common DB format (assumes UTC or server local) - Use UTC
            "2006-01-02",           // Date only (time defaults to 00:00:00 UTC)
        }
        for _, format := range formats {
             t, err := time.Parse(format, v)
             if err == nil {
                // If parsed without timezone, assume UTC
                if t.Location() == time.UTC || t.Location() == nil {
                     return t.UTC(), nil
                }
                return t, nil // Return with original location if parsed
             }
        }
		return time.Time{}, fmt.Errorf("string '%s' is not in any supported time format", v)
	case json.Number: // From DB timestamps stored as numbers
		// Try integer seconds first
		i, err := v.Int64()
		if err == nil {
			return time.Unix(i, 0).UTC(), nil // Assume seconds since epoch, UTC
		}
        // Try float seconds (potentially with fractional part)
        f, err := v.Float64()
        if err == nil {
            sec, dec := math.Modf(f)
            // Handle potential millisecond or nanosecond precision from float
            nsec := int64(dec * 1e9)
            return time.Unix(int64(sec), nsec).UTC(), nil // UTC
        }
		return time.Time{}, fmt.Errorf("json.Number '%s' could not be converted to int64/float64 timestamp", v.String())
    case int:
        return time.Unix(int64(v), 0).UTC(), nil // Assume seconds, UTC
    case int64:
         // Could be seconds or milliseconds? Assume seconds for standard Unix time.
        return time.Unix(v, 0).UTC(), nil
    case float64:
        sec, dec := math.Modf(v)
        nsec := int64(dec * 1e9)
        return time.Unix(int64(sec), nsec).UTC(), nil // UTC
	default:
		return time.Time{}, fmt.Errorf("unsupported type %T cannot be converted to time.Time", val)
	}
}