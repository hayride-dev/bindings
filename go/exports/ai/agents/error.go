package agents

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

func init() {
	agent.Exports.Error.Code = code
	agent.Exports.Error.Data = data
}

func code(rep cm.Rep) (result agent.ErrorCode) {
	return agent.ErrorCode(rep)
}

func data(rep cm.Rep) (result string) {
	switch rep {
	case cm.Rep(agent.ErrorCodeInvokeError):
		return "invoke error"
	default:
		return "unknown error"
	}
}
