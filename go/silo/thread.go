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

func (p *Thread) Spawn(path string, function string, args ...string) (string, error) {
	list := cm.ToList(args)
	result := threads.Spawn(path, function, list)
	if result.IsErr() {
		return "", fmt.Errorf("failed to spawn thread: %v", result.Err())
	}
	return *result.OK(), nil
}

func (p *Thread) Wait(threadID string) ([]byte, error) {
	result := threads.Wait(threadID)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to wait for thread: %v", result.Err())
	}

	return result.OK().Slice(), nil
}

func (p *Thread) Kill(threadID string) error {
	result := threads.Kill(threadID)
	if result.IsErr() {
		return fmt.Errorf("failed ot kill thread: %v", result.Err())
	}

	return nil
}
