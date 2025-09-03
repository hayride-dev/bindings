package mcp

import (
	"github.com/hayride-dev/bindings/go/internal/gen/types/hayride/mcp/types"
	"go.bytecodealliance.org/cm"
)

type None = struct{}

type Tool = types.Tool
type CallToolParams = types.CallToolParams
type CallToolResult = types.CallToolResult
type ListToolsResult = types.ListToolsResult

type Prompt = types.Prompt
type GetPromptParams = types.GetPromptParams
type GetPromptResult = types.GetPromptResult
type ListPromptsResult = types.ListPromptsResult

type Resource = types.McpResource
type ReadResourceParams = types.ReadResourceParams
type ReadResourceResult = types.ReadResourceResult
type ListResourcesResult = types.ListResourcesResult
type ListResourceTemplatesResult = types.ListResourceTemplatesResult

type ToolSchema = types.ToolSchema

type Content = types.Content
type TextContent = types.TextContent
type ImageContent = types.ImageContent
type AudioContent = types.AudioContent
type ResourceLinkContent = types.ResourceLinkContent
type EmbeddedResourceContent = types.EmbeddedResourceContent

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

type ResourceContents = types.ResourceContents
type TextResourceContents = types.TextResourceContents
type BlobResourceContents = types.BlobResourceContents

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
