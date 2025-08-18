package runner

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/wasi/streams"

	"go.bytecodealliance.org/cm"
)

var _ Runner = (*runnerImpl)(nil)

type Runner interface {
	Invoke(message types.Message, agent agents.Agent, format models.Format, stream graph.GraphExecutionContextStream, writer io.Writer) ([]types.Message, error)
}

type runnerImpl struct{}

func New() Runner {
	return &runnerImpl{}
}

func (r *runnerImpl) Invoke(message types.Message, agent agents.Agent, format models.Format, stream graph.GraphExecutionContextStream, writer io.Writer) ([]types.Message, error) {
	a, ok := agent.(agents.AgentResource)
	if !ok {
		return nil, fmt.Errorf("agent does not implement hayride ai agent resource")
	}
	agentResource := cm.Reinterpret[runner.Agent](a)

	f, ok := format.(models.FormatResource)
	if !ok {
		return nil, fmt.Errorf("format does not implement hayride ai format resource")
	}
	formatResource := cm.Reinterpret[runner.Format](f)

	graphExecCtxStream, ok := stream.(graph.GraphExecCtxStream)
	if !ok {
		return nil, fmt.Errorf("stream does not implement graph.GraphExecCtxStream")
	}
	graphExecutionContextStream := cm.Reinterpret[runner.GraphExecutionContextStream](graphExecCtxStream)

	// Optional output stream
	var outputOption cm.Option[runner.OutputStream]
	if writer != nil {
		w, ok := writer.(streams.Writer)
		if !ok {
			return nil, fmt.Errorf("writer does not implement wasi io outputstream resource")
		}
		agentOutputStream := cm.Reinterpret[runner.OutputStream](w)
		outputOption = cm.Some(agentOutputStream)
	} else {
		outputOption = cm.None[runner.OutputStream]()
	}

	result := runner.Invoke(cm.Reinterpret[runner.Message](message), agentResource, formatResource, graphExecutionContextStream, outputOption)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent: %s", result.Err().Data())
	}

	msgs := result.OK().Slice()
	return cm.Reinterpret[[]types.Message](msgs), nil
}
