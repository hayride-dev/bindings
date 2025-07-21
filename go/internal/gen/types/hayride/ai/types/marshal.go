package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.bytecodealliance.org/cm"
)

func (m Message) MarshalJSON() ([]byte, error) {
	var content []json.RawMessage
	for _, c := range m.Content.Slice() {
		raw, err := json.Marshal(c)
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

func (m *Message) UnmarshalJSON(data []byte) error {
	var aux struct {
		Role    string            `json:"role"`
		Content []json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Decode role
	if err := m.Role.UnmarshalText([]byte(strings.ToLower(aux.Role))); err != nil {
		return fmt.Errorf("invalid role: %w", err)
	}

	// Decode content
	var content []MessageContent
	for _, raw := range aux.Content {
		var cw MessageContent
		if err := json.Unmarshal(raw, &cw); err != nil {
			return fmt.Errorf("failed to unmarshal content item: %w", err)
		}
		content = append(content, cw)
	}
	m.Content = cm.ToList(content)

	return nil
}

func (c MessageContent) MarshalJSON() ([]byte, error) {
	var contentType string
	var value interface{}

	switch c.Tag() {
	case 1:
		if v := c.Text(); v != nil {
			contentType = "text"
			value = v
		}
	case 2:
		if v := c.Blob(); v != nil {
			contentType = "blob"
			value = v
		}
	case 3:
		if v := c.Tools(); v != nil {
			contentType = "tools"
			value = v
		}
	case 4:
		if v := c.ToolInput(); v != nil {
			contentType = "tool-input"
			value = v
		}
	case 5:
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

func (c *MessageContent) UnmarshalJSON(data []byte) error {
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
			var text string
			if err := json.Unmarshal(raw, &text); err != nil {
				return err
			}
			*c = MessageContentText(text)
		case "blob":
			var blob cm.List[uint8]
			if err := json.Unmarshal(raw, &blob); err != nil {
				return err
			}
			*c = MessageContentBlob(blob)
		case "tools":
			var tools []Tool
			if err := json.Unmarshal(raw, &tools); err != nil {
				return err
			}
			*c = MessageContentTools(cm.ToList(tools))
		case "tool-input":
			var input CallToolParams
			if err := json.Unmarshal(raw, &input); err != nil {
				return err
			}
			*c = MessageContentToolInput(input)
		case "tool-output":
			var output CallToolResult
			if err := json.Unmarshal(raw, &output); err != nil {
				return err
			}
			*c = MessageContentToolOutput(output)
		default:
			return fmt.Errorf("unknown content variant: %s", key)
		}
	}
	return nil
}
