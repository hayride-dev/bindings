package ctx

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/gen/domain/hayride/ai/types"
)

// TODO :: Remove from bindings

var _ Context = (*inMemoryContext)(nil)

type inMemoryContext struct {
	context []types.Message
}

func NewInmemoryCtx() *inMemoryContext {
	return &inMemoryContext{
		context: make([]types.Message, 0),
	}
}

func (c *inMemoryContext) Push(messages ...types.Message) error {
	for _, m := range messages {
		if m.Role == types.RoleSystem {
			c.context[0] = m
			continue
		}
		c.context = append(c.context, m)
	}
	return nil
}

func (c *inMemoryContext) Messages() ([]types.Message, error) {
	return c.context, nil
}

func (c *inMemoryContext) Next() (types.Message, error) {
	if len(c.context) == 0 {
		return types.Message{}, fmt.Errorf("missing messages")
	}
	return c.context[len(c.context)-1], nil
}
