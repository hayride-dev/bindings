package export

import (
	"github.com/hayride-dev/bindings/go/internal/gen/wasi/config/store"
	"github.com/hayride-dev/bindings/go/wasi/config"
	"go.bytecodealliance.org/cm"
)

type readerInstance struct {
	config.Reader
}

func Export(r config.Reader) {
	instance := &readerInstance{
		Reader: r,
	}

	store.Exports.Get = instance.get
	store.Exports.GetAll = instance.getAll
}

func (ri *readerInstance) get(key string) (result cm.Result[store.OptionStringShape_, cm.Option[string], store.Error]) {
	v, err := ri.Get(key)
	if err != nil {
		s := store.ErrorUpstream("error getting config value" + err.Error())
		return cm.Err[cm.Result[store.OptionStringShape_, cm.Option[string], store.Error]](s)
	}
	return cm.OK[cm.Result[store.OptionStringShape_, cm.Option[string], store.Error]](cm.Some(v))
}

func (ri *readerInstance) getAll() (result cm.Result[store.ErrorShape_, cm.List[[2]string], store.Error]) {
	values, err := ri.GetAll()
	if err != nil {
		s := store.ErrorUpstream("error getting config value" + err.Error())
		return cm.Err[cm.Result[store.ErrorShape_, cm.List[[2]string], store.Error]](s)
	}

	tuples := make([][2]string, 0)
	for k, v := range values {
		tuples = append(tuples, [2]string{k, v})
	}

	return cm.OK[cm.Result[store.ErrorShape_, cm.List[[2]string], store.Error]](cm.ToList(tuples))
}
