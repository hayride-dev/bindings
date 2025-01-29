package rag

import "github.com/hayride-dev/bindings/go/ai/gen/imports/hayride/ai/rag"

type ragErr struct {
	e *rag.Error
}

func (err *ragErr) Error() string {
	return err.e.Code().String()
}
