package morph

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/morph/gen/wasi/io/streams"
)

type morphReadCloser struct {
	stream streams.InputStream
}

// create an io.Reader from the input stream
func NewReader(s streams.InputStream) io.Reader {
	return &morphReadCloser{
		stream: s,
	}
}

func NewReadCloser(s streams.InputStream) io.ReadCloser {
	return &morphReadCloser{
		stream: s,
	}
}

func (w *morphReadCloser) Read(p []byte) (int, error) {
	readResult := w.stream.Read(uint64(len(p)))
	if readResult.IsErr() {
		readErr := readResult.Err()
		if readErr.Closed() {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to read from InputStream %s", readErr.LastOperationFailed().ToDebugString())
	}

	readList := readResult.OK()
	copy(p, readList.Slice())
	return int(len(p)), nil
}

func (w *morphReadCloser) Close() error {
	w.stream.ResourceDrop()
	return nil
}
