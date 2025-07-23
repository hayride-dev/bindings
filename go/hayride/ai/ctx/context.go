package ctx

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported context resource.
*/
import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

type Context interface {
	Push(messages ...types.Message) error
	Messages() ([]types.Message, error)
}

type ContextResource cm.Resource

// Create the resource
func New() (Context, error) {
	return ContextResource(context.NewContext()), nil
}

// Push take a list of messages, convert them to a list of wit Messages
// and call imported context push
func (c ContextResource) Push(messages ...types.Message) error {
	witContext := cm.Reinterpret[context.Context](c)
	// Convert types.Message to context.Message and push
	for _, msg := range messages {
		result := witContext.Push(cm.Reinterpret[context.Message](msg))
		if result.IsErr() {
			return fmt.Errorf("failed to push message: %s", result.Err().Data())
		}
	}
	return nil
}

// Messages returns the list of messages in the context
func (c ContextResource) Messages() ([]types.Message, error) {
	witContext := cm.Reinterpret[context.Context](c)
	result := witContext.Messages()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get messages: %s", result.Err().Data())
	}
	msgs := result.OK().Slice()
	return cm.Reinterpret[[]types.Message](msgs), nil
}
