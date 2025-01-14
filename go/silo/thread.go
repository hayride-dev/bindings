package silo

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/silo/gen/hayride/silo/threads"
	"go.bytecodealliance.org/cm"
)

type Thread struct {
}

func NewThread() *Thread {
	return &Thread{}
}

func (p *Thread) Spawn(path string, function string, args ...string) (int32, error) {
	list := cm.ToList(args)
	result := threads.Spawn(path, function, list)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to spawn thread: %v", result.Err())
	}
	return *result.OK(), nil
}

func (p *Thread) Wait(threadID uint32) ([]byte, error) {
	result := threads.Wait(threadID)
	if result.IsErr() {
		return nil, fmt.Errorf("thread error: %v", result.Err())
	}

	return result.OK().Slice(), nil
}
