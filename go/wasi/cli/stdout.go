package cli

import (
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/cli/stdout"
	"github.com/hayride-dev/bindings/go/wasi/streams"
)

// GetStdout returns a writer for stdout.
// If block is true, writer will block until ready to write.
func GetStdout(block bool) io.Writer {
	if block {
		return streams.Writer(stdout.GetStdout())
	}
	return nonBlockingWriter(stdout.GetStdout())
}
