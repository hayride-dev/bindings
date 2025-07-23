package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witTools "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/tools"
	"go.bytecodealliance.org/cm"
)

type Constructor func() tools.Tools

var toolsConstructor Constructor

type resources struct {
	tools map[cm.Rep]tools.Tools
}

var resourceTable = &resources{
	tools: make(map[cm.Rep]tools.Tools),
}

func init() {
}

func Export(c Constructor) {
	toolsConstructor = c

	witTools.Exports.Tools.Constructor = constructor
	witTools.Exports.Tools.CallTool = call
	witTools.Exports.Tools.ListTools = list
	witTools.Exports.Tools.Destructor = destructor
}

func constructor() witTools.Tools {
	toolbox := toolsConstructor()

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&toolbox))))
	v := witTools.ToolsResourceNew(key)
	resourceTable.tools[key] = toolbox
	return v
}

func call(self cm.Rep, params witTools.CallToolParams) cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error] {
	tool, ok := resourceTable.tools[self]
	if !ok {
		wasiErr := witTools.ErrorResourceNew(cm.Rep(witTools.ErrorCodeToolCallFailed))
		return cm.Err[cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error]](wasiErr)
	}

	result, err := tool.Call(cm.Reinterpret[types.CallToolParams](params))
	if err != nil {
		wasiErr := witTools.ErrorResourceNew(cm.Rep(witTools.ErrorCodeToolCallFailed))
		return cm.Err[cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error]](cm.Reinterpret[witTools.CallToolResult](result))
}

func list(self cm.Rep, cursor string) cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error] {
	tool, ok := resourceTable.tools[self]
	if !ok {
		wasiErr := witTools.ErrorResourceNew(cm.Rep(witTools.ErrorCodeToolNotFound))
		return cm.Err[cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error]](wasiErr)
	}

	result, err := tool.List(cursor)
	if err != nil {
		wasiErr := witTools.ErrorResourceNew(cm.Rep(witTools.ErrorCodeToolNotFound))
		return cm.Err[cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error]](cm.Reinterpret[witTools.ListToolsResult](result))
}

func destructor(self cm.Rep) {
	delete(resourceTable.tools, self)
}
