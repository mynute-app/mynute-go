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

func NewPolicyEngine(db *gorm.DB) *Policy {
	return &Policy{DB: db}
}

func (p *Policy) CanAccess(subject, resource map[string]any, policy *model.PolicyRule) (bool, error) {
	root, err := policy.GetConditionsNode()
	if err != nil {
		return false, err
	}
	result, err := evalNode(root, subject, resource)
	if err != nil {
		return false, err
	}
	if policy.Effect == "Allow" {
		return result, nil
	}
	return !result, nil
}

func evalNode(node model.ConditionNode, subject, resource map[string]interface{}) (bool, error) {
	if node.Leaf != nil {
		return evalLeaf(*node.Leaf, subject, resource)
	}

	results := make([]bool, len(node.Children))
	for i, child := range node.Children {
		res, err := evalNode(child, subject, resource)
		if err != nil {
			return false, err
		}
		results[i] = res
	}

	switch node.LogicType {
	case "AND":
		for _, r := range results {
			if !r {
				return false, nil
			}
		}
		return true, nil
	case "OR":
		for _, r := range results {
			if r {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown logic type: %s", node.LogicType)
	}
}

func evalLeaf(leaf model.ConditionLeaf, subject, resource map[string]interface{}) (bool, error) {
	left, err := resolveAttr(leaf.Attr, subject, resource)
	if err != nil {
		return false, err
	}

	var right interface{}
	if leaf.ValueSourceAttr != "" {
		right, err = resolveAttr(leaf.ValueSourceAttr, subject, resource)
		if err != nil {
			return false, err
		}
	} else if leaf.Value != nil {
		if err := json.Unmarshal(leaf.Value, &right); err != nil {
			return false, fmt.Errorf("invalid static value: %w", err)
		}
	}

	switch leaf.Op {
	case "Equals":
		return reflect.DeepEqual(left, right), nil
	case "IsNull":
		return left == nil, nil
	case "IsNotNull":
		return left != nil, nil
	case "Contains":
		val := reflect.ValueOf(left)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return false, fmt.Errorf("attribute %s is not a slice or array", leaf.Attr)
		}
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), right) {
				return true, nil
			}
		}
		return false, nil
	case "GreaterThan", "GreaterThanOrEqual", "LessThan", "LessThanOrEqual":
		return compareNumbers(left, right, leaf.Op)
	case "StartsWith", "EndsWith", "Includes":
		return compareStrings(left, right, leaf.Op)
	case "Before", "After":
		return compareTimes(left, right, leaf.Op)
	default:
		return false, fmt.Errorf("unsupported operation: %s", leaf.Op)
	}
}

func resolveAttr(attr string, subject, resource map[string]interface{}) (interface{}, error) {
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

func compareNumbers(left, right interface{}, op string) (bool, error) {
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

func toFloat(val interface{}) (float64, bool) {
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

func compareStrings(left, right interface{}, op string) (bool, error) {
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

func compareTimes(left, right interface{}, op string) (bool, error) {
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

func toTime(val interface{}) (time.Time, error) {
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
