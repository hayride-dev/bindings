package morph

import "github.com/hayride-dev/bindings/go/morph/gen/hayride/morph/errors"

type morphErr struct {
	e *errors.Error
}

func (err *morphErr) Error() string {
	return err.e.Code().String()
}
