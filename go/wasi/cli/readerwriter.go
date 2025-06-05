package cli

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/wasi/io/streams"
	"go.bytecodealliance.org/cm"
)

type nonBlockingReader cm.Resource

func (r nonBlockingReader) Read(p []byte) (n int, err error) {
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

func (r nonBlockingReader) Close() error {
	resource := cm.Reinterpret[streams.InputStream](r)
	resource.ResourceDrop()
	return nil
}

type nonBlockingWriter cm.Resource

func (w nonBlockingWriter) Write(p []byte) (n int, err error) {
	resource := cm.Reinterpret[streams.OutputStream](w)
	contents := cm.ToList(p)
	writeResult := resource.Write(contents)
	if writeResult.IsErr() {
		if writeResult.Err().Closed() {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to write to stream: %s", writeResult.Err().LastOperationFailed().ToDebugString())
	}
	resource.BlockingFlush()
	return int(contents.Len()), nil
}
