package tools

import (
	"github.com/hayride-dev/bindings/go/hayride/domain"
)

type ToolBox interface {
	Call(input domain.ToolInput) (*domain.ToolOutput, error)
	Capabilities() ([]domain.ToolSchema, error)
}
