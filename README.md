# Go Workflow Engine 🚀

A powerful mini workflow engine similar to n8n, built with Go. This project demonstrates advanced Go concepts including goroutines, channels, interfaces, and concurrent execution.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Workflow JSON Structure](#workflow-json-structure)
- [Node Types](#node-types)
- [Examples](#examples)
- [Advanced Features](#advanced-features)
- [Go Concepts Demonstrated](#go-concepts-demonstrated)
- [Project Structure](#project-structure)
- [Future Enhancements](#future-enhancements)

---

## ✨ Features

- **4 Node Types**: Start, Condition, MongoDB Insert, MongoDB Find
- **Parallel Execution**: Automatically executes multiple branches concurrently using goroutines
- **Template Variables**: Dynamic data resolution with `{{variableName}}` syntax
- **Thread-Safe Context**: Concurrent data access using RWMutex
- **REST API**: HTTP endpoints for workflow execution
- **Graph-Based Execution**: Follows edges with output-based routing
- **MongoDB Integration**: Native support for database operations
- **Error Handling**: Comprehensive error propagation and logging

---

## 🏗️ Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Workflow Engine                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐      ┌──────────────┐                     │
│  │  JSON Parser │─────▶│    Engine    │                     │
│  └──────────────┘      └───────┬──────┘                     │
│                                 │                             │
│                                 ▼                             │
│                        ┌────────────────┐                    │
│                        │  Node Factory  │                    │
│                        └────────┬───────┘                    │
│                                 │                             │
│                ┌────────────────┼────────────────┐           │
│                ▼                ▼                ▼           │
│         ┌──────────┐     ┌──────────┐    ┌──────────┐      │
│         │StartNode │     │Condition │    │ MongoDB  │      │
│         │          │     │   Node   │    │  Nodes   │      │
│         └──────────┘     └──────────┘    └──────────┘      │
│                                                               │
│         ┌─────────────────────────────────────────┐         │
│         │  Thread-Safe Context (RWMutex)          │         │
│         │  Data flows through workflow             │         │
│         └─────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### Execution Flow

1. **Load Workflow**: Parse JSON into workflow definition
2. **Build Nodes**: Convert node definitions to executable nodes using factory pattern
3. **Execute**: Start from the start node and follow edges
4. **Parallel Processing**: When multiple edges share the same output, execute in parallel using goroutines
5. **Context Management**: Data flows through nodes via thread-safe context

---

## 🚀 Quick Start

### Prerequisites

- Go 1.19+
- MongoDB 4.4+ (running locally or Docker)

### Installation

```bash
# Clone the repository
git clone https://github.com/arjun/go-workflow-engine.git
cd go-workflow-engine

# Install dependencies
go mod tidy

# Start MongoDB (if using Docker)
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Run the server
go run main.go
```

The server will start on `http://localhost:3002`

### Environment Variables

```bash
# MongoDB connection string (optional, defaults to localhost)
export MONGO_URI="mongodb://localhost:27017"
```

---

## 📡 API Documentation

### Health Check

**GET** `/health`

Check if the workflow engine is running.

**Response:**
```json
{
  "status": "healthy",
  "message": "Workflow engine is running"
}
```

---

### Execute Workflow

**POST** `/execute-workflow`

Execute a workflow by sending its JSON definition.

**Request Body:**
```json
{
  "id": "workflow-1",
  "name": "My Workflow",
  "nodes": [...],
  "edges": [...]
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Workflow executed successfully",
  "workflow_id": "workflow-1",
  "workflow_name": "My Workflow"
}
```

**Error Response (400/500):**
```json
{
  "error": "Error type",
  "details": "Detailed error message"
}
```

---

## 📝 Workflow JSON Structure

### Basic Structure

```json
{
  "id": "unique-workflow-id",
  "name": "Workflow Name",
  "nodes": [
    {
      "id": "node-1",
      "type": "start|condition|mongodb_insert|mongodb_find",
      "config": {
        // Node-specific configuration
      }
    }
  ],
  "edges": [
    {
      "from": "source-node-id",
      "to": "target-node-id",
      "output": "default|true|false"
    }
  ]
}
```

### Edge Outputs

- **`default`**: Standard output (Start, Insert, Find nodes)
- **`true`**: Condition evaluated to true
- **`false`**: Condition evaluated to false

Multiple edges with the same `from` and `output` will execute **in parallel**.

---

## 🔧 Node Types

### 1. Start Node

Entry point of the workflow. Initializes data into the context.

**Configuration:**
```json
{
  "id": "start-1",
  "type": "start",
  "config": {
    "initialData": {
      "key": "value",
      "age": 25
    }
  }
}
```

**Output:** `"default"`

---

### 2. Condition Node

Conditional branching based on comparison operations.

**Configuration:**
```json
{
  "id": "condition-1",
  "type": "condition",
  "config": {
    "lhs": "{{age}}",
    "operator": ">=",
    "rhs": "18"
  }
}
```

**Supported Operators:**
- `==` (equals)
- `!=` (not equals)
- `>` (greater than)
- `<` (less than)
- `>=` (greater than or equal)
- `<=` (less than or equal)

**Output:** `"true"` or `"false"`

**Template Variables:**
- `{{variableName}}` - Resolved from context
- Literal values - Used as-is (numbers, strings)

---

### 3. MongoDB Insert Node

Inserts a document into a MongoDB collection.

**Configuration:**
```json
{
  "id": "insert-1",
  "type": "mongodb_insert",
  "config": {
    "database": "mydb",
    "collection": "users",
    "document": {
      "name": "{{name}}",
      "age": "{{age}}",
      "status": "active"
    }
  }
}
```

**Context Updates:**
- Adds `insertedID` to context with the MongoDB ObjectID

**Output:** `"default"`

---

### 4. MongoDB Find Node

Queries documents from a MongoDB collection.

**Configuration:**
```json
{
  "id": "find-1",
  "type": "mongodb_find",
  "config": {
    "database": "mydb",
    "collection": "users",
    "filter": {
      "age": "{{minAge}}",
      "status": "active"
    },
    "limit": 10,
    "outputKey": "users"
  }
}
```

**Parameters:**
- `limit` (optional): Max documents to return (default: 10)
- `outputKey` (optional): Key name for results in context (default: "results")

**Context Updates:**
- Adds `{outputKey}` with array of documents
- Adds `{outputKey}Count` with number of results

**Output:** `"default"`

**Example Context After Execution:**
```json
{
  "users": [
    {"name": "John", "age": 25},
    {"name": "Alice", "age": 30}
  ],
  "usersCount": 2
}
```

---

## 📚 Examples

### Example 1: Simple User Registration

```json
{
  "id": "user-registration",
  "name": "User Registration Workflow",
  "nodes": [
    {
      "id": "start",
      "type": "start",
      "config": {
        "initialData": {
          "name": "John Doe",
          "age": 25,
          "email": "john@example.com"
        }
      }
    },
    {
      "id": "check-age",
      "type": "condition",
      "config": {
        "lhs": "{{age}}",
        "operator": ">=",
        "rhs": "18"
      }
    },
    {
      "id": "register",
      "type": "mongodb_insert",
      "config": {
        "database": "workflow_db",
        "collection": "users",
        "document": {
          "name": "{{name}}",
          "age": "{{age}}",
          "email": "{{email}}",
          "status": "active"
        }
      }
    }
  ],
  "edges": [
    {"from": "start", "to": "check-age", "output": "default"},
    {"from": "check-age", "to": "register", "output": "true"}
  ]
}
```

**Execution:**
```bash
curl -X POST http://localhost:3002/execute-workflow \
  -H "Content-Type: application/json" \
  -d @examples/simple_workflow.json
```

---

### Example 2: Find and Report

```json
{
  "id": "find-and-report",
  "name": "Find Users and Create Report",
  "nodes": [
    {
      "id": "start",
      "type": "start",
      "config": {
        "initialData": {"minAge": 18, "country": "USA"}
      }
    },
    {
      "id": "find",
      "type": "mongodb_find",
      "config": {
        "database": "workflow_db",
        "collection": "users",
        "filter": {"age": "{{minAge}}", "country": "{{country}}"},
        "limit": 10,
        "outputKey": "foundUsers"
      }
    },
    {
      "id": "check-count",
      "type": "condition",
      "config": {
        "lhs": "{{foundUsersCount}}",
        "operator": ">",
        "rhs": "0"
      }
    },
    {
      "id": "create-report",
      "type": "mongodb_insert",
      "config": {
        "database": "workflow_db",
        "collection": "reports",
        "document": {
          "totalFound": "{{foundUsersCount}}",
          "criteria": {"minAge": "{{minAge}}", "country": "{{country}}"}
        }
      }
    }
  ],
  "edges": [
    {"from": "start", "to": "find", "output": "default"},
    {"from": "find", "to": "check-count", "output": "default"},
    {"from": "check-count", "to": "create-report", "output": "true"}
  ]
}
```

More examples in the `examples/` directory!

---

## 🚀 Advanced Features

### 1. Parallel Execution

When multiple edges have the same `from` node and `output` value, they execute **concurrently** using goroutines.

```json
{
  "edges": [
    {"from": "start", "to": "task-a", "output": "default"},
    {"from": "start", "to": "task-b", "output": "default"},
    {"from": "start", "to": "task-c", "output": "default"}
  ]
}
```

**Execution:**
```
Start Node
    ├─→ Task A (goroutine 1)
    ├─→ Task B (goroutine 2)  } Execute in parallel
    └─→ Task C (goroutine 3)
```

**Implementation:**
- Uses `sync.WaitGroup` to wait for all goroutines
- Uses channels to collect errors
- Thread-safe context with `sync.RWMutex`

---

### 2. Template Variables

Use `{{variableName}}` to dynamically resolve values from context.

**Supported in:**
- Condition node: `lhs`, `rhs`
- MongoDB Insert: All fields in `document`
- MongoDB Find: All fields in `filter`

**Example:**
```json
// Context: {name: "John", age: 25}

"document": {
  "name": "{{name}}",     // Resolves to "John"
  "age": "{{age}}",       // Resolves to 25
  "status": "active"      // Literal value
}
```

**Nested Objects:**
```json
"document": {
  "user": {
    "name": "{{name}}",
    "details": {
      "age": "{{age}}"
    }
  }
}
```

---

### 3. Context Data Flow

Data flows through the workflow via a thread-safe context:

```
Start Node:
  Context = {name: "John", age: 25}

Condition Node:
  Reads: age
  Context = {name: "John", age: 25} (unchanged)

Insert Node:
  Inserts document
  Context = {name: "John", age: 25, insertedID: "abc123"}

Find Node:
  Finds documents
  Context = {name: "John", age: 25, insertedID: "abc123", users: [...], usersCount: 5}
```

Each node:
1. Receives a **copy** of the current context
2. Performs its operation
3. **Adds** new data to the context
4. Passes updated context to next node

---

## 🎓 Go Concepts Demonstrated

### 1. Interfaces & Polymorphism
```go
type Node interface {
    Execute(ctx map[string]interface{}) (NodeResult, error)
}
```
All nodes implement the same interface, enabling polymorphic execution.

### 2. Goroutines & Channels
```go
go func() {
    result := node.Execute(ctx)
    errChan <- result.Error
}()
```
Parallel execution using goroutines with error collection via channels.

### 3. Sync Primitives
- `sync.WaitGroup`: Wait for multiple goroutines
- `sync.RWMutex`: Thread-safe read/write operations
- Buffered channels: Non-blocking error collection

### 4. Error Handling
```go
if err != nil {
    return fmt.Errorf("failed to execute: %w", err)
}
```
Error wrapping with `%w` for error chains.

### 5. Factory Pattern
```go
func CreateNode(def NodeDefinition) (Node, error) {
    switch def.Type {
    case "start": return NewStartNode(def)
    case "condition": return NewConditionNode(def)
    }
}
```
Dynamic node creation based on type.

### 6. Dependency Injection
```go
var NodeFactory func(NodeDefinition) (Node, error)
```
Breaking import cycles with dependency injection.

### 7. Type Assertions
```go
if value, ok := data.(string); ok {
    // Use value safely
}
```
Safe type conversion with comma-ok idiom.

### 8. Recursion
```go
func (e *Engine) executeNode(nodeID string, ctx *Context) error {
    // Execute node
    // Find next nodes
    return e.executeNode(nextNodeID, ctx)  // Recursive call
}
```
Graph traversal using recursion.

---

## 📁 Project Structure

```
go-workflow-engine/
├── main.go                          # HTTP server & initialization
├── go.mod                           # Go module dependencies
├── README.md                        # This file
│
├── workflow/                        # Core workflow engine
│   ├── types.go                    # Node interface, Workflow struct
│   ├── context.go                  # Thread-safe execution context
│   ├── engine.go                   # Execution engine with parallel support
│   │
│   └── nodes/                      # Node implementations
│       ├── factory.go              # Node factory pattern
│       ├── start.go                # Start node
│       ├── condition.go            # Condition node with comparisons
│       ├── mongodb.go              # MongoDB connection & helpers
│       ├── mongodb_insert.go       # MongoDB insert node
│       └── mongodb_find.go         # MongoDB find node
│
└── examples/                       # Sample workflows
    ├── simple_workflow.json        # Basic insert workflow
    ├── mongodb_workflow.json       # MongoDB operations
    ├── simple_find_workflow.json   # Simple find example
    ├── find_and_insert_workflow.json  # Find + report workflow
    └── complete_workflow.json      # Complex multi-step workflow
```

---

## 🔮 Future Enhancements

### Potential Node Types
- [ ] **HTTP Request Node**: Make API calls
- [ ] **MongoDB Update Node**: Update documents
- [ ] **MongoDB Delete Node**: Delete documents
- [ ] **Transform Node**: Data transformation/mapping
- [ ] **Delay Node**: Wait for specified time
- [ ] **Loop Node**: Iterate over arrays
- [ ] **Email Node**: Send emails
- [ ] **Webhook Node**: Trigger external webhooks

### Advanced Features
- [ ] **Workflow Persistence**: Save workflow state to resume later
- [ ] **Workflow Scheduling**: Cron-based execution
- [ ] **Sub-workflows**: Call other workflows as nodes
- [ ] **Error Handling Paths**: Dedicated error output edges
- [ ] **Retry Logic**: Automatic retry on failure
- [ ] **Metrics & Monitoring**: Execution time, success rate
- [ ] **Visual Editor**: Web UI for workflow creation
- [ ] **Workflow Versioning**: Track workflow changes
- [ ] **Conditional Edges**: Edge-level conditions
- [ ] **Join Node**: Wait for multiple parallel branches

### Performance Optimizations
- [ ] **Connection Pooling**: Reuse MongoDB connections
- [ ] **Caching**: Cache frequently accessed data
- [ ] **Streaming**: Handle large datasets with cursors
- [ ] **Timeout Management**: Per-node execution timeouts

---

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./workflow/nodes
```

### Manual Testing
```bash
# Start the server
go run main.go

# In another terminal, test workflows
curl -X POST http://localhost:3002/execute-workflow \
  -H "Content-Type: application/json" \
  -d @examples/simple_workflow.json
```

---

## 🤝 Contributing

This is a learning project, but contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

---

## 📝 License

MIT License - Feel free to use this for learning!

---

## 🙏 Acknowledgments

Built as a learning project to understand:
- Go concurrency patterns
- Workflow engine architecture
- Graph-based execution
- MongoDB integration
- REST API design

Inspired by workflow automation tools like n8n, Zapier, and Apache Airflow.

---

## 📞 Contact

For questions or feedback, please open an issue on GitHub.

---

**Happy Workflow Building!** 🚀
