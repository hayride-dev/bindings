package transport

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/http/types"
	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/io/streams"
	"go.bytecodealliance.org/cm"
)

type bodyConsumer interface {
	Consume() (result cm.Result[types.IncomingBody, types.IncomingBody, struct{}])
	Headers() (result types.Fields)
}

type inputStreamReader struct {
	consumer    bodyConsumer
	body        *types.IncomingBody
	stream      *streams.InputStream
	trailerLock sync.Mutex
	trailers    http.Header
	trailerOnce sync.Once
}

func (r *inputStreamReader) Close() error {
	r.trailerOnce.Do(r.parseTrailers)

	if r.stream != nil {
		r.stream.ResourceDrop()
	}

	if r.body != nil {
		r.body.ResourceDrop()
		r.body = nil
	}

	return nil
}

func (r *inputStreamReader) parseTrailers() {
	r.trailerLock.Lock()
	defer r.trailerLock.Unlock()

	// if we got this far, then we release ownership from body, otherwise it is our responsibility to drop it
	r.stream.ResourceDrop()
	r.stream = nil

	futureTrailers := types.IncomingBodyFinish(*r.body)
	defer futureTrailers.ResourceDrop()

	trailersResult := futureTrailers.Get()
	r.body = nil

	// unroll the future
	if trailersResult.None() {
		return
	}
	if trailersResult.Some().IsErr() {
		return
	}
	if trailersResult.Some().OK().IsErr() {
		return
	}
	maybeWasiTrailers := trailersResult.Some().OK().OK()

	if maybeWasiTrailers.None() {
		return
	}

	wasiTrailers := maybeWasiTrailers.Some()
	for _, kv := range wasiTrailers.Entries().Slice() {
		r.trailers.Add(string(kv.F0), string(kv.F1.Slice()))
	}

	wasiTrailers.ResourceDrop()
}

func (r *inputStreamReader) Read(p []byte) (n int, err error) {
	pollable := r.stream.Subscribe()
	pollable.Block()
	pollable.ResourceDrop()

	readResult := r.stream.Read(uint64(len(p)))
	if err := readResult.Err(); err != nil {
		if err.Closed() {
			r.trailerOnce.Do(r.parseTrailers)
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to read from InputStream %s", err.LastOperationFailed().ToDebugString())
	}

	readList := *readResult.OK()
	copy(p, readList.Slice())
	return int(readList.Len()), nil
}

// NewIncomingBodyTrailer takes a [BodyConsumer] and parses it into corresponding [io.ReadCloser] and [net/http.Header].
func newIncomingBodyTrailer(consumer bodyConsumer) (io.ReadCloser, http.Header, error) {
	consumeResult := consumer.Consume()
	if consumeResult.IsErr() {
		return nil, nil, errors.New("failed to consume incoming request")
	}

	body := consumeResult.OK()
	streamResult := body.Stream()
	if streamResult.IsErr() {
		return nil, nil, errors.New("failed to consume incoming request body stream")
	}

	stream := streamResult.OK()

	trailers := http.Header{}
	return &inputStreamReader{
		consumer: consumer,
		trailers: trailers,
		body:     body,
		stream:   stream,
	}, trailers, nil
}

type outgoingBody struct {
	body   *types.OutgoingBody
	stream *streams.OutputStream
}

// newOutgoingBody takes a [types.OutgoingBody] and returns a [io.WriteCloser] encapsulating it.
func newOutgoingBody(body *types.OutgoingBody) (io.WriteCloser, error) {
	stream := body.Write()
	if stream.IsErr() {
		return nil, errors.New("failed to acquire resource handle to request body")
	}
	return &outgoingBody{
		body:   body,
		stream: stream.OK(),
	}, nil
}

func (r *outgoingBody) Close() error {
	r.stream.ResourceDrop()
	return nil
}

func (r *outgoingBody) Write(p []byte) (n int, err error) {
	// Send in chunks of max 4096 bytes
	written := uint(0)
	total := uint(len(p))

	for written < total {
		chunkSize := uint(4096)
		if chunkSize > (total - written) {
			chunkSize = total - written
		}

		chunk := p[written : written+chunkSize]
		_, err, isErr := r.stream.BlockingWriteAndFlush(cm.ToList(chunk)).Result()
		if isErr {
			// NOTE: Continue even if the stream is closed
			// Refer to https://github.com/WebAssembly/wasi-io/issues/109 for more details.
			if !err.Closed() {
				return 0, fmt.Errorf("failed to write to response body's stream: %s", err.LastOperationFailed().ToDebugString())
			}
		}

		written += chunkSize
	}

	return len(p), nil
}
