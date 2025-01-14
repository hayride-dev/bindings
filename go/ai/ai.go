package ai

import (
	"github.com/hayride-dev/bindings/go/ai/gen/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

func init() {
	errorView := &errorView{}
	agent.Exports.Initialize = initalize
	agent.Exports.Error.Destructor = errorView.Destructor
	agent.Exports.Error.Code = errorView.Code
	agent.Exports.Error.Data = errorView.Data

}

type A interface {
	Description() string
	Capabilities() ([]string, error)
	Enhance(components []string) error
}

func Initialize(a A) error {
	bb := b{a: a}
	agent.Exports.Agent.Constructor = bb.constructor
	return nil
}

func initalize() (result cm.Result[agent.Agent, agent.Agent, agent.Error]) {
	return
}

type agentError struct {
}

// Resource destructor.
func destructor(self cm.Rep) {
	//e.code.
}

// Code represents the caller-defined, exported method "code".
//
// return the error code.
//
//	code: func() -> error-code
func code(self cm.Rep) (result agent.ErrorCode) {
	switch self {
	case cm.Rep(agent.ErrorCodeInvalidArgument):
		return agent.ErrorCodeInvalidArgument
	case cm.Rep(agent.ErrorCodeMissingCapability):
		return agent.ErrorCodeMissingCapability
	case cm.Rep(agent.ErrorCodeRuntimeError):
		return agent.ErrorCodeRuntimeError
	default:
		return agent.ErrorCodeUnknown
	}
}

// Data represents the caller-defined, exported method "data".
//
// errors can propagated with backend specific status through a string value.
//
//	data: func() -> string
func data(self cm.Rep) (result string) {
	switch self {
	case cm.Rep(agent.ErrorCodeInvalidArgument):
		return agent.ErrorCodeInvalidArgument.String()
	case cm.Rep(agent.ErrorCodeMissingCapability):
		return agent.ErrorCodeMissingCapability.String()
	case cm.Rep(agent.ErrorCodeRuntimeError):
		return agent.ErrorCodeRuntimeError.String()
	default:
		return agent.ErrorCodeUnknown.String()
	}
}

type agentView struct {
	a        A
	resource agent.Agent
}

func (b *agentView) constructor(description string, components cm.List[string]) (result agent.Agent) {
	b.resource = agent.AgentResourceNew(cm.Rep(0))
	agent.AgentResourceNew(cm.Rep(0))
	return b.resource
}

func (b *agentView) destructor(self cm.Rep) {
	b.resource.ResourceDrop()
}

func (b *agentView) description(self cm.Rep) string {
	return b.a.Description()
}

func (b *agentView) enhance(self cm.Rep, components cm.List[string]) (result cm.Result[struct{}, agent.Error, agent.Error]) {
	if err := b.a.Enhance(components.Slice()); err != nil {
		ae := errorView{}
		e := agent.ErrorResourceNew(cm.Rep(0))

		return cm.Err[struct{}, agent.Error, agent.Error](e)
	}

	return cm.Ok[agent.Error, struct{}, agent.Error](struct{}{})
}
