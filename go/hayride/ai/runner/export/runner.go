package export

import (
	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witRunner "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/runner"
	witAgents "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agents"
	witStreams "github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/io/streams"
	"github.com/hayride-dev/bindings/go/wasi/streams"
	"go.bytecodealliance.org/cm"
)

var r runner.Runner

// Runner interface only defines error resources
type resources struct {
	errors map[cm.Rep]errorResource
}

var resourceTable = &resources{
	errors: make(map[cm.Rep]errorResource),
}

func Runner(runner runner.Runner) {
	r = runner

	witRunner.Exports.Invoke = invoke
	witRunner.Exports.InvokeStream = invokeStream

	witRunner.Exports.Error.Code = errorCode
	witRunner.Exports.Error.Data = errorData
	witRunner.Exports.Error.Destructor = errorDestructor
}

func invoke(message witRunner.Message, agent cm.Rep) cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error] {
	agentResource := cm.Reinterpret[agents.AgentResource](agent)
	defer witAgents.Agent(agent).ResourceDrop()

	msg, err := r.Invoke(cm.Reinterpret[types.Message](message), agentResource)
	if err != nil {
		wasiErr := createError(witRunner.ErrorCodeInvokeError, err.Error())
		return cm.Err[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](cm.Reinterpret[cm.List[witRunner.Message]](msg))
}

func invokeStream(message witRunner.Message, writer cm.Rep, agent cm.Rep) cm.Result[witRunner.Error, struct{}, witRunner.Error] {
	w := cm.Reinterpret[streams.Writer](writer)
	defer witStreams.OutputStream(writer).ResourceDrop()

	agentResource := cm.Reinterpret[agents.AgentResource](agent)
	defer witAgents.Agent(agent).ResourceDrop()

	err := r.InvokeStream(cm.Reinterpret[types.Message](message), w, agentResource)
	if err != nil {
		wasiErr := createError(witRunner.ErrorCodeInvokeError, err.Error())
		return cm.Err[cm.Result[witRunner.Error, struct{}, witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witRunner.Error, struct{}, witRunner.Error]](struct{}{})
}
