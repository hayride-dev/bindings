package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/ai/graph"
	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	witAgents "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agents"
	"go.bytecodealliance.org/cm"
)

type Constructor func(name string, instruction string, tools tools.Tools, context ctx.Context, format models.Format, graph graph.GraphExecutionContextStream) agents.Agent

var agentConstructor Constructor

type resources struct {
	agents map[cm.Rep]agents.Agent
}

var resourceTable = &resources{
	agents: make(map[cm.Rep]agents.Agent),
}

func init() {
}
func Export(c Constructor) {
	agentConstructor = c

	witAgents.Exports.Agent.Constructor = constructor
	witAgents.Exports.Agent.Destructor = destructor
	witAgents.Exports.Agent.Name = name
	witAgents.Exports.Agent.Instruction = instruction
	witAgents.Exports.Agent.Context = contextFunc
	witAgents.Exports.Agent.Graph = graphFunc
	witAgents.Exports.Agent.Tools = toolsFunc
	witAgents.Exports.Agent.Format = formatFunc
}

func constructor(name string, instruction string, t witAgents.Tools, context witAgents.Context, format witAgents.Format, g witAgents.GraphExecutionContextStream) witAgents.Agent {
	agent := agentConstructor(name, instruction, cm.Reinterpret[tools.ToolResource](t), cm.Reinterpret[ctx.ContextResource](context), cm.Reinterpret[models.FormatResource](format), cm.Reinterpret[graph.GraphExecCtxStream](g))

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&agent))))
	v := witAgents.AgentResourceNew(key)
	resourceTable.agents[key] = agent
	return v
}

func destructor(self cm.Rep) {
	delete(resourceTable.agents, self)
}

func name(self cm.Rep) string {
	agent, ok := resourceTable.agents[self]
	if !ok {
		return ""
	}
	return agent.Name()
}

func instruction(self cm.Rep) string {
	agent, ok := resourceTable.agents[self]
	if !ok {
		return ""
	}

	return agent.Instruction()
}

func toolsFunc(self cm.Rep) witAgents.Tools {
	agent, ok := resourceTable.agents[self]
	if !ok {
		return cm.ResourceNone
	}
	return cm.Reinterpret[witAgents.Tools](agent.Tools())
}

func contextFunc(self cm.Rep) witAgents.Context {
	agent, ok := resourceTable.agents[self]
	if !ok {
		return cm.ResourceNone
	}
	return cm.Reinterpret[witAgents.Context](agent.Context())
}

func formatFunc(self cm.Rep) witAgents.Format {
	agent, ok := resourceTable.agents[self]
	if !ok {
		return cm.ResourceNone
	}
	return cm.Reinterpret[witAgents.Format](agent.Format())
}

func graphFunc(self cm.Rep) witAgents.GraphExecutionContextStream {
	agent, ok := resourceTable.agents[self]
	if !ok {
		return cm.ResourceNone
	}
	return cm.Reinterpret[witAgents.GraphExecutionContextStream](agent.Graph())
}
