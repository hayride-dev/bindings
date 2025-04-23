package ctx

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

func init() {
	context.Exports.Error.Code = code
	context.Exports.Error.Data = data
}

func code(rep cm.Rep) (result context.ErrorCode) {
	return context.ErrorCode(rep)
}

func data(rep cm.Rep) (result string) {
	switch rep {
	case cm.Rep(context.ErrorCodeMessageNotFound):
		return context.ErrorCodeMessageNotFound.String()
	case cm.Rep(context.ErrorCodePushError):
		return context.ErrorCodePushError.String()
	case cm.Rep(context.ErrorCodeUnexpectedMessageType):
		return context.ErrorCodeUnexpectedMessageType.String()
	default:
		return context.ErrorCodeUnknown.String()
	}
}
