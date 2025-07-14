package version

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/core/version"
)

func Latest() (string, error) {
	result := version.Latest()
	if result.IsErr() {
		return "", fmt.Errorf("failed to get version: %v", result.Err().Data())
	}
	return *result.OK(), nil
}
