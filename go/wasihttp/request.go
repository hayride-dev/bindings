package wasihttp

import (
	"fmt"
	"io"
	"net/http"

	"github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/types"
)

func requestFromWASIIncomingRequest(incoming types.IncomingRequest) (*http.Request, error) {
	method, err := wasiMethodToString(incoming.Method())
	if err != nil {
		return nil, fmt.Errorf("failed to convert wasi/http/types.Method to string: %s", err)
	}

	var url string
	if pathWithQuery := incoming.PathWithQuery(); pathWithQuery.None() {
		url = ""
	} else {
		url = *pathWithQuery.Some()
	}

	var body io.Reader
	if consumeResult := incoming.Consume(); consumeResult.IsErr() {
		return nil, fmt.Errorf("failed to consume incoming request %s", *consumeResult.Err())
	} else if streamResult := consumeResult.OK().Stream(); streamResult.IsErr() {
		return nil, fmt.Errorf("failed to consume incoming request stream %s", streamResult.Err())
	} else {
		body = newReadCloser(*streamResult.OK())
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = wasiHeadertoHeader(incoming.Headers())

	return req, nil
}

// WASItoHTTPRequest takes an [IncomingRequest] and returns a [net/http.Request] representation of it.
func WASItoHTTPRequest(incoming types.IncomingRequest) (req *http.Request, err error) {
	method, err := wasiMethodToString(incoming.Method())
	if err != nil {
		return nil, err
	}

	authority := "localhost"
	if auth := incoming.Authority(); !auth.None() {
		authority = *auth.Some()
	}

	pathWithQuery := "/"
	if p := incoming.PathWithQuery(); !p.None() {
		pathWithQuery = *p.Some()
	}

	body, trailers, err := NewIncomingBodyTrailer(incoming)
	if err != nil {
		switch method {
		case http.MethodGet,
			http.MethodHead,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace:
		default:
			return nil, fmt.Errorf("failed to consume incoming request: %w", err)
		}
	}

	url := fmt.Sprintf("http://%s%s", authority, pathWithQuery)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Trailer = trailers

	headers := incoming.Headers()

	req.Header = wasiHeadertoHeader(headers)
	headers.ResourceDrop()

	req.Host = authority
	req.URL.Host = authority
	req.RequestURI = pathWithQuery

	return req, nil
}
