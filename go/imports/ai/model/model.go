package model

import (
	"fmt"

	witGraph "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	witModel "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type Model interface {
	Compute(messages []*ai.Message) (*ai.Message, error)
}

type wacModel struct {
	m witModel.Model
}

func (i *wacModel) Compute(messages []*ai.Message) (*ai.Message, error) {
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
				return nil, fmt.Errorf("unknown content type: %s", c.Type())
			}
		}
		msgs = append(msgs, witTypes.Message{
			Role:    witTypes.Role(message.Role),
			Content: cm.ToList(content),
		})
	}
	// compute will use the output writer to stream data, result is the "complete" message
	result := i.m.Compute(cm.ToList(msgs))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compute")
	}

	// message should always be a model response ( aka assistant)
	witMsg := result.OK()
	if witMsg.Role != witTypes.RoleAssistant {
		return nil, fmt.Errorf("expected assistant role, got %v", witMsg.Role)
	}

	content := make([]ai.Content, 0)
	for _, c := range witMsg.Content.Slice() {
		switch c.String() {
		case "text":
			content = append(content, &ai.TextContent{
				Text:        c.Text().Text,
				ContentType: c.Text().ContentType,
			})
		case "tool-input":
			content = append(content, &ai.ToolInput{
				ContentType: c.ToolInput().ContentType,
				ID:          c.ToolInput().ID,
				Name:        c.ToolInput().Name,
				Input:       c.ToolInput().Input,
			})
		default:
			return nil, fmt.Errorf("unknown assistant content type: %s", c.String())
		}
	}

	response := &ai.Message{
		Role:    ai.Role(witMsg.Role),
		Content: content,
	}

	return response, nil
}

func New(options ...Option[*ModelOptions]) (Model, error) {
	opts := defaultModelOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	result := witGraph.LoadByName(opts.name)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to load graph")
	}
	model := result.OK()

	resultCtxStream := model.InitExecutionContextStream()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to init execution graph context stream")
	}
	stream := *resultCtxStream.OK()
	format := witModel.NewFormat()

	return &wacModel{
		m: witModel.NewModel(format, stream)}, nil
}
