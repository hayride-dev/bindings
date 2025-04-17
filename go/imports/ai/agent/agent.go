package agent

import (
	"fmt"

	witAgent "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agent"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"

	"go.bytecodealliance.org/cm"
)

type Agent cm.Resource

func NewAgent() Agent {
	return Agent(witAgent.NewAgent())
}

func (a Agent) Invoke(messages []*ai.Message) ([]*ai.Message, error) {
	wa := cm.Reinterpret[witAgent.Agent](a)

	witMsgs := make([]types.Message, 0)
	for _, m := range messages {
		content := make([]types.Content, 0)
		for _, c := range m.Content {
			switch c.Type() {
			case "text":
				textContent := c.(*ai.TextContent)
				content = append(content, types.ContentText(types.TextContent{
					Text:        textContent.Text,
					ContentType: textContent.ContentType,
				}))
			case "tool-schema":
				toolSchema := c.(*ai.ToolSchema)
				content = append(content, types.ContentToolSchema(types.ToolSchema{
					ID:           toolSchema.ID,
					Name:         toolSchema.Name,
					Description:  toolSchema.Description,
					ParamsSchema: toolSchema.ParamsSchema,
				}))
			case "tool-input":
				toolContent := c.(*ai.ToolInput)
				content = append(content, types.ContentToolInput(types.ToolInput{
					ContentType: toolContent.ContentType,
					ID:          toolContent.ID,
					Name:        toolContent.Name,
					Input:       toolContent.Input,
				}))
			case "tool-output":
				toolResult := c.(*ai.ToolOutput)
				content = append(content, types.ContentToolOutput(types.ToolOutput{
					ContentType: toolResult.ContentType,
					ID:          toolResult.ID,
					Name:        toolResult.Name,
					Output:      toolResult.Output,
				}))
			default:
				return nil, fmt.Errorf("unknown content type: %s", c.Type())
			}
		}
		witMsgs = append(witMsgs, types.Message{
			Role:    types.Role(m.Role),
			Content: cm.ToList(content),
		})
	}

	result := wa.Invoke(cm.ToList(witMsgs))
	if result.IsErr() {
		// TODO: handle error
		return nil, fmt.Errorf("failed to invoke agent")
	}

	msgs := make([]*ai.Message, 0)
	for _, m := range result.OK().Slice() {
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
					Name:        c.ToolOutput().Name,
					Output:      c.ToolOutput().Output,
				})
			default:
				return nil, fmt.Errorf("unknown content type: %s", c.String())
			}
		}
		msgs = append(msgs, &ai.Message{
			Role:    ai.Role(m.Role),
			Content: content,
		})
	}

	return msgs, nil
}
