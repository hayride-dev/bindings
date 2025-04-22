package model

import (
	witModel "github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/model"

	"go.bytecodealliance.org/cm"
)

type Formatter interface {
	Encode(...witModel.Message) ([]byte, error)
	Decode([]byte) (witModel.Message, error)
}

type formatResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]Formatter
}

func (f *formatResourceTable) constructor() witModel.Format {
	f.rep++
	f.resources[f.rep] = formatter
	return witModel.FormatResourceNew(f.rep)
}

func (f *formatResourceTable) destructor(self cm.Rep) {
	delete(f.resources, self)
}

func (f *formatResourceTable) encode(self cm.Rep, messages cm.List[witModel.Message]) cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextEncode)))
	}

	data, err := f.resources[self].Encode(messages.Slice()...)
	if err != nil {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextEncode)))
	}
	return cm.OK[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](cm.ToList(data))
}

func (f *formatResourceTable) decode(self cm.Rep, raw cm.List[uint8]) cm.Result[witModel.MessageShape, witModel.Message, witModel.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextDecode)))
	}

	message, err := f.resources[self].Decode(raw.Slice())
	if err != nil {
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextDecode)))
	}

	return cm.OK[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](message)
}
