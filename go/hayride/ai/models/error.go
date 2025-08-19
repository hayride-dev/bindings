package models

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
)

type ContextError struct {
	Data string
}

func (e *ContextError) Error() string {
	return fmt.Sprintf("context error: data: %s, code: %d", e.Data, model.ErrorCodeContextError)
}

type ContextEncodeError struct {
	Data string
}

func (e *ContextEncodeError) Error() string {
	return fmt.Sprintf("context encode error: data: %s, code: %d", e.Data, model.ErrorCodeContextEncode)
}

type ContextDecodeError struct {
	Data string
}

func (e *ContextDecodeError) Error() string {
	return fmt.Sprintf("context decode error: data: %s, code: %d", e.Data, model.ErrorCodeContextDecode)
}

type ComputeError struct {
	Data string
}

func (e *ComputeError) Error() string {
	return fmt.Sprintf("compute error: data: %s, code: %d", e.Data, model.ErrorCodeComputeError)
}

type PartialDecodeError struct {
	Data string
}

func (e *PartialDecodeError) Error() string {
	return fmt.Sprintf("partial decode error: data: %s, code: %d", e.Data, model.ErrorCodePartialDecode)
}

type UnknownError struct {
	Data string
}

func (e *UnknownError) Error() string {
	return fmt.Sprintf("unknown error: data: %s, code: %d", e.Data, model.ErrorCodeUnknown)
}

func newError(err *model.Error) error {
	switch err.Code() {
	case model.ErrorCodeContextError:
		return &ContextError{
			Data: err.Data(),
		}
	case model.ErrorCodeContextEncode:
		return &ContextEncodeError{
			Data: err.Data(),
		}
	case model.ErrorCodeContextDecode:
		return &ContextDecodeError{
			Data: err.Data(),
		}
	case model.ErrorCodeComputeError:
		return &ComputeError{
			Data: err.Data(),
		}
	case model.ErrorCodePartialDecode:
		return &PartialDecodeError{
			Data: err.Data(),
		}

	default:
		return &UnknownError{
			Data: err.Data(),
		}
	}
}
