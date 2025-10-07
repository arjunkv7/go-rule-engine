package nodes

import (
	"fmt"

	"github.com/arjun/go-workflow-engine/workflow"
)

type StartNode struct {
	ID          string
	InitialData map[string]interface{}
}

func NewStartNode(def workflow.NodeDefinition) (*StartNode, error) {
	initialData, ok := def.Config["initialData"]
	if !ok {
		initialData = make(map[string]interface{})
	}

	dataMap, ok := initialData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("initialData must be an object")
	}
	var newNode StartNode = StartNode{
		ID:          def.ID,
		InitialData: dataMap,
	}

	return &newNode, nil
}

func (n *StartNode) Execute(ctx map[string]interface{}) (workflow.NodeResult, error) {

	result := workflow.NodeResult{
		Output: "default",
		Data:   n.InitialData,
	}

	return result, nil
}
