package wasinn

import (
	"github.com/hayride-dev/bindings/go/ai/gen/imports/wasi/io/streams"
	"github.com/hayride-dev/bindings/go/ai/gen/imports/wasi/nn/errors"
)

type wasinnErr struct {
	e *errors.Error
}

func (err *wasinnErr) Error() string {
	return err.e.Code().String()
}

type streamErr struct {
	e *streams.StreamError
}

func (err *streamErr) Error() string {
	return err.e.String()
}
