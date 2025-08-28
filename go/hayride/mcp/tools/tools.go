package tools

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/mcp"
	witTools "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/mcp/tools"

	"go.bytecodealliance.org/cm"
)

var _ Tools = (*ToolResource)(nil)

type Tools interface {
	Call(params mcp.CallToolParams) (*mcp.CallToolResult, error)
	List(cursor string) (*mcp.ListToolsResult, error)
}

type ToolResource cm.Resource

func New() (ToolResource, error) {
	return ToolResource(witTools.NewTools()), nil
}

func (t ToolResource) Call(params mcp.CallToolParams) (*mcp.CallToolResult, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.CallTool(cm.Reinterpret[witTools.CallToolParams](params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to call tool: %s", result.Err().Data())
	}

	return cm.Reinterpret[*mcp.CallToolResult](result.OK()), nil
}

func (t ToolResource) List(cursor string) (*mcp.ListToolsResult, error) {
	witToolsToolbox := cm.Reinterpret[witTools.Tools](t)

	result := witToolsToolbox.ListTools(cursor)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get tools: %s", result.Err().Data())
	}

	return cm.Reinterpret[*mcp.ListToolsResult](result.OK()), nil
}
