package streams

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/wasi/io/streams"
	"go.bytecodealliance.org/cm"
)

var _ io.ReadCloser = ReaderCloser(0)

type ReaderCloser cm.Resource

func (r ReaderCloser) Read(p []byte) (n int, err error) {
	resource := cm.Reinterpret[streams.InputStream](r)

	resource.Subscribe().Block()
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

func (r ReaderCloser) Close() error {
	resource := cm.Reinterpret[streams.InputStream](r)
	resource.ResourceDrop()
	return nil
}
