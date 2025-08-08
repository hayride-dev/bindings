package server

import (
	"github.com/hayride-dev/bindings/go/hayride/mcp/auth"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
)

type MCPServerOptions struct {
	name         string
	version      string
	toolbox      tools.Tools
	authProvider auth.Provider
}

type OptionType interface {
	*MCPServerOptions
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

func WithName(name string) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.name = name
		return nil
	})
}

func WithVersion(version string) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.version = version
		return nil
	})
}

func WithTools(tools tools.Tools) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.toolbox = tools
		return nil
	})
}

func WithAuthProvider(provider auth.Provider) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.authProvider = provider
		return nil
	})
}

func defaultMCPServerOptions() *MCPServerOptions {
	return &MCPServerOptions{
		name:         "MCP Server",
		version:      "1.0.0",
		toolbox:      nil,
		authProvider: nil,
	}
}
