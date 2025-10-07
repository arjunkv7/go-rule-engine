package main

import (
	"fmt"

	"github.com/arjun/go-workflow-engine/workflow"
	"github.com/arjun/go-workflow-engine/workflow/nodes"
)

func main() {
	fmt.Println("Go Workflow Engine")
	
	workflow.NodeFactory = nodes.CreateNode

	engine := workflow.NewEngine()

	engine.LoadWorkflow("./examples/simple_workflow.json")
	engine.BuildNodes()
	engine.Execute()
}
