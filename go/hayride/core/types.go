package core

import (
	"github.com/hayride-dev/bindings/go/hayride/x/silo"
	"github.com/hayride-dev/bindings/go/internal/gen/types/hayride/core/types"
	"go.bytecodealliance.org/cm"
)

type Unknown = struct{}

type Cast = types.Cast
type SessionID string
type Version string
type Generate = types.Generate
type Path = string
type RequestData = types.RequestData
type ResponseData = types.ResponseData
type Request = types.Request
type Response = types.Response

// Variant is a type constraint
type RequestDataVariant interface {
	Unknown | Cast | Generate | SessionID
}

func NewRequestData[T RequestDataVariant](data T) RequestData {
	switch any(data).(type) {
	case Cast:
		return cm.New[RequestData](1, data)
	case SessionID:
		return cm.New[RequestData](2, data)
	case Generate:
		return cm.New[RequestData](3, data)
	default:
		return cm.New[RequestData](0, struct{}{})
	}
}

// Variant is a type constraint
type ResponseDataVariant interface {
	Unknown | cm.List[types.Message] | SessionID | Path | cm.List[silo.ThreadMetadata] | silo.ThreadStatus | cm.List[string] | Version
}

func NewResponseData[T ResponseDataVariant](data T) ResponseData {
	switch any(data).(type) {
	case cm.List[silo.ThreadMetadata]:
		return cm.New[ResponseData](1, data)
	case SessionID:
		return cm.New[ResponseData](2, data)
	case silo.ThreadStatus:
		return cm.New[ResponseData](3, data)
	case cm.List[types.Message]:
		return cm.New[ResponseData](4, data)
	case Path:
		return cm.New[ResponseData](5, data)
	case cm.List[string]:
		return cm.New[ResponseData](6, data)
	case Version:
		return cm.New[ResponseData](7, data)
	default:
		return cm.New[ResponseData](0, struct{}{})
	}
}
