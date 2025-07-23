package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/ai/ctx"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

type Constructor func() ctx.Context

var ctxConstructor Constructor

type resources struct {
	ctx map[cm.Rep]ctx.Context
}

var resourceTable = &resources{
	ctx: make(map[cm.Rep]ctx.Context),
}

func init() {
}

func Export(c Constructor) {
	ctxConstructor = c

	context.Exports.Context.Constructor = constructor
	context.Exports.Context.Push = push
	context.Exports.Context.Messages = messages
	context.Exports.Context.Destructor = destructor
}

func constructor() context.Context {
	ctx := ctxConstructor()

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&ctx))))
	v := context.ContextResourceNew(key)
	resourceTable.ctx[key] = ctx
	return v
}

func destructor(self cm.Rep) {
	delete(resourceTable.ctx, self)
}

func push(self cm.Rep, msg context.Message) cm.Result[context.Error, struct{}, context.Error] {
	ctx, ok := resourceTable.ctx[self]
	if !ok {
		wasiErr := context.ErrorResourceNew(cm.Rep(context.ErrorCodePushError))
		return cm.Err[cm.Result[context.Error, struct{}, context.Error]](wasiErr)
	}

	m := cm.Reinterpret[types.Message](msg)

	if err := ctx.Push(m); err != nil {
		wasiErr := context.ErrorResourceNew(cm.Rep(context.ErrorCodePushError))
		return cm.Err[cm.Result[context.Error, struct{}, context.Error]](wasiErr)
	}
	return cm.Result[context.Error, struct{}, context.Error]{}
}

func messages(self cm.Rep) (result cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]) {
	ctx, ok := resourceTable.ctx[self]
	if !ok {
		wasiErr := context.ErrorResourceNew(cm.Rep(context.ErrorCodeMessageNotFound))
		return cm.Err[cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]](wasiErr)
	}

	messages, err := ctx.Messages()
	if err != nil {
		wasiErr := context.ErrorResourceNew(cm.Rep(context.ErrorCodeMessageNotFound))
		return cm.Err[cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]](wasiErr)
	}

	msgs := cm.Reinterpret[cm.List[context.Message]](messages)

	return cm.OK[cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]](msgs)
}
