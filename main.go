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

	// CORS middleware for UI
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve static UI files
	router.Static("/ui", "./ui")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})

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
	log.Println("📱 UI available at: http://localhost:3002/ui")
	if err := router.Run(":3002"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
