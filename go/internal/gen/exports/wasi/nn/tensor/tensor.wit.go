// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package tensor represents the imported interface "wasi:nn/tensor@0.2.0-rc-2024-10-28".
//
// All inputs and outputs to an ML inference are represented as `tensor`s.
package tensor

import (
	"go.bytecodealliance.org/cm"
)

// TensorDimensions represents the list "wasi:nn/tensor@0.2.0-rc-2024-10-28#tensor-dimensions".
//
// The dimensions of a tensor.
//
// The array length matches the tensor rank and each element in the array describes
// the size of
// each dimension
//
//	type tensor-dimensions = list<u32>
type TensorDimensions cm.List[uint32]

// TensorType represents the enum "wasi:nn/tensor@0.2.0-rc-2024-10-28#tensor-type".
//
// The type of the elements in a tensor.
//
//	enum tensor-type {
//		FP16,
//		FP32,
//		FP64,
//		BF16,
//		U8,
//		I32,
//		I64
//	}
type TensorType uint8

const (
	TensorTypeFP16 TensorType = iota
	TensorTypeFP32
	TensorTypeFP64
	TensorTypeBF16
	TensorTypeU8
	TensorTypeI32
	TensorTypeI64
)

var _TensorTypeStrings = [7]string{
	"FP16",
	"FP32",
	"FP64",
	"BF16",
	"U8",
	"I32",
	"I64",
}

// String implements [fmt.Stringer], returning the enum case name of e.
func (e TensorType) String() string {
	return _TensorTypeStrings[e]
}

// MarshalText implements [encoding.TextMarshaler].
func (e TensorType) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], unmarshaling into an enum
// case. Returns an error if the supplied text is not one of the enum cases.
func (e *TensorType) UnmarshalText(text []byte) error {
	return _TensorTypeUnmarshalCase(e, text)
}

var _TensorTypeUnmarshalCase = cm.CaseUnmarshaler[TensorType](_TensorTypeStrings[:])

// TensorData represents the list "wasi:nn/tensor@0.2.0-rc-2024-10-28#tensor-data".
//
// The tensor data.
//
// Initially conceived as a sparse representation, each empty cell would be filled
// with zeros
// and the array length must match the product of all of the dimensions and the number
// of bytes
// in the type (e.g., a 2x2 tensor with 4-byte f32 elements would have a data array
// of length
// 16). Naturally, this representation requires some knowledge of how to lay out data
// in
// memory--e.g., using row-major ordering--and could perhaps be improved.
//
//	type tensor-data = list<u8>
type TensorData cm.List[uint8]

// Tensor represents the imported resource "wasi:nn/tensor@0.2.0-rc-2024-10-28#tensor".
//
//	resource tensor
type Tensor cm.Resource

// ResourceDrop represents the imported resource-drop for resource "tensor".
//
// Drops a resource handle.
//
//go:nosplit
func (self Tensor) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_TensorResourceDrop((uint32)(self0))
	return
}

// NewTensor represents the imported constructor for resource "tensor".
//
//	constructor(dimensions: tensor-dimensions, ty: tensor-type, data: tensor-data)
//
//go:nosplit
func NewTensor(dimensions TensorDimensions, ty TensorType, data TensorData) (result Tensor) {
	dimensions0, dimensions1 := cm.LowerList(dimensions)
	ty0 := (uint32)(ty)
	data0, data1 := cm.LowerList(data)
	result0 := wasmimport_NewTensor((*uint32)(dimensions0), (uint32)(dimensions1), (uint32)(ty0), (*uint8)(data0), (uint32)(data1))
	result = cm.Reinterpret[Tensor]((uint32)(result0))
	return
}

// Data represents the imported method "data".
//
// Return the tensor data.
//
//	data: func() -> tensor-data
//
//go:nosplit
func (self Tensor) Data() (result TensorData) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_TensorData((uint32)(self0), &result)
	return
}

// Dimensions represents the imported method "dimensions".
//
// Describe the size of the tensor (e.g., 2x2x2x2 -> [2, 2, 2, 2]). To represent a
// tensor
// containing a single value, use `[1]` for the tensor dimensions.
//
//	dimensions: func() -> tensor-dimensions
//
//go:nosplit
func (self Tensor) Dimensions() (result TensorDimensions) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_TensorDimensions((uint32)(self0), &result)
	return
}

// Ty represents the imported method "ty".
//
// Describe the type of element in the tensor (e.g., `f32`).
//
//	ty: func() -> tensor-type
//
//go:nosplit
func (self Tensor) Ty() (result TensorType) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_TensorTy((uint32)(self0))
	result = (TensorType)((uint32)(result0))
	return
}
