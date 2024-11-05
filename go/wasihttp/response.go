package wasihttp

import (
	"fmt"
	"net/http"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/types"
	"github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/io/streams"
)

var _ http.ResponseWriter = &wasiResponseWriter{}

type wasiResponseWriter struct {
	outparam    types.ResponseOutparam
	response    types.OutgoingResponse
	wasiHeaders types.Fields
	httpHeaders http.Header
	body        *types.OutgoingBody
	stream      *streams.OutputStream
	statuscode  int
}

func newWasiResponseWriter(out types.ResponseOutparam) *wasiResponseWriter {
	return &wasiResponseWriter{
		outparam:    out,
		httpHeaders: http.Header{},
		wasiHeaders: types.NewFields(),
	}
}

func (w *wasiResponseWriter) Header() http.Header {
	return w.httpHeaders
}

func (w *wasiResponseWriter) Write(buf []byte) (int, error) {
	if w.body == nil {
		bodyResult := w.response.Body()
		if bodyResult.IsErr() {
			return 0, fmt.Errorf("failed to acquire resource handle to response body: %s", bodyResult.Err())
		}
		w.body = bodyResult.OK()

		writeResult := w.body.Write()
		if writeResult.IsErr() {
			return 0, fmt.Errorf("failed to acquire resource handle for response body's stream: %s", writeResult.Err())
		}
		w.stream = writeResult.OK()
	}

	contents := cm.ToList(buf)
	writeResult := w.stream.Write(contents)

	if writeResult.IsErr() {
		if writeResult.Err().Closed() {
			return 0, fmt.Errorf("failed to write to response body's stream: closed")
		}
		return 0, fmt.Errorf("failed to write to response body's stream: %s", writeResult.Err().LastOperationFailed().ToDebugString())
	}

	result := cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](w.response)
	types.ResponseOutparamSet(w.outparam, result)

	return int(contents.Len()), nil
}

func (w *wasiResponseWriter) WriteHeader(statusCode int) {
	w.statuscode = statusCode
	w.reconcile()
}

func (w *wasiResponseWriter) reconcile() {
	w.wasiHeaders = headerToWASIHeader(w.httpHeaders)

	//setting headers after this cause panic
	w.response = types.NewOutgoingResponse(w.wasiHeaders)

	//set status code
	w.response.SetStatusCode(types.StatusCode(w.statuscode))
}
