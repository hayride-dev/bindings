package agent

import (
	"fmt"

	witAgent "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/agent"

	"go.bytecodealliance.org/cm"
)

type Agent cm.Resource

func NewAgent() Agent {
	return Agent(witAgent.NewAgent())
}

func (a Agent) Invoke(messages []witAgent.Message) ([]witAgent.Message, error) {
	wa := cm.Reinterpret[witAgent.Agent](a)
	result := wa.Invoke(cm.ToList(messages))
	if result.IsErr() {
		// TODO: handle error
		return nil, fmt.Errorf("failed to invoke agent")
	}

	return result.OK().Slice(), nil
}
