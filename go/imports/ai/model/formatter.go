package model

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported format resource. It is only used for
constrcuting the format resources
*/

import (
	"fmt"

	witModel "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type Formatter cm.Resource

func (f Formatter) Encode(messages ...witModel.Message) ([]byte, error) {
	return nil, nil
}

func (f Formatter) Decode(data []byte) (*witModel.Message, error) {
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
		return nil, fmt.Errorf("expected assistant role, got %v", witMsg.Role)
	}

	return witMsg, nil
}

func NewFormatter() Formatter {
	return Formatter(witModel.NewFormat())
}
