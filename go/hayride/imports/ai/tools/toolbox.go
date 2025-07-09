package tools

import (
	"github.com/hayride-dev/bindings/go/hayride/types/hayride/ai/types"
)

type ToolBox interface {
	Call(input types.ToolInput) (*types.ToolOutput, error)
	Capabilities() ([]types.ToolSchema, error)
}
