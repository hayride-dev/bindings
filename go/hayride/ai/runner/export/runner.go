package export

import (
	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witRunner "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/wasi/streams"
	"go.bytecodealliance.org/cm"
)

var r runner.Runner

// Runner interface only defines error resources
type resources struct {
	errors map[cm.Rep]error
}

var resourceTable = &resources{
	errors: make(map[cm.Rep]error),
}

func Runner(runner runner.Runner) {
	r = runner

	witRunner.Exports.Invoke = invoke
	witRunner.Exports.InvokeStream = invokeStream

	witRunner.Exports.Error.Code = errorCode
	witRunner.Exports.Error.Data = errorData
	witRunner.Exports.Error.Destructor = errorDestructor
}

func invoke(message witRunner.Message, agent witRunner.Agent) cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error] {
	msg, err := r.Invoke(cm.Reinterpret[types.Message](message), cm.Reinterpret[agents.AgentResource](agent))
	if err != nil {
		wasiErr := createError(witRunner.ErrorCodeInvokeError, err.Error())
		return cm.Err[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](cm.Reinterpret[cm.List[witRunner.Message]](msg))
}

func invokeStream(message witRunner.Message, writer witRunner.OutputStream, agent witRunner.Agent) cm.Result[witRunner.Error, struct{}, witRunner.Error] {
	w := cm.Reinterpret[streams.Writer](writer)

	err := r.InvokeStream(cm.Reinterpret[types.Message](message), w, cm.Reinterpret[agents.AgentResource](agent))
	if err != nil {
		wasiErr := createError(witRunner.ErrorCodeInvokeError, err.Error())
		return cm.Err[cm.Result[witRunner.Error, struct{}, witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witRunner.Error, struct{}, witRunner.Error]](struct{}{})
}
