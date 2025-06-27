package agents

type AgentOptions struct {
	name        string
	instruction string
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

func WithName(name string) Option[*AgentOptions] {
	return newFuncOption(func(m *AgentOptions) error {
		m.name = name
		return nil
	})
}

func WithInstruction(instruction string) Option[*AgentOptions] {
	return newFuncOption(func(m *AgentOptions) error {
		m.instruction = instruction
		return nil
	})
}

func defaultAgentOptions() *AgentOptions {
	return &AgentOptions{
		name:        "default-agent",
		instruction: "default-instruction",
	}
}
