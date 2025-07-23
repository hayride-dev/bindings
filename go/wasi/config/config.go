package config

import (
	"errors"
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/wasip2/wasi/config/store"
)

type Reader interface {
	Get(key string) (string, error)
	GetAll() (map[string]string, error)
}

type readerImpl struct{}

func NewReader() Reader {
	return &readerImpl{}
}

func (r *readerImpl) Get(key string) (string, error) {
	result := store.Get(key)
	if result.IsErr() {
		switch result.Err().String() {
		case "upstream":
			return "", fmt.Errorf("upstream error: %s", *result.Err().Upstream())
		case "io":
			return "", fmt.Errorf("I/O error: %s", *result.Err().IO())
		default:
			return "", fmt.Errorf("failed to get key %s: %s", key, result.Err().String())
		}
	}
	return result.OK().Value(), nil
}

func (r *readerImpl) GetAll() (map[string]string, error) {
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
