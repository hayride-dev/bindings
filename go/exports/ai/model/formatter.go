package model

import "github.com/hayride-dev/bindings/go/exports/ai/types"

type Encode interface {
	Encode(...*types.Message) ([]byte, error)
}

type Decode interface {
	Decode([]byte) (*types.Message, error)
}

type Formatter interface {
	Encode
	Decode
}
