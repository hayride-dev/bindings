package morph

import (
	"context"
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/morph/gen/hayride/morph/spawn"
	"github.com/hayride-dev/bindings/go/morph/gen/wasi/io/streams"
)

type morph struct {
}

func Spawn() *morph {
	return nil
}

func (m *morph) Execute(ctx context.Context, name string, args []byte) (io.ReadCloser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var output streams.OutputStream
		w := newWriterCloser(output)
		_, err := w.Write(args)
		if err != nil {
			return nil, err
		}

		result := spawn.Exec(name, output)
		if result.IsErr() {
			return nil, fmt.Errorf("failed to read from InputStream %s", result.Err().ToDebugString())
		}
		stream := *result.OK()
		return NewReadCloser(stream), nil
	}
}
