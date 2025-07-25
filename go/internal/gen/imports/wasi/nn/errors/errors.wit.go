// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package errors represents the imported interface "wasi:nn/errors@0.2.0-rc-2024-10-28".
//
// TODO: create function-specific errors (https://github.com/WebAssembly/wasi-nn/issues/42)
package errors

import (
	"go.bytecodealliance.org/cm"
)

// ErrorCode represents the enum "wasi:nn/errors@0.2.0-rc-2024-10-28#error-code".
//
//	enum error-code {
//		invalid-argument,
//		invalid-encoding,
//		timeout,
//		runtime-error,
//		unsupported-operation,
//		too-large,
//		not-found,
//		security,
//		unknown
//	}
type ErrorCode uint8

const (
	// Caller module passed an invalid argument.
	ErrorCodeInvalidArgument ErrorCode = iota

	// Invalid encoding.
	ErrorCodeInvalidEncoding

	// The operation timed out.
	ErrorCodeTimeout

	// Runtime Error.
	ErrorCodeRuntimeError

	// Unsupported operation.
	ErrorCodeUnsupportedOperation

	// Graph is too large.
	ErrorCodeTooLarge

	// Graph not found.
	ErrorCodeNotFound

	// The operation is insecure or has insufficient privilege to be performed.
	// e.g., cannot access a hardware feature requested
	ErrorCodeSecurity

	// The operation failed for an unspecified reason.
	ErrorCodeUnknown
)

var _ErrorCodeStrings = [9]string{
	"invalid-argument",
	"invalid-encoding",
	"timeout",
	"runtime-error",
	"unsupported-operation",
	"too-large",
	"not-found",
	"security",
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

// Error represents the imported resource "wasi:nn/errors@0.2.0-rc-2024-10-28#error".
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
// Return the error code.
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
// Errors can propagated with backend specific status through a string value.
//
//	data: func() -> string
//
//go:nosplit
func (self Error) Data() (result string) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ErrorData((uint32)(self0), &result)
	return
}
