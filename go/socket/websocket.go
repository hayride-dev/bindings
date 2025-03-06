package socket

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/io"
	"github.com/hayride-dev/bindings/go/socket/gen/hayride/socket/websocket"
)

func init() {
	websocket.Exports.Handle = websocketHandle
}

type Handler interface {
	Handle(msg string, writer io.Writer)
}

type defaulthandler struct{}

func (dh *defaulthandler) Handle(msg string, writer io.Writer) {
	writer.Write([]byte("websocket handler undefined"))
}

var handler Handler = &defaulthandler{}

func Handle(h Handler) {
	handler = h
}

func websocketHandle(text string, out websocket.OutputStream) {
	w := wasiio.NewWriter(uint32(out))
	handler.Handle(text, w)
}
