package workflow

import "sync"

type WorkflowContext struct {
	Data  map[string]interface{}
	mutex sync.RWMutex
}

func NewWorkflowContext(initialData map[string]interface{}) *WorkflowContext {
	if initialData == nil {
		initialData = make(map[string]interface{})
	}
	return &WorkflowContext{
		Data:  initialData,
		mutex: sync.RWMutex{},
	}
}

func (ctx *WorkflowContext) Get(key string) (interface{}, bool) {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()

	value, ok := ctx.Data[key]
	return value, ok
}

func (ctx *WorkflowContext) Set(key string, value interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.Data[key] = value
}

func (ctx *WorkflowContext) Delete(key string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	delete(ctx.Data, key)
}

func (ctx *WorkflowContext) GetAll() map[string]interface{} {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()

	// Return a copy to prevent external modifications
	result := make(map[string]interface{})
	for k, v := range ctx.Data {
		result[k] = v
	}
	return result
}
