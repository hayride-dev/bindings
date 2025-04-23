package socket

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/socket/websocket"
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

func Handle(h Handler) {
	handler = h
}

func websocketHandle(input websocket.InputStream, output websocket.OutputStream) {
	w := wasiio.Clone(uint32(output))
	reader := wasiio.CloneReader(uint32(input))
	handler.Handle(reader, w)
}
