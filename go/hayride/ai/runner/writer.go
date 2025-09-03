package runner

import (
	"fmt"
	"io"
	"strings"

	"github.com/hayride-dev/bindings/go/hayride/ai"
)

var _ io.Writer = &Writer{}

type Writer struct {
	writerType ai.WriterType
	w          io.Writer
}

func NewWriter(writerType ai.WriterType, w io.Writer) *Writer {
	return &Writer{
		writerType: writerType,
		w:          w,
	}
}

func (mw *Writer) Write(p []byte) (int, error) {
	switch mw.writerType {
	case ai.WriterTypeSse:
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

		return len(p), nil
	case ai.WriterTypeRaw:
		if _, err := mw.w.Write(p); err != nil {
			return 0, fmt.Errorf("failed to write raw message: %w", err)
		}
		return len(p), nil

	default:
		return 0, fmt.Errorf("unknown writer type: %s", mw.writerType)
	}
}
