package cli

import (
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/cli/stdin"

	"github.com/hayride-dev/bindings/go/wasi/streams"
)

// GetStdin returns a read closer for stdin.
// If block is true, reader will block until data is available.
func GetStdin(block bool) io.ReadCloser {
	if block {
		return streams.ReaderCloser(stdin.GetStdin())
	}
	return nonBlockingReader(stdin.GetStdin())
}
