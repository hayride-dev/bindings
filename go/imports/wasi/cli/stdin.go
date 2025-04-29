package cli

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/cli/stdin"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/io/streams"
	"go.bytecodealliance.org/cm"

	wasiio "github.com/hayride-dev/bindings/go/imports/wasi/io"
)

// GetStdin returns a blocking reader for stdin.
func GetStdin() io.ReadCloser {
	r := wasiio.ReaderCloser(stdin.GetStdin())
	return r
}

// GetNonBlockingStdin returns a non-blocking reader for stdin.
func GetNonBlockingStdin() io.ReadCloser {
	return NonBlockingReader(stdin.GetStdin())
}

type NonBlockingReader cm.Resource

func (r NonBlockingReader) Read(p []byte) (n int, err error) {
	resource := cm.Reinterpret[streams.InputStream](r)

	readResult := resource.Read(uint64(len(p)))
	if readResult.IsErr() {
		readErr := readResult.Err()
		if readErr.Closed() {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to read from InputStream %s", readErr.LastOperationFailed().ToDebugString())
	}

	data := readResult.OK().Slice()
	copy(p, data)
	return int(len(data)), nil
}

func (r NonBlockingReader) Close() error {
	resource := cm.Reinterpret[streams.InputStream](r)
	resource.ResourceDrop()
	return nil
}
