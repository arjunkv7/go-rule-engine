package workflow

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	Workflow Workflow
	Nodes    map[string]Node
	Context  *WorkflowContext
}

func NewEngine() *Engine {
	return &Engine{
		Workflow: Workflow{},
		Nodes:    map[string]Node{},
	}
}

func (e *Engine) LoadWorkflowFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var workflow Workflow
	if err := json.Unmarshal(data, &workflow); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	e.Workflow = workflow
	return nil
}

func (e *Engine) LoadWorkflowFromPayload(c *gin.Context) error {
	var workflow Workflow

	if err := c.BindJSON(&workflow); err != nil {
        return fmt.Errorf("invalid workflow JSON: %w", err)
    }
    e.Workflow = workflow
    return nil
}

var NodeFactory func(NodeDefinition) (Node, error)

func (e *Engine) BuildNodes() error {
	for _, nodeDef := range e.Workflow.Nodes {
		node, err := NodeFactory(nodeDef)
		if err != nil {
			return fmt.Errorf("failed to create node %s: %w", nodeDef.ID, err)
		}
		e.Nodes[nodeDef.ID] = node
	}
	return nil
}

func (e *Engine) Execute() error {
	ctx := NewWorkflowContext(nil)
	e.Context = ctx
	startNode := e.findStartNode()
	if startNode == nil {
		return fmt.Errorf("no start node found")
	}
	log.Printf("Starting workflow: %s", e.Workflow.Name)
	return e.executeNode(startNode.ID, ctx)

}

func (e *Engine) findStartNode() *NodeDefinition {
	for _, nodeDef := range e.Workflow.Nodes {
		if nodeDef.Type == "start" {
			return &nodeDef
		}
	}
	return nil
}

func (e *Engine) executeNode(nodeId string, ctx *WorkflowContext) error {
	node, exist := e.Nodes[nodeId]
	if !exist {
		return fmt.Errorf("node not found")
	}

	response, err := node.Execute(ctx.GetAll())
	if err != nil {
		return fmt.Errorf("error executing node %s: %w", nodeId, err)
	}

	log.Printf("Node %s executed. Output: %s", nodeId, response.Output)

	// Update context with result data - CRITICAL!
	for key, value := range response.Data {
		ctx.Set(key, value)
	}

	log.Printf("Context after node %s: %v", nodeId, ctx.GetAll())

	nextNodes := e.findNextNodes(nodeId, response.Output)

	if len(nextNodes) == 0 {
		log.Printf("No more nodes to execute. Workflow complete")
		return nil
	}

	if len(nextNodes) == 1 {
		return e.executeNode(nextNodes[0], ctx)
	}

	// Multiple paths - execute in parallel
	log.Printf("Executing %d nodes in parallel", len(nextNodes))
	return e.executeNodeParallel(nextNodes, ctx)

}

func (e *Engine) findNextNodes(fromNode string, output string) []string {
	var nextNodes []string

	for _, edge := range e.Workflow.Edges {
		if edge.From == fromNode && edge.Output == output {
			nextNodes = append(nextNodes, edge.To)
		}
	}
	return nextNodes
}

func (e *Engine) executeNodeParallel(nodeIds []string, ctx *WorkflowContext) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(nodeIds))
	for _, nodeId := range nodeIds {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			log.Printf("Started executing id: %s ", id)
			err := e.executeNode(id, ctx)
			if err != nil {
				log.Printf("Execution error for the node: %s", id)
				errChan <- err
			}
		}(nodeId)

	}

	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
