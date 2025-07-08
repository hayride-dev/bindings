package tools

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/domain"
	witTools "github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/tools"

	"go.bytecodealliance.org/cm"
)

type Tools interface {
	Call(input domain.ToolInput) (*domain.ToolOutput, error)
	Capabilities() ([]domain.ToolSchema, error)
}

type Toolbox cm.Resource

func New() (Toolbox, error) {
	return Toolbox(witTools.NewTools()), nil
}

func (t Toolbox) Call(input domain.ToolInput) (*domain.ToolOutput, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.Call(cm.Reinterpret[witTools.ToolInput](input))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to call tool: %s", result.Err().String())
	}

	return cm.Reinterpret[*domain.ToolOutput](result.OK()), nil
}

func (t Toolbox) Capabilities() ([]domain.ToolSchema, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.Capabilities()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get capabilities: %s", result.Err().Data())
	}

	schemas := cm.Reinterpret[cm.List[domain.ToolSchema]](result.OK())
	return schemas.Slice(), nil
}
