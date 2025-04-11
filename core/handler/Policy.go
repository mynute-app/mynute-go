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
		if !decision.Allowed {     // Access denied if Deny conditions are TRUE
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

// This function evaluates a condition node recursively, passing the depth for indentation
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
		leafDescription = fmt.Sprintf("'%s' ", leaf.Description)
	}

	// --- Resolve Right-Hand Value FIRST (needed for comparison later) ---
	var right any
	var rightSourceDesc string
	var err error
	if leaf.ResourceAttribute != "" {
		rightSourceDesc = fmt.Sprintf("resource attribute '%s'", leaf.ResourceAttribute)
		right, err = resolveAttr(leaf.ResourceAttribute, subject, resource)
		if err != nil {
			reason := fmt.Sprintf("%s- Condition %sfailed: could not resolve %s: %v)", indent, leafDescription, rightSourceDesc, err)
			return false, reason, fmt.Errorf("%sFailed to resolve %s for condition %s: %w", indent, rightSourceDesc, leafDescription, err)
		}
	} else if leaf.Value != nil {
		if err := json.Unmarshal(leaf.Value, &right); err != nil {
			reason := fmt.Sprintf("%s- Condition %sfailed: invalid static JSON value '%s': %v)", indent, leafDescription, string(leaf.Value), err)
			return false, reason, fmt.Errorf("%sInvalid static JSON value '%s' for condition %s: %w", indent, string(leaf.Value), leafDescription, err)
		}
		valueStr := strings.Trim(string(leaf.Value), `"`)
		rightSourceDesc = fmt.Sprintf("static value '%s'", valueStr)
	} else if leaf.Operator != "IsNull" && leaf.Operator != "IsNotNull" {
		reason := fmt.Sprintf("%s- Condition %sfailed: operator '%s' requires 'value' or 'resource_attribute'", indent, leafDescription, leaf.Operator)
		return false, reason, fmt.Errorf("%sCondition %srequires 'value' or 'resource_attribute' for operator '%s'", indent, leafDescription, leaf.Operator)
	} else {
		rightSourceDesc = "(not applicable)"
	}

	// We resolve the LEFT value LATER, only if needed and based on syntax for 'Contains'

	// --- Helper for creating indented failure reason ---
	// Note: Left value formatting is deferred until we know how 'Contains' uses it
	failReason := func(opResult bool, leftValFormatted string, comparisonDesc string, rightValFormatted string) (bool, string) { // Add rightValFormatted
		if opResult {
			return true, ""
		}
		// Construct the core comparison part of the message
		var reasonCore string
		// Add the right value details into the comparisonDesc where appropriate
		switch leaf.Operator {
		case "Equals":
			// Example: "does not equal resource attribute 'resource.company_id' (value: <nil>)"
			reasonCore = fmt.Sprintf("%s: %s does not equal %s (value: %s))", leaf.Attribute, leftValFormatted, rightSourceDesc, rightValFormatted)
		case "NotEquals":
			reasonCore = fmt.Sprintf("%s: %s equals %s (value: %s))", leaf.Attribute, leftValFormatted, rightSourceDesc, rightValFormatted)
		case "IsNull":
			reasonCore = fmt.Sprintf("%s: %s is not null)", leaf.Attribute, leftValFormatted) // No right value here
		case "IsNotNull":
			reasonCore = fmt.Sprintf("%s: %s is null)", leaf.Attribute, leftValFormatted) // No right value here
		case "GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual":
			opSymbols := map[string]string{">": ">", ">=": ">=", "<": "<", "<=": "<="}
			reasonCore = fmt.Sprintf("%s: %s is not %s %s (value: %s))", leaf.Attribute, leftValFormatted, opSymbols[leaf.Operator], rightSourceDesc, rightValFormatted)
		case "StartsWith", "EndsWith", "Includes":
			opDesc := map[string]string{"StartsWith": "start with", "EndsWith": "end with", "Includes": "include"}
			reasonCore = fmt.Sprintf("%s: %s does not %s %s (value: %s))", leaf.Attribute, leftValFormatted, opDesc[leaf.Operator], rightSourceDesc, rightValFormatted)
		case "Before", "After":
			opDesc := map[string]string{"Before": "before", "After": "after"}[leaf.Operator]
			// Time formatting is already complex, maybe just add the raw value if comparison failed
			reasonCore = fmt.Sprintf("%s: %s is not %s %s (value: %s))", leaf.Attribute, leftValFormatted, opDesc, rightSourceDesc, rightValFormatted)
		case "Contains": // Contains needs special handling based on syntax used within its logic
			// The 'comparisonDescription' calculated earlier for Contains already includes details
			reasonCore = fmt.Sprintf("%s: %s %s)", leaf.Attribute, leftValFormatted, comparisonDesc)
		default:
			reasonCore = fmt.Sprintf("%s: %s comparison (%s) failed against %s (value: %s))", leaf.Attribute, leftValFormatted, leaf.Operator, rightSourceDesc, rightValFormatted) // Fallback
		}

		// Construct the final reason string
		reason := fmt.Sprintf("%s- Condition %s(%s", indent, leafDescription, reasonCore)
		return false, reason
	}

	// --- Perform Comparison ---
	var res bool
	var compareErr error
	var comparisonDescription string
	var left any       // Declare left here, resolved within specific cases
	var leftStr string // Formatted left value

	left, err = resolveAttr(leaf.Attribute, subject, resource)
	if err != nil {
		reason := fmt.Sprintf("%s- Condition %sfailed: could not resolve attribute '%s': %v)", indent, leafDescription, leaf.Attribute, err)
		return false, reason, fmt.Errorf("%sFailed to resolve attribute '%s' for condition %s: %w", indent, leaf.Attribute, leafDescription, err)
	}
	leftStr = formatValueForLog(left) // Format it now

	switch leaf.Operator {
	case "Equals", "NotEquals", "IsNull", "IsNotNull",
		"GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual",
		"StartsWith", "EndsWith", "Includes", "Before", "After":
		// --- Resolve Left value for non-Contains operators ---

		// --- Handle the specific operator ---
		switch leaf.Operator {
		case "Equals":
			res = robustCompareEquals(left, right)
			comparisonDescription = fmt.Sprintf("does not equal %s)", rightSourceDesc)
		case "NotEquals":
			res = !robustCompareEquals(left, right)
			comparisonDescription = fmt.Sprintf("equals %s)", rightSourceDesc)
		case "IsNull":
			leftValCheck := reflect.ValueOf(left)
			res = left == nil || (leftValCheck.IsValid() && leftValCheck.Kind() == reflect.Ptr && leftValCheck.IsNil())
			comparisonDescription = "is not null)"
		case "IsNotNull":
			leftValCheck := reflect.ValueOf(left)
			res = left != nil && !(leftValCheck.IsValid() && leftValCheck.Kind() == reflect.Ptr && leftValCheck.IsNil())
			comparisonDescription = "is null)"
		// ... other simple operators ...
		case "GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual":
			res, compareErr = compareNumbers(left, right, leaf.Operator)
			if compareErr == nil {
				opSymbols := map[string]string{">": ">", ">=": ">=", "<": "<", "<=": "<="}
				comparisonDescription = fmt.Sprintf("is not %s %s)", opSymbols[leaf.Operator], rightSourceDesc)
			} else {
				comparisonDescription = fmt.Sprintf("could not be compared with %s: %v)", rightSourceDesc, compareErr)
			}
		case "StartsWith", "EndsWith", "Includes":
			res, compareErr = compareStrings(left, right, leaf.Operator)
			if compareErr == nil {
				opDesc := map[string]string{"StartsWith": "start with", "EndsWith": "end with", "Includes": "include"}
				comparisonDescription = fmt.Sprintf("does not %s %s)", opDesc[leaf.Operator], rightSourceDesc)
			} else {
				comparisonDescription = fmt.Sprintf("could not be compared with %s: %v)", rightSourceDesc, compareErr)
			}
		case "Before", "After":
			res, compareErr = compareTimes(left, right, leaf.Operator)
			if compareErr == nil {
				opDesc := map[string]string{"Before": "before", "After": "after"}[leaf.Operator]
				lTime, l_err := toTime(left)
				rTime, r_err := toTime(right)
				lValDesc := formatValueForLog(left)
				if l_err == nil {
					lValDesc = lTime.Format(time.RFC3339)
				}
				rValDesc := rightSourceDesc
				if r_err == nil {
					rValDesc = fmt.Sprintf("%s (%s)", rightSourceDesc, rTime.Format(time.RFC3339))
				}
				comparisonDescription = fmt.Sprintf("(%s) is not %s %s)", lValDesc, opDesc, rValDesc)
			} else {
				comparisonDescription = fmt.Sprintf("could not be compared with %s: %v)", rightSourceDesc, compareErr)
			}
		}

		// --- SYNTAX-DRIVEN 'Contains' Operator ---
	case "Contains":
		targetValue := right // Target value is already resolved

		// Parse the attribute string for the special syntax
		syntaxParts := strings.SplitN(leaf.Attribute, "[*].", 2)
		isObjectSyntax := len(syntaxParts) == 2

		var basePath string
		var fieldToExtract string
		if isObjectSyntax {
			basePath = syntaxParts[0]
			fieldToExtract = syntaxParts[1] // Assume FieldName follows [*].
			if fieldToExtract == "" {
				compareErr = fmt.Errorf("invalid 'Contains' syntax in attribute '%s': missing field name after [*]", leaf.Attribute)
			}
		} else {
			basePath = leaf.Attribute // Use the full attribute string as the path
		}

		if compareErr == nil {
			// Resolve the base path (e.g., "subject.roles[*].ID" or the full "subject.tags")
			left, err = resolveAttr(basePath, subject, resource)
			if err != nil {
				// Specific error for Contains path resolution
				reason := fmt.Sprintf("%s- Condition %sfailed: could not resolve path '%s' from attribute '%s': %v)", indent, leafDescription, basePath, leaf.Attribute, err)
				return false, reason, fmt.Errorf("%sFailed to resolve path '%s' for attribute '%s' for condition %s: %w", indent, basePath, leaf.Attribute, leafDescription, err)
			}
			leftStr = formatValueForLog(left) // Format the resolved collection value

			// Ensure the resolved 'left' value is a slice/array
			collValue := reflect.ValueOf(left)
			if !collValue.IsValid() || (collValue.Kind() != reflect.Slice && collValue.Kind() != reflect.Array) {
				compareErr = fmt.Errorf("resolved path '%s' from attribute '%s' is type %T (value: %s), not a slice/array, cannot perform 'Contains'", basePath, leaf.Attribute, left, leftStr)
			} else {
				// Now iterate and compare based on syntax used
				found := false
				if isObjectSyntax {
					// --- Behavior A: Object Syntax ("path[*].Field") ---
					comparisonDescription = fmt.Sprintf("does not contain an object where field '%s' matches %s)", fieldToExtract, rightSourceDesc) // Default failure reason

					for i := range collValue.Len() {
						element := collValue.Index(i)
						if !element.IsValid() {
							continue
						}

						// Extract the specified field by name
						extractedValue, fieldFound := extractFieldByName(element, fieldToExtract)
						if fieldFound {
							boolResult := robustCompareEquals(extractedValue, targetValue)
							if boolResult {
								found = true
								break
							}
						} // else: Field not found or accessible on this element, skip it.
					}
				} else {
					// --- Behavior B: Simple Syntax ("path") ---
					comparisonDescription = fmt.Sprintf("does not contain element matching %s)", rightSourceDesc) // Default failure reason

					for i := range collValue.Len() {
						element := collValue.Index(i)
						if !element.IsValid() {
							continue
						}
						if !element.CanInterface() {
							continue
						}
						elementInterface := element.Interface()

						if robustCompareEquals(elementInterface, targetValue) {
							found = true
							break
						}
					}
				}
				res = found // Set final result
			} // End slice/array check
		} // End compareErr check

	default:
		compareErr = fmt.Errorf("unsupported operator '%s'", leaf.Operator)
	}

	// --- Final Result and Reason ---
	if compareErr != nil {
		// Format error reason consistently. Note leftStr might be empty if resolution failed early.
		rightStrFmt := formatValueForLog(right) // Format right value for error
		reason := fmt.Sprintf("%s- Condition %s(%s: %s failed comparison with %s (value: %s): %v)", indent, leafDescription, leaf.Attribute, leftStr, rightSourceDesc, rightStrFmt, compareErr)
		return false, reason, fmt.Errorf("%sComparison Error for condition %s: %w", indent, leafDescription, compareErr)
	}

	// Make sure leftStr is set correctly before calling failReason
	if leftStr == "" && left != nil { // If left was resolved successfully but leftStr wasn't formatted yet
		leftStr = formatValueForLog(left)
	} else if left == nil && leaf.Operator != "IsNull" && leaf.Operator != "IsNotNull" && !strings.Contains(fmt.Sprintf("%v", compareErr), "resolve") {
		leftStr = "<nil>"
	}

	rightStrFmt := formatValueForLog(right) // Format right value for success/fail reason message

	// Call the updated failReason helper
	finalRes, finalReason := failReason(res, leftStr, comparisonDescription, rightStrFmt) // Pass formatted right value
	return finalRes, finalReason, nil
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
		// Show type and length, maybe first few elements?
		// Avoid printing large slices.
		if v.Len() > 5 {
			return fmt.Sprintf("%T (len:%d) [preview omitted]", val, v.Len())
		}
		// Try marshalling small slices to JSON? Or just use default %v?
		return fmt.Sprintf("%v", val) // Default Go representation for smaller slices
	case reflect.Map, reflect.Struct:
		// Avoid printing large/complex maps/structs
		// Just show the type? Or a very simple representation?
		return fmt.Sprintf("%T (value omitted)", val) // Safe fallback
	default:
		// Use standard formatting for simple types (string, int, bool, float, uuid.UUID)
		return fmt.Sprintf("%v", val)
	}
}

func robustCompareEquals(left, right any) bool {
	if left == nil || right == nil {
		return left == right // Nils are only equal to nil
	}

	// Attempt specific type comparisons, especially for UUIDs
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)

	leftUUID, leftIsUUID := left.(uuid.UUID)
	rightUUID, rightIsUUID := right.(uuid.UUID)

	// Case 1: Both UUIDs
	if leftIsUUID && rightIsUUID {
		return leftUUID == rightUUID
	}

	// Case 2: One UUID, one string - Try parsing the string
	if leftIsUUID && rightIsStr {
		parsedRightUUID, err := uuid.Parse(rightStr)
		return err == nil && leftUUID == parsedRightUUID
	}
	if leftIsStr && rightIsUUID {
		parsedLeftUUID, err := uuid.Parse(leftStr)
		return err == nil && parsedLeftUUID == rightUUID
	}

	// Case 3: Both Strings - check direct equality and also if they represent the same UUID
	if leftIsStr && rightIsStr {
		if leftStr == rightStr {
			return true
		} // Direct string equality
		// Try parsing both as UUIDs as a fallback
		parsedLeftUUID, err1 := uuid.Parse(leftStr)
		parsedRightUUID, err2 := uuid.Parse(rightStr)
		if err1 == nil && err2 == nil {
			return parsedLeftUUID == parsedRightUUID
		}
		return false // String contents differ and don't parse to comparable UUIDs
	}

	// Case 4: Numeric types? Try converting to float64 (can lose precision)
	leftFloat, leftIsNum := toFloat(left)    // Assumes toFloat helper exists
	rightFloat, rightIsNum := toFloat(right) // Assumes toFloat helper exists
	if leftIsNum && rightIsNum {
		// Be cautious with float comparisons!
		return leftFloat == rightFloat // Direct float comparison
	}

	// Case 5: Maybe boolean?
	leftBool, leftIsBool := left.(bool)
	rightBool, rightIsBool := right.(bool)
	if leftIsBool && rightIsBool {
		return leftBool == rightBool
	}

	// Default Fallback: Use DeepEqual for complex types or unhandled pairs
	// Be aware of its limitations (unexported fields, type strictness).
	return reflect.DeepEqual(left, right)
}

func extractFieldByName(element reflect.Value, fieldName string) (any, bool) {
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
		if key == "" {
			return nil, errors.New("invalid subject attribute key (e.g., 'subject.id')")
		}
	case strings.HasPrefix(attr, "resource."):
		sourceMap = resource
		key = strings.TrimPrefix(attr, "resource.")
		contextName = "resource"
		if key == "" {
			return nil, errors.New("invalid resource attribute key (e.g., 'resource.owner_id')")
		}
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
			return nil, fmt.Errorf("attribute '%s' found non-map value at '%s' while resolving '%s'", attr, strings.Join(parts[:i], "."), parts[i])
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
