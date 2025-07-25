package runner

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/hayride/ai/agents"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/wasi/streams"

	"go.bytecodealliance.org/cm"
)

var _ Runner = (*runnerImpl)(nil)

type Runner interface {
	Invoke(message types.Message, agent agents.Agent) ([]types.Message, error)
	InvokeStream(message types.Message, writer io.Writer, agent agents.Agent) error
}

type runnerImpl struct{}

func New() Runner {
	return &runnerImpl{}
}

func (r *runnerImpl) Invoke(message types.Message, agent agents.Agent) ([]types.Message, error) {
	a, ok := agent.(agents.AgentResource)
	if !ok {
		return nil, fmt.Errorf("agent does not implement hayride ai agent resource")
	}
	agentResource := cm.Reinterpret[runner.Agent](a)

	result := runner.Invoke(cm.Reinterpret[runner.Message](message), agentResource)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent: %s", result.Err().Data())
	}

	msgs := result.OK().Slice()
	return cm.Reinterpret[[]types.Message](msgs), nil
}

func (r *runnerImpl) InvokeStream(message types.Message, writer io.Writer, agent agents.Agent) error {
	w, ok := writer.(streams.Writer)
	if !ok {
		return fmt.Errorf("writer does not implement wasi io outputstream resource")
	}

	a, ok := agent.(agents.AgentResource)
	if !ok {
		return fmt.Errorf("agent does not implement hayride ai agent resource")
	}

	agentMessage := cm.Reinterpret[runner.Message](message)
	agentOutputStream := cm.Reinterpret[runner.OutputStream](w)
	agentResource := cm.Reinterpret[runner.Agent](a)

	result := runner.InvokeStream(agentMessage, agentOutputStream, agentResource)
	if result.IsErr() {
		return fmt.Errorf("failed to invoke agent: %s", result.Err().Data())
	}

	return nil
}
