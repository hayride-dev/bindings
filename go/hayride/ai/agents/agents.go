package agents

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/ai/tools"
	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/agents"
	graphstream "github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/graph-stream"
	"github.com/hayride-dev/bindings/go/wasi/streams"

	"go.bytecodealliance.org/cm"
)

type Agent interface {
	Invoke(message types.Message) ([]types.Message, error)
	InvokeStream(message types.Message, writer io.Writer) error
}

type agent cm.Resource

func New(toolbox tools.Tools, context ctx.Context, format models.Format, stream graph.GraphExecutionContextStream, options ...Option[*AgentOptions]) (Agent, error) {
	opts := defaultAgentOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	tb, ok := toolbox.(tools.Toolbox)
	if !ok {
		return nil, fmt.Errorf("toolbox does not implement tools.Toolbox")
	}

	c, ok := context.(ctx.Ctx)
	if !ok {
		return nil, fmt.Errorf("context does not implement ctx.Context")
	}

	f, ok := format.(models.Fmt)
	if !ok {
		return nil, fmt.Errorf("format does not implement models.Format")
	}

	graphExecCtxStream, ok := stream.(graph.GraphExecCtxStream)
	if !ok {
		return nil, fmt.Errorf("stream does not implement graph.GraphExecCtxStream")
	}

	wa := agents.NewAgent(opts.name, opts.instruction,
		agents.Tools(tb),
		agents.Context(c),
		agents.Format(f),
		graphstream.GraphExecutionContextStream(graphExecCtxStream),
	)

	return agent(wa), nil
}

func (a agent) Invoke(message types.Message) ([]types.Message, error) {
	wa := cm.Reinterpret[agents.Agent](a)

	result := wa.Invoke(cm.Reinterpret[agents.Message](message))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent")
	}

	msgs := result.OK().Slice()
	return cm.Reinterpret[[]types.Message](msgs), nil
}

func (a agent) InvokeStream(message types.Message, writer io.Writer) error {
	wa := cm.Reinterpret[agents.Agent](a)

	w, ok := writer.(streams.Writer)
	if !ok {
		return fmt.Errorf("writer does not implement wasi io outputstream resource")
	}

	agentMessage := cm.Reinterpret[agents.Message](message)
	agentOutputStream := cm.Reinterpret[agents.OutputStream](w)

	result := wa.InvokeStream(agentMessage, agentOutputStream)
	if result.IsErr() {
		return fmt.Errorf("failed to invoke agent")
	}

	return nil
}
