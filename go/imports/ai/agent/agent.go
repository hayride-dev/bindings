package agent

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agent"

	"go.bytecodealliance.org/cm"
)

type Agent cm.Resource

func NewAgent() Agent {
	return Agent(agent.NewAgent())
}

func (a Agent) Invoke(messages []agent.Message) ([]agent.Message, error) {
	wa := cm.Reinterpret[agent.Agent](a)
	result := wa.Invoke(cm.ToList(messages))
	if result.IsErr() {
		// TODO: handle error
		return nil, fmt.Errorf("failed to invoke agent")
	}

	return result.OK().Slice(), nil
}
