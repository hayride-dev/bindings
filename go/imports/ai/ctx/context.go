package ctx

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported context resource.
*/
import (
	"fmt"

	witContext "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/context"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type Context interface {
	Push(messages ...*ai.Message) error
	Messages() ([]*ai.Message, error)
	Next() (*ai.Message, error)
}

type wacContext struct {
	ctx witContext.Context
}

// Push take a list of messages, convert them to a list of wit Messages
// and call imported context push
func (c *wacContext) Push(messages ...*ai.Message) error {
	msgs := make([]witTypes.Message, 0)
	for _, message := range messages {
		content := make([]witTypes.Content, 0)
		for _, c := range message.Content {
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
				return fmt.Errorf("unknown content type: %s", c.Type())
			}
		}
		msgs = append(msgs, witTypes.Message{
			Role:    witTypes.Role(message.Role),
			Content: cm.ToList(content),
		})
	}
	result := c.ctx.Push(cm.ToList(msgs))
	if result.IsErr() {
		// TODO: handle error result
		return fmt.Errorf("failed to push message")
	}
	return nil
}

// Messages take a list of messages, convert them to a list of wit Messages
// and call imported context push
func (c *wacContext) Messages() ([]*ai.Message, error) {
	msgs := make([]*ai.Message, 0)
	result := c.ctx.Messages()
	if result.IsErr() {
		// TODO: handle error result
		return nil, fmt.Errorf("failed to get messages")
	}
	witMessages := result.OK()
	for _, message := range witMessages.Slice() {
		content := make([]ai.Content, 0)
		for _, c := range message.Content.Slice() {
			if !c.None() {
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
						ID:          c.ToolInput().ID,
						Name:        c.ToolInput().Name,
						Input:       c.ToolInput().Input,
						ContentType: c.ToolInput().ContentType,
					})
				case "tool-output":
					content = append(content, &ai.ToolOutput{
						ID:          c.ToolOutput().ID,
						Name:        c.ToolOutput().Name,
						Output:      c.ToolOutput().Output,
						ContentType: c.ToolOutput().ContentType,
					})
				default:
					return nil, fmt.Errorf("unknown content type: %s", c.String())
				}
			}
		}
		msgs = append(msgs, &ai.Message{
			Role:    ai.Role(message.Role),
			Content: content,
		})
	}
	return msgs, nil
}

func (c *wacContext) Next() (*ai.Message, error) {
	result := c.ctx.Next()
	if result.IsErr() {
		// TODO : handle error result
		return nil, fmt.Errorf("failed to get next message")
	}
	content := make([]ai.Content, 0)

	for _, c := range result.OK().Content.Slice() {
		if !c.None() {
			switch c.String() {
			case "text":
				content = append(content, ai.TextContent{
					Text:        c.Text().Text,
					ContentType: c.Text().ContentType,
				})

			case "tool-schema":
				content = append(content, ai.ToolSchema{
					ID:           c.ToolSchema().ID,
					Name:         c.ToolSchema().Name,
					Description:  c.ToolSchema().Description,
					ParamsSchema: c.ToolSchema().ParamsSchema,
				})
			case "tool-input":
				content = append(content, ai.ToolInput{
					ID:          c.ToolInput().ID,
					Name:        c.ToolInput().Name,
					Input:       c.ToolInput().Input,
					ContentType: c.ToolInput().ContentType,
				})
			case "tool-output":
				content = append(content, ai.ToolOutput{
					ID:          c.ToolOutput().ID,
					Name:        c.ToolOutput().Name,
					Output:      c.ToolOutput().Output,
					ContentType: c.ToolOutput().ContentType,
				})
			default:
				return nil, fmt.Errorf("unknown content type: %s", c.String())
			}
		}
	}
	return &ai.Message{
		Role:    ai.Role(result.OK().Role),
		Content: content,
	}, nil
}

func NewContext() *wacContext {
	return &wacContext{
		// connect the wac'd context component
		ctx: witContext.NewContext(),
	}
}
