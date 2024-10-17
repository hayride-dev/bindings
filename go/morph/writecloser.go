package morph

import (
	"fmt"
	"io"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/morph/gen/wasi/io/streams"
)

type morphWriteCloser struct {
	stream streams.OutputStream
}

func newWriter(s streams.OutputStream) io.Writer {
	return &morphWriteCloser{
		stream: s,
	}
}

func newWriterCloser(s streams.OutputStream) io.WriteCloser {
	return &morphWriteCloser{
		stream: s,
	}
}

func (w *morphWriteCloser) Write(p []byte) (int, error) {
	content := cm.ToList(p)
	// TODO:: eval if should we do a blocking write and flush
	writeResult := w.stream.Write(content)
	if writeResult.IsErr() {
		writeErr := writeResult.Err()
		if writeErr.Closed() {
			return 0, io.ErrClosedPipe
		}
		return 0, fmt.Errorf("failed to write to OutputStream %s", writeErr.LastOperationFailed().ToDebugString())
	}
	// after writing flush the stream
	flushResult := w.stream.Flush()
	if flushResult.IsErr() {
		flushErr := flushResult.Err()
		if flushErr.Closed() {
			return 0, io.ErrClosedPipe
		}
		return 0, fmt.Errorf("failed to flush OutputStream %s", flushErr.LastOperationFailed().ToDebugString())
	}
	return int(len(p)), nil
}

func (w *morphWriteCloser) Close() error {
	w.stream.ResourceDrop()
	return nil
}
