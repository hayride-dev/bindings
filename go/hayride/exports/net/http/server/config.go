package server

import (
	"fmt"
	"net/http"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/config"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/types"
	"github.com/hayride-dev/bindings/go/wasi/net/http/handle"
	"go.bytecodealliance.org/cm"
)

type Config = config.ServerConfig

var cfg = config.ServerConfig{
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

func get() (result cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]) {
	if cfg.Address == "" {
		wasiErr := config.ErrorResourceNew(cm.Rep(types.ErrorCodeInvalid))
		return cm.Err[cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]](wasiErr)
	}
	return cm.OK[cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]](cfg)
}
