package workflow

type Node interface {
	Execute(ctx map[string]interface{}) (NodeResult, error)
}

type NodeDefinition struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type NodeResult struct {
	Output string
	Data   map[string]interface{}
}

type Workflow struct {
	ID    string           `json:"id"`
	Name  string           `json:"name"`
	Nodes []NodeDefinition `json:"nodes"`
	Edges []Edge           `json:"edges"`
}

type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Output string `json:"output"`
}
