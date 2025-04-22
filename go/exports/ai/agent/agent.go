package agent

import (
	"unsafe"

	witAgent "github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/agent"
	witTypes "github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type invokeFunc func(messages []witAgent.Message) ([]witAgent.Message, error)

type resource struct {
	name       string
	invokeFunc invokeFunc
}

var agent *resource

func init() {
	agent = &resource{}
	witAgent.Exports.Agent.Constructor = agent.constructor
	witAgent.Exports.Agent.Invoke = agent.invoke
}

func Export(name string, f func(messages []witAgent.Message) ([]witAgent.Message, error)) {
	agent.name = name
	agent.invokeFunc = f
}

func (a *resource) constructor() witAgent.Agent {
	return witAgent.AgentResourceNew(cm.Rep(uintptr(unsafe.Pointer(&agent))))
}

func (a *resource) invoke(self cm.Rep, messages cm.List[witTypes.Message]) (result cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]) {
	msgs, err := a.invokeFunc(messages.Slice())
	if err != nil {
		wasiErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeInvokeError))
		return cm.Err[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](cm.ToList(msgs))
}
