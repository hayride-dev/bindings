package runner

import (
	"fmt"
	"io"
	"net/http"

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

var _ Runner = (*RunnerResource)(nil)

type Runner interface {
	Invoke(message types.Message, agent agents.Agent, format models.Format, stream graph.GraphExecutionContextStream, writer io.Writer) ([]types.Message, error)
}

type RunnerResource cm.Rep

func New(options types.RunnerOptions) (Runner, error) {
	return RunnerResource(runner.NewRunner(cm.Reinterpret[runner.RunnerOptions](options))), nil
}

func (r RunnerResource) Invoke(message types.Message, agent agents.Agent, format models.Format, stream graph.GraphExecutionContextStream, writer io.Writer) ([]types.Message, error) {
	runnerResource := cm.Reinterpret[runner.Runner](r)

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
			// Ensure the headers are set for the response writer
			// NOTE:: This is required to initialize the output stream for the response writer
			// TODO:: Heavily follow wasi-http 0.3 updates for better output stream handling
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)
			agentOutputStream := cm.Reinterpret[runner.OutputStream](w.Writer)
			outputOption = cm.Some(agentOutputStream)
		}
	}

	result := runnerResource.Invoke(
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
