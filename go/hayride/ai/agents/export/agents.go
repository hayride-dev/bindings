package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witAgents "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agents"
	"go.bytecodealliance.org/cm"
)

type Constructor func(name string, instruction string, tools tools.Tools, context ctx.Context) (agents.Agent, error)

var agentConstructor Constructor

type resources struct {
	agents map[cm.Rep]agents.Agent
	errors map[cm.Rep]errorResource
}

var resourceTable = &resources{
	agents: make(map[cm.Rep]agents.Agent),
	errors: make(map[cm.Rep]errorResource),
}

func Agent(c Constructor) {
	agentConstructor = c

	witAgents.Exports.Agent.Constructor = constructor
	witAgents.Exports.Agent.Destructor = destructor
	witAgents.Exports.Agent.Name = name
	witAgents.Exports.Agent.Instruction = instruction
	witAgents.Exports.Agent.Capabilities = capabilities
	witAgents.Exports.Agent.Push = push
	witAgents.Exports.Agent.Context = context
	witAgents.Exports.Agent.Execute = execute

	witAgents.Exports.Error.Code = errorCode
	witAgents.Exports.Error.Data = errorData
	witAgents.Exports.Error.Destructor = errorDestructor
}

func constructor(name string, instruction string, tools_ cm.Option[witAgents.Tools], context cm.Option[witAgents.Context]) witAgents.Agent {
	t := tools_.Some()
	var toolResource tools.Tools
	if t != nil {
		toolResource = cm.Reinterpret[tools.ToolResource](*t)
	}

	c := context.Some()
	var contextResource ctx.Context
	if c != nil {
		contextResource = cm.Reinterpret[ctx.ContextResource](*c)
	}

	agent, err := agentConstructor(name, instruction, toolResource, contextResource)
	if err != nil {
		return cm.ResourceNone
	}

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

func capabilities(self cm.Rep) cm.Result[cm.List[witAgents.Tool], cm.List[witAgents.Tool], witAgents.Error] {
	agent, ok := resourceTable.agents[self]
	if !ok {
		wasiErr := createError(witAgents.ErrorCodeCapabilitiesError, "failed to find agent resource")
		return cm.Err[cm.Result[cm.List[witAgents.Tool], cm.List[witAgents.Tool], witAgents.Error]](wasiErr)
	}

	capabilities, err := agent.Capabilities()
	if err != nil {
		wasiErr := createError(witAgents.ErrorCodeCapabilitiesError, err.Error())
		return cm.Err[cm.Result[cm.List[witAgents.Tool], cm.List[witAgents.Tool], witAgents.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witAgents.Tool], cm.List[witAgents.Tool], witAgents.Error]](cm.Reinterpret[cm.List[witAgents.Tool]](cm.ToList(capabilities)))
}

func context(self cm.Rep) cm.Result[cm.List[witAgents.Message], cm.List[witAgents.Message], witAgents.Error] {
	agent, ok := resourceTable.agents[self]
	if !ok {
		wasiErr := createError(witAgents.ErrorCodeContextError, "failed to find agent resource")
		return cm.Err[cm.Result[cm.List[witAgents.Message], cm.List[witAgents.Message], witAgents.Error]](wasiErr)
	}

	msg, err := agent.Context()
	if err != nil {
		wasiErr := createError(witAgents.ErrorCodeContextError, err.Error())
		return cm.Err[cm.Result[cm.List[witAgents.Message], cm.List[witAgents.Message], witAgents.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[witAgents.Message], cm.List[witAgents.Message], witAgents.Error]](cm.Reinterpret[cm.List[witAgents.Message]](cm.ToList(msg)))
}

func push(self cm.Rep, msg witAgents.Message) cm.Result[witAgents.Error, struct{}, witAgents.Error] {
	agent, ok := resourceTable.agents[self]
	if !ok {
		wasiErr := createError(witAgents.ErrorCodePushError, "failed to find agent resource")
		return cm.Err[cm.Result[witAgents.Error, struct{}, witAgents.Error]](wasiErr)
	}
	message := cm.Reinterpret[types.Message](msg)
	err := agent.Push(message)
	if err != nil {
		wasiErr := createError(witAgents.ErrorCodePushError, err.Error())
		return cm.Err[cm.Result[witAgents.Error, struct{}, witAgents.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witAgents.Error, struct{}, witAgents.Error]](struct{}{})
}

func execute(self cm.Rep, params witAgents.CallToolParams) cm.Result[witAgents.CallToolResultShape, witAgents.CallToolResult, witAgents.Error] {
	agent, ok := resourceTable.agents[self]
	if !ok {
		wasiErr := createError(witAgents.ErrorCodeExecuteError, "failed to find agent resource")
		return cm.Err[cm.Result[witAgents.CallToolResultShape, witAgents.CallToolResult, witAgents.Error]](wasiErr)
	}

	p := cm.Reinterpret[types.CallToolParams](params)

	result, err := agent.Execute(p)
	if err != nil {
		wasiErr := createError(witAgents.ErrorCodeExecuteError, err.Error())
		return cm.Err[cm.Result[witAgents.CallToolResultShape, witAgents.CallToolResult, witAgents.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witAgents.CallToolResultShape, witAgents.CallToolResult, witAgents.Error]](cm.Reinterpret[witAgents.CallToolResult](*result))
}
