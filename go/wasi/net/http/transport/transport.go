package transport

import (
	"fmt"
	"io"
	"net/http"
	"time"

	monotonicclock "github.com/hayride-dev/bindings/go/internal/gen/wasi/clocks/monotonic-clock"
	outgoinghandler "github.com/hayride-dev/bindings/go/internal/gen/wasi/http/outgoing-handler"
	"github.com/hayride-dev/bindings/go/internal/gen/wasi/http/types"
	"go.bytecodealliance.org/cm"
)

var _ http.RoundTripper = (*transport)(nil)

type transport struct {
	ConnectTimeout time.Duration
}

func New() http.RoundTripper {
	return &transport{
		ConnectTimeout: 30 * time.Second,
	}
}

func (r *transport) requestOptions() types.RequestOptions {
	options := types.NewRequestOptions()
	options.SetConnectTimeout(cm.Some(monotonicclock.Duration(r.ConnectTimeout)))
	return options
}

// RoundTrip implements the [net/http.RoundTripper] interface.
func (r *transport) RoundTrip(incomingRequest *http.Request) (*http.Response, error) {
	var err error

	outHeaders := headerToWASIHeader(incomingRequest.Header)

	outRequest := types.NewOutgoingRequest(outHeaders)

	outRequest.SetAuthority(cm.Some(incomingRequest.Host))
	outRequest.SetMethod(methodToWASIMethod(incomingRequest.Method))

	pathWithQuery := incomingRequest.URL.Path
	if incomingRequest.URL.RawQuery != "" {
		pathWithQuery = pathWithQuery + "?" + incomingRequest.URL.Query().Encode()
	}
	outRequest.SetPathWithQuery(cm.Some(pathWithQuery))

	switch incomingRequest.URL.Scheme {
	case "http":
		outRequest.SetScheme(cm.Some(types.SchemeHTTP()))
	case "https":
		outRequest.SetScheme(cm.Some(types.SchemeHTTPS()))
	default:
		outRequest.SetScheme(cm.Some(types.SchemeOther(incomingRequest.URL.Scheme)))
	}

	var adaptedBody io.WriteCloser
	var body *types.OutgoingBody
	if incomingRequest.Body != nil {
		bodyRes := outRequest.Body()
		if bodyRes.IsErr() {
			return nil, fmt.Errorf("failed to acquire resource handle to request body: %s", bodyRes.Err())
		}
		body = bodyRes.OK()
		adaptedBody, err = newOutgoingBody(body)
		if err != nil {
			return nil, fmt.Errorf("failed to adapt body: %s", err)
		}
	}

	if body != nil {
		if _, err := io.Copy(adaptedBody, incomingRequest.Body); err != nil {
			return nil, fmt.Errorf("failed to copy body: %v", err)
		}

		if err := adaptedBody.Close(); err != nil {
			return nil, fmt.Errorf("failed to close body: %v", err)
		}

		outTrailers := headerToWASIHeader(incomingRequest.Trailer)

		maybeTrailers := cm.None[types.Fields]()
		if len(incomingRequest.Trailer) > 0 {
			maybeTrailers = cm.Some(outTrailers)
		}

		outFinish := types.OutgoingBodyFinish(*body, maybeTrailers)
		if outFinish.IsErr() {
			return nil, fmt.Errorf("failed to finish body: %v", outFinish.Err())
		}
	}

	handleResp := outgoinghandler.Handle(outRequest, cm.Some(r.requestOptions()))
	if handleResp.Err() != nil {
		return nil, fmt.Errorf("%v", handleResp.Err())
	}

	futureResponse := handleResp.OK()

	// wait until resp is returned
	futureResponse.Subscribe().Block()

	pollableOption := futureResponse.Get()
	if pollableOption.None() {
		return nil, fmt.Errorf("incoming resp is None")
	}

	pollableResult := pollableOption.Some()
	if pollableResult.IsErr() {
		return nil, fmt.Errorf("error is %v", pollableResult.Err())
	}

	resultOption := pollableResult.OK()
	if resultOption.IsErr() {
		return nil, fmt.Errorf("%v", resultOption.Err())
	}

	incomingResponse := resultOption.OK()
	incomingBody, incomingTrailers, err := newIncomingBodyTrailer(incomingResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to consume incoming request %s", err)
	}

	headers := incomingResponse.Headers()
	incomingHeaders := wasiHeadertoHeader(headers)
	headers.ResourceDrop()

	resp := &http.Response{
		StatusCode: int(incomingResponse.Status()),
		Status:     http.StatusText(int(incomingResponse.Status())),
		Request:    incomingRequest,
		Header:     incomingHeaders,
		Body:       incomingBody,
		Trailer:    incomingTrailers,
	}

	return resp, nil
}

func headerToWASIHeader(headers http.Header) types.Fields {
	fields := types.NewFields()
	for key, values := range headers {
		fieldValues := []types.FieldValue{}
		for _, v := range values {
			fieldValues = append(fieldValues, types.FieldValue(cm.ToList([]uint8(v))))
		}
		fields.Set(types.FieldKey(key), cm.ToList(fieldValues))
	}
	return fields
}

func wasiHeadertoHeader(fields types.Fields) http.Header {
	headers := http.Header{}
	for _, f := range fields.Entries().Slice() {
		key := string(f.F0)
		value := string(cm.List[uint8](f.F1).Slice())
		headers.Add(key, value)
	}
	return headers
}

func methodToWASIMethod(method string) types.Method {
	switch method {
	case "GET":
		return types.MethodGet()
	case "POST":
		return types.MethodPost()
	case "PUT":
		return types.MethodPut()
	case "DELETE":
		return types.MethodDelete()
	case "PATCH":
		return types.MethodPatch()
	case "HEAD":
		return types.MethodHead()
	case "OPTIONS":
		return types.MethodOptions()
	case "TRACE":
		return types.MethodTrace()
	case "CONNECT":
		return types.MethodConnect()
	default:
		return types.MethodOther(method)
	}
}
