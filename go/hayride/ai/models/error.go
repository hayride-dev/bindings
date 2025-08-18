package models

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
)

type ContextError struct {
	Code uint8
	Data string
}

func (e *ContextError) Error() string {
	return fmt.Sprintf("context error: data: %s, code: %d", e.Data, e.Code)
}

type ContextEncodeError struct {
	Code uint8
	Data string
}

func (e *ContextEncodeError) Error() string {
	return fmt.Sprintf("context encode error: data: %s, code: %d", e.Data, e.Code)
}

type ContextDecodeError struct {
	Code uint8
	Data string
}

func (e *ContextDecodeError) Error() string {
	return fmt.Sprintf("context decode error: data: %s, code: %d", e.Data, e.Code)
}

type ComputeError struct {
	Code uint8
	Data string
}

func (e *ComputeError) Error() string {
	return fmt.Sprintf("compute error: data: %s, code: %d", e.Data, e.Code)
}

type PartialDecodeError struct {
	Code uint8
	Data string
}

func (e *PartialDecodeError) Error() string {
	return fmt.Sprintf("partial decode error: data: %s, code: %d", e.Data, e.Code)
}

type UnknownError struct {
	Code uint8
	Data string
}

func (e *UnknownError) Error() string {
	return fmt.Sprintf("unknown error: data: %s, code: %d", e.Data, e.Code)
}

func newError(err *model.Error) error {
	switch err.Code() {
	case model.ErrorCodeContextError:
		return &ContextError{
			Code: uint8(err.Code()),
			Data: err.Data(),
		}
	case model.ErrorCodeContextEncode:
		return &ContextEncodeError{
			Code: uint8(err.Code()),
			Data: err.Data(),
		}
	case model.ErrorCodeContextDecode:
		return &ContextDecodeError{
			Code: uint8(err.Code()),
			Data: err.Data(),
		}
	case model.ErrorCodeComputeError:
		return &ComputeError{
			Code: uint8(err.Code()),
			Data: err.Data(),
		}
	case model.ErrorCodePartialDecode:
		return &PartialDecodeError{
			Code: uint8(err.Code()),
			Data: err.Data(),
		}

	default:
		return &UnknownError{
			Code: uint8(err.Code()),
			Data: err.Data(),
		}
	}
}
