package model

import (
	"io"

	"github.com/hayride-dev/bindings/go/imports/ai/nn"
	inferencestream "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/inference-stream"
	witModel "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/wasi/nn/tensor"

	"go.bytecodealliance.org/cm"
)

var formatter Formatter
var formatResourceTableInstance = formatResourceTable{rep: 0, resources: make(map[cm.Rep]Formatter)}
var modelResourceTableInstance = modelResourceTable{rep: 0, resources: make(map[cm.Rep]*modelResource)}

func init() {
	// TODO:: error exports
	// do we need a resource table?

	// format exports
	witModel.Exports.Format.Constructor = formatResourceTableInstance.constructor
	witModel.Exports.Format.Encode = formatResourceTableInstance.encode
	witModel.Exports.Format.Decode = formatResourceTableInstance.decode
	witModel.Exports.Format.Destructor = formatResourceTableInstance.destructor
	// model exports
	witModel.Exports.Model.ExportConstructor = modelResourceTableInstance.exportConstructor
	witModel.Exports.Model.Compute = modelResourceTableInstance.compute
}

type modelResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]*modelResource
}

type modelResource struct {
	format witModel.Format
	graph  witModel.GraphExecutionContextStream
}

func Export(f Formatter) {
	formatter = f
}

func (w *modelResourceTable) exportConstructor(f witModel.Format, graph witModel.GraphExecutionContextStream) witModel.Model {
	resource := &modelResource{f, graph}
	w.rep++
	w.resources[w.rep] = resource
	return witModel.ModelResourceNew(w.rep)
}

func (w *modelResourceTable) compute(self cm.Rep, messages cm.List[witModel.Message]) cm.Result[witModel.MessageShape, witModel.Message, witModel.Error] {
	if _, ok := w.resources[self]; !ok {
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeComputeError)))
	}
	resource := w.resources[self]
	result := formatResourceTableInstance.encode(cm.Rep(resource.format), messages)
	if result.IsErr() {
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextEncode)))
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
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeComputeError)))
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
			wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeComputeError))
			return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
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
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextDecode)))
	}
	return cm.OK[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](*decodeResult.OK())
}
