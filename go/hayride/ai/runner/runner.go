package runner

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/types"
	graphstream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/wasi/net/http/handle"
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

	graphExecCtxStream, ok := stream.(graph.GraphExecCtxStream)
	if !ok {
		return nil, fmt.Errorf("stream does not implement graph.GraphExecCtxStream")
	}

	// Optional output stream
	outputOption := cm.None[runner.OutputStream]()
	if writer != nil {
		switch w := writer.(type) {
		case streams.Writer:
			agentOutputStream := cm.Reinterpret[runner.OutputStream](w)
			outputOption = cm.Some(agentOutputStream)
		case *handle.WasiResponseWriter:
			stream, err := w.Stream()
			if err != nil {
				return nil, fmt.Errorf("failed to get stream from WasiResponseWriter: %s", err)
			}

			agentOutputStream := cm.Reinterpret[runner.OutputStream](*stream)
			outputOption = cm.Some(agentOutputStream)
		}
	}

	result := runner.Invoke(
		cm.Reinterpret[runner.Message](message),
		agentResource,
		runner.Format(f),
		graphstream.GraphExecutionContextStream(graphExecCtxStream),
		outputOption)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent: %s", result.Err().Data())
	}

	msgs := result.OK().Slice()
	return cm.Reinterpret[[]types.Message](msgs), nil
}
