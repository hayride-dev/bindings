package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

// errorResource represents the error resource
type errorResource struct {
	Code model.ErrorCode
	Data string
}

// createError creates a new error resource and stores it in the resource table.
func newErrorResource(e error) model.Error {
	switch e.(type) {
	case *models.PartialDecodeError:
		err := errorResource{
			Code: model.ErrorCodePartialDecode,
			Data: e.Error(),
		}
		key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&err))))
		resourceTable.errors[key] = err
		return model.ErrorResourceNew(key)
	default:
		err := errorResource{
			Code: model.ErrorCodeUnknown,
			Data: e.Error(),
		}
		key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&err))))
		resourceTable.errors[key] = err
		return model.ErrorResourceNew(key)
	}
}

func errorCode(self cm.Rep) model.ErrorCode {
	err, ok := resourceTable.errors[self]
	if !ok {
		return model.ErrorCodeUnknown
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
