package runner

import (
	"fmt"
	"io"
	"strings"

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

func (mw *Writer) Write(p []byte) (int, error) {
	switch mw.writerType {
	case types.WriterTypeSse:
		// SSE is text; normalize line endings to \n
		s := strings.ReplaceAll(string(p), "\r\n", "\n")

		// Split into lines; each becomes its own "data:" line.
		lines := strings.Split(s, "\n")
		for _, ln := range lines {
			// Note: per spec, either "data:<v>" or "data: <v>" is fine.
			if _, err := io.WriteString(mw.w, "data: "); err != nil {
				return 0, fmt.Errorf("failed to write SSE prefix: %w", err)
			}
			if _, err := io.WriteString(mw.w, ln); err != nil {
				return 0, fmt.Errorf("failed to write SSE line: %w", err)
			}
			if _, err := io.WriteString(mw.w, "\n"); err != nil {
				return 0, fmt.Errorf("failed to write SSE newline: %w", err)
			}
		}

		// Blank line terminates the event frame.
		if _, err := io.WriteString(mw.w, "\n"); err != nil {
			return 0, fmt.Errorf("failed to write SSE terminator: %w", err)
		}

		// Flush if supported to push the event immediately.
		if f, ok := mw.w.(interface{ Flush() }); ok {
			f.Flush()
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
