package tools

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	witTools "github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/tools"

	"go.bytecodealliance.org/cm"
)

var _ ToolBox = toolbox(0)

type toolbox cm.Resource

func New(tools ...types.ToolSchema) (toolbox, error) {
	witList := cm.ToList(tools)
	result := witTools.NewTools(cm.Reinterpret[cm.List[witTools.ToolSchema]](witList))
	return toolbox(result), nil
}

func (t toolbox) Call(input types.ToolInput) (*types.ToolOutput, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.Call(cm.Reinterpret[witTools.ToolInput](input))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to call tool: %s", result.Err().String())
	}

	return cm.Reinterpret[*types.ToolOutput](result.OK()), nil
}

func (t toolbox) Capabilities() ([]types.ToolSchema, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.Capabilities()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get capabilities: %s", result.Err().Data())
	}

	schemas := cm.Reinterpret[cm.List[types.ToolSchema]](result.OK())
	return schemas.Slice(), nil
}
