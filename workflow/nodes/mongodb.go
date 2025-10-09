package nodes

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongoDB(connectionString string) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Add this ping verification
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoClient = client
	log.Println("✅ Successfully connected to MongoDB")
	return nil
}

func ResolveMapValues(data map[string]interface{}, ctx map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range data {
		// Check if value is a string
		strValue, isString := value.(string)

		if isString {
			// Check if it's a template variable
			if strings.HasPrefix(strValue, "{{") && strings.HasSuffix(strValue, "}}") {
				// Extract variable name
				varName := strings.TrimPrefix(strValue, "{{")
				varName = strings.TrimSuffix(varName, "}}")
				varName = strings.TrimSpace(varName)

				// Get value from context
				ctxValue, exist := ctx[varName]
				if !exist {
					return nil, fmt.Errorf("variable %s not found in context", varName)
				}
				result[key] = ctxValue
			} else {
				// It's a regular string, keep as-is
				result[key] = strValue
			}
		} else if nestedMap, isMap := value.(map[string]interface{}); isMap {
			// It's a nested map, resolve recursively
			resolvedNested, err := ResolveMapValues(nestedMap, ctx)
			if err != nil {
				return nil, err
			}
			result[key] = resolvedNested
		} else {
			// It's a number, boolean, etc - keep as-is
			result[key] = value
		}
	}

	return result, nil
}
