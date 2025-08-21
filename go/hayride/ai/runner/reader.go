package runner

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/hayride-dev/bindings/go/hayride/types"
)

var _ io.Reader = &Reader{}

type Reader struct {
	writerType types.WriterType
	r          *bufio.Reader
}

func NewReader(writerType types.WriterType, r io.Reader) *Reader {
	return &Reader{
		writerType: writerType,
		r:          bufio.NewReader(r),
	}
}

// Read implements io.Reader interface - always reads raw bytes
func (r *Reader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// ReadLine reads a single line and handles SSE parsing if in SSE mode
func (r *Reader) ReadLine() ([]byte, error) {
	switch r.writerType {
	case types.WriterTypeSse:
		// For SSE mode, read one data line at a time
		for {
			line, err := r.r.ReadString('\n')
			if err != nil {
				return nil, err
			}

			// Remove only the trailing newline added by ReadString
			if len(line) > 0 && line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			// Also remove \r if present (for Windows line endings)
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}

			// Skip empty lines (SSE message separators)
			if line == "" {
				continue
			}

			// Parse SSE format: "data: content"
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				return []byte(data), nil
			}

			// Skip non-data lines (like comments, event types, etc.)
		}

	case types.WriterTypeRaw:
		// For raw mode, just read a single line
		line, err := r.r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		// Remove only the trailing newline added by ReadString
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		return []byte(line), nil

	default:
		return nil, fmt.Errorf("unknown writer type: %s", r.writerType)
	}
}

// Buffered returns the number of bytes that can be read from the current buffer
func (r *Reader) Buffered() int {
	return r.r.Buffered()
}

// Peek returns the next n bytes without advancing the reader
func (r *Reader) Peek(n int) ([]byte, error) {
	return r.r.Peek(n)
}
