package streams

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/wasi/io/streams"
	"go.bytecodealliance.org/cm"
)

var _ io.Writer = Writer(0)

type Writer cm.Resource

func (w Writer) Write(p []byte) (n int, err error) {
	resource := cm.Reinterpret[streams.OutputStream](w)
	contents := cm.ToList(p)
	resource.Subscribe().Block()
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
