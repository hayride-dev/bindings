package agent

import (
	witAgent "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

func init() {
	witAgent.Exports.Error.Code = code
	witAgent.Exports.Error.Data = data
}

func code(rep cm.Rep) (result witAgent.ErrorCode) {
	return witAgent.ErrorCode(rep)
}

func data(rep cm.Rep) (result string) {
	switch rep {
	case cm.Rep(witAgent.ErrorCodeInvokeError):
		return "invoke error"
	default:
		return "unknown error"
	}
}
