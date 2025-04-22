package ctx

import (
	witContext "github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

func init() {
	witContext.Exports.Error.Code = code
	witContext.Exports.Error.Data = data
}

func code(rep cm.Rep) (result witContext.ErrorCode) {
	return witContext.ErrorCode(rep)
}

func data(rep cm.Rep) (result string) {
	switch rep {
	case cm.Rep(witContext.ErrorCodeMessageNotFound):
		return witContext.ErrorCodeMessageNotFound.String()
	case cm.Rep(witContext.ErrorCodePushError):
		return witContext.ErrorCodePushError.String()
	case cm.Rep(witContext.ErrorCodeUnexpectedMessageType):
		return witContext.ErrorCodeUnexpectedMessageType.String()
	default:
		return witContext.ErrorCodeUnknown.String()
	}
}
