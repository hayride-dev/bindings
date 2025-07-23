package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/tools"
	"go.bytecodealliance.org/cm"
)

// errorResource represents the error resource
type errorResource struct {
	Code tools.ErrorCode
	Data string
}

// createError creates a new error resource and stores it in the resource table.
func createError(code tools.ErrorCode, data string) tools.Error {
	err := errorResource{
		Code: code,
		Data: data,
	}
	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&err))))
	resourceTable.errors[key] = err
	return tools.ErrorResourceNew(key)
}

func errorCode(self cm.Rep) tools.ErrorCode {
	err, ok := resourceTable.errors[self]
	if !ok {
		return tools.ErrorCodeUnknown
	}

	return err.Code
}

func errorData(self cm.Rep) string {
	err, ok := resourceTable.errors[self]
	if !ok {
		return ""
	}

	return err.Data
}

func errorDestructor(self cm.Rep) {
	delete(resourceTable.errors, self)
}
