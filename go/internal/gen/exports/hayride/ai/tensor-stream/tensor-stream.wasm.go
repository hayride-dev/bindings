// Code generated by wit-bindgen-go. DO NOT EDIT.

package tensorstream

import (
	"go.bytecodealliance.org/cm"
)

// This file contains wasmimport and wasmexport declarations for "hayride:ai@0.0.61".

//go:wasmimport hayride:ai/tensor-stream@0.0.61 [resource-drop]tensor-stream
//go:noescape
func wasmimport_TensorStreamResourceDrop(self0 uint32)

//go:wasmimport hayride:ai/tensor-stream@0.0.61 [constructor]tensor-stream
//go:noescape
func wasmimport_NewTensorStream(dimensions0 *uint32, dimensions1 uint32, ty0 uint32, data0 *uint8, data1 uint32) (result0 uint32)

//go:wasmimport hayride:ai/tensor-stream@0.0.61 [method]tensor-stream.dimensions
//go:noescape
func wasmimport_TensorStreamDimensions(self0 uint32, result *TensorDimensions)

//go:wasmimport hayride:ai/tensor-stream@0.0.61 [method]tensor-stream.read
//go:noescape
func wasmimport_TensorStreamRead(self0 uint32, len0 uint64, result *cm.Result[TensorData, TensorData, StreamError])

//go:wasmimport hayride:ai/tensor-stream@0.0.61 [method]tensor-stream.subscribe
//go:noescape
func wasmimport_TensorStreamSubscribe(self0 uint32) (result0 uint32)

//go:wasmimport hayride:ai/tensor-stream@0.0.61 [method]tensor-stream.ty
//go:noescape
func wasmimport_TensorStreamTy(self0 uint32) (result0 uint32)
