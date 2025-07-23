package tools

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/types"
	witTools "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/mcp/tools"

	"go.bytecodealliance.org/cm"
)

var _ Tools = (*ToolResource)(nil)

type Tools interface {
	Call(params types.CallToolParams) (*types.CallToolResult, error)
	List(cursor string) (*types.ListToolsResult, error)
}

type ToolResource cm.Resource

func New() (ToolResource, error) {
	return ToolResource(witTools.NewTools()), nil
}

func (t ToolResource) Call(params types.CallToolParams) (*types.CallToolResult, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.CallTool(cm.Reinterpret[witTools.CallToolParams](params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to call tool: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.CallToolResult](result.OK()), nil
}

func (t ToolResource) List(cursor string) (*types.ListToolsResult, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.ListTools(cursor)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get tools: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.ListToolsResult](result.OK()), nil
}
