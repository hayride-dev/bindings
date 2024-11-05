package wasihttp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	incominghandler "github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/incoming-handler"
	"github.com/hayride-dev/bindings/go/wasihttp/gen/wasi/http/types"
)

func init() {
	incominghandler.Exports.Handle = wasiHandle
}

type defaulthandler struct{}

func (dh *defaulthandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	fmt.Fprintln(os.Stderr, "http handler undefined")
}

var handler http.Handler = &defaulthandler{}

func Handle(h http.Handler) {
	handler = h
}

func wasiHandle(request types.IncomingRequest, responseOut types.ResponseOutparam) {
	// construct the http.Request and http.ResponseWriter from wasi types
	req, err := requestFromWASIIncomingRequest(request)
	if err != nil {
		fmt.Printf("failed to convert wasi/http/types.IncomingRequest to http.Request: %s\n", err)
		return
	}
	resp := newWasiResponseWriter(responseOut)
	// call the go handler
	handler.ServeHTTP(resp, req)
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

func wasiMethodToString(m types.Method) (string, error) {
	if m.Connect() {
		return http.MethodConnect, nil
	} else if m.Delete() {
		return http.MethodDelete, nil
	} else if m.Get() {
		return http.MethodGet, nil
	} else if m.Head() {
		return http.MethodHead, nil
	} else if m.Options() {
		return http.MethodOptions, nil
	} else if m.Patch() {
		return http.MethodPatch, nil
	} else if m.Post() {
		return http.MethodPost, nil
	} else if m.Put() {
		return http.MethodPut, nil
	} else if m.Trace() {
		return http.MethodTrace, nil
	} else if other := m.Other(); other != nil {
		return *other, fmt.Errorf("unknown http method '%s'", *other)
	}
	return "", fmt.Errorf("failed to convert http method")
}

func schemeToWASIScheme(scheme string) types.Scheme {
	switch scheme {
	case "http":
		return types.SchemeHTTP()
	case "https":
		return types.SchemeHTTPS()
	default:
		return types.SchemeOther(scheme)
	}
}
