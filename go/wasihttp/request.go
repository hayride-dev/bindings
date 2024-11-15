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
