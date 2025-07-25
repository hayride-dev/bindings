// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package context represents the imported interface "hayride:ai/context@0.0.61".
package context

import (
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

// Message represents the type alias "hayride:ai/context@0.0.61#message".
//
// See [types.Message] for more information.
type Message = types.Message

// ErrorCode represents the enum "hayride:ai/context@0.0.61#error-code".
//
//	enum error-code {
//		unexpected-message-type,
//		push-error,
//		message-not-found,
//		unknown
//	}
type ErrorCode uint8

const (
	ErrorCodeUnexpectedMessageType ErrorCode = iota
	ErrorCodePushError
	ErrorCodeMessageNotFound
	ErrorCodeUnknown
)

var _ErrorCodeStrings = [4]string{
	"unexpected-message-type",
	"push-error",
	"message-not-found",
	"unknown",
}

// String implements [fmt.Stringer], returning the enum case name of e.
func (e ErrorCode) String() string {
	return _ErrorCodeStrings[e]
}

// MarshalText implements [encoding.TextMarshaler].
func (e ErrorCode) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], unmarshaling into an enum
// case. Returns an error if the supplied text is not one of the enum cases.
func (e *ErrorCode) UnmarshalText(text []byte) error {
	return _ErrorCodeUnmarshalCase(e, text)
}

var _ErrorCodeUnmarshalCase = cm.CaseUnmarshaler[ErrorCode](_ErrorCodeStrings[:])

// Error represents the imported resource "hayride:ai/context@0.0.61#error".
//
//	resource error
type Error cm.Resource

// ResourceDrop represents the imported resource-drop for resource "error".
//
// Drops a resource handle.
//
//go:nosplit
func (self Error) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ErrorResourceDrop((uint32)(self0))
	return
}

// Code represents the imported method "code".
//
// return the error code.
//
//	code: func() -> error-code
//
//go:nosplit
func (self Error) Code() (result ErrorCode) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_ErrorCode((uint32)(self0))
	result = (ErrorCode)((uint32)(result0))
	return
}

// Data represents the imported method "data".
//
// errors can propagated with backend specific status through a string value.
//
//	data: func() -> string
//
//go:nosplit
func (self Error) Data() (result string) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ErrorData((uint32)(self0), &result)
	return
}

// Context represents the imported resource "hayride:ai/context@0.0.61#context".
//
//	resource context
type Context cm.Resource

// ResourceDrop represents the imported resource-drop for resource "context".
//
// Drops a resource handle.
//
//go:nosplit
func (self Context) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ContextResourceDrop((uint32)(self0))
	return
}

// NewContext represents the imported constructor for resource "context".
//
//	constructor()
//
//go:nosplit
func NewContext() (result Context) {
	result0 := wasmimport_NewContext()
	result = cm.Reinterpret[Context]((uint32)(result0))
	return
}

// Messages represents the imported method "messages".
//
//	messages: func() -> result<list<message>, error>
//
//go:nosplit
func (self Context) Messages() (result cm.Result[cm.List[Message], cm.List[Message], Error]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ContextMessages((uint32)(self0), &result)
	return
}

// Push represents the imported method "push".
//
//	push: func(msg: message) -> result<_, error>
//
//go:nosplit
func (self Context) Push(msg Message) (result cm.Result[Error, struct{}, Error]) {
	self0 := cm.Reinterpret[uint32](self)
	msg0, msg1, msg2 := lower_Message(msg)
	wasmimport_ContextPush((uint32)(self0), (uint32)(msg0), (*types.MessageContent)(msg1), (uint32)(msg2), &result)
	return
}
