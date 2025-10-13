package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arjun/go-workflow-engine/workflow"
	"github.com/arjun/go-workflow-engine/workflow/nodes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Go Workflow Engine")

	// Set up node factory
	workflow.NodeFactory = nodes.CreateNode

	// Initialize MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	log.Printf("Connecting to MongoDB: %s", mongoURI)
	if err := nodes.InitMongoDB(mongoURI); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Workflow engine is running",
		})
	})

	// Execute workflow endpoint
	router.POST("/execute-workflow", ExecuteWorkflowHandler)

	router.POST("/execute-workflow-by-id", ExecuteWorkflowByIdHandler)

	router.POST("/create-workflow", CreateWorkflowHandler)

	// Start server
	log.Println("🚀 Starting workflow engine API on :3002")
	if err := router.Run(":3002"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}


