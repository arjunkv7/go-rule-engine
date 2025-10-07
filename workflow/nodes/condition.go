package nodes

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/arjun/go-workflow-engine/workflow"
)

type ConditionNode struct {
	ID       string
	LHS      string
	RHS      string // may need to change to interface
	Operator string
}

func NewConditionNode(def workflow.NodeDefinition) (*ConditionNode, error) {
	operator, ok := def.Config["operator"]
	if !ok {
		return nil, fmt.Errorf("operator is required")
	}

	lhs, ok := def.Config["lhs"]
	if !ok {
		return nil, fmt.Errorf("lhs is required")
	}

	rhs, ok := def.Config["rhs"]
	if !ok {
		return nil, fmt.Errorf("rhs is required")
	}

	lhsStr, ok := lhs.(string)
	if !ok {
		return nil, fmt.Errorf("lhs must be a string")
	}

	operatorStr, ok := operator.(string)
	if !ok {
		return nil, fmt.Errorf("operator must be a string")
	}

	rhsStr, ok := rhs.(string)
	if !ok {
		return nil, fmt.Errorf("rhs must be a string")
	}

	return &ConditionNode{
		ID:       def.ID,
		LHS:      lhsStr,
		RHS:      rhsStr,
		Operator: operatorStr,
	}, nil

}

func (n *ConditionNode) Execute(ctx map[string]interface{}) (workflow.NodeResult, error) {
	var output string = "false"

	lhsValue, err := resolveValue(n.LHS, ctx)
	if err != nil {
		return workflow.NodeResult{}, err
	}

	rhsValue, err := resolveValue(n.RHS, ctx)
	if err != nil {
		return workflow.NodeResult{}, err
	}

	result, err := compareValues(lhsValue, rhsValue, n.Operator)
	if err != nil {
		return workflow.NodeResult{}, err
	}

	if result {
		output = "true"
	}
	return workflow.NodeResult{
		Output: output,
		Data:   ctx,
	}, nil
}

func compareValues(lhs interface{}, rhs interface{}, operator string) (bool, error) {
	lhsFloat, lhsIsNum := toFloat64(lhs)
	rhsFloat, rhsIsNum := toFloat64(rhs)

	// If both are numbers, do numeric comparison
	if lhsIsNum && rhsIsNum {
		switch operator {
		case "==":
			return lhsFloat == rhsFloat, nil
		case "!=":
			return lhsFloat != rhsFloat, nil
		case ">":
			return lhsFloat > rhsFloat, nil
		case "<":
			return lhsFloat < rhsFloat, nil
		case ">=":
			return lhsFloat >= rhsFloat, nil
		case "<=":
			return lhsFloat <= rhsFloat, nil
		default:
			return false, fmt.Errorf("unknown operator: %s", operator)
		}
	}

	// For non-numeric values, only == and != make sense
	switch operator {
	case "==":
		return lhs == rhs, nil
	case "!=":
		return lhs != rhs, nil
	default:
		return false, fmt.Errorf("operator %s not supported for non-numeric values", operator)
	}

}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func resolveValue(template string, ctx map[string]interface{}) (interface{}, error) {
	if strings.HasPrefix(template, "{{") && strings.HasSuffix(template, "}}") {
		varName := strings.TrimPrefix(template, "{{")
		varName = strings.TrimSuffix(varName, "}}")
		varName = strings.TrimSpace(varName)

		value, exists := ctx[varName]
		if !exists {
			return nil, fmt.Errorf("variable %s not found in context", varName)
		}
		return value, nil
	}

	if numValue, err := strconv.ParseFloat(template, 64); err == nil {
		return numValue, nil
	}

	return template, nil
}
