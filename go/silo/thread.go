package silo

import (
	"fmt"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/silo/gen/hayride/silo/threads"
)

type Thread struct {
}

func NewThread() *Thread {
	return &Thread{}
}

func (p *Thread) Spawn(path string, args ...string) (int32, error) {
	list := cm.ToList(args)
	result := threads.Spawn(path, list)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to spawn thread: %v", result.Err())
	}
	return *result.OK(), nil
}

func (p *Thread) Wait(threadID uint32) (int32, error) {
	result := threads.Wait(threadID)
	if result.IsErr() {
		return -1, fmt.Errorf("thread error: %v", result.Err())
	}
	return 0, nil
}
