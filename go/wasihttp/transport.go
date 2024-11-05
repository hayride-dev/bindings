package wasihttp

import (
	"fmt"
	"net/http"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	outgoinghandler "github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/outgoing-handler"
	"github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/types"
)

var _ http.RoundTripper = &RoundTrip{}

type RoundTrip struct {
}

func NewWasiRoundTripper() *RoundTrip {
	return &RoundTrip{}
}

func (r *RoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	wasiRequest, err := wasiOutGoingFromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert http.Request to wasi/http/types.OutgoingRequest: %s", err)
	}

	result := outgoinghandler.Handle(wasiRequest, cm.None[types.RequestOptions]())
	if result.IsErr() {
		return nil, fmt.Errorf("error %v", result.Err())
	}

	if result.IsOK() {
		result.OK().Subscribe().Block()
		incomingResponse := result.OK().Get()
		if incomingResponse.Some().IsErr() {
			return nil, fmt.Errorf("error %v", incomingResponse.Some().Err())
		}
		if incomingResponse.Some().OK().IsErr() {
			return nil, fmt.Errorf("error %v", incomingResponse.Some().OK().Err())
		}
		ok := incomingResponse.Some().OK().OK()
		consume := ok.Consume()
		if consume.IsErr() {
			return nil, fmt.Errorf("error %v", consume.Err())
		}

		stream := consume.OK().Stream()
		if stream.IsErr() {
			return nil, fmt.Errorf("error %v", stream.Err())
		}

		body := newReadCloser(*stream.OK())

		response := &http.Response{
			StatusCode:    int(ok.Status()),
			Status:        http.StatusText(int(ok.Status())),
			ContentLength: 0,
			Body:          body,
			Request:       req,
		}
		return response, nil
	}
	return nil, fmt.Errorf("failed to get response")
}
