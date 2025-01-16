package agents

type AgentOptions struct {
	agents []*Agent
}

type OptionType interface {
	*AgentOptions
}

type Option[T OptionType] interface {
	Apply(T) error
}

type funcOption[T OptionType] struct {
	f func(T) error
}

func (fo *funcOption[T]) Apply(opt T) error {
	return fo.f(opt)
}

func newFuncOption[T OptionType](f func(T) error) *funcOption[T] {
	return &funcOption[T]{
		f: f,
	}
}

func WithAgents(agents ...*Agent) Option[*AgentOptions] {
	return newFuncOption(func(a *AgentOptions) error {
		a.agents = append(a.agents, agents...)
		return nil
	})
}

func defaultAgentOptions() *AgentOptions {
	return &AgentOptions{
		agents: make([]*Agent, 0),
	}
}
