// Code generated by wit-bindgen-go. DO NOT EDIT.

package runner

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

// This file contains wasmimport and wasmexport declarations for "hayride:ai@0.0.61".

//go:wasmimport [export]hayride:ai/runner@0.0.61 [resource-new]error
//go:noescape
func wasmimport_ErrorResourceNew(rep0 uint32) (result0 uint32)

//go:wasmimport [export]hayride:ai/runner@0.0.61 [resource-rep]error
//go:noescape
func wasmimport_ErrorResourceRep(self0 uint32) (result0 uint32)

//go:wasmimport [export]hayride:ai/runner@0.0.61 [resource-drop]error
//go:noescape
func wasmimport_ErrorResourceDrop(self0 uint32)

//go:wasmexport hayride:ai/runner@0.0.61#[dtor]error
//export hayride:ai/runner@0.0.61#[dtor]error
func wasmexport_ErrorDestructor(self0 uint32) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	Exports.Error.Destructor(self)
	return
}

//go:wasmexport hayride:ai/runner@0.0.61#[method]error.code
//export hayride:ai/runner@0.0.61#[method]error.code
func wasmexport_ErrorCode(self0 uint32) (result0 uint32) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result := Exports.Error.Code(self)
	result0 = (uint32)(result)
	return
}

//go:wasmexport hayride:ai/runner@0.0.61#[method]error.data
//export hayride:ai/runner@0.0.61#[method]error.data
func wasmexport_ErrorData(self0 uint32) (result *string) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result_ := Exports.Error.Data(self)
	result = &result_
	return
}

//go:wasmexport hayride:ai/runner@0.0.61#invoke
//export hayride:ai/runner@0.0.61#invoke
func wasmexport_Invoke(message0 uint32, message1 *types.MessageContent, message2 uint32, agent0 uint32) (result *cm.Result[cm.List[Message], cm.List[Message], Error]) {
	message := lift_Message((uint32)(message0), (*types.MessageContent)(message1), (uint32)(message2))
	agent := cm.Reinterpret[cm.Rep]((uint32)(agent0))
	result_ := Exports.Invoke(message, agent)
	result = &result_
	return
}

//go:wasmexport hayride:ai/runner@0.0.61#invoke-stream
//export hayride:ai/runner@0.0.61#invoke-stream
func wasmexport_InvokeStream(message0 uint32, message1 *types.MessageContent, message2 uint32, writer0 uint32, agent0 uint32) (result *cm.Result[Error, struct{}, Error]) {
	message := lift_Message((uint32)(message0), (*types.MessageContent)(message1), (uint32)(message2))
	writer := cm.Reinterpret[cm.Rep]((uint32)(writer0))
	agent := cm.Reinterpret[cm.Rep]((uint32)(agent0))
	result_ := Exports.InvokeStream(message, writer, agent)
	result = &result_
	return
}
