package models

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"

	"go.bytecodealliance.org/cm"
)

type Formatter interface {
	Encode(...model.Message) ([]byte, error)
	Decode([]byte) (model.Message, error)
}

type formatResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]Formatter
}

func (f *formatResourceTable) constructor() model.Format {
	f.rep++
	f.resources[f.rep] = formatter
	return model.FormatResourceNew(f.rep)
}

func (f *formatResourceTable) destructor(self cm.Rep) {
	delete(f.resources, self)
}

func (f *formatResourceTable) encode(self cm.Rep, messages cm.List[model.Message]) cm.Result[cm.List[uint8], cm.List[uint8], model.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextEncode)))
	}

	data, err := f.resources[self].Encode(messages.Slice()...)
	if err != nil {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextEncode)))
	}
	return cm.OK[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](cm.ToList(data))
}

func (f *formatResourceTable) decode(self cm.Rep, raw cm.List[uint8]) cm.Result[model.MessageShape, model.Message, model.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextDecode)))
	}

	message, err := f.resources[self].Decode(raw.Slice())
	if err != nil {
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextDecode)))
	}

	return cm.OK[cm.Result[model.MessageShape, model.Message, model.Error]](message)
}
