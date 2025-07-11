package cli

import (
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/cli/stderr"
	"github.com/hayride-dev/bindings/go/wasi/streams"
)

// GetStderr returns a writer for stderr.
// If block is true, writer will block until ready to write.
func GetStderr(block bool) io.Writer {
	if block {
		return streams.Writer(stderr.GetStderr())
	}
	return nonBlockingWriter(stderr.GetStderr())
}
