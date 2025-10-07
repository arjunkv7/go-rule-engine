package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arjun/go-workflow-engine/workflow"
	"github.com/arjun/go-workflow-engine/workflow/nodes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Go Workflow Engine")

	workflow.NodeFactory = nodes.CreateNode

	router := gin.Default()

	router.POST("/execute-workflow", executeWorkflowHandler)

	// Start server
	log.Println("🚀 Starting workflow engine API on :3002")
	if err := router.Run(":3002"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func executeWorkflowHandler(c *gin.Context) {
	// Create new engine for this request
	engine := workflow.NewEngine()

	// Load workflow from request body
	if err := engine.LoadWorkflowFromPayload(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid workflow definition",
			"details": err.Error(),
		})
		return
	}

	// Build nodes
	if err := engine.BuildNodes(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to build workflow nodes",
			"details": err.Error(),
		})
		return
	}

	// Execute workflow
	log.Printf("=== Executing workflow: %s ===", engine.Workflow.Name)
	if err := engine.Execute(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Workflow execution failed",
			"details": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "Workflow executed successfully",
		"workflow_id":   engine.Workflow.ID,
		"workflow_name": engine.Workflow.Name,
	})
}
