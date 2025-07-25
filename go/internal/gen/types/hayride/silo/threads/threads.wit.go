// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package threads represents the imported interface "hayride:silo/threads@0.0.61".
package threads

import (
	"github.com/hayride-dev/bindings/go/internal/gen/types/hayride/silo/types"
	"go.bytecodealliance.org/cm"
)

// ErrNo represents the type alias "hayride:silo/threads@0.0.61#err-no".
//
// See [types.ErrNo] for more information.
type ErrNo = types.ErrNo

// ThreadMetadata represents the type alias "hayride:silo/threads@0.0.61#thread-metadata".
//
// See [types.ThreadMetadata] for more information.
type ThreadMetadata = types.ThreadMetadata

// ThreadStatus represents the type alias "hayride:silo/threads@0.0.61#thread-status".
//
// See [types.ThreadStatus] for more information.
type ThreadStatus = types.ThreadStatus

// Thread represents the imported resource "hayride:silo/threads@0.0.61#thread".
//
//	resource thread
type Thread cm.Resource

// ResourceDrop represents the imported resource-drop for resource "thread".
//
// Drops a resource handle.
//
//go:nosplit
func (self Thread) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ThreadResourceDrop((uint32)(self0))
	return
}

// ID represents the imported method "id".
//
//	id: func() -> result<string, err-no>
//
//go:nosplit
func (self Thread) ID() (result cm.Result[string, string, ErrNo]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ThreadID((uint32)(self0), &result)
	return
}

// Wait represents the imported method "wait".
//
//	wait: func() -> result<list<u8>, err-no>
//
//go:nosplit
func (self Thread) Wait() (result cm.Result[cm.List[uint8], cm.List[uint8], ErrNo]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_ThreadWait((uint32)(self0), &result)
	return
}

// Spawn represents the imported function "spawn".
//
//	spawn: func(pkg: string, function: string, args: list<string>) -> result<thread,
//	err-no>
//
//go:nosplit
func Spawn(pkg string, function string, args cm.List[string]) (result cm.Result[Thread, Thread, ErrNo]) {
	pkg0, pkg1 := cm.LowerString(pkg)
	function0, function1 := cm.LowerString(function)
	args0, args1 := cm.LowerList(args)
	wasmimport_Spawn((*uint8)(pkg0), (uint32)(pkg1), (*uint8)(function0), (uint32)(function1), (*string)(args0), (uint32)(args1), &result)
	return
}

// Status represents the imported function "status".
//
//	status: func(id: string) -> result<thread-metadata, err-no>
//
//go:nosplit
func Status(id string) (result cm.Result[ThreadMetadataShape, ThreadMetadata, ErrNo]) {
	id0, id1 := cm.LowerString(id)
	wasmimport_Status((*uint8)(id0), (uint32)(id1), &result)
	return
}

// Kill represents the imported function "kill".
//
// get metadata about a single thread
//
//	kill: func(id: string) -> result<_, err-no>
//
//go:nosplit
func Kill(id string) (result cm.Result[ErrNo, struct{}, ErrNo]) {
	id0, id1 := cm.LowerString(id)
	wasmimport_Kill((*uint8)(id0), (uint32)(id1), &result)
	return
}

// Group represents the imported function "group".
//
//	group: func() -> result<list<thread-metadata>, err-no>
//
//go:nosplit
func Group() (result cm.Result[cm.List[ThreadMetadata], cm.List[ThreadMetadata], ErrNo]) {
	wasmimport_Group(&result)
	return
}
