package ctx

type CtxOptions struct {
	ctx Context
}

type OptionType interface {
	*CtxOptions
}

type Option[T OptionType] interface {
	Apply(T) error
}

type funcOption[T OptionType] struct {
	f func(T) error
}

func (fo *funcOption[T]) Apply(opt T) error {
	return fo.f(opt)
}

func newFuncOption[T OptionType](f func(T) error) *funcOption[T] {
	return &funcOption[T]{
		f: f,
	}
}

func WithContext(ctx Context) Option[*CtxOptions] {
	return newFuncOption(func(c *CtxOptions) error {
		c.ctx = ctx
		return nil
	})
}

func defaultModelOptions() *CtxOptions {
	return &CtxOptions{
		ctx: newInmemoryCtx(),
	}
}
