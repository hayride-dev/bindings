package threads

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/hayride/silo/threads"
	"go.bytecodealliance.org/cm"
)

type Thread interface {
	ID() (string, error)
	Wait(threadID string) ([]byte, error)
}

type thread cm.Resource

// ID returns the ID of the thread.
func (t thread) ID() (string, error) {
	witThread := cm.Reinterpret[threads.Thread](t)

	result := witThread.ID()
	if result.IsErr() {
		return "", fmt.Errorf("failed to get thread ID: %v", result.Err())
	}

	return *result.OK(), nil
}

// Wait waits for the thread to finish and returns the result.
func (t thread) Wait(threadID string) ([]byte, error) {
	witThread := cm.Reinterpret[threads.Thread](t)

	result := witThread.Wait()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to wait for thread: %v", result.Err())
	}

	return result.OK().Slice(), nil
}
