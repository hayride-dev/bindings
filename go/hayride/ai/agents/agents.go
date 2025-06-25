package agents

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/ai/tools"
	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/agents"
	graphstream "github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/graph-stream"
	"github.com/hayride-dev/bindings/go/wasi/streams"

	"go.bytecodealliance.org/cm"
)

type Agent interface {
	Invoke(message types.Message) (*types.Message, error)
	InvokeStream(message types.Message, writer io.Writer) error
}

type agent cm.Resource

func New(options ...Option[*AgentOptions]) (Agent, error) {
	opts := defaultAgentOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	tools, err := tools.New(opts.tools...)
	if err != nil {
		return nil, fmt.Errorf("failed to create tools: %w", err)
	}

	ctx := ctx.New()

	format, err := models.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create format: %w", err)
	}

	// host provides a graph stream
	result := graphstream.LoadByName(opts.model)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to load graph")
	}
	graph := result.OK()
	resultCtxStream := graph.InitExecutionContextStream()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to init execution graph context stream")
	}
	stream := *resultCtxStream.OK()

	wa := agents.NewAgent(opts.name, opts.instruction,
		agents.Tools(tools),
		agents.Context(ctx),
		agents.Format(format),
		stream,
	)

	return agent(wa), nil
}

func (a agent) Invoke(message types.Message) (*types.Message, error) {
	wa := cm.Reinterpret[agents.Agent](a)

	result := wa.Invoke(cm.Reinterpret[agents.Message](message))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent")
	}

	return cm.Reinterpret[*types.Message](result.OK()), nil
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
