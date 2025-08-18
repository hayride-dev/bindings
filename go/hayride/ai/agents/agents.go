package agents

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agents"

	"go.bytecodealliance.org/cm"
)

var _ Agent = (*AgentResource)(nil)

type Agent interface {
	Name() string
	Instruction() string
	Capabilities() ([]types.Tool, error)
	Context() ([]types.Message, error)
	Push(message types.Message) error
	Execute(params types.CallToolParams) (*types.CallToolResult, error)
}

type AgentResource cm.Resource

func New(format models.Format, stream graph.GraphExecutionContextStream, options ...Option[*AgentOptions]) (Agent, error) {
	opts := defaultAgentOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	var toolsOption cm.Option[agents.Tools]
	if opts.toolbox != nil {
		tb, ok := opts.toolbox.(tools.ToolResource)
		if !ok {
			return nil, fmt.Errorf("toolbox does not implement tools.Toolbox")
		}

		toolsOption = cm.Some(agents.Tools(tb))
	} else {
		toolsOption = cm.None[agents.Tools]()
	}

	var ctxOption cm.Option[agents.Context]
	if opts.context != nil {
		c, ok := opts.context.(ctx.ContextResource)
		if !ok {
			return nil, fmt.Errorf("context does not implement ctx.Context")
		}
		ctxOption = cm.Some(agents.Context(c))
	} else {
		ctxOption = cm.None[agents.Context]()
	}

	wa := agents.NewAgent(opts.name, opts.instruction,
		toolsOption,
		ctxOption,
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

func (a AgentResource) Push(message types.Message) error {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Push(cm.Reinterpret[agents.Message](message))
	if result.IsErr() {
		return fmt.Errorf("failed to push message: %s", result.Err().Data())
	}

	return nil
}

func (a AgentResource) Execute(params types.CallToolParams) (*types.CallToolResult, error) {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Execute(cm.Reinterpret[agents.CallToolParams](params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to execute tool: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.CallToolResult](result.OK()), nil
}
