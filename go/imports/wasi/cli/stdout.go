package cli

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stdout"
)

// GetStdout returns a writer for stdout.
// If block is true, writer will block until ready to write.
func GetStdout(block bool) io.Writer {
	if block {
		return wasiio.Writer(stdout.GetStdout())
	}
	return nonBlockingWriter(stdout.GetStdout())
}
