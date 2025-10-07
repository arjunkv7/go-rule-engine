package nodes

import (
	"fmt"

	"github.com/arjun/go-workflow-engine/workflow"
)

func CreateNode(def workflow.NodeDefinition) (workflow.Node, error) {
	switch def.Type {
	case "start":
		return NewStartNode(def)
	case "condition":
		return NewConditionNode(def)
	case "mongodb_insert":
		return NewMongoDBInsertNode(def)

	default:
		return nil, fmt.Errorf("unknown node type: %s", def.Type)
	}
}
