package process

import (
	"fmt"

	witProcess "github.com/hayride-dev/bindings/go/internal/gen/hayride/silo/process"
	"go.bytecodealliance.org/cm"
)

// key value tuple for environment variables
type Tuple [][2]string

func Spawn(path string, args []string, envs Tuple) (int32, error) {
	argList := cm.ToList(args)
	envList := cm.ToList(envs)
	result := witProcess.Spawn(path, argList, envList)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to spawn process: %v", result.Err())
	}
	return *result.OK(), nil
}

func Wait(pid uint32) (int32, error) {
	result := witProcess.Wait(pid)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to wait for process: %v", result.Err())
	}
	return 0, nil
}

func Status(pid uint32) (bool, error) {
	result := witProcess.Status(pid)
	if result.IsErr() {
		return false, fmt.Errorf("failed to get process status: %v", result.Err())
	}
	return *result.OK(), nil
}

func Kill(pid uint32, sig int32) (int32, error) {
	result := witProcess.Kill(pid, sig)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to kill process: %v", result.Err())
	}
	return 0, nil
}
