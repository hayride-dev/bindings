package config

import (
	"errors"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/wasi/config/store"
)

func Get(key string) (string, error) {
	result := store.Get(key)
	if result.IsErr() {
		return "", errors.New(result.Err().String())
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
