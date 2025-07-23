package graph

import (
	"fmt"

	graphstream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	inferencestream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/inference-stream"
	"go.bytecodealliance.org/cm"
)

var _ InferenceStream = (*GraphStream)(nil)

type InferenceStream interface {
	InitExecutionContextStream() (GraphExecutionContextStream, error)
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
	Compute(namedTensors []inferencestream.NamedTensor) ([]inferencestream.NamedTensor, error)
}

type GraphExecCtxStream cm.Resource

func (g GraphExecCtxStream) Compute(namedTensors []inferencestream.NamedTensor) ([]inferencestream.NamedTensor, error) {
	witGraphExecCtxStream := cm.Reinterpret[graphstream.GraphExecutionContextStream](g)

	witList := cm.ToList(namedTensors)
	result := witGraphExecCtxStream.Compute(cm.Reinterpret[cm.List[inferencestream.NamedTensor]](witList))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compute: %s", result.Err().Data())
	}

	namedTensorsResult := cm.Reinterpret[cm.List[inferencestream.NamedTensor]](result.OK())
	return namedTensorsResult.Slice(), nil
}
