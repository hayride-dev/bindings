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

	// Capitalize the role string to match wit-bindgen's enum serialization format.
	roleStr := m.Role.String()
	if len(roleStr) > 0 {
		roleStr = strings.ToUpper(roleStr[:1]) + roleStr[1:]
	}

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
	var content []Content
	for _, raw := range aux.Content {
		var cw Content
		if err := json.Unmarshal(raw, &cw); err != nil {
			return fmt.Errorf("failed to unmarshal content item: %w", err)
		}
		content = append(content, cw)
	}
	m.Content = cm.ToList(content)

	return nil
}

func (c Content) MarshalJSON() ([]byte, error) {
	var contentType string
	var value interface{}

	switch c.Tag() {
	case 1:
		if v := c.Text(); v != nil {
			contentType = "Text"
			value = v
		}
	case 2:
		if v := c.ToolSchema(); v != nil {
			contentType = "ToolSchema"
			value = v
		}
	case 3:
		if v := c.ToolInput(); v != nil {
			contentType = "ToolInput"
			value = v
		}
	case 4:
		if v := c.ToolOutput(); v != nil {
			contentType = "ToolOutput"
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
		case "Text":
			var text TextContent
			if err := json.Unmarshal(raw, &text); err != nil {
				return err
			}
			*c = ContentText(text)
		case "ToolSchema":
			var schema ToolSchema
			if err := json.Unmarshal(raw, &schema); err != nil {
				return err
			}
			*c = ContentToolSchema(schema)
		case "ToolInput":
			var input ToolInput
			if err := json.Unmarshal(raw, &input); err != nil {
				return err
			}
			*c = ContentToolInput(input)
		case "ToolOutput":
			var output ToolOutput
			if err := json.Unmarshal(raw, &output); err != nil {
				return err
			}
			*c = ContentToolOutput(output)
		default:
			return fmt.Errorf("unknown content variant: %s", key)
		}
	}
	return nil
}
