package agent

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/imports/ai/ctx"
	"github.com/hayride-dev/bindings/go/imports/ai/model"
	witAgent "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type invokeFunc func(ctx ctx.Context, model model.Model) ([]*ai.Message, error)

type wacAgent struct {
	name        string
	instruction string
	invokeFunc  invokeFunc
}

var agent *wacAgent

func init() {
	agent = &wacAgent{}
	witAgent.Exports.Agent.Constructor = agent.wacConstructorfunc
	witAgent.Exports.Agent.Invoke = agent.wacInvoke
}

func Export(name string, instruction string, f func(ctx ctx.Context, model model.Model) ([]*ai.Message, error)) {
	agent.name = name
	agent.instruction = instruction
	agent.invokeFunc = f
}

func (a *wacAgent) wacConstructorfunc() witAgent.Agent {
	return witAgent.AgentResourceNew(cm.Rep(uintptr(unsafe.Pointer(&agent))))
}

func (a *wacAgent) wacInvoke(self cm.Rep, ctx_ cm.Rep, model_ cm.Rep) (result cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]) {
	context := ctx.Context(ctx_)
	model := model.Model(model_)

	systemprompt := &ai.Message{
		Role:    ai.RoleSystem,
		Content: []ai.Content{&ai.TextContent{Text: a.instruction}},
	}

	// TODO: eval a way to avoid setting this on each invoke
	context.Push(systemprompt)

	msgs, err := a.invokeFunc(context, model)
	if err != nil {
		wasiErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeInvokeError))
		return cm.Err[cm.Result[cm.List[witTypes.Message], cm.List[witTypes.Message], witAgent.Error]](wasiErr)
	}

	context.Push(msgs...)

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
