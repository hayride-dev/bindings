package morph

import (
	"context"
	"fmt"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/morph/gen/hayride/morph/spawn"
)

type morph struct {
}

func Spawn() *morph {
	return &morph{}
}

func (m *morph) Execute(ctx context.Context, name string, args []string) (string, error) {
	// one-shot context check
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	// Guest -> args -> Writer (reader - to get bytes ) -> Host
	list := cm.ToList([]string(args))
	result := spawn.Exec(name, list)
	if result.IsErr() {
		return "", fmt.Errorf("failed to exec morph: %s", result.Err().ToDebugString())
	}
	return *result.OK(), nil
}
