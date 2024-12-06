package socket

import (
	"fmt"
	"io"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/socket/gen/hayride/socket/websocket"
	"github.com/hayride-dev/bindings/go/socket/gen/wasi/io/streams"
)

func init() {
	websocket.Exports.Handle = websocketHandle
}

type Handler interface {
	Handle(msg string, writer Writer)
}

type Writer interface {
	Write(buf []byte) (int, error)
}

type websocketWriter struct {
	stream *streams.OutputStream
}

func (w *websocketWriter) Write(buf []byte) (int, error) {
	contents := cm.ToList(buf)
	writeResult := w.stream.Write(contents)
	if writeResult.IsErr() {
		if writeResult.Err().Closed() {
			return 0, io.EOF
		}

		return 0, fmt.Errorf("failed to write to response body's stream: %s", writeResult.Err().LastOperationFailed().ToDebugString())
	}
	w.stream.BlockingFlush()
	return int(contents.Len()), nil
}

type defaulthandler struct{}

func (dh *defaulthandler) Handle(msg string, writer Writer) {
	writer.Write([]byte("websocket handler undefined"))
}

var handler Handler = &defaulthandler{}

func Handle(h Handler) {
	handler = h
}

func websocketHandle(text string, out websocket.OutputStream) {
	writer := &websocketWriter{stream: &out}
	handler.Handle(text, writer)
}
