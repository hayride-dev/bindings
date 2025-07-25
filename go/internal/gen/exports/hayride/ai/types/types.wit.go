// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package types represents the imported interface "hayride:ai/types@0.0.61".
package types

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/types"
	"go.bytecodealliance.org/cm"
)

// Tool represents the type alias "hayride:ai/types@0.0.61#tool".
//
// See [types.Tool] for more information.
type Tool = types.Tool

// CallToolParams represents the type alias "hayride:ai/types@0.0.61#call-tool-params".
//
// See [types.CallToolParams] for more information.
type CallToolParams = types.CallToolParams

// CallToolResult represents the type alias "hayride:ai/types@0.0.61#call-tool-result".
//
// See [types.CallToolResult] for more information.
type CallToolResult = types.CallToolResult

// Role represents the enum "hayride:ai/types@0.0.61#role".
//
//	enum role {
//		user,
//		assistant,
//		system,
//		tool,
//		unknown
//	}
type Role uint8

const (
	RoleUser Role = iota
	RoleAssistant
	RoleSystem
	RoleTool
	RoleUnknown
)

var _RoleStrings = [5]string{
	"user",
	"assistant",
	"system",
	"tool",
	"unknown",
}

// String implements [fmt.Stringer], returning the enum case name of e.
func (e Role) String() string {
	return _RoleStrings[e]
}

// MarshalText implements [encoding.TextMarshaler].
func (e Role) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], unmarshaling into an enum
// case. Returns an error if the supplied text is not one of the enum cases.
func (e *Role) UnmarshalText(text []byte) error {
	return _RoleUnmarshalCase(e, text)
}

var _RoleUnmarshalCase = cm.CaseUnmarshaler[Role](_RoleStrings[:])

// MessageContent represents the variant "hayride:ai/types@0.0.61#message-content".
//
//	variant message-content {
//		none,
//		text(string),
//		blob(list<u8>),
//		tools(list<tool>),
//		tool-input(call-tool-params),
//		tool-output(call-tool-result),
//	}
type MessageContent cm.Variant[uint8, CallToolResultShape, CallToolResult]

// MessageContentNone returns a [MessageContent] of case "none".
func MessageContentNone() MessageContent {
	var data struct{}
	return cm.New[MessageContent](0, data)
}

// None returns true if [MessageContent] represents the variant case "none".
func (self *MessageContent) None() bool {
	return self.Tag() == 0
}

// MessageContentText returns a [MessageContent] of case "text".
func MessageContentText(data string) MessageContent {
	return cm.New[MessageContent](1, data)
}

// Text returns a non-nil *[string] if [MessageContent] represents the variant case "text".
func (self *MessageContent) Text() *string {
	return cm.Case[string](self, 1)
}

// MessageContentBlob returns a [MessageContent] of case "blob".
func MessageContentBlob(data cm.List[uint8]) MessageContent {
	return cm.New[MessageContent](2, data)
}

// Blob returns a non-nil *[cm.List[uint8]] if [MessageContent] represents the variant case "blob".
func (self *MessageContent) Blob() *cm.List[uint8] {
	return cm.Case[cm.List[uint8]](self, 2)
}

// MessageContentTools returns a [MessageContent] of case "tools".
func MessageContentTools(data cm.List[Tool]) MessageContent {
	return cm.New[MessageContent](3, data)
}

// Tools returns a non-nil *[cm.List[Tool]] if [MessageContent] represents the variant case "tools".
func (self *MessageContent) Tools() *cm.List[Tool] {
	return cm.Case[cm.List[Tool]](self, 3)
}

// MessageContentToolInput returns a [MessageContent] of case "tool-input".
func MessageContentToolInput(data CallToolParams) MessageContent {
	return cm.New[MessageContent](4, data)
}

// ToolInput returns a non-nil *[CallToolParams] if [MessageContent] represents the variant case "tool-input".
func (self *MessageContent) ToolInput() *CallToolParams {
	return cm.Case[CallToolParams](self, 4)
}

// MessageContentToolOutput returns a [MessageContent] of case "tool-output".
func MessageContentToolOutput(data CallToolResult) MessageContent {
	return cm.New[MessageContent](5, data)
}

// ToolOutput returns a non-nil *[CallToolResult] if [MessageContent] represents the variant case "tool-output".
func (self *MessageContent) ToolOutput() *CallToolResult {
	return cm.Case[CallToolResult](self, 5)
}

var _MessageContentStrings = [6]string{
	"none",
	"text",
	"blob",
	"tools",
	"tool-input",
	"tool-output",
}

// String implements [fmt.Stringer], returning the variant case name of v.
func (v MessageContent) String() string {
	return _MessageContentStrings[v.Tag()]
}

// Message represents the record "hayride:ai/types@0.0.61#message".
//
//	record message {
//		role: role,
//		content: list<message-content>,
//	}
type Message struct {
	_       cm.HostLayout           `json:"-"`
	Role    Role                    `json:"role"`
	Content cm.List[MessageContent] `json:"content"`
}
