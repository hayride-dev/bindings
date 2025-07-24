package agents

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agents"
	graphstream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"

	"go.bytecodealliance.org/cm"
)

var _ Agent = (*AgentResource)(nil)

type Agent interface {
	Name() string
	Instruction() string
	Capabilities() ([]types.Tool, error)
	Context() ([]types.Message, error)
	Compute(message types.Message) (*types.Message, error)
}

type AgentResource cm.Resource

func New(toolbox tools.Tools, context ctx.Context, format models.Format, stream graph.GraphExecutionContextStream, options ...Option[*AgentOptions]) (Agent, error) {
	opts := defaultAgentOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	tb, ok := toolbox.(tools.ToolResource)
	if !ok {
		return nil, fmt.Errorf("toolbox does not implement tools.Toolbox")
	}

	c, ok := context.(ctx.ContextResource)
	if !ok {
		return nil, fmt.Errorf("context does not implement ctx.Context")
	}

	f, ok := format.(models.FormatResource)
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

	return AgentResource(wa), nil
}

func (a AgentResource) Name() string {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Name()

	return result
}

func (a AgentResource) Instruction() string {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Instruction()

	return result
}

func (a AgentResource) Capabilities() ([]types.Tool, error) {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Capabilities()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get capabilities: %s", result.Err().Data())
	}

	return cm.Reinterpret[[]types.Tool](result.OK().Slice()), nil
}

func (a AgentResource) Context() ([]types.Message, error) {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Context()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get context: %s", result.Err().Data())
	}

	return cm.Reinterpret[[]types.Message](result.OK().Slice()), nil
}

func (a AgentResource) Compute(message types.Message) (*types.Message, error) {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Compute(cm.Reinterpret[agents.Message](message))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compute message: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.Message](result.OK()), nil
}
