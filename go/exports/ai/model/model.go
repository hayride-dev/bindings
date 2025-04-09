package model

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/hayride-dev/bindings/go/exports/ai/types"
	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	inferencestream "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/inference-stream"
	witModel "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	tensorstream "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/tensor-stream"

	witTypes "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/wasi/nn/tensor"
	"go.bytecodealliance.org/cm"
)

var impl *wacModel

func init() {
	impl = &wacModel{
		history: make([]*types.Message, 0),
	}
	witModel.Exports.Model.Constructor = impl.wacConstructorfunc
	witModel.Exports.Model.Push = impl.wacPush
	witModel.Exports.Model.Compute = impl.wacCompute
}

// TODO :: this likely should be moved - or figure out how to make it a io.reader that can
// be used as the wasi stream wit
type tensorStream struct {
	stream tensorstream.TensorStream
}

// Read will read the next `len` bytes from the stream
// will return empty byte slice if the stream is closed.
// blocks until the data is available
func (t tensorStream) Read(p []byte) (int, error) {
	t.stream.Subscribe().Block()
	data := t.stream.Read(uint64(len(p)))
	if data.IsErr() {
		if data.Err().Closed() {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to read from stream: %v", data.Err())
	}
	n := copy(p, data.OK().Slice())
	p = p[:n]
	return len(p), nil
}

func Model(m Formatter) error {
	impl.formatter = m
	return nil
}

type wacModel struct {
	formatter Formatter
	model     witModel.GraphExecutionContextStream
	history   []*types.Message
	ref       cm.Rep
}

func (w *wacModel) wacConstructorfunc(graph witModel.GraphExecutionContextStream) (result witModel.Model) {
	w.model = graph
	w.ref = cm.Rep(uintptr(unsafe.Pointer(&impl)))
	return witModel.ModelResourceNew(w.ref)
}

func (w *wacModel) wacCompute(self cm.Rep, output cm.Rep) (result cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]) {
	writer := wasiio.Clone(uint32(output))
	defer witModel.OutputStream(cm.Rep(uint32(output))).ResourceDrop()

	ctx, err := w.formatter.Encode(w.history...)
	if err != nil {
		wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextEncode))
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
	}

	d := tensor.TensorDimensions(cm.ToList([]uint32{1}))
	td := tensor.TensorData(cm.ToList(ctx))
	t := tensor.NewTensor(d, tensor.TensorTypeU8, td)

	// TODO :: add support for named configuration
	inputs := []inferencestream.NamedTensor{
		{
			F0: "user",
			F1: t,
		},
	}

	inputList := cm.ToList(inputs)
	computeResult := w.model.Compute(inputList)
	if computeResult.IsErr() {
		wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeComputeError))
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
	}

	stream := computeResult.OK().F1
	ts := &tensorStream{
		stream: stream,
	}

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

		if _, err := writer.Write(p[:len]); err != nil {
			wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeComputeError))
			return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
		}
		text = append(text, p[:len]...)
	}

	response, err := w.formatter.Decode([]byte(text))
	if err != nil {
		wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeMessageDecode))
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
	}

	content := make([]witTypes.Content, 0)
	for _, c := range response.Content {
		switch c.Type() {
		case "text":
			c := c.(*types.TextContent)
			content = append(content, witTypes.ContentText(witTypes.TextContent{
				Text:        c.Text,
				ContentType: c.ContentType,
			}))
		case "tool-input":
			c := c.(*types.ToolInput)
			content = append(content, witTypes.ContentToolInput(witTypes.ToolInput{
				ContentType: c.ContentType,
				ID:          c.ID,
				Name:        c.Name,
				Input:       c.Input,
			}))
		default:
			wasiErr := witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeUnexpectedMessageType))
			return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](wasiErr)
		}
	}

	resp := witModel.Message{
		Role:    witTypes.RoleAssistant,
		Content: cm.ToList(content),
	}

	// store model response in model context
	w.wacPush(w.ref, cm.ToList([]witModel.Message{resp}))

	return cm.OK[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](resp)
}
