package tools

import (
	"github.com/hayride-dev/bindings/go/hayride/ai"
)

type ToolBox interface {
	Call(input ai.ToolInput) (*ai.ToolOutput, error)
	Capabilities() ([]ai.ToolSchema, error)
}
