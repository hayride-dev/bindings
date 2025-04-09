package model

import (
	witModel "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"

	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type Formatter interface {
	Encode(...*ai.Message) ([]byte, error)
	Decode([]byte) (*ai.Message, error)
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

func (f *formatResourceTable) encode(self cm.Rep, messages cm.List[witModel.Message]) cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextEncode)))
	}
	f.resources[self].Encode(nil)
	return cm.OK[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](cm.List[uint8]{})
}

func (f *formatResourceTable) decode(self cm.Rep, raw cm.List[uint8]) cm.Result[witModel.MessageShape, witModel.Message, witModel.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextDecode)))
	}
	f.resources[self].Decode(nil)
	msg := witTypes.Message{}
	return cm.OK[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](msg)
}

func (f *formatResourceTable) destructor(self cm.Rep) {
	delete(f.resources, self)
}
