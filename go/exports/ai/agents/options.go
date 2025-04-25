package agents

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
)

type AgentsOptions struct {
	name             string
	invokeFunc       func(messages []types.Message) ([]types.Message, error)
	invokeStreamFunc func(messages []types.Message, writer io.Writer) error
}

type OptionType interface {
	*AgentsOptions
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

func WithName(name string) Option[*AgentsOptions] {
	return newFuncOption(func(m *AgentsOptions) error {
		m.name = name
		return nil
	})
}

func WithInvokeFunc(f func(messages []types.Message) ([]types.Message, error)) Option[*AgentsOptions] {
	return newFuncOption(func(m *AgentsOptions) error {

		m.invokeFunc = f
		return nil
	})
}

func WithInvokeStreamFunc(f func(messages []types.Message, writer io.Writer) error) Option[*AgentsOptions] {
	return newFuncOption(func(m *AgentsOptions) error {

		m.invokeStreamFunc = f
		return nil
	})
}

func defaultAgentOptions() *AgentsOptions {
	return &AgentsOptions{
		name:             "default",
		invokeFunc:       func(messages []types.Message) ([]types.Message, error) { return nil, fmt.Errorf("invokeFunc not set") },
		invokeStreamFunc: func(messages []types.Message, writer io.Writer) error { return fmt.Errorf("invokeStreamFunc not set") },
	}
}
