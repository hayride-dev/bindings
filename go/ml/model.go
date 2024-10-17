package ml

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/ml/gen/wasi/nn/graph"
	"github.com/hayride-dev/bindings/go/ml/gen/wasi/nn/inference"
	"github.com/hayride-dev/bindings/go/ml/gen/wasi/nn/tensor"
)

type Model struct {
	graphExecCtx *inference.GraphExecutionContext
	inputTensor  *tensor.Tensor
	options      *ModelOptions
}

func New(options ...Option[*ModelOptions]) (*Model, error) {
	opts := defaultModelOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	graphResult := graph.LoadByName(opts.name)
	if graphResult.IsErr() {
		e := &mlErr{graphResult.Err()}
		return nil, e
	}
	graph := graphResult.OK()

	execCtxResult := graph.InitExecutionContext()
	if execCtxResult.IsErr() {
		return nil, &mlErr{execCtxResult.Err()}
	}
	execCtx := execCtxResult.OK()

	return &Model{graphExecCtx: execCtx, options: opts}, nil
}

func (m *Model) Input(name string, text string) error {
	d := tensor.TensorDimensions(cm.ToList([]uint32{1}))
	td := tensor.TensorData(cm.ToList([]uint8(text)))
	t := tensor.NewTensor(d, tensor.TensorTypeU8, td)
	// TODO :: validate name ?
	inputResult := m.graphExecCtx.SetInput(name, t)
	if inputResult.IsErr() {
		return &mlErr{inputResult.Err()}
	}
	m.inputTensor = &t
	return nil
}

func (m *Model) Output(name string) (string, error) {
	// TODO :: validate name ?
	outputResult := m.graphExecCtx.GetOutput(name)
	if outputResult.IsErr() {
		return "", &mlErr{outputResult.Err()}
	}
	tensor := outputResult.OK()
	return string(tensor.Data().Slice()), nil
}

func (m *Model) Compute() error {
	computeResult := m.graphExecCtx.Compute()
	if computeResult.IsErr() {
		return &mlErr{computeResult.Err()}
	}
	return nil
}
