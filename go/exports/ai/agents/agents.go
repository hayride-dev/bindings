package agents

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

type invokeFunc func(messages []types.Message) ([]types.Message, error)

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

func Export(name string, f func(messages []types.Message) ([]types.Message, error)) {
	agentResource.name = name
	agentResource.invokeFunc = f
}

func (a *resource) constructor() agent.Agent {
	return agent.AgentResourceNew(cm.Rep(uintptr(unsafe.Pointer(&agentResource))))
}

func (a *resource) invoke(self cm.Rep, messages cm.List[agent.Message]) cm.Result[cm.List[agent.Message], cm.List[agent.Message], agent.Error] {
	msgs := make([]types.Message, len(messages.Slice()))
	for i, msg := range messages.Slice() {
		msgs[i] = cm.Reinterpret[types.Message](msg)
	}

	msgs, err := a.invokeFunc(msgs)
	if err != nil {
		wasiErr := agent.ErrorResourceNew(cm.Rep(agent.ErrorCodeInvokeError))
		return cm.Err[cm.Result[cm.List[agent.Message], cm.List[agent.Message], agent.Error]](wasiErr)
	}

	result := make([]agent.Message, len(msgs))
	for i, msg := range msgs {
		result[i] = cm.Reinterpret[agent.Message](msg)
	}

	return cm.OK[cm.Result[cm.List[agent.Message], cm.List[agent.Message], agent.Error]](cm.ToList(result))
}
