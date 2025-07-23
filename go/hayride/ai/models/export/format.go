package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/ai/models"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

type Constructor func() models.Format

var formatConstructor Constructor

type resources struct {
	format map[cm.Rep]models.Format
}

var resourceTable = &resources{
	format: make(map[cm.Rep]models.Format),
}

func init() {
}

func Export(c Constructor) {
	formatConstructor = c

	model.Exports.Format.Constructor = constructor
	model.Exports.Format.Decode = decode
	model.Exports.Format.Encode = encode
	model.Exports.Format.Destructor = destructor
}

func constructor() model.Format {
	formatter := formatConstructor()

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&formatter))))
	v := model.FormatResourceNew(key)
	resourceTable.format[key] = formatter
	return v
}

func destructor(self cm.Rep) {
	delete(resourceTable.format, self)
}

func decode(self cm.Rep, raw cm.List[uint8]) (result cm.Result[model.MessageShape, model.Message, model.Error]) {
	m, ok := resourceTable.format[self]
	if !ok {
		wasiErr := model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextDecode))
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](wasiErr)
	}
	msg, err := m.Decode(raw.Slice())
	if err != nil {
		wasiErr := model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextDecode))
		return cm.Err[cm.Result[model.MessageShape, model.Message, model.Error]](wasiErr)
	}

	message := cm.Reinterpret[model.Message](msg)

	return cm.OK[cm.Result[model.MessageShape, model.Message, model.Error]](message)
}

func encode(self cm.Rep, messages cm.List[model.Message]) (result cm.Result[cm.List[uint8], cm.List[uint8], model.Error]) {
	m, ok := resourceTable.format[self]
	if !ok {
		wasiErr := model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextEncode))
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](wasiErr)
	}

	msgs := cm.Reinterpret[cm.List[types.Message]](messages)

	msg, err := m.Encode(msgs.Slice()...)
	if err != nil {
		wasiErr := model.ErrorResourceNew(cm.Rep(model.ErrorCodeContextEncode))
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[uint8], cm.List[uint8], model.Error]](cm.ToList(msg))
}
