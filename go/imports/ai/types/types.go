package types

import "github.com/hayride-dev/bindings/go/internal/shared/domain/ai/msg"

type Message = msg.Message
type Role = msg.Role
type Content = msg.Content
type TextContent = msg.TextContent
type ToolInput = msg.ToolInput
type ToolOutput = msg.ToolOutput

const (
	RoleUser Role = iota
	RoleAssistant
	RoleSystem
	RoleTool
	RoleUnknown
)
