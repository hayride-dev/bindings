package ctx

import (
	"fmt"

	witContext "github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/context"
	"github.com/hayride-dev/bindings/go/gen/exports/hayride/ai/types"
)

// TODO :: Remove from bindings

var _ Context = (*inMemoryContext)(nil)

type inMemoryContext struct {
	context []witContext.Message
}

func NewInmemoryCtx() *inMemoryContext {
	return &inMemoryContext{
		context: make([]witContext.Message, 0),
	}
}

func (c *inMemoryContext) Push(messages ...witContext.Message) error {
	for _, m := range messages {
		if m.Role == types.RoleSystem {
			c.context[0] = m
			continue
		}
		c.context = append(c.context, m)
	}
	return nil
}

func (c *inMemoryContext) Messages() ([]witContext.Message, error) {
	return c.context, nil
}

func (c *inMemoryContext) Next() (witContext.Message, error) {
	if len(c.context) == 0 {
		return witContext.Message{}, fmt.Errorf("missing messages")
	}
	return c.context[len(c.context)-1], nil
}
