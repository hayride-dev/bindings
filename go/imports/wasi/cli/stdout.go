package cli

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stdout"
)

func GetStdout() io.Writer {
	w := wasiio.Writer(stdout.GetStdout())
	return w
}
