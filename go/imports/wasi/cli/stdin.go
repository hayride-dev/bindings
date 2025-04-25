package cli

import (
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stdin"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
)

func GetStdin() io.Reader {
	r := wasiio.ReaderCloser(stdin.GetStdin())
	return r
}
