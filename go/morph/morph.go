package morph

import (
	"context"
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/morph/gen/hayride/morph/spawn"
)

type morph struct {
}

func Spawn() *morph {
	return &morph{}
}

func (m *morph) Execute(ctx context.Context, name string, args []byte) (io.ReadCloser, error) {
	// one-shot context check
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	writer := spawn.NewWriter()
	stream := writer.Stream()
	if stream.IsErr() {
		return nil, fmt.Errorf("failed to create OutputStream: %s", stream.Err().ToDebugString())
	}

	w := newWriterCloser(*stream.OK())
	defer w.Close()

	if _, err := w.Write(args); err != nil {
		return nil, err
	}

	result := spawn.Exec(name, *stream.OK())
	if result.IsErr() {
		return nil, fmt.Errorf("failed to read from InputStream: %s", result.Err().ToDebugString())
	}

	return NewReadCloser(*result.OK()), nil
}
