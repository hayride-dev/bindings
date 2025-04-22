package ctx

import (
	witContext "github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

var ctx Context
var ctxResourceTableInstance = ctxResourceTable{rep: 0, resources: make(map[cm.Rep]Context)}

type Context interface {
	Push(messages ...witContext.Message) error
	Messages() ([]witContext.Message, error)
	Next() (witContext.Message, error)
}

type ctxResourceTable struct {
	rep       cm.Rep
	resources map[cm.Rep]Context
}

func init() {
	witContext.Exports.Context.Constructor = ctxResourceTableInstance.constructor
	witContext.Exports.Context.Push = ctxResourceTableInstance.push
	witContext.Exports.Context.Messages = ctxResourceTableInstance.messages
	witContext.Exports.Context.Next = ctxResourceTableInstance.next
	witContext.Exports.Context.Destructor = ctxResourceTableInstance.destructor

}

func Export(c Context) error {
	ctx = c
	return nil
}

func (c *ctxResourceTable) constructor() witContext.Context {
	c.rep++
	c.resources[c.rep] = ctx
	return witContext.ContextResourceNew(c.rep)
}

func (c *ctxResourceTable) destructor(self cm.Rep) {
	delete(c.resources, self)
}

func (c *ctxResourceTable) push(self cm.Rep, messages cm.List[witContext.Message]) (result cm.Result[witContext.Error, struct{}, witContext.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[witContext.Error, struct{}, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodePushError)))
	}

	err := c.resources[self].Push(messages.Slice()...)
	if err != nil {
		return cm.Err[cm.Result[witContext.Error, struct{}, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodePushError)))
	}
	return cm.OK[cm.Result[witContext.Error, struct{}, witContext.Error]](struct{}{})
}

func (c *ctxResourceTable) messages(self cm.Rep) (result cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
	}
	msgs, err := c.resources[self].Messages()
	if err != nil {
		return cm.Err[cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
	}

	return cm.OK[cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]](cm.ToList(msgs))
}

func (c *ctxResourceTable) next(self cm.Rep) (result cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
	}
	msg, err := c.resources[self].Next()
	if err != nil {
		return cm.Err[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
	}

	return cm.OK[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](msg)
}
