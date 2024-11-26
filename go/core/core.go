package core

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/core/gen/hayride/core/config"
	"github.com/hayride-dev/bindings/go/core/gen/hayride/core/types"
)

type Config = types.Config
type Logging = types.Logging
type Http = types.HTTP
type Runtime = types.Runtime

func SetConfig(cfg Config) error {
	result := config.Set(cfg)
	if result.IsErr() {
		return fmt.Errorf("failed to set config: %v", result.Err())
	}
	return nil
}

func GetConfig() (*Config, error) {
	result := config.Get()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get config: %v", result.Err())
	}
	return result.OK(), nil
}
