package config

import (
	"fmt"
	"os"

	"github.com/hayride-dev/bindings/go/internal/gen/exports/wasi/config/store"
	"go.bytecodealliance.org/cm"
)

type Reader interface {
	Get(key string) (string, error)
	GetAll() (map[string]string, error)
}

func init() {
	store.Exports.Get = get
	store.Exports.GetAll = getAll
}

type defaultReader struct{}

func (r *defaultReader) Get(key string) (string, error) {
	if key == "" {
		return "", os.ErrInvalid
	}
	return "", fmt.Errorf("config reader not set")
}
func (r *defaultReader) GetAll() (map[string]string, error) {
	return nil, fmt.Errorf("config reader not set")
}

var readerInstance Reader = &defaultReader{}

func Export(r Reader) {
	readerInstance = r
}

func get(key string) (result cm.Result[store.OptionStringShape, cm.Option[string], store.Error]) {
	v, err := readerInstance.Get(key)
	if err != nil {
		s := store.ErrorUpstream("error getting config value" + err.Error())
		return cm.Err[cm.Result[store.OptionStringShape, cm.Option[string], store.Error]](s)
	}
	return cm.OK[cm.Result[store.OptionStringShape, cm.Option[string], store.Error]](cm.Some(v))
}

func getAll() (result cm.Result[store.ErrorShape, cm.List[[2]string], store.Error]) {
	values, err := readerInstance.GetAll()
	if err != nil {
		s := store.ErrorUpstream("error getting config value" + err.Error())
		return cm.Err[cm.Result[store.ErrorShape, cm.List[[2]string], store.Error]](s)
	}

	tuples := make([][2]string, 0)
	for k, v := range values {
		tuples = append(tuples, [2]string{k, v})
	}

	return cm.OK[cm.Result[store.ErrorShape, cm.List[[2]string], store.Error]](cm.ToList(tuples))
}
