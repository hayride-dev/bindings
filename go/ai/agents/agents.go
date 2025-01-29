package agents

import (
	wasiagents "github.com/hayride-dev/bindings/go/ai/gen/exports/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/ai/gen/exports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type Tool = types.Tool

type Agent struct {
	Name         string
	Description  string
	Capabilities []Tool
}

type store struct {
	agents map[string]*Agent
}

func New(options ...Option[*AgentOptions]) error {
	opts := defaultAgentOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return err
		}
	}

	s := &store{agents: make(map[string]*Agent)}

	for _, agent := range opts.agents {
		s.agents[agent.Name] = agent
	}

	wasiagents.Exports.Set = s.wasiSet
	wasiagents.Exports.Get = s.wasiGet
	wasiagents.Exports.Enhance = s.wasiEnhance

	wasiagents.Exports.Error.Code = code
	wasiagents.Exports.Error.Data = data

	return nil
}

func code(self cm.Rep) (result wasiagents.ErrorCode) {
	switch self {
	case cm.Rep(wasiagents.ErrorCodeEnhanceError):
		return wasiagents.ErrorCodeEnhanceError
	default:
		return wasiagents.ErrorCodeUnknown
	}
}

func data(self cm.Rep) string {
	switch self {
	case cm.Rep(wasiagents.ErrorCodeEnhanceError):
		return wasiagents.ErrorCodeEnhanceError.String()
	default:
		return wasiagents.ErrorCodeUnknown.String()
	}
}

func (a *store) wasiSet(wasiAgent types.Agent) (result cm.Result[wasiagents.Error, struct{}, wasiagents.Error]) {
	agent := &Agent{
		Name:         wasiAgent.Name,
		Description:  wasiAgent.Description,
		Capabilities: wasiAgent.Capabilities.Slice(),
	}
	// duplicate agent
	if _, ok := a.agents[agent.Name]; ok {
		resourceErr := wasiagents.ErrorResourceNew(cm.Rep(wasiagents.ErrorCodeUnknown))
		return cm.Err[cm.Result[wasiagents.Error, struct{}, wasiagents.Error]](resourceErr)
	}
	a.agents[agent.Name] = agent

	ok := struct{}{}
	return cm.OK[cm.Result[wasiagents.Error, struct{}, wasiagents.Error]](ok)
}

func (a *store) wasiGet(name string) (result cm.Result[wasiagents.AgentShape, types.Agent, wasiagents.Error]) {
	if _, ok := a.agents[name]; !ok {
		resourceErr := wasiagents.ErrorResourceNew(cm.Rep(wasiagents.ErrorCodeUnknown))
		return cm.Err[cm.Result[wasiagents.AgentShape, types.Agent, wasiagents.Error]](resourceErr)
	}
	agent := a.agents[name]
	wasiAgent := types.Agent{
		Name:         agent.Name,
		Description:  agent.Description,
		Capabilities: cm.ToList[[]types.Tool](agent.Capabilities),
	}
	return cm.OK[cm.Result[wasiagents.AgentShape, types.Agent, wasiagents.Error]](wasiAgent)
}

func (a *store) wasiEnhance(agent types.Agent, tools cm.List[Tool]) (result cm.Result[wasiagents.Error, struct{}, wasiagents.Error]) {
	if _, ok := a.agents[agent.Name]; !ok {
		resourceErr := wasiagents.ErrorResourceNew(cm.Rep(wasiagents.ErrorCodeEnhanceError))
		return cm.Err[cm.Result[wasiagents.Error, struct{}, wasiagents.Error]](resourceErr)
	}
	a.agents[agent.Name].Capabilities = append(a.agents[agent.Name].Capabilities, tools.Slice()...)
	ok := struct{}{}
	return cm.OK[cm.Result[wasiagents.Error, struct{}, wasiagents.Error]](ok)
}
