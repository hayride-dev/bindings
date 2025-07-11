package threads

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/silo/threads"
	"go.bytecodealliance.org/cm"
)

type Thread interface {
	ID() (string, error)
	Wait(threadID string) ([]byte, error)
}

type thread cm.Resource

// ID returns the ID of the thread.
func (t thread) ID() (string, error) {
	threadResource := cm.Reinterpret[threads.Thread](t)

	result := threadResource.ID()
	if result.IsErr() {
		return "", fmt.Errorf("failed to get thread ID: %v", result.Err())
	}

	return *result.OK(), nil
}

// Wait waits for the thread to finish and returns the result.
func (t thread) Wait(threadID string) ([]byte, error) {
	threadResource := cm.Reinterpret[threads.Thread](t)

	result := threadResource.Wait()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to wait for thread: %v", result.Err())
	}

	return result.OK().Slice(), nil
}
