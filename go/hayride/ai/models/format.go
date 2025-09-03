package models

import (
	"github.com/hayride-dev/bindings/go/hayride/ai"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
	"go.bytecodealliance.org/cm"
)

var _ Format = (*FormatResource)(nil)

type Format interface {
	Encode(messages ...ai.Message) ([]byte, error)
	Decode(b []byte) (*ai.Message, error)
}

type FormatResource cm.Resource

func New() (Format, error) {
	return FormatResource(model.NewFormat()), nil
}

func (f FormatResource) Encode(messages ...ai.Message) ([]byte, error) {
	witFormat := cm.Reinterpret[model.Format](f)

	witList := cm.ToList(messages)
	result := witFormat.Encode(cm.Reinterpret[cm.List[model.Message]](witList))
	if result.IsErr() {
		err := result.Err()

		return nil, newError(err)
	}

	return result.OK().Slice(), nil
}

func (f FormatResource) Decode(b []byte) (*ai.Message, error) {
	witFormat := cm.Reinterpret[model.Format](f)

	data := cm.ToList(b)
	result := witFormat.Decode(data)
	if result.IsErr() {
		err := result.Err()

		return nil, newError(err)
	}

	return cm.Reinterpret[*ai.Message](result.OK()), nil
}
