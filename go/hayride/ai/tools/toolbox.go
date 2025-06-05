package tools

import "github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"

type ToolBox interface {
	Call(input types.ToolInput) (*types.ToolOutput, error)
	Capabilities() ([]types.ToolSchema, error)
}
