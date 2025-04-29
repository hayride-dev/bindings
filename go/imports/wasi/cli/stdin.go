package cli

import (
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stdin"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
)

// GetStdin returns a read closer for stdin.
// If block is true, reader will block until data is available.
func GetStdin(block bool) io.ReadCloser {
	if block {
		return wasiio.ReaderCloser(stdin.GetStdin())
	}
	return nonBlockingReader(stdin.GetStdin())
}
