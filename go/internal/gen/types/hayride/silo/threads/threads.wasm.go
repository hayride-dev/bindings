// Code generated by wit-bindgen-go. DO NOT EDIT.

package threads

import (
	"go.bytecodealliance.org/cm"
)

// This file contains wasmimport and wasmexport declarations for "hayride:silo@0.0.61".

//go:wasmimport hayride:silo/threads@0.0.61 [resource-drop]thread
//go:noescape
func wasmimport_ThreadResourceDrop(self0 uint32)

//go:wasmimport hayride:silo/threads@0.0.61 [method]thread.id
//go:noescape
func wasmimport_ThreadID(self0 uint32, result *cm.Result[string, string, ErrNo])

//go:wasmimport hayride:silo/threads@0.0.61 [method]thread.wait
//go:noescape
func wasmimport_ThreadWait(self0 uint32, result *cm.Result[cm.List[uint8], cm.List[uint8], ErrNo])

//go:wasmimport hayride:silo/threads@0.0.61 spawn
//go:noescape
func wasmimport_Spawn(pkg0 *uint8, pkg1 uint32, function0 *uint8, function1 uint32, args0 *string, args1 uint32, result *cm.Result[Thread, Thread, ErrNo])

//go:wasmimport hayride:silo/threads@0.0.61 status
//go:noescape
func wasmimport_Status(id0 *uint8, id1 uint32, result *cm.Result[ThreadMetadataShape, ThreadMetadata, ErrNo])

//go:wasmimport hayride:silo/threads@0.0.61 kill
//go:noescape
func wasmimport_Kill(id0 *uint8, id1 uint32, result *cm.Result[ErrNo, struct{}, ErrNo])

//go:wasmimport hayride:silo/threads@0.0.61 group
//go:noescape
func wasmimport_Group(result *cm.Result[cm.List[ThreadMetadata], cm.List[ThreadMetadata], ErrNo])
