package silo

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/silo/process"
	"go.bytecodealliance.org/cm"
)

type Process struct {
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) Spawn(path string, args ...string) (int32, error) {
	list := cm.ToList(args)
	result := process.Spawn(path, list)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to spawn process: %v", result.Err())
	}
	return *result.OK(), nil
}

func (p *Process) Wait(pid uint32) (int32, error) {
	result := process.Wait(pid)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to wait for process: %v", result.Err())
	}
	return 0, nil
}

func (p *Process) Status(pid uint32) (bool, error) {
	result := process.Status(pid)
	if result.IsErr() {
		return false, fmt.Errorf("failed to get process status: %v", result.Err())
	}
	return *result.OK(), nil
}

func (p *Process) Kill(pid uint32, sig int32) (int32, error) {
	result := process.Kill(pid, sig)
	if result.IsErr() {
		return -1, fmt.Errorf("failed to kill process: %v", result.Err())
	}
	return 0, nil
}
