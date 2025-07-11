package export

import (
	"fmt"
	"net/http"

	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/config"
	types_ "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/http/types"

	"github.com/hayride-dev/bindings/go/wasi/net/http/handle"
	"go.bytecodealliance.org/cm"
)

var cfg = types.ServerConfig{
	Address:        "localhost:80",
	MaxHeaderBytes: http.DefaultMaxHeaderBytes,
}

func init() {
	config.Exports.Get = get
}

func ServerConfig(h http.Handler, c types.ServerConfig) error {
	handle.Handler(h)
	if c.Address == "" {
		return fmt.Errorf("invalid address: %s", c.Address)
	}
	cfg = c
	return nil
}

func get() (result cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]) {
	if cfg.Address == "" {
		wasiErr := config.ErrorResourceNew(cm.Rep(types_.ErrorCodeInvalid))
		return cm.Err[cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]](wasiErr)
	}
	return cm.OK[cm.Result[config.ServerConfigShape, config.ServerConfig, config.Error]](cm.Reinterpret[config.ServerConfig](cfg))
}
