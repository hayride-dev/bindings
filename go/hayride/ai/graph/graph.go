package graph

import (
	"fmt"

	graphstream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	inferencestream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/inference-stream"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/nn/tensor"
	"go.bytecodealliance.org/cm"
)

var _ InferenceStream = (*GraphStream)(nil)

type InferenceStream interface {
	InitExecutionContextStream() (GraphExecutionContextStream, error)
}

type NamedTensor = inferencestream.NamedTensor
type NamedTensorStream = inferencestream.NamedTensorStream

type TensorDimensions = tensor.TensorDimensions
type TensorType = tensor.TensorType

const TensorTypeFP16 = tensor.TensorTypeFP16
const TensorTypeFP32 = tensor.TensorTypeFP32
const TensorTypeFP64 = tensor.TensorTypeFP64
const TensorTypeBF16 = tensor.TensorTypeBF16
const TensorTypeU8 = tensor.TensorTypeU8
const TensorTypeI32 = tensor.TensorTypeI32
const TensorTypeI64 = tensor.TensorTypeI64

type TensorData = tensor.TensorData
type Tensor = tensor.Tensor

func NewTensor(dimensions TensorDimensions, ty TensorType, data TensorData) Tensor {
	return tensor.NewTensor(dimensions, ty, data)
}

func LoadByName(name string) (InferenceStream, error) {
	witGraphStream := graphstream.LoadByName(name)
	if witGraphStream.IsErr() {
		return nil, fmt.Errorf("failed to load graph stream: %s", witGraphStream.Err().Data())
	}

	return GraphStream(*witGraphStream.OK()), nil
}

type GraphStream cm.Resource

func (g GraphStream) InitExecutionContextStream() (GraphExecutionContextStream, error) {
	witGraphStream := cm.Reinterpret[graphstream.GraphStream](g)

	result := witGraphStream.InitExecutionContextStream()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to init execution context stream: %s", result.Err().Data())
	}

	return GraphExecCtxStream(*result.OK()), nil
}

var _ GraphExecutionContextStream = (*GraphExecCtxStream)(nil)

type GraphExecutionContextStream interface {
	Compute(namedTensors []NamedTensor) (inferencestream.NamedTensorStream, error)
}

type GraphExecCtxStream cm.Resource

func (g GraphExecCtxStream) Compute(namedTensors []NamedTensor) (NamedTensorStream, error) {
	witGraphExecCtxStream := cm.Reinterpret[graphstream.GraphExecutionContextStream](g)

	witList := cm.ToList(namedTensors)
	result := witGraphExecCtxStream.Compute(cm.Reinterpret[cm.List[inferencestream.NamedTensor]](witList))
	if result.IsErr() {
		return NamedTensorStream{}, fmt.Errorf("failed to compute: %s", result.Err().Data())
	}

	namedTensorStream := cm.Reinterpret[NamedTensorStream](result.OK())
	return namedTensorStream, nil
}

type TensorStream cm.Resource

// Read will read the next `len` bytes from the stream
// will return empty byte slice if the stream is closed.
// blocks until the data is available
func (t TensorStream) Read(p []byte) (int, error) {
	ts := cm.Reinterpret[inferencestream.TensorStream](t)
	ts.Subscribe().Block()
	data := ts.Read(uint64(len(p)))
	if data.IsErr() {
		if data.Err().Closed() {
			return 0, nil
		}
		return 0, fmt.Errorf("%s", data.Err().String())
	}
	n := copy(p, data.OK().Slice())
	p = p[:n]
	return len(p), nil
}
