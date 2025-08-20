package handle

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/http/types"
	wasistreams "github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/io/streams"
	"github.com/hayride-dev/bindings/go/wasi/streams"
	"go.bytecodealliance.org/cm"
)

var _ http.ResponseWriter = &WasiResponseWriter{}

type WasiResponseWriter struct {
	streams.Writer

	outparam    types.ResponseOutparam
	response    types.OutgoingResponse
	wasiHeaders types.Fields
	httpHeaders http.Header
	body        *types.OutgoingBody

	headerOnce sync.Once
	headerErr  error

	statuscode int
}

func newWasiResponseWriter(out types.ResponseOutparam) *WasiResponseWriter {
	return &WasiResponseWriter{
		outparam:    out,
		httpHeaders: http.Header{},
		wasiHeaders: types.NewFields(),
		statuscode:  http.StatusOK,
	}
}

func (w *WasiResponseWriter) Header() http.Header {
	return w.httpHeaders
}

func (w *WasiResponseWriter) Write(buf []byte) (int, error) {
	// NOTE: If this is the first write, make sure we set the headers/statuscode
	w.headerOnce.Do(w.reconcile)
	if w.headerErr != nil {
		return 0, w.headerErr
	}

	return w.Writer.Write(buf)
}

func (w *WasiResponseWriter) WriteHeader(statusCode int) {
	w.headerOnce.Do(func() {
		w.statuscode = statusCode
		w.reconcile()
	})
}

func (w *WasiResponseWriter) reconcile() {
	for key, vals := range w.httpHeaders {
		fieldVals := []types.FieldValue{}
		for _, val := range vals {
			fieldVals = append(fieldVals, types.FieldValue(cm.ToList([]uint8(val))))
		}
		w.wasiHeaders.Set(types.FieldKey(key), cm.ToList(fieldVals))
	}

	// NOTE: once headers are written we clear them out so they can emit http trailers
	w.httpHeaders = http.Header{}

	w.response = types.NewOutgoingResponse(w.wasiHeaders)
	w.response.SetStatusCode(types.StatusCode(w.statuscode))

	bodyResult := w.response.Body()
	if bodyResult.IsErr() {
		w.headerErr = fmt.Errorf("failed to acquire resource handle to response body: %s", bodyResult.Err())
		return
	}
	w.body = bodyResult.OK()

	writeResult := w.body.Write()
	if writeResult.IsErr() {
		w.headerErr = fmt.Errorf("failed to acquire resource handle for response body's stream: %s", writeResult.Err())
		return
	}
	w.Writer = streams.Writer(*writeResult.OK())
	result := cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](w.response)
	types.ResponseOutparamSet(w.outparam, result)
}

// Close closes out the underlying stream by flushing the response and making
// sure that the underlying resource handle is dropped.
func (w *WasiResponseWriter) Close() error {
	if w.Writer != cm.ResourceNone {
		return nil
	}

	stream := cm.Reinterpret[wasistreams.OutputStream](w.Writer)

	stream.BlockingFlush()
	stream.ResourceDrop()

	w.Writer = cm.ResourceNone

	var maybeTrailers cm.Option[types.Fields]
	wasiTrailers := types.NewFields()
	for key, vals := range w.httpHeaders {
		fieldVals := []types.FieldValue{}
		for _, val := range vals {
			fieldVals = append(fieldVals, types.FieldValue(cm.ToList([]uint8(val))))
		}

		if result := wasiTrailers.Set(types.FieldKey(key), cm.ToList(fieldVals)); result.IsErr() {
			return fmt.Errorf("failed to set trailer %s: %s", key, result.Err())
		}
	}
	if len(w.httpHeaders) > 0 {
		maybeTrailers = cm.Some(wasiTrailers)
	} else {
		maybeTrailers = cm.None[types.Fields]()
	}

	res := types.OutgoingBodyFinish(*w.body, maybeTrailers)
	if res.IsErr() {
		return fmt.Errorf("failed to set trailer: %v", res.Err())
	}
	return nil
}
