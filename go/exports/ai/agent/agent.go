package agent

import (
	"unsafe"

	witAgent "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type invokeFunc func(messages []*ai.Message) ([]*ai.Message, error)

type resource struct {
	name       string
	invokeFunc invokeFunc
}

var agent *resource

func init() {
	agent = &resource{}
	witAgent.Exports.Agent.Constructor = agent.constructor
	witAgent.Exports.Agent.Invoke = agent.invoke
}

func Export(name string, f func(messages []*ai.Message) ([]*ai.Message, error)) {
	agent.name = name
	agent.invokeFunc = f
}

func (a *resource) constructor() witAgent.Agent {
	return witAgent.AgentResourceNew(cm.Rep(uintptr(unsafe.Pointer(&agent))))
}

func (a *resource) invoke(self cm.Rep, messages cm.List[witTypes.Message]) (result cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]) {

	messageList := []*ai.Message{}
	for _, m := range messages.Slice() {
		content := make([]ai.Content, 0)
		for _, c := range m.Content.Slice() {
			switch c.String() {
			case "text":
				content = append(content, &ai.TextContent{
					Text:        c.Text().Text,
					ContentType: c.Text().ContentType,
				})
			case "tool-schema":
				content = append(content, &ai.ToolSchema{
					ID:           c.ToolSchema().ID,
					Name:         c.ToolSchema().Name,
					Description:  c.ToolSchema().Description,
					ParamsSchema: c.ToolSchema().ParamsSchema,
				})
			case "tool-input":
				content = append(content, &ai.ToolInput{
					ContentType: c.ToolInput().ContentType,
					ID:          c.ToolInput().ID,
					Name:        c.ToolInput().Name,
					Input:       c.ToolInput().Input,
				})

			case "tool-output":
				content = append(content, &ai.ToolOutput{
					ContentType: c.ToolOutput().ContentType,
					ID:          c.ToolOutput().ID,

					Name:   c.ToolOutput().Name,
					Output: c.ToolOutput().Output,
				})
			default:
				return cm.Err[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeUnknown)))

			}
		}

		messageList = append(messageList, &ai.Message{
			Role:    ai.Role(m.Role),
			Content: content,
		})
	}

	msgs, err := a.invokeFunc(messageList)
	if err != nil {
		wasiErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeInvokeError))
		return cm.Err[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](wasiErr)
	}

	witMsgs := make([]witTypes.Message, 0)
	for _, m := range msgs {
		content := make([]witTypes.Content, 0)
		for _, c := range m.Content {
			switch c.Type() {
			case "text":
				textContent := c.(*ai.TextContent)
				content = append(content, witTypes.ContentText(witTypes.TextContent{
					Text:        textContent.Text,
					ContentType: textContent.ContentType,
				}))
			case "tool-schema":
				toolSchema := c.(*ai.ToolSchema)
				content = append(content, witTypes.ContentToolSchema(witTypes.ToolSchema{
					ID:           toolSchema.ID,
					Name:         toolSchema.Name,
					Description:  toolSchema.Description,
					ParamsSchema: toolSchema.ParamsSchema,
				}))
			case "tool-input":
				toolContent := c.(*ai.ToolInput)
				content = append(content, witTypes.ContentToolInput(witTypes.ToolInput{
					ContentType: toolContent.ContentType,
					ID:          toolContent.ID,
					Name:        toolContent.Name,
					Input:       toolContent.Input,
				}))
			case "tool-output":
				toolResult := c.(*ai.ToolOutput)
				content = append(content, witTypes.ContentToolOutput(witTypes.ToolOutput{
					ContentType: toolResult.ContentType,
					ID:          toolResult.ID,
					Name:        toolResult.Name,
					Output:      toolResult.Output,
				}))
			default:
				return cm.Err[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeUnknown)))
			}
		}
		witMsgs = append(witMsgs, witTypes.Message{
			Role:    witTypes.Role(m.Role),
			Content: cm.ToList(content),
		})
	}

	return cm.OK[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](cm.ToList(witMsgs))
}
