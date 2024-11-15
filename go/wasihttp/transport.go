package wasihttp

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	monotonicclock "github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/clocks/monotonic-clock"
	outgoinghandler "github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/outgoing-handler"
	"github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/types"
)

var _ http.RoundTripper = (*Transport)(nil)

type Transport struct {
	ConnectTimeout time.Duration
}

func NewWasiRoundTripper() *Transport {
	return &Transport{
		ConnectTimeout: 30 * time.Second,
	}
}

func (r *Transport) requestOptions() types.RequestOptions {
	options := types.NewRequestOptions()
	options.SetConnectTimeout(cm.Some(monotonicclock.Duration(r.ConnectTimeout)))
	return options
}

// RoundTrip implements the [net/http.RoundTripper] interface.
func (r *Transport) RoundTrip(incomingRequest *http.Request) (*http.Response, error) {
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
		fmt.Println("Acquiring body")
		bodyRes := outRequest.Body()
		if bodyRes.IsErr() {
			return nil, fmt.Errorf("failed to acquire resource handle to request body: %s", bodyRes.Err())
		}
		body = bodyRes.OK()
		adaptedBody, err = NewOutgoingBody(body)
		if err != nil {
			return nil, fmt.Errorf("failed to adapt body: %s", err)
		}
	}

	handleResp := outgoinghandler.Handle(outRequest, cm.Some(r.requestOptions()))
	if handleResp.Err() != nil {
		return nil, fmt.Errorf("%v", handleResp.Err())
	}

	if body != nil {
		fmt.Println("Copying body")
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
		fmt.Println("Finished body")
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
	incomingBody, incomingTrailers, err := NewIncomingBodyTrailer(incomingResponse)
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
