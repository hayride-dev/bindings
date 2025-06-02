package models

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type Format cm.Resource

func New() (Format, error) {
	return Format(model.NewFormat()), nil
}

func (f Format) Encode(messages ...types.Message) (string, error) {
	witFormat := cm.Reinterpret[model.Format](f)

	witList := cm.ToList(messages)
	result := witFormat.Encode(cm.Reinterpret[cm.List[model.Message]](witList))
	if result.IsErr() {
		return "", fmt.Errorf("error encoding: %s", result.Err().Code().String())
	}

	return cm.Reinterpret[string](result.OK()), nil
}

func (f Format) Decode(b []byte) (*types.Message, error) {
	witFormat := cm.Reinterpret[model.Format](f)

	data := cm.ToList(b)
	result := witFormat.Decode(data)
	if result.IsErr() {
		return nil, fmt.Errorf("error decoding: %s", result.Err().Code().String())
	}

	message := cm.Reinterpret[types.Message](result.OK())
	return &message, nil
}
