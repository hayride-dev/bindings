package ctx

import (
	"unsafe"

	witContext "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/context"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

var impl *wacContext

func init() {
	impl = &wacContext{}
	witContext.Exports.Context.Constructor = impl.wacConstructorfunc
	witContext.Exports.Context.Push = impl.wacPush
	witContext.Exports.Context.Messages = impl.wacMessages
	witContext.Exports.Context.Next = impl.wacNext

}

type Context interface {
	Push(messages ...*ai.Message) error
	Messages() ([]*ai.Message, error)
	Next() (*ai.Message, error)
}

type wacContext struct {
	ctx Context
}

func New(options ...Option[*CtxOptions]) error {
	opts := defaultModelOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return err
		}
	}
	impl.ctx = opts.ctx
	return nil
}

func (c *wacContext) wacConstructorfunc() witContext.Context {
	return witContext.ContextResourceNew(cm.Rep(uintptr(unsafe.Pointer(c))))
}

func (c *wacContext) wacPush(self cm.Rep, messages cm.List[witContext.Message]) (result cm.Result[witContext.Error, struct{}, witContext.Error]) {

	c.ctx.Push(nil)
	// Implementation of wacPush
	return cm.OK[cm.Result[witContext.Error, struct{}, witContext.Error]](struct{}{})
}

func (c *wacContext) wacMessages(self cm.Rep) (result cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]) {
	c.ctx.Messages()
	return cm.OK[cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]](cm.List[witContext.Message]{})
}

func (c *wacContext) wacNext(self cm.Rep) (result cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]) {
	c.ctx.Next()
	return cm.OK[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witContext.Message{})
}
