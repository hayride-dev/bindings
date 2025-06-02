package agents

import "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/tools"

type AgentOptions struct {
	name        string
	instruction string
	tools       tools.Tools
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

func WithModel(name string) Option[*AgentOptions] {
	return newFuncOption(func(m *AgentOptions) error {
		m.name = name
		return nil
	})
}

func defaultAgentOptions() *AgentOptions {
	return &AgentOptions{}
}
