package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/runner"
	"go.bytecodealliance.org/cm"
)

// error represents the error resource
type error struct {
	Code runner.ErrorCode
	Data string
}

// createError creates a new error resource and stores it in the resource table.
func createError(code runner.ErrorCode, data string) runner.Error {
	err := error{
		Code: code,
		Data: data,
	}
	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&err))))
	resourceTable.errors[key] = err
	return runner.ErrorResourceNew(key)
}

func errorCode(self cm.Rep) runner.ErrorCode {
	err, ok := resourceTable.errors[self]
	if !ok {
		return runner.ErrorCodeUnknown
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
