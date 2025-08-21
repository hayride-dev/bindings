package types

import (
	ai "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/ai/types"
	core "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/core/types"
	http "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/http/types"
	mcp "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/mcp/types"
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

type None = struct{}
type SessionID string
type Version string
type Text string

type Message = ai.Message
type Role = ai.Role
type MessageContent = ai.MessageContent

const (
	MessageContentNone       = "none"
	MessageContentText       = "text"
	MessageContentBlob       = "blob"
	MessageContentTools      = "tools"
	MessageContentToolInput  = "tool-input"
	MessageContentToolOutput = "tool-output"
)

type WriterType = ai.WriterType
type RunnerOptions = ai.RunnerOptions

const (
	WriterTypeSse = ai.WriterTypeSse
	WriterTypeRaw = ai.WriterTypeRaw
)

type Tool = mcp.Tool
type CallToolParams = mcp.CallToolParams
type CallToolResult = mcp.CallToolResult
type ListToolsResult = mcp.ListToolsResult

type Prompt = mcp.Prompt
type GetPromptParams = mcp.GetPromptParams
type GetPromptResult = mcp.GetPromptResult
type ListPromptsResult = mcp.ListPromptsResult

type Resource = mcp.McpResource
type ReadResourceParams = mcp.ReadResourceParams
type ReadResourceResult = mcp.ReadResourceResult
type ListResourcesResult = mcp.ListResourcesResult
type ListResourceTemplatesResult = mcp.ListResourceTemplatesResult

const (
	RoleUser      = ai.RoleUser
	RoleAssistant = ai.RoleAssistant
	RoleSystem    = ai.RoleSystem
	RoleTool      = ai.RoleTool
	RoleUnknown   = ai.RoleUnknown
)

type MessageContentType interface {
	None | Text | cm.List[uint8] | cm.List[Tool] | CallToolParams | CallToolResult
}

func NewMessageContent[T MessageContentType](data T) MessageContent {
	switch any(data).(type) {
	case Text:
		return cm.New[MessageContent](1, data)
	case cm.List[uint8]:
		return cm.New[MessageContent](2, data)
	case cm.List[Tool]:
		return cm.New[MessageContent](3, data)
	case CallToolParams:
		return cm.New[MessageContent](4, data)
	case CallToolResult:
		return cm.New[MessageContent](5, data)
	default:
		return cm.New[MessageContent](0, struct{}{})
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

type ToolSchema = mcp.ToolSchema

type Content = mcp.Content
type TextContent = mcp.TextContent
type ImageContent = mcp.ImageContent
type AudioContent = mcp.AudioContent
type ResourceLinkContent = mcp.ResourceLinkContent
type EmbeddedResourceContent = mcp.EmbeddedResourceContent

type ContentType interface {
	None | TextContent | ImageContent | AudioContent | ResourceLinkContent | EmbeddedResourceContent
}

func NewContent[T ContentType](data T) Content {
	switch any(data).(type) {
	case TextContent:
		return cm.New[Content](1, data)
	case ImageContent:
		return cm.New[Content](2, data)
	case AudioContent:
		return cm.New[Content](3, data)
	case ResourceLinkContent:
		return cm.New[Content](4, data)
	case EmbeddedResourceContent:
		return cm.New[Content](5, data)
	default:
		return cm.New[Content](0, struct{}{})
	}
}

type ResourceContents = mcp.ResourceContents
type TextResourceContents = mcp.TextResourceContents
type BlobResourceContents = mcp.BlobResourceContents

type ResourceContentsType interface {
	None | TextResourceContents | BlobResourceContents
}

func NewResourceContents[T ResourceContentsType](data T) ResourceContents {
	switch any(data).(type) {
	case TextResourceContents:
		return cm.New[ResourceContents](1, data)
	case BlobResourceContents:
		return cm.New[ResourceContents](2, data)
	default:
		return cm.New[ResourceContents](0, struct{}{})
	}
}
