package io

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/io/streams"
	"go.bytecodealliance.org/cm"
)

// WasiWriter is used to convert a io.Writer to a WasiWriter. This is used when
// a io.Writer is passed to a go wrapper function that calls a imported function
// that requires an wasi io output stream. The additional Ptr() method is used to
// get the pointer to the output stream.
type WasiWriter interface {
	io.Writer
	Ptr() uint32
}

func Clone(ptr uint32) WasiWriter {
	outstream := cm.Reinterpret[streams.OutputStream](ptr)
	return &wasiWriter{
		output: outstream,
		ref:    ptr}
}

type wasiWriter struct {
	output streams.OutputStream
	ref    uint32
}

func (w *wasiWriter) Write(p []byte) (n int, err error) {
	contents := cm.ToList(p)
	writeResult := w.output.Write(contents)
	if writeResult.IsErr() {
		if writeResult.Err().Closed() {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to write to stream: %s", writeResult.Err().LastOperationFailed().ToDebugString())
	}
	w.output.BlockingFlush()
	return int(contents.Len()), nil
}

func (w *wasiWriter) Ptr() uint32 {
	return w.ref
}
