package agents

import (
	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
)

type AgentOptions struct {
	name        string
	instruction string
	toolbox     tools.Tools
	context     ctx.Context
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

func WithTools(tools tools.Tools) Option[*AgentOptions] {
	return newFuncOption(func(m *AgentOptions) error {
		m.toolbox = tools
		return nil
	})
}

func WithContext(context ctx.Context) Option[*AgentOptions] {
	return newFuncOption(func(m *AgentOptions) error {
		m.context = context
		return nil
	})
}

func defaultAgentOptions() *AgentOptions {
	return &AgentOptions{
		name:        "default-agent",
		instruction: "default-instruction",
		toolbox:     nil,
		context:     nil,
	}
}
