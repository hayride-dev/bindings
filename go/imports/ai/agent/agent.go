package agent

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/imports/ai/ctx"
	"github.com/hayride-dev/bindings/go/imports/ai/model"
	witAgent "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agent"
	witContext "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/context"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"

	"go.bytecodealliance.org/cm"
)

type Agent cm.Resource

func NewAgent() Agent {
	return Agent(witAgent.NewAgent())
}

func (a Agent) Invoke(ctx ctx.Context, model model.Model) ([]*ai.Message, error) {
	wa := cm.Reinterpret[witAgent.Agent](a)
	wctx := cm.Reinterpret[witContext.Context](ctx)
	wmodel := cm.Reinterpret[witAgent.Model](model)
	result := wa.Invoke(wctx, wmodel)
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
