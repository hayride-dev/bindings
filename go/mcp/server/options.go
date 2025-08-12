package server

import (
	"time"

	"github.com/hayride-dev/bindings/go/hayride/mcp/auth"
	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
)

type MCPServerOptions struct {
	name         string
	version      string
	toolbox      tools.Tools
	authProvider auth.Provider

	// Cors Options
	corsEnabled      bool
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	exposedHeaders   []string
	allowCredentials bool
	maxAge           time.Duration
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

func WithCorsEnabled(enabled bool) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.corsEnabled = enabled
		return nil
	})
}

func WithAllowedOrigins(origins []string) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.allowedOrigins = origins
		return nil
	})
}

func WithAllowedMethods(methods []string) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.allowedMethods = methods
		return nil
	})
}

func WithAllowedHeaders(headers []string) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.allowedHeaders = headers
		return nil
	})
}

func WithExposedHeaders(headers []string) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.exposedHeaders = headers
		return nil
	})
}

func WithAllowCredentials(allow bool) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.allowCredentials = allow
		return nil
	})
}

func WithMaxAge(maxAge time.Duration) Option[*MCPServerOptions] {
	return newFuncOption(func(m *MCPServerOptions) error {
		m.maxAge = maxAge
		return nil
	})
}

func defaultMCPServerOptions() *MCPServerOptions {
	return &MCPServerOptions{
		name:         "MCP Server",
		version:      "1.0.0",
		toolbox:      nil,
		authProvider: nil,
		// Default CORS options
		corsEnabled:      true,
		allowedOrigins:   []string{"*"},
		allowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		allowedHeaders:   []string{"Content-Type", "Authorization"},
		exposedHeaders:   nil,
		allowCredentials: false,
		maxAge:           24 * time.Hour,
	}
}
