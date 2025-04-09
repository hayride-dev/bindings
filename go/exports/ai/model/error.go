package model

import (
	witModel "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

func init() {
	witModel.Exports.Error.Code = code
	witModel.Exports.Error.Data = data
}

func code(rep cm.Rep) (result witModel.ErrorCode) {
	return witModel.ErrorCode(rep)
}

func data(rep cm.Rep) (result string) {
	switch rep {
	case cm.Rep(witModel.ErrorCodeComputeError):
		return witModel.ErrorCodeComputeError.String()
	default:
		return witModel.ErrorCodeUnknown.String()
	}
}
