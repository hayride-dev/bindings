package cli

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stderr"
)

func GetStderr() io.Writer {
	w := wasiio.Writer(stderr.GetStderr())
	return w
}
