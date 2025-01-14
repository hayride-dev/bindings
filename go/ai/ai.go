package ai

import (
	"unsafe"

	wasiagent "github.com/hayride-dev/bindings/go/ai/gen/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

var defaultAgent = &agentResource{}

func init() {
	wasiagent.Exports.Initialize = initialize
	wasiagent.Exports.Agent.Description = description
	wasiagent.Exports.Agent.Enhance = enhance
	wasiagent.Exports.Agent.Capabilities = capabilities
}

type Agent interface {
	Description() string
	Capabilities() []string
	Enhance(components []string) error
}

func New(agent Agent) error {
	defaultAgent.agent = agent
	return nil
}

func initialize() cm.Result[wasiagent.Agent, wasiagent.Agent, wasiagent.Error] {
	ptr := unsafe.Pointer(defaultAgent)
	// Convert the pointer address to uint32 (truncated on 64-bit systems)
	address := uint32(uintptr(ptr))
	rep := cm.Rep(address)

	a := wasiagent.AgentResourceNew(rep)
	result := cm.OK[cm.Result[wasiagent.Agent, wasiagent.Agent, wasiagent.Error]](a)
	return result
}

// Component representation of an agent resource
type agentResource struct {
	agent Agent
}

type agentError struct {
	code wasiagent.ErrorCode
}

func description(rep cm.Rep) string {
	// Convert unsafe pointer rep to agentResource
	ptr := unsafe.Pointer(uintptr(rep))
	a := (*agentResource)(ptr)

	return a.agent.Description()
}

func enhance(rep cm.Rep, components cm.List[string]) (result cm.Result[wasiagent.Error, struct{}, wasiagent.Error]) {
	ptr := unsafe.Pointer(uintptr(rep))
	a := (*agentResource)(ptr)

	err := a.agent.Enhance(components.Slice())
	if err != nil {
		resourceErr := wasiagent.ErrorResourceNew(cm.Rep(wasiagent.ErrorCodeRuntimeError))
		return cm.Err[cm.Result[wasiagent.Error, struct{}, wasiagent.Error]](resourceErr)
	}

	ok := struct{}{}
	return cm.OK[cm.Result[wasiagent.Error, struct{}, wasiagent.Error]](ok)
}

func capabilities(rep cm.Rep) cm.Result[cm.List[string], cm.List[string], wasiagent.Error] {
	ptr := unsafe.Pointer(uintptr(rep))
	a := (*agentResource)(ptr)
	caps := a.agent.Capabilities()
	list := cm.ToList[[]string](caps)
	result := cm.OK[cm.Result[cm.List[string], cm.List[string], wasiagent.Error]](list)
	return result
}
