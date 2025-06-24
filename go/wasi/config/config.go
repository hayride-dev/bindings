package config

import (
	"errors"
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/wasi/config/store"
)

type Reader interface {
	Get(key string) (string, error)
	GetAll() (map[string]string, error)
}

func Get(key string) (string, error) {
	result := store.Get(key)
	if result.IsErr() {
		switch result.Err().String() {
		case "upstream":
			return *result.Err().Upstream(), nil
		case "io":
			return *result.Err().IO(), nil
		default:
			return "", fmt.Errorf("failed to get key %s: %s", key, result.Err().String())
		}
	}
	return result.OK().Value(), nil
}

func GetAll() (map[string]string, error) {
	result := store.GetAll()
	if result.IsErr() {
		return nil, errors.New(result.Err().String())
	}

	tuples := result.OK().Slice()
	values := make(map[string]string, 0)
	for _, tuple := range tuples {
		values[tuple[0]] = tuple[1]
	}

	return values, nil
}
