package nodes

import (
	"context"
	"fmt"
	"log"

	"github.com/arjun/go-workflow-engine/workflow"
)

type MongoDBInsertNode struct {
	ID         string
	Database   string
	Collection string
	Document   map[string]interface{}
}

func NewMongoDBInsertNode(def workflow.NodeDefinition) (*MongoDBInsertNode, error) {
	// Extract and validate database
	database, ok := def.Config["database"].(string)
	if !ok {
		return nil, fmt.Errorf("database must be a string")
	}

	// Extract and validate collection
	collection, ok := def.Config["collection"].(string)
	if !ok {
		return nil, fmt.Errorf("collection must be a string")
	}

	// Extract and validate document
	document, ok := def.Config["document"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("document must be an object")
	}

	return &MongoDBInsertNode{
		ID:         def.ID,
		Database:   database,
		Collection: collection,
		Document:   document,
	}, nil
}

func (n *MongoDBInsertNode) Execute(ctx map[string]interface{}) (workflow.NodeResult, error) {
	resolvedDoc, err := ResolveMapValues(n.Document, ctx)
	if err != nil {
		return workflow.NodeResult{}, fmt.Errorf("failed to resolve document values: %w", err)
	}

	collection := MongoClient.Database(n.Database).Collection(n.Collection)
	result, err := collection.InsertOne(context.Background(), resolvedDoc)
	if err != nil {
		return workflow.NodeResult{}, fmt.Errorf("failed to insert document: %w", err)
	}

	log.Printf("✅ Inserted document with ID: %v", result.InsertedID)

	ctx["insertedID"] = result.InsertedID

	return workflow.NodeResult{
		Output: "default",
		Data:   ctx,
	}, nil
}
