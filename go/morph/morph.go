package morph

import (
	"fmt"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/morph/gen/hayride/silo/threads"
)

type morph struct {
}

func New() *morph {
	return &morph{}
}

func (p *morph) Spawn(path string, args ...string) (int32, error) {
	list := cm.ToList(args)
	result := threads.Spawn(path, list)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to spawn thread: %v", result.Err())
	}
	return *result.OK(), nil
}

func (p *morph) Wait(threadID uint32) (int32, error) {
	result := threads.Wait(threadID)
	if result.IsErr() {
		return -1, fmt.Errorf("thread error: %v", result.Err())
	}
	return 0, nil
}
