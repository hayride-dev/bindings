package export

import (
	"github.com/hayride-dev/bindings/go/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witRunner "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/wasi/streams"
	"go.bytecodealliance.org/cm"
)

var r runner.Runner

func init() {
}

func Export(runner runner.Runner) {
	r = runner

	witRunner.Exports.Invoke = invoke
	witRunner.Exports.InvokeStream = invokeStream
}

func invoke(message witRunner.Message, agent witRunner.Agent) cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error] {
	msg, err := r.Invoke(cm.Reinterpret[types.Message](message), cm.Reinterpret[runner.Agent](agent))
	if err != nil {
		wasiErr := witRunner.ErrorResourceNew(cm.Rep(witRunner.ErrorCodeInvokeError))
		return cm.Err[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](cm.Reinterpret[cm.List[witRunner.Message]](msg))
}

func invokeStream(message witRunner.Message, writer witRunner.OutputStream, agent witRunner.Agent) cm.Result[witRunner.Error, struct{}, witRunner.Error] {
	w := cm.Reinterpret[streams.Writer](writer)

	err := r.InvokeStream(cm.Reinterpret[types.Message](message), w, cm.Reinterpret[runner.Agent](agent))
	if err != nil {
		wasiErr := witRunner.ErrorResourceNew(cm.Rep(witRunner.ErrorCodeInvokeError))
		return cm.Err[cm.Result[witRunner.Error, struct{}, witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witRunner.Error, struct{}, witRunner.Error]](struct{}{})
}
