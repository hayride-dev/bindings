// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package config represents the exported interface "hayride:http/config@0.0.61".
package config

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/types"
	"go.bytecodealliance.org/cm"
)

// ServerConfig represents the type alias "hayride:http/config@0.0.61#server-config".
//
// See [types.ServerConfig] for more information.
type ServerConfig = types.ServerConfig

// ErrorCode represents the type alias "hayride:http/config@0.0.61#error-code".
//
// See [types.ErrorCode] for more information.
type ErrorCode = types.ErrorCode

// Error represents the exported resource "hayride:http/config@0.0.61#error".
//
//	resource error
type Error cm.Resource

// ErrorResourceNew represents the imported resource-new for resource "error".
//
// Creates a new resource handle.
//
//go:nosplit
func ErrorResourceNew(rep cm.Rep) (result Error) {
	rep0 := cm.Reinterpret[uint32](rep)
	result0 := wasmimport_ErrorResourceNew((uint32)(rep0))
	result = cm.Reinterpret[Error]((uint32)(result0))
	return
}

// ResourceRep represents the imported resource-rep for resource "error".
//
// Returns the underlying resource representation.
//
//go:nosplit
func (self Error) ResourceRep() (result cm.Rep) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_ErrorResourceRep((uint32)(self0))
	result = cm.Reinterpret[cm.Rep]((uint32)(result0))
	return
}

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

func init() {
	Exports.Error.Destructor = func(self cm.Rep) {}
}
