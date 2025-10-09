package nodes

import (
	"context"
	"fmt"
	"log"

	"github.com/arjun/go-workflow-engine/workflow"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBFindNode struct {
	ID         string
	Database   string
	Collection string
	Query      map[string]interface{}
	Limit      int64
	OutputKey  string
}

func NewMongoDBFindNode(def workflow.NodeDefinition) (*MongoDBFindNode, error) {
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
	filter, ok := def.Config["filter"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("filter must be an object")
	}

	var limit int64 = 10
	if limitValue, exists := def.Config["limit"]; exists {
		if limitFloat, ok := limitValue.(float64); ok {
			limit = int64(limitFloat)
		} else {
			return nil, fmt.Errorf("limit must be a number")
		}
	}

	outputKey := "results" // default
	if keyValue, exists := def.Config["outputKey"]; exists {
		if key, ok := keyValue.(string); ok {
			outputKey = key
		}
	}

	return &MongoDBFindNode{
		ID:         def.ID,
		Database:   database,
		Collection: collection,
		Query:      filter,
		Limit:      limit,
		OutputKey:  outputKey,
	}, nil
}

func (n *MongoDBFindNode) Execute(ctx map[string]interface{}) (workflow.NodeResult, error) {
	resolvedQuery, err := ResolveMapValues(n.Query, ctx)
	if err != nil {
		return workflow.NodeResult{}, fmt.Errorf("failed to resolve query: %w", err)
	}

	log.Printf("Finding documents in %s.%s with query: %v", n.Database, n.Collection, resolvedQuery)

	
	collection := MongoClient.Database(n.Database).Collection(n.Collection)
	cursor, err := collection.Find(context.Background(), resolvedQuery, options.Find().SetLimit(n.Limit))
	if err != nil {
		return workflow.NodeResult{}, fmt.Errorf("failed to find documents: %w", err)
	}

	defer cursor.Close(context.Background())

	results := make([]map[string]interface{}, 0)
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			return workflow.NodeResult{}, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, result)
	}

	ctx[n.OutputKey] = results
	ctx[n.OutputKey+"Count"] = len(results)

	log.Printf("Found %d documents in %s.%s", len(results), n.Database, n.Collection)

	return workflow.NodeResult{
		Output: "default",
		Data:   ctx,
	}, nil
}
