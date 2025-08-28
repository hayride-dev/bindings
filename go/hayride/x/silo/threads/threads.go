package threads

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/x/silo"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/silo/threads"
	"go.bytecodealliance.org/cm"
)

// Spawn spawns a new thread with the given path, function, and arguments.
func Spawn(path string, function string, args []string, envs [][2]string) (Thread, error) {
	argList := cm.ToList(args)
	envList := cm.ToList(envs)
	result := threads.Spawn(path, function, argList, envList)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to spawn thread: %v", result.Err())
	}

	return thread(*result.OK()), nil
}

// Status returns the status of the thread with the given ID.
func Status(threadID string) (*silo.ThreadMetadata, error) {
	result := threads.Status(threadID)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get thread status: %v", result.Err())
	}

	return cm.Reinterpret[*silo.ThreadMetadata](result.OK()), nil
}

// Kill kills the thread with the given ID.
func Kill(threadID string) error {
	result := threads.Kill(threadID)
	if result.IsErr() {
		return fmt.Errorf("failed to kill thread: %v", result.Err())
	}

	return nil
}

// Group returns a list of all threads with their metadata.
func Group() ([]silo.ThreadMetadata, error) {
	result := threads.Group()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to list threads: %v", result.Err())
	}

	return cm.Reinterpret[[]silo.ThreadMetadata](result.OK().Slice()), nil
}
