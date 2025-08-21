package runner

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/hayride/types"
)

var _ io.Writer = &Writer{}

type Writer struct {
	writerType types.WriterType
	w          io.Writer
}

func NewWriter(writerType types.WriterType, w io.Writer) *Writer {
	return &Writer{
		writerType: writerType,
		w:          w,
	}
}

func (mw *Writer) Write(p []byte) (n int, err error) {
	switch mw.writerType {
	case types.WriterTypeSse:
		content := []byte(fmt.Sprintf("data: %s\n\n", p))
		if _, err := mw.w.Write(content); err != nil {
			return 0, fmt.Errorf("failed to write SSE message: %w", err)
		}
		return len(p), nil
	case types.WriterTypeRaw:
		if _, err := mw.w.Write(p); err != nil {
			return 0, fmt.Errorf("failed to write raw message: %w", err)
		}
		return len(p), nil
	default:
		return 0, fmt.Errorf("unknown writer type: %s", mw.writerType)
	}
}
