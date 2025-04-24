package agents

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agent"

	"go.bytecodealliance.org/cm"
)

type Agent cm.Resource

func NewAgent() Agent {
	return Agent(agent.NewAgent())
}

func (a Agent) Invoke(messages []types.Message) ([]types.Message, error) {
	wa := cm.Reinterpret[agent.Agent](a)

	msgs := make([]agent.Message, len(messages))
	for i, msg := range messages {
		msgs[i] = cm.Reinterpret[agent.Message](msg)
	}

	result := wa.Invoke(cm.ToList(msgs))
	if result.IsErr() {
		// TODO: handle error
		return nil, fmt.Errorf("failed to invoke agent")
	}

	witMessages := make([]types.Message, len(result.OK().Slice()))
	for i, msg := range result.OK().Slice() {
		witMessages[i] = cm.Reinterpret[types.Message](msg)
	}

	return witMessages, nil
}

func (a Agent) InvokeStream(messages []types.Message, writer io.Writer) error {
	wa := cm.Reinterpret[agent.Agent](a)

	msgs := make([]agent.Message, len(messages))
	for i, msg := range messages {
		msgs[i] = cm.Reinterpret[agent.Message](msg)
	}

	result := wa.InvokeStream(cm.ToList(msgs), cm.Reinterpret[agent.OutputStream](writer))
	if result.IsErr() {
		return fmt.Errorf("failed to invoke agent")
	}

	return nil
}
