package silo

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/silo/threads"
	"go.bytecodealliance.org/cm"
)

type Thread cm.Resource

// ID returns the ID of the thread.
func (t *Thread) ID() (string, error) {
	witThread := cm.Reinterpret[threads.Thread](t)

	result := witThread.ID()
	if result.IsErr() {
		return "", fmt.Errorf("failed to get thread ID: %v", result.Err())
	}

	return *result.OK(), nil
}

// Wait waits for the thread to finish and returns the result.
func (t *Thread) Wait(threadID string) ([]byte, error) {
	witThread := cm.Reinterpret[threads.Thread](t)

	result := witThread.Wait()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to wait for thread: %v", result.Err())
	}

	return result.OK().Slice(), nil
}

// Spawn spawns a new thread with the given path, function, and arguments.
func Spawn(path string, function string, args ...string) (Thread, error) {
	list := cm.ToList(args)
	result := threads.Spawn(path, function, list)
	if result.IsErr() {
		return 0, fmt.Errorf("failed to spawn thread: %v", result.Err())
	}

	thread := cm.Reinterpret[Thread](result.OK())
	return thread, nil
}

// Status returns the status of the thread with the given ID.
func Status(threadID string) (*threads.ThreadMetadata, error) {
	result := threads.Status(threadID)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get thread status: %v", result.Err())
	}

	return result.OK(), nil
}

// Kill kills the thread with the given ID.
func Kill(threadID string) error {
	result := threads.Kill(threadID)
	if result.IsErr() {
		return fmt.Errorf("failed ot kill thread: %v", result.Err())
	}

	return nil
}

// Group returns a list of all threads with their metadata.
func Group() ([]threads.ThreadMetadata, error) {
	result := threads.Group()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to list threads: %v", result.Err())
	}

	return result.OK().Slice(), nil
}
