package ctx

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/shared/domain/ai"
)

var _ Context = (*inMemoryContext)(nil)

type inMemoryContext struct {
	context []*ai.Message
}

func newInmemoryCtx() *inMemoryContext {
	return &inMemoryContext{
		context: make([]*ai.Message, 0),
	}
}

func (c *inMemoryContext) Push(messages ...*ai.Message) error {
	for _, m := range messages {
		if m.Role == ai.RoleSystem {
			c.context[0] = m
			continue
		}
		c.context = append(c.context, m)
	}
	return nil
}

func (c *inMemoryContext) Messages() ([]*ai.Message, error) {
	return c.context, nil
}

func (c *inMemoryContext) Next() (*ai.Message, error) {
	if len(c.context) == 0 {
		return nil, fmt.Errorf("missing messages")
	}
	return c.context[len(c.context)-1], nil
}
