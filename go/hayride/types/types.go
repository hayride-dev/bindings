package types

import (
	ai "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/ai/types"
	core "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/core/types"
	http "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/http/types"
	silo "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/silo/types"
	"go.bytecodealliance.org/cm"
)

type ServerConfig = http.ServerConfig

type Unknown = struct{}

type ThreadMetadata = silo.ThreadMetadata
type ThreadStatus = silo.ThreadStatus

const (
	ThreadStatusUnknown    = silo.ThreadStatusUnknown
	ThreadStatusProcessing = silo.ThreadStatusProcessing
	ThreadStatusExited     = silo.ThreadStatusExited
	ThreadStatusKilled     = silo.ThreadStatusKilled
)

type SessionID string
type Version string

type Message = ai.Message
type Role = ai.Role
type TextContent = ai.TextContent
type ToolSchema = ai.ToolSchema
type ToolInput = ai.ToolInput
type ToolOutput = ai.ToolOutput
type Content = ai.Content
type None = struct{}

const (
	RoleUser      = ai.RoleUser
	RoleAssistant = ai.RoleAssistant
	RoleSystem    = ai.RoleSystem
	RoleTool      = ai.RoleTool
	RoleUnknown   = ai.RoleUnknown
)

type ContentType interface {
	None | TextContent | ToolSchema | ToolInput | ToolOutput
}

func NewContent[T ContentType](data T) Content {
	switch any(data).(type) {
	case TextContent:
		return cm.New[Content](1, data)
	case ToolSchema:
		return cm.New[Content](2, data)
	case ToolInput:
		return cm.New[Content](3, data)
	case ToolOutput:
		return cm.New[Content](4, data)
	default:
		return cm.New[Content](0, struct{}{})
	}
}

type Cast = core.Cast
type Generate = core.Generate
type Path = string
type RequestData = core.RequestData
type ResponseData = core.ResponseData
type Request = core.Request
type Response = core.Response

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
	Unknown | cm.List[Message] | SessionID | Path | cm.List[ThreadMetadata] | ThreadStatus | cm.List[string] | Version
}

func NewResponseData[T ResponseDataVariant](data T) ResponseData {
	switch any(data).(type) {
	case cm.List[ThreadMetadata]:
		return cm.New[ResponseData](1, data)
	case SessionID:
		return cm.New[ResponseData](2, data)
	case ThreadStatus:
		return cm.New[ResponseData](3, data)
	case cm.List[Message]:
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
