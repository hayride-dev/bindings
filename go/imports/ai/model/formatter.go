package model

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported format resource. It is only used for
constrcuting the format resources
*/

import (
	"fmt"

	witModel "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/shared/domain/ai"
	"go.bytecodealliance.org/cm"
)

type Formatter cm.Resource

func (f Formatter) Encode(messages ...*ai.Message) ([]byte, error) {
	return nil, nil
}

func (f Formatter) Decode(data []byte) (*ai.Message, error) {
	if len(data) == 0 {
		return nil, nil
	}

	wformat := cm.Reinterpret[witModel.Format](f)
	result := wformat.Decode(cm.ToList(data))
	if result.IsErr() {
		return nil, fmt.Errorf("decode error: %s", result.Err().Data())
	}
	witMsg := result.OK()

	// NOTE: message should always be a model response ( aka assistant).
	// Should we make this assumption?
	if witMsg.Role != witTypes.RoleAssistant {
		return nil, nil
	}

	content := make([]ai.Content, 0)
	for _, c := range witMsg.Content.Slice() {
		switch c.String() {
		case "text":
			content = append(content, &ai.TextContent{
				Text:        c.Text().Text,
				ContentType: c.Text().ContentType,
			})
		default:
			return nil, nil
		}
	}

	return &ai.Message{
		Role:    ai.RoleAssistant,
		Content: content,
	}, nil
}

func NewFormatter() Formatter {
	return Formatter(witModel.NewFormat())
}
