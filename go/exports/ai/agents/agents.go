package agents

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

type invokeFunc func(messages []agent.Message) ([]agent.Message, error)

type resource struct {
	name       string
	invokeFunc invokeFunc
}

var agentResource *resource

func init() {
	agentResource = &resource{}
	agent.Exports.Agent.Constructor = agentResource.constructor
	agent.Exports.Agent.Invoke = agentResource.invoke
}

func Export(name string, f func(messages []agent.Message) ([]agent.Message, error)) {
	agentResource.name = name
	agentResource.invokeFunc = f
}

func (a *resource) constructor() agent.Agent {
	return agent.AgentResourceNew(cm.Rep(uintptr(unsafe.Pointer(&agentResource))))
}

func (a *resource) invoke(self cm.Rep, messages cm.List[agent.Message]) (result cm.Result[cm.List[agent.Message], cm.List[agent.Message], agent.Error]) {
	msgs, err := a.invokeFunc(messages.Slice())
	if err != nil {
		wasiErr := agent.ErrorResourceNew(cm.Rep(agent.ErrorCodeInvokeError))
		return cm.Err[cm.Result[cm.List[agent.Message], cm.List[agent.Message], agent.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[agent.Message], cm.List[agent.Message], agent.Error]](cm.ToList(msgs))
}
