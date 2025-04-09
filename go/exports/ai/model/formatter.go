package model

import (
	"unsafe"

	witModel "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"

	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type Formatter interface {
	Encode(...*ai.Message) ([]byte, error)
	Decode([]byte) (*ai.Message, error)
}

type wacFormat struct {
	f Formatter
}

var format *wacFormat

func init() {
	format = &wacFormat{}
	witModel.Exports.Format.Constructor = format.wacConstructorfunc
}

func NewFormat(formatter Formatter) {
	format.f = formatter
}

func (f *wacFormat) wacConstructorfunc() witModel.Format {
	return witModel.FormatResourceNew(cm.Rep(uintptr(unsafe.Pointer(f))))
}

func (f *wacFormat) wacEncode(self cm.Rep, messages cm.List[witModel.Message]) (result cm.Result[witModel.Error, []byte, witModel.Error]) {
	return cm.OK[cm.Result[witModel.Error, []byte, witModel.Error]](nil)
}

func (f *wacFormat) wacDecode(self cm.Rep, data []byte) (result cm.Result[witModel.Error, witModel.Message, witModel.Error]) {
	msg := witTypes.Message{}
	return cm.OK[cm.Result[witModel.Error, witModel.Message, witModel.Error]](msg)
}
