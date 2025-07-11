package wac

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/wac/wac"
	"go.bytecodealliance.org/cm"
)

func Compose(path string) ([]byte, error) {
	result := wac.Compose(path)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compose: %v", result.Err().Data())
	}
	return result.OK().Slice(), nil
}

func Plug(socket string, plugs []string) ([]byte, error) {
	result := wac.Plug(socket, cm.ToList(plugs))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to plug: %v", result.Err().Data())
	}
	return result.OK().Slice(), nil
}
