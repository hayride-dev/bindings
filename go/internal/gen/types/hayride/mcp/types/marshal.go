package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

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
		if v := c.Image(); v != nil {
			contentType = "image"
			value = v
		}
	case 3:
		if v := c.Audio(); v != nil {
			contentType = "audio"
			value = v
		}
	case 4:
		if v := c.ResourceLink(); v != nil {
			contentType = "resource-link"
			value = v
		}
	case 5:
		if v := c.ResourceContent(); v != nil {
			contentType = "resource-content"
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
			var text TextContent
			if err := json.Unmarshal(raw, &text); err != nil {
				return err
			}
			*c = ContentText(text)
		case "image":
			var image ImageContent
			if err := json.Unmarshal(raw, &image); err != nil {
				return err
			}
			*c = ContentImage(image)
		case "audio":
			var audio AudioContent
			if err := json.Unmarshal(raw, &audio); err != nil {
				return err
			}
			*c = ContentAudio(audio)
		case "resource-link":
			var output ResourceLinkContent
			if err := json.Unmarshal(raw, &output); err != nil {
				return err
			}
			*c = ContentResourceLink(output)
		case "resource-content":
			var resource EmbeddedResourceContent
			if err := json.Unmarshal(raw, &resource); err != nil {
				return err
			}
			*c = ContentResourceContent(resource)
		default:
			return fmt.Errorf("unknown content variant: %s", key)
		}
	}
	return nil
}

func (r *ResourceContents) MarshalJSON() ([]byte, error) {
	var contentType string
	var value interface{}

	switch r.Tag() {
	case 1:
		if v := r.Text(); v != nil {
			contentType = "text"
			value = v
		}
	case 2:
		if v := r.Blob(); v != nil {
			contentType = "blob"
			value = v
		}
	default:
		return nil, fmt.Errorf("unsupported resource contents tag: %d", r.Tag())
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]json.RawMessage{
		contentType: raw,
	})
}

func (r *ResourceContents) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if len(temp) != 1 {
		return errors.New("invalid resource contents format")
	}
	for key, raw := range temp {
		switch key {
		case "text":
			var text TextResourceContents
			if err := json.Unmarshal(raw, &text); err != nil {
				return err
			}
			*r = ResourceContentsText(text)
		case "blob":
			var blob BlobResourceContents
			if err := json.Unmarshal(raw, &blob); err != nil {
				return err
			}
			*r = ResourceContentsBlob(blob)
		default:
			return fmt.Errorf("unknown resource contents variant: %s", key)
		}
	}
	return nil
}
