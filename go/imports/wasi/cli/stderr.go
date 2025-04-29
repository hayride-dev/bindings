package cli

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stderr"
)

// GetStderr returns a writer for stderr.
// If block is true, writer will block until ready to write.
func GetStderr(block bool) io.Writer {
	if block {
		return wasiio.Writer(stderr.GetStderr())
	}
	return nonBlockingWriter(stderr.GetStderr())
}
