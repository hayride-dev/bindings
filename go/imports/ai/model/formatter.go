package models

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported format resource. It is only used for
constrcuting the format resources
*/

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/gen/domain/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

type Formatter cm.Resource

func (f Formatter) Encode(messages ...types.Message) ([]byte, error) {
	return nil, nil
}

func (f Formatter) Decode(data []byte) (*types.Message, error) {
	if len(data) == 0 {
		return nil, nil
	}

	wformat := cm.Reinterpret[model.Format](f)
	result := wformat.Decode(cm.ToList(data))
	if result.IsErr() {
		return nil, fmt.Errorf("decode error: %s", result.Err().Data())
	}
	witMsg := result.OK()

	msg := cm.Reinterpret[*types.Message](witMsg)

	return msg, nil
}

func NewFormatter() Formatter {
	return Formatter(model.NewFormat())
}
