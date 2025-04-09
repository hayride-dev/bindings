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

func (f *formatResourceTable) destructor(self cm.Rep) {
	delete(f.resources, self)
}

func (f *formatResourceTable) encode(self cm.Rep, messages cm.List[witModel.Message]) cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error] {
	if _, ok := f.resources[self]; !ok {
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeContextEncode)))
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
				return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeUnknown)))
			}
		}
		msgs = append(msgs, &ai.Message{
			Role:    ai.Role(msg.Role),
			Content: content,
		})
	}
	data, err := f.resources[self].Encode(msgs...)
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
	content := make([]witTypes.Content, 0)
	for _, c := range message.Content {
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
			return cm.Err[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](witModel.ErrorResourceNew(cm.Rep(witModel.ErrorCodeUnknown)))
		}
	}
	msg := witTypes.Message{
		Role:    witTypes.Role(message.Role),
		Content: cm.ToList(content),
	}
	return cm.OK[cm.Result[witModel.MessageShape, witModel.Message, witModel.Error]](msg)
}
