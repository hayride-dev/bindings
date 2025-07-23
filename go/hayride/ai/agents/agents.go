package agents

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agents"
	graphstream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"

	"go.bytecodealliance.org/cm"
)

type Agent interface {
	Name() string
	Instruction() string
	Tools() tools.Tools
	Context() ctx.Context
	Format() models.Format
	Graph() graph.GraphExecutionContextStream
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

func (a agent) Name() string {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Name()

	return result
}

func (a agent) Instruction() string {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Instruction()

	return result
}

func (a agent) Tools() tools.Tools {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Tools()

	return cm.Reinterpret[tools.Tools](result)
}

func (a agent) Context() ctx.Context {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Context()

	return cm.Reinterpret[ctx.Context](result)
}

func (a agent) Format() models.Format {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Format()

	return cm.Reinterpret[models.Format](result)
}

func (a agent) Graph() graph.GraphExecutionContextStream {
	wa := cm.Reinterpret[agents.Agent](a)
	result := wa.Graph()

	return cm.Reinterpret[graph.GraphExecutionContextStream](result)
}
