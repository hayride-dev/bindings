package models

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

func init() {
	model.Exports.Error.Code = code
	model.Exports.Error.Data = data
}

func code(rep cm.Rep) (result model.ErrorCode) {
	return model.ErrorCode(rep)
}

func data(rep cm.Rep) (result string) {
	switch rep {
	case cm.Rep(model.ErrorCodeComputeError):
		return model.ErrorCodeComputeError.String()
	default:
		return model.ErrorCodeUnknown.String()
	}
}
