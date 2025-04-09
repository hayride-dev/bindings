package ctx

import (
	witContext "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/context"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

var ctx Context
var ctxResourceTableInstance = ctxResourceTable{rep: 0, resources: make(map[cm.Rep]Context)}

type Context interface {
	Push(messages ...*ai.Message) error
	Messages() ([]*ai.Message, error)
	Next() (*ai.Message, error)
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

func (c *ctxResourceTable) push(self cm.Rep, messages cm.List[witContext.Message]) (result cm.Result[witContext.Error, struct{}, witContext.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[witContext.Error, struct{}, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodePushError)))
	}
	msgs := make([]*ai.Message, 0)
	for _, msg := range messages.Slice() {
		content := make([]ai.Content, 0)
		for _, c := range msg.Content.Slice() {
			switch c.String() {
			case "text":
				content = append(content, &ai.TextContent{
					ContentType: c.Text().ContentType,
					Text:        c.Text().Text})
			case "tool-schema":
				content = append(content, &ai.ToolSchema{
					ID:           c.ToolSchema().ID,
					Name:         c.ToolSchema().Name,
					Description:  c.ToolSchema().Description,
					ParamsSchema: c.ToolSchema().ParamsSchema,
				})
			case "tool-input":
				content = append(content, &ai.ToolInput{
					ID:          c.ToolInput().ID,
					Name:        c.ToolInput().Name,
					ContentType: c.ToolInput().ContentType,
					Input:       c.ToolInput().Input,
				})
			case "tool-output":
				content = append(content, &ai.ToolOutput{
					ContentType: c.ToolOutput().ContentType,
					ID:          c.ToolOutput().ID,
					Name:        c.ToolOutput().Name,
					Output:      c.ToolOutput().Output,
				})
			default:
				return cm.Err[cm.Result[witContext.Error, struct{}, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeUnknown)))
			}
		}
		msgs = append(msgs, &ai.Message{
			Role:    ai.Role(msg.Role),
			Content: content,
		})
	}
	err := c.resources[self].Push(msgs...)
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
	witMessages := make([]witTypes.Message, 0)
	for _, m := range msgs {
		content := make([]witTypes.Content, 0)
		for _, c := range m.Content {
			switch c.Type() {
			case "text":
				c := c.(*ai.TextContent)
				content = append(content, witTypes.ContentText(witTypes.TextContent{
					Text:        c.Text,
					ContentType: c.ContentType,
				}))
			case "tool-schema":
				c := c.(*ai.ToolSchema)
				content = append(content, witTypes.ContentToolSchema(witTypes.ToolSchema{
					ID:           c.ID,
					Name:         c.Name,
					Description:  c.Description,
					ParamsSchema: c.ParamsSchema,
				}))
			case "tool-input":
				c := c.(*ai.ToolInput)
				content = append(content, witTypes.ContentToolInput(witTypes.ToolInput{
					ID:          c.ID,
					Name:        c.Name,
					ContentType: c.ContentType,
					Input:       c.Input,
				}))
			case "tool-output":
				c := c.(*ai.ToolOutput)
				content = append(content, witTypes.ContentToolOutput(witTypes.ToolOutput{
					ContentType: c.ContentType,
					ID:          c.ID,
					Name:        c.Name,
					Output:      c.Output,
				}))
			default:
				return cm.Err[cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeUnknown)))
			}
		}
		witMessages = append(witMessages, witTypes.Message{
			Role:    witTypes.Role(m.Role),
			Content: cm.ToList(content),
		})
	}

	return cm.OK[cm.Result[cm.List[witContext.Message], cm.List[witContext.Message], witContext.Error]](cm.ToList(witMessages))
}

func (c *ctxResourceTable) next(self cm.Rep) (result cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]) {
	if _, ok := c.resources[self]; !ok {
		return cm.Err[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
	}
	msg, err := c.resources[self].Next()
	if err != nil {
		return cm.Err[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
	}
	content := make([]witTypes.Content, 0)
	for _, c := range msg.Content {
		switch c.Type() {
		case "text":
			c := c.(*ai.TextContent)
			content = append(content, witTypes.ContentText(witTypes.TextContent{
				Text:        c.Text,
				ContentType: c.ContentType,
			}))
		case "tool-schema":
			c := c.(*ai.ToolSchema)
			content = append(content, witTypes.ContentToolSchema(witTypes.ToolSchema{
				ID:           c.ID,
				Name:         c.Name,
				Description:  c.Description,
				ParamsSchema: c.ParamsSchema,
			}))
		case "tool-input":
			c := c.(*ai.ToolInput)
			content = append(content, witTypes.ContentToolInput(witTypes.ToolInput{
				ID:          c.ID,
				Name:        c.Name,
				ContentType: c.ContentType,
				Input:       c.Input,
			}))
		case "tool-output":
			c := c.(*ai.ToolOutput)
			content = append(content, witTypes.ContentToolOutput(witTypes.ToolOutput{
				ContentType: c.ContentType,
				ID:          c.ID,
				Name:        c.Name,
				Output:      c.Output,
			}))
		default:
			return cm.Err[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witContext.ErrorResourceNew(cm.Rep(witContext.ErrorCodeMessageNotFound)))
		}
	}
	witMessage := witTypes.Message{
		Role:    witTypes.Role(msg.Role),
		Content: cm.ToList(content),
	}
	return cm.OK[cm.Result[witContext.MessageShape, witContext.Message, witContext.Error]](witMessage)
}
