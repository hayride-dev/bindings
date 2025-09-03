package ai

import (
	"github.com/hayride-dev/bindings/go/hayride/mcp"
	"github.com/hayride-dev/bindings/go/internal/gen/types/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type Unknown = struct{}

type None = struct{}
type SessionID string
type Text string

type Message = types.Message
type Role = types.Role

const (
	RoleUser      = types.RoleUser
	RoleAssistant = types.RoleAssistant
	RoleSystem    = types.RoleSystem
	RoleTool      = types.RoleTool
	RoleUnknown   = types.RoleUnknown
)

type MessageContent = types.MessageContent

const (
	MessageContentNone       = "none"
	MessageContentText       = "text"
	MessageContentBlob       = "blob"
	MessageContentTools      = "tools"
	MessageContentToolInput  = "tool-input"
	MessageContentToolOutput = "tool-output"
)

type MessageContentType interface {
	None | Text | cm.List[uint8] | cm.List[mcp.Tool] | mcp.CallToolParams | mcp.CallToolResult
}

func NewMessageContent[T MessageContentType](data T) MessageContent {
	switch any(data).(type) {
	case Text:
		return cm.New[MessageContent](1, data)
	case cm.List[uint8]:
		return cm.New[MessageContent](2, data)
	case cm.List[mcp.Tool]:
		return cm.New[MessageContent](3, data)
	case mcp.CallToolParams:
		return cm.New[MessageContent](4, data)
	case mcp.CallToolResult:
		return cm.New[MessageContent](5, data)
	default:
		return cm.New[MessageContent](0, struct{}{})
	}
}

type WriterType = types.WriterType
type RunnerOptions = types.RunnerOptions

const (
	WriterTypeSse = types.WriterTypeSse
	WriterTypeRaw = types.WriterTypeRaw
)
