package export

import (
	"io"

	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witRunner "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/runner"
	witAgents "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agents"
	witGraph "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	witModel "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
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

	witRunner.Exports.Error.Code = errorCode
	witRunner.Exports.Error.Data = errorData
	witRunner.Exports.Error.Destructor = errorDestructor
}

func invoke(message witRunner.Message, agent cm.Rep, format cm.Rep, stream cm.Rep, writer cm.Option[cm.Rep]) cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error] {
	agentResource := cm.Reinterpret[agents.AgentResource](agent)
	defer witAgents.Agent(agent).ResourceDrop()

	formatResource := cm.Reinterpret[models.FormatResource](format)
	defer witModel.Format(format).ResourceDrop()

	graphStreamResource := cm.Reinterpret[graph.GraphExecutionContextStream](stream)
	defer witGraph.GraphExecutionContextStream(stream).ResourceDrop()

	o := writer.Some()
	var outputStream io.Writer
	if o != nil {
		w := cm.Reinterpret[streams.Writer](o)
		defer witStreams.OutputStream(*o).ResourceDrop()

		outputStream = w
	}

	msg, err := r.Invoke(cm.Reinterpret[types.Message](message), agentResource, formatResource, graphStreamResource, outputStream)
	if err != nil {
		wasiErr := createError(witRunner.ErrorCodeInvokeError, err.Error())
		return cm.Err[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witRunner.Message], cm.List[witRunner.Message], witRunner.Error]](cm.Reinterpret[cm.List[witRunner.Message]](msg))
}
