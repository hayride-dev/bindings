package handle

import (
	"fmt"
	"net/http"

	"github.com/hayride-dev/bindings/go/gen/exports/wasi/http/types"
)

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

	body, trailers, err := newIncomingBodyTrailer(incoming)
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
