package ai

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/ai/gen/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

var defaultAgent = &agentResource{}

func init() {
	agent.Exports.Initialize = initialize
	agent.Exports.Agent.Constructor = constructor
	// agent.Exports.Agent.Destructor = destructor
	agent.Exports.Agent.Description = description
}

func New(description string, components []string) (*agentResource, error) {
	defaultAgent.description = description
	defaultAgent.components = components
	return defaultAgent, nil
}

func initialize() cm.Result[agent.Agent, agent.Agent, agent.Error] {
	ptr := unsafe.Pointer(defaultAgent)
	// Convert the pointer address to uint32 (truncated on 64-bit systems)
	address := uint32(uintptr(ptr))
	rep := cm.Rep(address)

	a := agent.AgentResourceNew(rep)
	result := cm.OK[cm.Result[agent.Agent, agent.Agent, agent.Error]](a)
	return result
}

// Component representation of an agent resource
type agentResource struct {
	description string
	components  []string
}

type agentError struct {
	code agent.ErrorCode
}

func constructor(description string, components cm.List[string]) (result agent.Agent) {
	// Create a new resource.
	a := &agentResource{
		description: description,
		components:  components.Slice(),
	}

	ptr := unsafe.Pointer(a)
	// Convert the pointer address to uint32 (truncated on 64-bit systems)
	address := uint32(uintptr(ptr))
	rep := cm.Rep(address)

	return agent.AgentResourceNew(rep)
}

/*
func destructor(self cm.Rep) {
	agent := agent.AgentResourceNew(self)
	agent.ResourceDrop()
}
*/

func description(rep cm.Rep) string {
	// Convert unsafe pointer rep to agentResource
	ptr := unsafe.Pointer(uintptr(rep))
	a := (*agentResource)(ptr)

	return a.description
}

/*
func enhance(rep cm.Rep, components cm.List[string]) (result cm.Result[agent.Error, struct{}, agent.Error]) {
	ptr := unsafe.Pointer(uintptr(rep))
	a := (*agentResource)(ptr)

	a.components = append(a.components, components.Slice()...)

	// TODO: Check for duplicates
	// result = cm.OK[agent.Error, struct{}, agent.Error]()
	result = cm.Err()

	return
}

// Code represents the caller-defined, exported method "code".
//
// return the error code.
//
//	code: func() -> error-code
func (b *agentView) code(self cm.Rep) (result agent.ErrorCode) {
	// TODO: Handle the case where the resource is not an agentError.
	err := b.table[self].(agentError)
	return err.code
}

// Data represents the caller-defined, exported method "data".
//
// errors can propagated with backend specific status through a string value.
//
//	data: func() -> string
func (b *agentView) data(self cm.Rep) (result string) {
	// TODO: Handle the case where the resource is not an agentError.
	err := b.table[self].(agentError)
	return err.code.String()
}
*/
