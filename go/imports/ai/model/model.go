package model

import (
	"fmt"
	"io"

	"go.bytecodealliance.org/cm"

	"github.com/hayride-dev/bindings/go/imports/ai/types"
	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	witGraph "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	witModel "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
)

type Model interface {
	Push(messages ...*types.Message) error
	Compute(io.Writer) (*types.Message, error)
}

type wacModel struct {
	m witModel.Model
}

// Push take a list of messages, convert them to a list of imports.Message
// and call imported push
func (i *wacModel) Push(messages ...*types.Message) error {
	msgs := make([]witTypes.Message, 0)
	for _, message := range messages {
		content := make([]witTypes.Content, 0)
		for _, c := range message.Content {
			switch c.Type() {
			case "text":
				textContent := c.(*types.TextContent)
				content = append(content, witTypes.ContentText(witTypes.TextContent{
					Text:        textContent.Text,
					ContentType: textContent.ContentType,
				}))
			case "tool-schema":
				toolSchema := c.(*types.ToolSchema)
				content = append(content, witTypes.ContentToolSchema(witTypes.ToolSchema{
					ID:           toolSchema.ID,
					Name:         toolSchema.Name,
					Description:  toolSchema.Description,
					ParamsSchema: toolSchema.ParamsSchema,
				}))
			case "tool-input":
				toolContent := c.(*types.ToolInput)
				content = append(content, witTypes.ContentToolInput(witTypes.ToolInput{
					ContentType: toolContent.ContentType,
					ID:          toolContent.ID,
					Name:        toolContent.Name,
					Input:       toolContent.Input,
				}))
			case "tool-output":
				toolResult := c.(*types.ToolOutput)
				content = append(content, witTypes.ContentToolOutput(witTypes.ToolOutput{
					ContentType: toolResult.ContentType,
					ID:          toolResult.ID,
					Name:        toolResult.Name,
					Output:      toolResult.Output,
				}))
			}
			msgs = append(msgs, witTypes.Message{
				Role:    witTypes.Role(message.Role),
				Content: cm.ToList(content),
			})
		}
	}
	result := i.m.Push(cm.ToList(msgs))
	if result.IsErr() {
		return fmt.Errorf("failed to push message")
	}
	return nil
}

func (i *wacModel) Compute(w io.Writer) (*types.Message, error) {
	if _, ok := w.(wasiio.WasiWriter); !ok {
		return nil, fmt.Errorf("expected io.WasiWriter, got %T", w)
	}

	wasiStream := w.(wasiio.WasiWriter)
	ptr := wasiStream.Ptr()
	output := witModel.OutputStream(cm.Rep(ptr))

	// compute will use the output writer to stream data, result is the "complete" message
	result := i.m.Compute(output)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compute")
	}

	// message should always be a model response ( aka assistant)
	witMsg := result.OK()
	if witMsg.Role != witTypes.RoleAssistant {
		return nil, fmt.Errorf("expected assistant role, got %v", witMsg.Role)
	}

	content := make([]types.Content, 0)
	for _, c := range witMsg.Content.Slice() {
		switch c.String() {
		case "text":
			content = append(content, &types.TextContent{
				Text:        c.Text().Text,
				ContentType: c.Text().ContentType,
			})
		case "tool-input":
			content = append(content, &types.ToolInput{
				ContentType: c.ToolInput().ContentType,
				ID:          c.ToolInput().ID,
				Name:        c.ToolInput().Name,
				Input:       c.ToolInput().Input,
			})
		default:
			return nil, fmt.Errorf("unknown assistant content type: %s", c.String())
		}
	}

	response := &types.Message{
		Role:    types.Role(witMsg.Role),
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

	return &wacModel{m: witModel.NewModel(stream)}, nil
}
