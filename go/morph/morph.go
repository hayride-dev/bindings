package morph

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/morph/gen/hayride/morph/spawn"
)

type morph struct {
}

func Spawn() *morph {
	return &morph{}
}

func (m *morph) Sync(name string, args []string) (string, error) {
	list := cm.ToList([]string(args))
	result := spawn.Sync(name, list)
	if result.IsErr() {
		return "", &morphErr{result.Err()}
	}
	return *result.OK(), nil
}
