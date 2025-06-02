package agents

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agents"

	"go.bytecodealliance.org/cm"
)

type Agent cm.Resource

func New(options ...Option[*AgentOptions]) (Agent, error) {
	opts := defaultAgentOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return cm.ResourceNone, err
		}
	}
	return cm.ResourceNone, nil
}

func (a Agent) Invoke(message types.Message) (*types.Message, error) {
	wa := cm.Reinterpret[agents.Agent](a)

	result := wa.Invoke(cm.Reinterpret[agents.Message](message))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to invoke agent")
	}

	return cm.Reinterpret[*types.Message](result.OK()), nil
}

func (a Agent) InvokeStream(message types.Message, writer io.Writer) error {
	wa := cm.Reinterpret[agents.Agent](a)

	_, ok := writer.(wasiio.Writer)
	if !ok {
		return fmt.Errorf("writer does not implement wasi io outputstream resource")
	}

	result := wa.InvokeStream(cm.Reinterpret[agents.Message](message), cm.Reinterpret[agents.OutputStream](writer))
	if result.IsErr() {
		return fmt.Errorf("failed to invoke agent")
	}

	return nil
}
