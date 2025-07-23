package runner

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/runner"
	"github.com/hayride-dev/bindings/go/wasi/streams"

	"go.bytecodealliance.org/cm"
)

// TODO: What to use in place of runner.Agent?
type Agent = runner.Agent

type Runner interface {
	Invoke(message types.Message, agent Agent) ([]types.Message, error)
	InvokeStream(message types.Message, writer io.Writer, agent Agent) error
}

func Invoke(message types.Message, agent Agent) ([]types.Message, error) {
	result := runner.Invoke(cm.Reinterpret[runner.Message](message), cm.Reinterpret[runner.Agent](agent))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent")
	}

	msgs := result.OK().Slice()
	return cm.Reinterpret[[]types.Message](msgs), nil
}

func InvokeStream(message types.Message, writer io.Writer, agent Agent) error {
	w, ok := writer.(streams.Writer)
	if !ok {
		return fmt.Errorf("writer does not implement wasi io outputstream resource")
	}

	agentMessage := cm.Reinterpret[runner.Message](message)
	agentOutputStream := cm.Reinterpret[runner.OutputStream](w)

	result := runner.InvokeStream(agentMessage, agentOutputStream, agent)
	if result.IsErr() {
		return fmt.Errorf("failed to invoke agent")
	}

	return nil
}
