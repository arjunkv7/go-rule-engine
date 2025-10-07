# Go Workflow Engine

A mini workflow engine similar to n8n, built for learning Go.

## Features

- Start Node: Entry point for workflows
- Condition Node: Conditional branching based on expressions
- JSON-based workflow definitions
- Graph-based execution engine

## Project Structure

```
go-workflow-engine/
├── main.go                  # Entry point
├── workflow/
│   ├── types.go            # Core types and interfaces
│   ├── context.go          # Execution context
│   ├── engine.go           # Execution engine
│   └── nodes/
│       ├── start.go        # Start node implementation
│       └── condition.go    # Condition node implementation
└── examples/
    └── simple_workflow.json  # Sample workflow
```

## Usage

```bash
go run main.go examples/simple_workflow.json
```

## Learning Goals

- Interfaces & polymorphism
- JSON marshaling/unmarshaling
- Error handling patterns
- Graph traversal
- Expression evaluation

