package export

import (
	"fmt"
	"net/http"

	"github.com/hayride-dev/bindings/go/hayride/x/net/http/server"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/config"
	"github.com/hayride-dev/bindings/go/internal/gen/types/hayride/http/types"

	"github.com/hayride-dev/bindings/go/wasi/net/http/handle"
	"go.bytecodealliance.org/cm"
)

var cfg = server.ServerConfig{
	Address:        "localhost:80",
	MaxHeaderBytes: http.DefaultMaxHeaderBytes,
}

func init() {
	config.Exports.Get = get
}

func ServerConfig(h http.Handler, c server.ServerConfig) error {
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
	return cm.OK[cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]](cm.Reinterpret[config.ServerConfig](cfg))
}
