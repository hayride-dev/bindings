package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/mcp/prompts"
	"github.com/hayride-dev/bindings/go/hayride/types"
	witPrompts "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/prompts"
	"go.bytecodealliance.org/cm"
)

type Constructor func() (prompts.Prompts, error)

var promptsConstructor Constructor

type resources struct {
	prompts map[cm.Rep]prompts.Prompts
	errors  map[cm.Rep]errorResource
}

var resourceTable = &resources{
	prompts: make(map[cm.Rep]prompts.Prompts),
	errors:  make(map[cm.Rep]errorResource),
}

func Prompts(c Constructor) {
	promptsConstructor = c

	witPrompts.Exports.Prompts.Constructor = constructor
	witPrompts.Exports.Prompts.GetPrompt = get
	witPrompts.Exports.Prompts.ListPrompts = list
	witPrompts.Exports.Prompts.Destructor = destructor

	witPrompts.Exports.Error.Code = errorCode
	witPrompts.Exports.Error.Data = errorData
	witPrompts.Exports.Error.Destructor = errorDestructor
}

func constructor() witPrompts.Prompts {
	prompts, err := promptsConstructor()
	if err != nil {
		return cm.ResourceNone
	}

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&prompts))))
	v := witPrompts.PromptsResourceNew(key)
	resourceTable.prompts[key] = prompts
	return v
}

func get(self cm.Rep, params witPrompts.GetPromptParams) cm.Result[witPrompts.GetPromptResultShape, witPrompts.GetPromptResult, witPrompts.Error] {
	prompt, ok := resourceTable.prompts[self]
	if !ok {
		wasiErr := createError(witPrompts.ErrorCodePromptNotFound, "failed to find prompt resource")
		return cm.Err[cm.Result[witPrompts.GetPromptResultShape, witPrompts.GetPromptResult, witPrompts.Error]](wasiErr)
	}

	result, err := prompt.Get(cm.Reinterpret[types.GetPromptParams](params))
	if err != nil {
		wasiErr := createError(witPrompts.ErrorCodePromptNotFound, err.Error())
		return cm.Err[cm.Result[witPrompts.GetPromptResultShape, witPrompts.GetPromptResult, witPrompts.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witPrompts.GetPromptResultShape, witPrompts.GetPromptResult, witPrompts.Error]](cm.Reinterpret[witPrompts.GetPromptResult](*result))
}

func list(self cm.Rep, cursor string) cm.Result[witPrompts.ListPromptsResultShape, witPrompts.ListPromptsResult, witPrompts.Error] {
	prompt, ok := resourceTable.prompts[self]
	if !ok {
		wasiErr := createError(witPrompts.ErrorCodePromptNotFound, "failed to find prompt resource")
		return cm.Err[cm.Result[witPrompts.ListPromptsResultShape, witPrompts.ListPromptsResult, witPrompts.Error]](wasiErr)
	}

	result, err := prompt.List(cursor)
	if err != nil {
		wasiErr := createError(witPrompts.ErrorCodePromptNotFound, err.Error())
		return cm.Err[cm.Result[witPrompts.ListPromptsResultShape, witPrompts.ListPromptsResult, witPrompts.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witPrompts.ListPromptsResultShape, witPrompts.ListPromptsResult, witPrompts.Error]](cm.Reinterpret[witPrompts.ListPromptsResult](*result))
}

func destructor(self cm.Rep) {
	delete(resourceTable.prompts, self)
}
