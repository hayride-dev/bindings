package ctx

import (
	"github.com/hayride-dev/bindings/go/gen/domain/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

var ctx Context
var ctxResourceTableInstance = ctxResourceTable{rep: 0, resources: make(map[cm.Rep]Context)}

type Context interface {
	Push(messages ...types.Message) error
	Messages() ([]types.Message, error)
	Next() (types.Message, error)
}

type ctxResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]Context
}

func init() {
	context.Exports.Context.Constructor = ctxResourceTableInstance.constructor
	context.Exports.Context.Push = ctxResourceTableInstance.push
	context.Exports.Context.Messages = ctxResourceTableInstance.messages
	context.Exports.Context.Next = ctxResourceTableInstance.next
	context.Exports.Context.Destructor = ctxResourceTableInstance.destructor

}

func Export(c Context) error {
	ctx = c
	return nil
}

func (c *ctxResourceTable) constructor() context.Context {
	c.rep++
	c.resources[c.rep] = ctx
	return context.ContextResourceNew(c.rep)
}

func (c *ctxResourceTable) destructor(self cm.Rep) {
	delete(c.resources, self)
}

func (c *ctxResourceTable) push(self cm.Rep, messages cm.List[context.Message]) (result cm.Result[context.Error, struct{}, context.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[context.Error, struct{}, context.Error]](context.ErrorResourceNew(cm.Rep(context.ErrorCodePushError)))
	}

	// Convert context.Message to types.Message
	msgs := make([]types.Message, len(messages.Slice()))
	for i, msg := range messages.Slice() {
		msgs[i] = cm.Reinterpret[types.Message](msg)
	}

	err := c.resources[self].Push(msgs...)
	if err != nil {
		return cm.Err[cm.Result[context.Error, struct{}, context.Error]](context.ErrorResourceNew(cm.Rep(context.ErrorCodePushError)))
	}
	return cm.OK[cm.Result[context.Error, struct{}, context.Error]](struct{}{})
}

func (c *ctxResourceTable) messages(self cm.Rep) (result cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]](context.ErrorResourceNew(cm.Rep(context.ErrorCodeMessageNotFound)))
	}
	msgs, err := c.resources[self].Messages()
	if err != nil {
		return cm.Err[cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]](context.ErrorResourceNew(cm.Rep(context.ErrorCodeMessageNotFound)))
	}

	// convert to context.Messages
	vals := make([]context.Message, len(msgs))
	for i, msg := range msgs {
		vals[i] = cm.Reinterpret[context.Message](msg)
	}

	return cm.OK[cm.Result[cm.List[context.Message], cm.List[context.Message], context.Error]](cm.ToList(vals))
}

func (c *ctxResourceTable) next(self cm.Rep) cm.Result[context.MessageShape, context.Message, context.Error] {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[context.MessageShape, context.Message, context.Error]](context.ErrorResourceNew(cm.Rep(context.ErrorCodeMessageNotFound)))
	}
	msg, err := c.resources[self].Next()
	if err != nil {
		return cm.Err[cm.Result[context.MessageShape, context.Message, context.Error]](context.ErrorResourceNew(cm.Rep(context.ErrorCodeMessageNotFound)))
	}

	val := cm.Reinterpret[context.Message](msg)

	return cm.OK[cm.Result[context.MessageShape, context.Message, context.Error]](context.Message(val))
}
