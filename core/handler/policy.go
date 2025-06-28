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

	"github.com/google/uuid"
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

// CanAccess evaluates if the subject can perform an action on a resource,
// considering context from the request path, body, and query parameters.
// subject: Information about the entity performing the action.
// resource: Information about the entity being acted upon (e.g., fetched DB record). Can be nil.
// path: Data extracted from URL path parameters. /employee_id/{id}
// body: Data parsed from the request body.
// query: Data extracted from URL query string parameters.
// headers: Data extracted from request headers.
// policy: The specific PolicyRule to evaluate.
func (p *Policy) CanAccess(subject, resource, path, body, query, headers map[string]any, policy *model.PolicyRule) AccessDecision {
	decision := AccessDecision{}
	if policy == nil {
		decision.Error = errors.New("policy rule cannot be nil")
		return decision
	}

	policyIdentifier := fmt.Sprintf("ID %s", policy.ID.String())
	if policy.Name != "" {
		policyIdentifier = fmt.Sprintf("'%s' (ID %s)", policy.Name, policy.ID.String())
	}

	// Validation should now allow 'header.' prefix
	root, err := policy.GetConditionsNode()
	if err != nil {
		decision.Error = fmt.Errorf("failed to get/validate conditions node for policy %s: %w", policyIdentifier, err)
		return decision
	}

	// Start evaluation passing all distinct context maps, including headers
	result, reason, err := evalNode(root, subject, resource, path, body, query, headers, 0) // Pass headers
	if err != nil {
		decision.Error = fmt.Errorf("error evaluating conditions for policy %s: %w", policyIdentifier, err)
		return decision
	}

	effectDesc := fmt.Sprintf("[%s]", policy.Effect)

	if policy.Effect == "Allow" {
		decision.Allowed = result
		if !decision.Allowed {
			decision.Reason = fmt.Sprintf("Policy %s %s denied", effectDesc, policyIdentifier)
			if reason != "" {
				decision.Reason += ":\n" + reason
			} else {
				decision.Reason += " (conditions not met)."
			}
		}
	} else if policy.Effect == "Deny" {
		decision.Allowed = !result
		if !decision.Allowed { // Denied because Deny rule evaluated to TRUE
			decision.Reason = fmt.Sprintf("Policy %s %s enforced", effectDesc, policyIdentifier)
			if reason != "" {
				decision.Reason += ":\n" + reason
			} else {
				decision.Reason += " (conditions met)."
			}
		}
	} else {
		decision.Error = fmt.Errorf("unknown policy effect '%s' in policy %s", policy.Effect, policyIdentifier)
	}

	return decision
}

// evalNode evaluates a condition node recursively, passing all context maps.
func evalNode(node model.ConditionNode, subject, resource, path, body, query, headers map[string]any, depth int) (bool, string, error) {
	indent := strings.Repeat(indentStep, depth)
	nodeDescription := ""
	if node.Description != "" {
		nodeDescription = fmt.Sprintf("'%s'", node.Description)
	} else {
		nodeDescription = "(unnamed node)"
	}

	if node.Leaf != nil {
		// Evaluate leaf node, passing all maps
		res, reason, err := evalLeaf(*node.Leaf, subject, resource, path, body, query, headers, depth) // Pass headers
		if err != nil {
			return false, "", fmt.Errorf("%s Leaf Node %s evaluation failed: %w", indent, nodeDescription, err)
		}
		return res, reason, nil
	}

	if len(node.Children) == 0 {
		return false, "", fmt.Errorf("%s Node %s has logic type '%s' but no children", indent, nodeDescription, node.LogicType)
	}
	if node.LogicType != "AND" && node.LogicType != "OR" {
		return false, "", fmt.Errorf("%s Node %s has unknown logic type: %s", indent, nodeDescription, node.LogicType)
	}

	childResults := make([]bool, len(node.Children))
	childReasons := make([]string, len(node.Children))

	for i, child := range node.Children {
		// Evaluate child with increased depth, passing all maps
		res, reason, err := evalNode(child, subject, resource, path, body, query, headers, depth+1) // Pass headers
		if err != nil {
			return false, "", fmt.Errorf("%s Child evaluation failed within %s node: %w", indent, nodeDescription, err)
		}
		childResults[i] = res
		if !res {
			childReasons[i] = reason
		} else {
			childReasons[i] = ""
		}

		switch node.LogicType {
		case "AND":
			if !res {
				finalReason := fmt.Sprintf("%s - [%s Node: %s] failed:\n%s", indent, node.LogicType, nodeDescription, reason)
				return false, finalReason, nil
			}
		case "OR":
			if res {
				return true, "", nil
			}
		}
	}

	switch node.LogicType {
	case "AND":
		return true, "", nil
	case "OR":
		var failingReasons []string
		for _, r := range childReasons {
			if r != "" {
				failingReasons = append(failingReasons, r)
			}
		}
		combinedReason := strings.Join(failingReasons, "\n")
		finalReason := fmt.Sprintf("%s - [%s Node: %s] failed (all children false):\n%s", indent, node.LogicType, nodeDescription, combinedReason)
		return false, finalReason, nil
	default:
		return false, "", fmt.Errorf("%s Internal Error: Unhandled logic type %s", indent, node.LogicType)
	}
}

// evalLeaf evaluates a leaf condition, passing all context maps.
func evalLeaf(leaf model.ConditionLeaf, subject, resource, path, body, query, headers map[string]any, depth int) (bool, string, error) {
	indent := strings.Repeat(indentStep, depth)
	leafDescription := ""
	if leaf.Description != "" {
		leafDescription = fmt.Sprintf("'%s' ", leaf.Description)
	}

	// --- Resolve Right-Hand Value ---
	var right any
	var rightSourceDesc string
	var formattedRightVal string
	var err error
	if leaf.ResourceAttribute != "" {
		rightSourceDesc = fmt.Sprintf("attribute '%s'", leaf.ResourceAttribute)
		right, err = resolveAttr(leaf.ResourceAttribute, subject, resource, path, body, query, headers)
		if err != nil {
			reason := fmt.Sprintf("%s - Condition %s failed: resolve %s: %v", indent, leafDescription, rightSourceDesc, err)
			return false, reason, fmt.Errorf("%s Resolve attribute '%s' failed: %w", indent, leaf.ResourceAttribute, err) // Use full attribute name in error
		}
		formattedRightVal = formatValueForLog(right)
	} else if leaf.Value != nil {
		if err = json.Unmarshal(leaf.Value, &right); err != nil {
			reason := fmt.Sprintf("%s - Condition %s failed: invalid static JSON '%s': %v", indent, leafDescription, string(leaf.Value), err)
			return false, reason, fmt.Errorf("%s Invalid static JSON '%s': %w", indent, string(leaf.Value), err)
		}
		// Format the static value description better for logging
		formattedRightVal = formatValueForLog(right)
		rightSourceDesc = fmt.Sprintf("static value %s", formattedRightVal)
	} else if leaf.Operator != "IsNull" && leaf.Operator != "IsNotNull" {
		reason := fmt.Sprintf("%s - Condition %sfailed: operator '%s' requires value/resource_attribute", indent, leafDescription, leaf.Operator)
		return false, reason, fmt.Errorf("%s Op '%s' needs value/attribute", indent, leaf.Operator)
	} else {
		rightSourceDesc = "(comparison value not applicable)" // Clearer description
		formattedRightVal = "<not applicable>"
	}

	// --- Resolve Left-Hand Value ---
	var left any
	var leftStr string // Formatted string representation of the left value for logs/reasons
	left, err = resolveAttr(leaf.Attribute, subject, resource, path, body, query, headers)
	if err != nil {
		// Provide more context in the reason message
		reason := fmt.Sprintf("%s - Condition %s(attribute '%s') failed: could not resolve: %v", indent, leafDescription, leaf.Attribute, err)
		return false, reason, fmt.Errorf("%s Resolve attribute '%s' failed: %w", indent, leaf.Attribute, err)
	}
	leftStr = formatValueForLog(left)

	var comparisonTargetDesc string
	if leaf.ResourceAttribute != "" {
		// For attributes, show both the attribute name and its resolved value
		comparisonTargetDesc = fmt.Sprintf("%s: %s", rightSourceDesc, formattedRightVal)
	} else if leaf.Value != nil {
		// For static values, rightSourceDesc already includes the formatted value
		comparisonTargetDesc = rightSourceDesc
	} else {
		// For operators like IsNull/IsNotNull
		comparisonTargetDesc = rightSourceDesc
	}

	// --- Helper for creating indented failure reason ---
	// Uses comparisonDescription set by the specific operator logic below
	failReason := func(opResult bool, lValFmt, compDesc string) (bool, string) {
		if opResult {
			return true, ""
		}
		var reasonCore string
		// comparisonDesc should contain the explanation of failure (e.g., "is not null", "does not equal X")
		if compDesc != "" {
			reasonCore = fmt.Sprintf("%s: %s %s", leaf.Attribute, lValFmt, compDesc)
		} else {
			// Fallback if comparisonDescription wasn't set (should be avoided)
			reasonCore = fmt.Sprintf("%s: %s comparison (%s) failed", leaf.Attribute, lValFmt, leaf.Operator)
		}
		reason := fmt.Sprintf("%s - Condition %s (%s)", indent, leafDescription, reasonCore) // Removed extra closing parenthesis
		return false, reason
	}

	// --- Perform Comparison ---
	var res bool                     // The boolean result of the comparison
	var compareErr error             // Any error occurring during comparison (e.g., type mismatch)
	var comparisonDescription string // Human-readable explanation of *why* a comparison failed (if res == false)

	switch leaf.Operator {
	case "Equals":
		res = robustCompareEquals(left, right)
		if !res {
			comparisonDescription = fmt.Sprintf("does not equal %s", comparisonTargetDesc)
		}
	case "NotEquals":
		res = !robustCompareEquals(left, right)
		if !res {
			comparisonDescription = fmt.Sprintf("equals %s", comparisonTargetDesc)
		}
	case "IsNull":
		// Consider pointer nils as well as interface nil
		leftValCheck := reflect.ValueOf(left)
		isConsideredNil := left == nil || (leftValCheck.IsValid() && leftValCheck.Kind() == reflect.Ptr && leftValCheck.IsNil())
		res = isConsideredNil
		if !res {
			comparisonDescription = "is not null"
		}
	case "IsNotNull":
		leftValCheck := reflect.ValueOf(left)
		isConsideredNil := left == nil || (leftValCheck.IsValid() && leftValCheck.Kind() == reflect.Ptr && leftValCheck.IsNil())
		res = !isConsideredNil
		if !res {
			comparisonDescription = "is null"
		}
	case "GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual":
		res, compareErr = compareNumbers(left, right, leaf.Operator)
		if compareErr != nil {
			// Error during comparison
			comparisonDescription = fmt.Sprintf("could not be compared numerically with %s: %v", comparisonTargetDesc, compareErr) // USE comparisonTargetDesc
		} else if !res {
			// Comparison successful but returned false
			opSymbols := map[string]string{">": ">", ">=": ">=", "<": "<", "<=": "<="}
			comparisonDescription = fmt.Sprintf("is not %s %s", opSymbols[leaf.Operator], comparisonTargetDesc) // USE comparisonTargetDesc
		}
	case "StartsWith", "EndsWith", "Includes":
		res, compareErr = compareStrings(left, right, leaf.Operator)
		if compareErr != nil {
			comparisonDescription = fmt.Sprintf("could not be compared as strings with %s: %v", comparisonTargetDesc, compareErr) // USE comparisonTargetDesc
		} else if !res {
			opDesc := map[string]string{"StartsWith": "start with", "EndsWith": "end with", "Includes": "include"}
			comparisonDescription = fmt.Sprintf("does not %s %s", opDesc[leaf.Operator], comparisonTargetDesc) // USE comparisonTargetDesc
		}
	case "Before", "After":
		res, compareErr = compareTimes(left, right, leaf.Operator)
		if compareErr != nil {
			comparisonDescription = fmt.Sprintf("could not be compared as times with %s: %v", comparisonTargetDesc, compareErr) // USE comparisonTargetDesc
		} else if !res {
			opDesc := map[string]string{"Before": "before", "After": "after"}[leaf.Operator]
			lValDesc := leftStr // Already formatted left value
			// Time formatting is handled within formatValueForLog, so comparisonTargetDesc is sufficient
			comparisonDescription = fmt.Sprintf("(value: %s) is not %s %s", lValDesc, opDesc, comparisonTargetDesc) // USE comparisonTargetDesc
		}
	case "Contains":
		targetValue := right // Use resolved 'right' value directly as the target
		var collectionPath string
		var fieldToExtract string
		isObjectSyntax := false

		if strings.Contains(leaf.Attribute, "[*]") {
			parts := strings.SplitN(leaf.Attribute, "[*]", 2)
			collectionPath = parts[0]
			if len(parts) > 1 && strings.HasPrefix(parts[1], ".") {
				isObjectSyntax = true
				fieldToExtract = strings.TrimPrefix(parts[1], ".")
				if fieldToExtract == "" {
					compareErr = fmt.Errorf("invalid 'Contains' object syntax: missing field name after [*]. in '%s'", leaf.Attribute)
					comparisonDescription = fmt.Sprintf("invalid syntax: missing field name after [*]. in '%s'", leaf.Attribute)
				}
			} else if len(parts) > 1 && parts[1] != "" {
				compareErr = fmt.Errorf("invalid 'Contains' syntax: unexpected characters after [*] in '%s'", leaf.Attribute)
				comparisonDescription = fmt.Sprintf("invalid syntax: unexpected characters after [*] in '%s'", leaf.Attribute)
			}
		} else {
			collectionPath = leaf.Attribute
			isObjectSyntax = false
		}

		var leftCollection any
		if compareErr == nil {
			leftCollection, err = resolveAttr(collectionPath, subject, resource, path, body, query, headers)
			if err != nil {
				compareErr = fmt.Errorf("resolve collection path '%s' failed: %w", collectionPath, err)
				comparisonDescription = fmt.Sprintf("could not resolve collection path '%s': %v", collectionPath, err)
				leftStr = "<unresolved collection>"
			}
		}

		if compareErr == nil {
			leftStr = formatValueForLog(leftCollection) // Format collection *after* potential resolution error handled
			collValue := reflect.ValueOf(leftCollection)

			if !collValue.IsValid() || (collValue.Kind() != reflect.Slice && collValue.Kind() != reflect.Array) {
				res = false
				comparisonDescription = fmt.Sprintf("resolved path '%s' is type %T (value: %s), not a slice/array, cannot perform 'Contains'", collectionPath, leftCollection, leftStr)
			} else {
				found := false
				if isObjectSyntax {
					for i := range collValue.Len() {
						element := collValue.Index(i)
						extractedValue, fieldFound := extractFieldByName(element, fieldToExtract)
						if fieldFound && robustCompareEquals(extractedValue, targetValue) {
							found = true
							break
						}
					}
					if !found {
						// USE comparisonTargetDesc for the value being looked for
						comparisonDescription = fmt.Sprintf("in path '%s', no object found where field '%s' equals %s", collectionPath, fieldToExtract, comparisonTargetDesc)
					}
				} else {
					for i := range collValue.Len() {
						element := collValue.Index(i)
						if !element.IsValid() || !element.CanInterface() {
							continue
						}
						elementInterface := element.Interface()
						if robustCompareEquals(elementInterface, targetValue) {
							found = true
							break
						}
					}
					if !found {
						// USE comparisonTargetDesc for the value being looked for
						comparisonDescription = fmt.Sprintf("collection at path '%s' does not contain an element equal to %s", collectionPath, comparisonTargetDesc)
					}
				}
				res = found
				if res {
					comparisonDescription = "" // Clear on success
				}
			}
		}

	default:
		compareErr = fmt.Errorf("unsupported operator '%s'", leaf.Operator)
	}

	// --- Final Result and Reason ---
	if compareErr != nil {
		// An error occurred *during* the comparison itself (e.g., type mismatch, invalid syntax)
		// Generate reason based on the comparison error. comparisonDescription might also contain info.
		reasonMsg := fmt.Sprintf("Comparison error: %v", compareErr)
		if comparisonDescription != "" {
			reasonMsg = comparisonDescription // Use more specific description if available
		}
		// Use failReason, passing false and the detailed error description
		_, finalReason := failReason(false, leftStr, reasonMsg)
		// Return the error as well for detailed logging
		return false, finalReason, fmt.Errorf("%s Comparison Error on attribute '%s': %w", indent, leaf.Attribute, compareErr)
	}

	// If no comparison error occurred, use the standard failReason helper
	// comparisonDescription will be populated if 'res' is false
	finalRes, finalReason := failReason(res, leftStr, comparisonDescription)
	return finalRes, finalReason, nil // No error occurred during comparison itself
}

// resolveAttr resolves attribute path against ALL provided context maps, including headers.
func resolveAttr(attr string, subject, resource, path, body, query, headers map[string]any) (any, error) {
	if attr == "" {
		return nil, errors.New("attribute name cannot be empty")
	}

	var sourceMap map[string]any
	var key string
	var contextName string

	switch {
	case strings.HasPrefix(attr, "subject."):
		sourceMap, key, contextName = subject, strings.TrimPrefix(attr, "subject."), "subject"
	case strings.HasPrefix(attr, "resource."):
		sourceMap, key, contextName = resource, strings.TrimPrefix(attr, "resource."), "resource"
	case strings.HasPrefix(attr, "path."):
		sourceMap, key, contextName = path, strings.TrimPrefix(attr, "path."), "path"
	case strings.HasPrefix(attr, "body."):
		sourceMap, key, contextName = body, strings.TrimPrefix(attr, "body."), "body"
	case strings.HasPrefix(attr, "query."):
		sourceMap, key, contextName = query, strings.TrimPrefix(attr, "query."), "query"
	case strings.HasPrefix(attr, "header."): // Added Header case
		sourceMap, key, contextName = headers, strings.TrimPrefix(attr, "header."), "header"
		// Note: Header keys in sourceMap should ideally be canonicalized (e.g., "Content-Type").
	default:
		return nil, fmt.Errorf("invalid attribute format: '%s' (must start with 'subject.', 'resource.', 'path.', 'body.', 'query.', or 'header.')", attr)
	}

	if key == "" {
		return nil, fmt.Errorf("invalid %s attribute key (e.g., '%s.id')", contextName, contextName)
	}

	// Allow context map to be nil (e.g., no headers sent, no body, etc.)
	if sourceMap == nil {
		return nil, nil
	}

	// Handle nested keys (less common for headers/path/query, but possible for subject/resource/body)
	parts := strings.Split(key, ".")
	currentVal, ok := sourceMap[parts[0]]
	if !ok {
		return nil, nil /* Top-level key not found */
	}

	for i := 1; i < len(parts); i++ {
		currentMap, isMap := currentVal.(map[string]any)
		if !isMap {
			nestedPath := strings.Join(parts[:i], ".")
			return nil, fmt.Errorf("attr '%s': value at '%s.%s' is %T, not map for nested access", attr, contextName, nestedPath, currentVal)
		}
		currentVal, ok = currentMap[parts[i]]
		if !ok {
			return nil, nil /* Nested key not found */
		}
	}

	return currentVal, nil
}

// --- Helper functions (resolveAttr, compare*, to*) ---
// Minor improvements & error messages

// Helper to format values for log messages (prevents overly verbose maps/structs)
func formatValueForLog(val any) string {
	if val == nil {
		return "<nil>"
	}

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "<nil>"
		}
		// Dereference pointer for further checks
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() > 5 {
			return fmt.Sprintf("%T (len:%d) [preview omitted]", val, v.Len())
		}
		marshaled, err := json.Marshal(val)
		if err == nil && len(marshaled) < 100 {
			return string(marshaled)
		}
		return fmt.Sprintf("%T (len:%d, %v)", val, v.Len(), val)
	case reflect.Map, reflect.Struct:
		if v.Kind() == reflect.Map {
			marshaled, err := json.Marshal(val)
			if err == nil && len(marshaled) < 150 {
				return string(marshaled)
			}
		}
		return fmt.Sprintf("%T (details omitted)", val)
	default:
		if str, ok := val.(string); ok {
			return fmt.Sprintf(`"%s"`, str)
		}
		if t, ok := val.(time.Time); ok {
			return t.Format(time.RFC3339)
		}
		return fmt.Sprintf("%v", val)
	}
}

func robustCompareEquals(left, right any) bool {
	if left == nil || right == nil {
		return left == right
	}
	leftFloat, leftIsNum := toFloat(left)
	rightFloat, rightIsNum := toFloat(right)
	if leftIsNum && rightIsNum {
		return leftFloat == rightFloat
	}
	leftStr, leftIsStr := toString(left)
	rightStr, rightIsStr := toString(right)
	leftUUID, leftIsUUID := toUUID(left)
	rightUUID, rightIsUUID := toUUID(right)
	if leftUUID != uuid.Nil && rightUUID != uuid.Nil {
		return leftUUID == rightUUID
	}
	if leftUUID != uuid.Nil && rightIsStr {
		parsedRight, _ := uuid.Parse(rightStr)
		return leftUUID == parsedRight
	}
	if rightUUID != uuid.Nil && leftIsStr {
		parsedLeft, _ := uuid.Parse(leftStr)
		return rightUUID == parsedLeft
	}
	if (leftIsUUID && !rightIsUUID && rightUUID == uuid.Nil) || (rightIsUUID && !leftIsUUID && leftUUID == uuid.Nil) {
		return false
	}
	if leftIsBool, ok := left.(bool); ok {
		if rightIsBool, okR := right.(bool); okR {
			return leftIsBool == rightIsBool
		}
	}
	if leftIsStr && rightIsStr {
		return leftStr == rightStr
	} // Direct string compare after others
	if reflect.TypeOf(left) == reflect.TypeOf(right) {
		return reflect.DeepEqual(left, right)
	}
	return false
}

func extractFieldByName(element reflect.Value, fieldName string) (any, bool) {
	for element.Kind() == reflect.Ptr || element.Kind() == reflect.Interface {
		if element.IsNil() {
			return nil, false
		}
		element = element.Elem()
		if !element.IsValid() {
			return nil, false
		}
	}

	// 1. Initial Validity Check
	if !element.IsValid() {
		return nil, false
	}

	// 2. Handle Interface Wrapping
	if element.Kind() == reflect.Interface {
		if element.IsNil() {
			return nil, false // Interface is nil, contains no value
		}
		// Replace 'element' with the value inside the interface
		element = element.Elem()
		// Need to re-check validity in case the elem inside interface was invalid (though less common)
		if !element.IsValid() {
			return nil, false
		}
	}

	// 3. Handle Pointers
	if element.Kind() == reflect.Ptr {
		if element.IsNil() {
			return nil, false // Pointer is nil
		}
		// Replace 'element' with the value the pointer points to
		element = element.Elem()
		// Re-check validity after dereferencing
		if !element.IsValid() {
			return nil, false
		}
	}

	// 4. Try Map Access (Most common for map[string]any)
	// At this point, 'element' should represent the actual map or struct (not an interface or pointer).
	if element.Kind() == reflect.Map {
		// Ensure map keys are strings for direct lookup
		if element.Type().Key().Kind() == reflect.String {
			// Use reflect.ValueOf to create a key of the correct type for lookup
			key := reflect.ValueOf(fieldName)
			val := element.MapIndex(key)
			// Check if the key existed and the value is accessible
			if val.IsValid() && val.CanInterface() {
				return val.Interface(), true
			}
		}
	}

	// 5. Try Struct Field Access
	if element.Kind() == reflect.Struct {
		// FieldByName requires the EXACT Go struct field name (usually CamelCase and exported)
		field := element.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() {
			// Found an exported field with the exact name
			return field.Interface(), true
		}
	}

	// 6. Not Found
	// If it's not a map or struct, or the field/key wasn't found/accessible
	return nil, false
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

func toString(val any) (string, bool) { /* ... Copy from previous correct version ... */
	if val == nil {
		return "", false
	}
	if s, ok := val.(string); ok {
		return s, true
	}
	if stringer, ok := val.(fmt.Stringer); ok {
		return stringer.String(), true
	}
	return "", false
}

func toUUID(val any) (uuid.UUID, bool) { /* ... Copy from previous correct version ... */
	if val == nil {
		return uuid.Nil, false
	}
	if u, ok := val.(uuid.UUID); ok {
		return u, true
	}
	if s, ok := val.(string); ok {
		u, err := uuid.Parse(s)
		return u, err == nil
	}
	if b, ok := val.([]byte); ok {
		u, err := uuid.FromBytes(b)
		return u, err == nil
	}
	return uuid.Nil, false
}

func toFloat(val any) (float64, bool) {
	if val == nil {
		return 0, false
	} // Explicitly handle nil
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true // Potential precision loss but generally okay
	case float32:
		return float64(v), true
	case float64:
		return v, true
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
	if !ok1 {
		return false, fmt.Errorf("cannot compare: left value %v (%T) is not a string", left, left)
	}
	rStr, ok2 := right.(string)
	if !ok2 {
		return false, fmt.Errorf("cannot compare: right value %v (%T) is not a string", right, right)
	}

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
	if val == nil {
		return time.Time{}, errors.New("cannot convert nil to time.Time")
	}

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
			"2006-01-02T15:04:05Z07:00",  // RFC3339 slightly simplified
			"2006-01-02 15:04:05 Z07:00", // Space separation
			"2006-01-02 15:04:05",        // Common DB format (assumes UTC or server local) - Use UTC
			"2006-01-02",                 // Date only (time defaults to 00:00:00 UTC)
		}
		for _, format := range formats {
			t, err := time.Parse(format, v)
			if err == nil {
				// If parsed without time_zone, assume UTC
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
