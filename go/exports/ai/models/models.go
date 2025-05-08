package models

import (
	"io"

	"github.com/hayride-dev/bindings/go/imports/ai/nn"
	inferencestream "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/inference-stream"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/wasi/nn/tensor"
	"go.bytecodealliance.org/cm"
)

var formatResourceTableInstance = formatResourceTable{rep: 0, resources: make(map[cm.Rep]Formatter)}
var modelResourceTableInstance = modelResourceTable{rep: 0, resources: make(map[cm.Rep]*modelResource)}

func init() {
	// format exports
	model.Exports.Format.Constructor = formatResourceTableInstance.constructor
	model.Exports.Format.Encode = formatResourceTableInstance.encode
	model.Exports.Format.Decode = formatResourceTableInstance.decode
	model.Exports.Format.Destructor = formatResourceTableInstance.destructor
	// model exports
	model.Exports.Model.ExportConstructor = modelResourceTableInstance.exportConstructor
	model.Exports.Model.Compute = modelResourceTableInstance.compute
	model.Exports.Model.Destructor = modelResourceTableInstance.destructor
}

type modelResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]*modelResource
}

type modelResource struct {
	format model.Format
	graph  model.GraphExecutionContextStream
}

func Export(f Formatter) {
	formatResourceTableInstance.resources[formatResourceTableInstance.rep] = f
}

func (w *modelResourceTable) exportConstructor(f model.Format, graph model.GraphExecutionContextStream) model.Model {
	resource := &modelResource{f, graph}
	w.rep++
	w.resources[w.rep] = resource
	return model.ModelResourceNew(w.rep)
}

func (w *modelResourceTable) destructor(self cm.Rep) {
	delete(w.resources, self)
}

func (w *modelResourceTable) compute(self cm.Rep, messages cm.List[model.Message]) cm.Result[model.MessageShape, model.Message, model.Error] {
	if _, ok := w.resources[self]; !ok {
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeComputeError)))
	}
	resource := w.resources[self]
	result := formatResourceTableInstance.encode(cm.Rep(resource.format), messages)
	if result.IsErr() {
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextEncode)))
	}
	d := tensor.TensorDimensions(cm.ToList([]uint32{1}))
	td := tensor.TensorData(cm.ToList(result.OK().Slice()))
	t := tensor.NewTensor(d, tensor.TensorTypeU8, td)
	inputs := []inferencestream.NamedTensor{
		{
			F0: "user",
			F1: t,
		},
	}
	graphResult := resource.graph.Compute(cm.ToList(inputs))
	if graphResult.IsErr() {
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeComputeError)))
	}

	stream := graphResult.OK().F1
	ts := nn.TensorStream(stream)
	// read the output from the stream
	text := make([]byte, 0)
	for {
		// Read up to 100 bytes from the output
		// to get any tokens that have been generated and push to socket
		p := make([]byte, 100)
		len, err := ts.Read(p)
		if len == 0 || err == io.EOF {
			break
		} else if err != nil {
			wasiErr := model.ErrorResourceNew(cm.Rep(model.ErrorCodeComputeError))
			return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](wasiErr)
		}
		/*
			// Removed:: example exposing raw tensor stream
			// this is commented out until there is a good way to return an AI message stream
			if _, err := writer.Write(p[:len]); err != nil {
				wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeComputeError))
				return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
			}
		*/
		text = append(text, p[:len]...)
	}

	decodeResult := formatResourceTableInstance.decode(cm.Rep(resource.format), cm.ToList(text))
	if decodeResult.IsErr() {
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextDecode)))
	}
	return cm.OK[cm.Result[model.MessageShape, model.Message, model.Error]](*decodeResult.OK())
}
