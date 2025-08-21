package runner

import (
	"bufio"
	"errors"
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
		var dataLines []string

		for {
			line, err := r.r.ReadString('\n') // includes '\n' if found
			if err != nil {
				// If EOF and we have buffered data, flush it as a final event
				if errors.Is(err, io.EOF) && len(dataLines) > 0 {
					return []byte(strings.Join(dataLines, "\n")), nil
				}
				return nil, err
			}

			// Trim CRLF (SSE commonly uses CRLF)
			line = strings.TrimRight(line, "\r\n")

			// Blank line => end of event frame
			if line == "" {
				if len(dataLines) == 0 {
					// Skip stray keep-alive blank lines
					continue
				}
				return []byte(strings.Join(dataLines, "\n")), nil
			}

			// Comments (ignore)
			if strings.HasPrefix(line, ":") {
				continue
			}

			// Parse field: "name: value" or "name:value"
			// We only care about data fields; ignore others.
			if i := strings.IndexByte(line, ':'); i >= 0 {
				field := line[:i]
				value := line[i+1:]
				// If there is a single leading space after ":", strip it (per spec).
				if len(value) > 0 && value[0] == ' ' {
					value = value[1:]
				}
				if field == "data" {
					dataLines = append(dataLines, value)
				}
				// ignore id/event/retry
				continue
			}

			// Field without ":" (non-standard)
		}

	case types.WriterTypeRaw:
		line, err := r.r.ReadString('\n')
		if err != nil {
			// If EOF and we got some bytes, return them (without forcing an error)
			if errors.Is(err, io.EOF) && len(line) > 0 {
				return []byte(strings.TrimRight(line, "\r\n")), nil
			}
			return nil, err
		}
		return []byte(strings.TrimRight(line, "\r\n")), nil

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
