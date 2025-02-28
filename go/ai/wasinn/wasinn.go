package wasinn

import (
	"errors"

	graphstream "github.com/hayride-dev/bindings/go/ai/gen/imports/hayride/ai/graph-stream"
	inferencestream "github.com/hayride-dev/bindings/go/ai/gen/imports/hayride/ai/inference-stream"
	tensorstream "github.com/hayride-dev/bindings/go/ai/gen/imports/hayride/ai/tensor-stream"
	"github.com/hayride-dev/bindings/go/ai/gen/imports/wasi/nn/graph"
	"github.com/hayride-dev/bindings/go/ai/gen/imports/wasi/nn/inference"
	"github.com/hayride-dev/bindings/go/ai/gen/imports/wasi/nn/tensor"
	"go.bytecodealliance.org/cm"
)

type Wasinn struct {
	graphExecCtx       *inference.GraphExecutionContext
	graphExecCtxStream *inferencestream.GraphExecutionContextStream
	options            *ModelOptions
}

func New(options ...Option[*ModelOptions]) (*Wasinn, error) {
	opts := defaultModelOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	var execCtx *inference.GraphExecutionContext
	var execCtxStream *inferencestream.GraphExecutionContextStream
	if opts.streaming {
		graphResult := graphstream.LoadByName(opts.name)
		if graphResult.IsErr() {
			return nil, &wasinnErr{graphResult.Err()}
		}
		graph := graphResult.OK()

		execCtxStreamResult := graph.InitExecutionContextStream()
		if execCtxStreamResult.IsErr() {
			return nil, &wasinnErr{execCtxStreamResult.Err()}
		}
		execCtxStream = execCtxStreamResult.OK()
	} else {
		graphResult := graph.LoadByName(opts.name)
		if graphResult.IsErr() {
			return nil, &wasinnErr{graphResult.Err()}
		}
		graph := graphResult.OK()

		execCtxResult := graph.InitExecutionContext()
		if execCtxResult.IsErr() {
			return nil, &wasinnErr{execCtxResult.Err()}
		}
		execCtx = execCtxResult.OK()
	}

	return &Wasinn{
		graphExecCtx:       execCtx,
		graphExecCtxStream: execCtxStream,
		options:            opts,
	}, nil
}

func (w *Wasinn) Compute(namedTensors map[string]string) (string, error) {
	inputs := []inference.NamedTensor{}
	for name, text := range namedTensors {
		d := tensor.TensorDimensions(cm.ToList([]uint32{1}))
		td := tensor.TensorData(cm.ToList([]uint8(text)))
		t := tensor.NewTensor(d, tensor.TensorTypeU8, td)

		namedTensor := inference.NamedTensor{
			F0: name,
			F1: t,
		}

		inputs = append(inputs, namedTensor)
	}
	inputList := cm.ToList(inputs)

	computeResult := w.graphExecCtx.Compute(inputList)
	if computeResult.IsErr() {
		return "", &wasinnErr{computeResult.Err()}
	}

	// Expecting only a single tensor as output
	// If multiple input is supported we could check the name of the named tensor for output matching
	outputs := computeResult.OK()
	if outputs.Len() != 1 {
		return "", errors.New("unexpected number of outputs")
	}

	return string(outputs.Data().F1.Data().Slice()), nil
}

func (w *Wasinn) ComputeStreaming(namedTensors map[string]string) (*TensorStream, error) {
	inputs := []inferencestream.NamedTensor{}
	for name, text := range namedTensors {
		d := tensor.TensorDimensions(cm.ToList([]uint32{1}))
		td := tensor.TensorData(cm.ToList([]uint8(text)))
		t := tensor.NewTensor(d, tensor.TensorTypeU8, td)

		namedTensor := inferencestream.NamedTensor{
			F0: name,
			F1: t,
		}

		inputs = append(inputs, namedTensor)
	}
	inputList := cm.ToList(inputs)

	computeResult := w.graphExecCtxStream.Compute(inputList)
	if computeResult.IsErr() {
		return nil, &wasinnErr{computeResult.Err()}
	}

	// Expecting only a single tensor as output
	// If multiple input is supported we could check the name of the named tensor for output matching
	tensorStream := &TensorStream{
		stream: computeResult.OK().F1,
	}

	return tensorStream, nil
}

type TensorStream struct {
	stream tensorstream.TensorStream
}

// Read will read the next `len` bytes from the stream
// will return empty byte slice if the stream is closed.
// blocks until the data is available
func (t TensorStream) Read(p []byte) (int, error) {
	t.stream.Subscribe().Block()
	data := t.stream.Read(uint64(len(p)))
	if data.IsErr() {
		if data.Err().Closed() {
			return 0, nil
		}

		return 0, &streamErr{data.Err()}
	}
	n := copy(p, data.OK().Slice())
	p = p[:n]
	return len(p), nil
}
