package server

import (
	"fmt"
	"net/http"

	"github.com/hayride-dev/bindings/go/exports/net/http/handle"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/config"
	"go.bytecodealliance.org/cm"
)

type Config = config.Server

var cfg = config.Server{
	Address:        "localhost:80",
	MaxHeaderBytes: http.DefaultMaxHeaderBytes,
}

func init() {
	config.Exports.Get = get
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
		wasiErr := config.ErrorResourceNew(cm.Rep(config.ErrorCodeInvalid))
		return cm.Err[cm.Result[config.ServerShape, config.Server, config.Error]](wasiErr)
	}
	return cm.OK[cm.Result[config.ServerShape, config.Server, config.Error]](cfg)
}
