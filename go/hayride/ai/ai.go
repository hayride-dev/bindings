package ai

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/types"
)

type TextContent = types.TextContent
type ToolSchema = types.ToolSchema
type ToolInput = types.ToolInput
type ToolOutput = types.ToolOutput
type Role = types.Role

type Message struct {
	types.Message
}

type Content struct {
	types.Content
}

func (m Message) MarshalJSON() ([]byte, error) {
	var content []json.RawMessage
	for _, c := range m.Content.Slice() {
		t := Content{c}
		raw, err := json.Marshal(t)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal content: %w", err)
		}
		content = append(content, raw)
	}

	roleStr := m.Role.String()
	return json.Marshal(struct {
		Role    string            `json:"role"`
		Content []json.RawMessage `json:"content"`
	}{
		Role:    roleStr,
		Content: content,
	})
}

func (c Content) MarshalJSON() ([]byte, error) {
	var contentType string
	var value interface{}

	switch c.Tag() {
	case 1:
		if v := c.Text(); v != nil {
			contentType = "text"
			value = v
		}
	case 2:
		if v := c.ToolSchema(); v != nil {
			contentType = "tool-schema"
			value = v
		}
	case 3:
		if v := c.ToolInput(); v != nil {
			contentType = "tool-input"
			value = v
		}
	case 4:
		if v := c.ToolOutput(); v != nil {
			contentType = "tool-output"
			value = v
		}
	default:
		return nil, fmt.Errorf("unsupported content tag: %d", c.Tag())
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]json.RawMessage{
		contentType: raw,
	})
}

func (c *Content) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if len(temp) != 1 {
		return errors.New("invalid content format")
	}
	for key, raw := range temp {
		switch key {
		case "text":
			var text types.TextContent
			if err := json.Unmarshal(raw, &text); err != nil {
				return err
			}
			*c = Content{types.ContentText(text)}
		case "tool-schema":
			var schema types.ToolSchema
			if err := json.Unmarshal(raw, &schema); err != nil {
				return err
			}
			*c = Content{types.ContentToolSchema(schema)}
		case "tool-input":
			var input types.ToolInput
			if err := json.Unmarshal(raw, &input); err != nil {
				return err
			}
			*c = Content{types.ContentToolInput(input)}
		case "tool-output":
			var output types.ToolOutput
			if err := json.Unmarshal(raw, &output); err != nil {
				return err
			}
			*c = Content{types.ContentToolOutput(output)}
		default:
			return fmt.Errorf("unknown content variant: %s", key)
		}
	}
	return nil
}
