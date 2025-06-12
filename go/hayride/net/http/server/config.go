package server

import (
	"fmt"
	"net/http"

	"github.com/hayride-dev/bindings/go/internal/gen/hayride/http/config"
	"github.com/hayride-dev/bindings/go/wasi/net/http/handle"
	"go.bytecodealliance.org/cm"
)

type Config = config.Server

var cfg = config.Server{
	Address:        "localhost:80",
	MaxHeaderBytes: http.DefaultMaxHeaderBytes,
}

type configError struct {
	code config.ErrorCode
	data string
}

var errors = make(map[cm.Rep]configError)
var errorCount cm.Rep = 0

func init() {
	config.Exports.Get = get

	config.Exports.Error.Code = code
	config.Exports.Error.Data = data
	config.Exports.Error.Destructor = errorDestructor
}

func Export(h http.Handler, c Config) error {
	handle.Handler(h)
	if c.Address == "" {
		return fmt.Errorf("invalid address: %s", c.Address)
	}
	cfg = c
	return nil
}

func get() (result cm.Result[config.ServerShape, config.Server, config.Error]) {
	if cfg.Address == "" {
		err := configError{
			code: config.ErrorCodeInvalid,
			data: "address cannot be empty",
		}

		// Store the error and return a reference to it
		errors[errorCount] = err
		wasiErr := config.ErrorResourceNew(errorCount)
		errorCount++

		return cm.Err[cm.Result[config.ServerShape, config.Server, config.Error]](wasiErr)
	}
	return cm.OK[cm.Result[config.ServerShape, config.Server, config.Error]](cfg)
}

func code(self cm.Rep) (result config.ErrorCode) {
	err, ok := errors[self]
	if !ok {
		return config.ErrorCodeUnknown
	}
	return err.code
}

func data(self cm.Rep) (result string) {
	err, ok := errors[self]
	if !ok {
		return ""
	}
	return err.data
}

func errorDestructor(self cm.Rep) {
	delete(errors, self)
}
