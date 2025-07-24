package models

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

var _ Format = (*FormatResource)(nil)

type Format interface {
	Encode(messages ...types.Message) ([]byte, error)
	Decode(b []byte) (*types.Message, error)
}

type FormatResource cm.Resource

func New() (Format, error) {
	return FormatResource(model.NewFormat()), nil
}

func (f FormatResource) Encode(messages ...types.Message) ([]byte, error) {
	witFormat := cm.Reinterpret[model.Format](f)

	witList := cm.ToList(messages)
	result := witFormat.Encode(cm.Reinterpret[cm.List[model.Message]](witList))
	if result.IsErr() {
		return nil, fmt.Errorf("error encoding: %s", result.Err().Code().String())
	}

	return result.OK().Slice(), nil
}

func (f FormatResource) Decode(b []byte) (*types.Message, error) {
	witFormat := cm.Reinterpret[model.Format](f)

	data := cm.ToList(b)
	result := witFormat.Decode(data)
	if result.IsErr() {
		return nil, fmt.Errorf("error decoding: %s", result.Err().Code().String())
	}

	return cm.Reinterpret[*types.Message](result.OK()), nil
}
