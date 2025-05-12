package models

import (
	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"

	"go.bytecodealliance.org/cm"
)

type Formatter interface {
	Encode(...types.Message) ([]byte, error)
	Decode([]byte) (types.Message, error)
}

type formatResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]Formatter
}

func (f *formatResourceTable) constructor() model.Format {
	value := model.FormatResourceNew(f.rep)
	f.resources[f.rep] = formatter
	f.rep++
	return value
}

func (f *formatResourceTable) destructor(self cm.Rep) {
	delete(f.resources, self)
}

func (f *formatResourceTable) encode(self cm.Rep, messages cm.List[model.Message]) cm.Result[cm.List[uint8], cm.List[uint8], model.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextEncode)))
	}

	messagesSlice := make([]types.Message, len(messages.Slice()))
	for i, msg := range messages.Slice() {
		messagesSlice[i] = cm.Reinterpret[types.Message](msg)
	}

	data, err := f.resources[self].Encode(messagesSlice...)
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

	result := cm.Reinterpret[model.Message](message)

	return cm.OK[cm.Result[model.MessageShape, model.Message, model.Error]](result)
}
