package format

import (
	"github.com/hayride-dev/bindings/go/shared/ai/msg"
)

type Encode interface {
	Encode(...*msg.Message) ([]byte, error)
}

type Decode interface {
	Decode([]byte) (*msg.Message, error)
}

type Formatter interface {
	Encode
	Decode
}
