package socket

import (
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/socket/websocket"
	"github.com/hayride-dev/bindings/go/wasi/streams"
)

func init() {
	websocket.Exports.Handle = websocketHandle
}

type Handler interface {
	Handle(reader io.ReadCloser, writer io.Writer)
}

type defaulthandler struct{}

func (dh *defaulthandler) Handle(reader io.ReadCloser, writer io.Writer) {
	writer.Write([]byte("websocket handler undefined"))
}

var handler Handler = &defaulthandler{}

func WebSocketHandler(h Handler) {
	handler = h
}

func websocketHandle(input websocket.InputStream, output websocket.OutputStream) {
	w := streams.Writer(output)
	reader := streams.ReaderCloser(input)
	handler.Handle(reader, w)
}
