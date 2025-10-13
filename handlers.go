package main

import (
	"context"
	"log"
	"net/http"

	"github.com/arjun/go-workflow-engine/workflow"
	"github.com/arjun/go-workflow-engine/workflow/nodes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExecuteWorkflowHandler(c *gin.Context) {
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
	if err := engine.Execute(nil); err != nil {
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
		"data":          engine.Context,
	})
}

func CreateWorkflowHandler(c *gin.Context) {
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

	engine.Workflow.ID = uuid.New().String()

	collection := nodes.MongoClient.Database("workflow_db").Collection("workflows")

	collection.InsertOne(context.Background(), engine.Workflow)

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "Workflow created successfully",
		"workflow_id":   engine.Workflow.ID,
		"workflow_name": engine.Workflow.Name,
	})
}

func ExecuteWorkflowByIdHandler(c *gin.Context) {
	workflowId := c.Query("workflow_id")
	engine := workflow.NewEngine()

	collection := nodes.MongoClient.Database("workflow_db").Collection("workflows")

	workflowData := collection.FindOne(context.Background(), bson.M{"id": workflowId})
	err := workflowData.Decode(&engine.Workflow)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Workflow not found",
			"details": "Workflow with id " + workflowId + " not found",
		})
		return
	}

	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to decode workflow data",
			"details": err.Error(),
		})
		return
	}

	if err := engine.BuildNodes(); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to build workflow nodes",
			"details": err.Error(),
		})
		return
	}

	inputData := make(map[string]interface{})
	if err := c.ShouldBindJSON(&inputData); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

	engine.Execute(inputData)
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "Workflow executed successfully",
		"workflow_id":   engine.Workflow.ID,
		"workflow_name": engine.Workflow.Name,
		"data":          engine.Context,
	})

}
