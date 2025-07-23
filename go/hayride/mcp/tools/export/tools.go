package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/mcp/tools"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witTools "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/tools"
	"go.bytecodealliance.org/cm"
)

type Constructor func() (tools.Tools, error)

var toolsConstructor Constructor

type resources struct {
	tools  map[cm.Rep]tools.Tools
	errors map[cm.Rep]errorResource
}

var resourceTable = &resources{
	tools:  make(map[cm.Rep]tools.Tools),
	errors: make(map[cm.Rep]errorResource),
}

func Tools(c Constructor) {
	toolsConstructor = c

	witTools.Exports.Tools.Constructor = constructor
	witTools.Exports.Tools.CallTool = call
	witTools.Exports.Tools.ListTools = list
	witTools.Exports.Tools.Destructor = destructor

	witTools.Exports.Error.Code = errorCode
	witTools.Exports.Error.Data = errorData
	witTools.Exports.Error.Destructor = errorDestructor
}

func constructor() witTools.Tools {
	toolbox, err := toolsConstructor()
	if err != nil {
		return cm.ResourceNone
	}

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&toolbox))))
	v := witTools.ToolsResourceNew(key)
	resourceTable.tools[key] = toolbox
	return v
}

func call(self cm.Rep, params witTools.CallToolParams) cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error] {
	tool, ok := resourceTable.tools[self]
	if !ok {
		wasiErr := createError(witTools.ErrorCodeToolCallFailed, "failed to find tool resource")
		return cm.Err[cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error]](wasiErr)
	}

	result, err := tool.Call(cm.Reinterpret[types.CallToolParams](params))
	if err != nil {
		wasiErr := createError(witTools.ErrorCodeToolCallFailed, err.Error())
		return cm.Err[cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witTools.CallToolResultShape, witTools.CallToolResult, witTools.Error]](cm.Reinterpret[witTools.CallToolResult](result))
}

func list(self cm.Rep, cursor string) cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error] {
	tool, ok := resourceTable.tools[self]
	if !ok {
		wasiErr := createError(witTools.ErrorCodeToolNotFound, "failed to find tool resource")
		return cm.Err[cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error]](wasiErr)
	}

	result, err := tool.List(cursor)
	if err != nil {
		wasiErr := createError(witTools.ErrorCodeToolNotFound, err.Error())
		return cm.Err[cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witTools.ListToolsResultShape, witTools.ListToolsResult, witTools.Error]](cm.Reinterpret[witTools.ListToolsResult](result))
}

func destructor(self cm.Rep) {
	delete(resourceTable.tools, self)
}
