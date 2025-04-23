package ctx

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported context resource.
*/
import (
	"fmt"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

type Context cm.Resource

// Push take a list of messages, convert them to a list of wit Messages
// and call imported context push
func (c Context) Push(messages ...types.Message) error {
	witContext := cm.Reinterpret[context.Context](c)

	// Convert types.Message to context.Message
	witMessages := make([]context.Message, len(messages))
	for i, msg := range messages {
		witMessages[i] = cm.Reinterpret[context.Message](msg)
	}

	result := witContext.Push(cm.ToList(witMessages))
	if result.IsErr() {
		// TODO: handle error result
		return fmt.Errorf("failed to push message")
	}
	return nil
}

// Messages returns the list of messages in the context
func (c Context) Messages() ([]types.Message, error) {
	witContext := cm.Reinterpret[context.Context](c)
	result := witContext.Messages()
	if result.IsErr() {
		// TODO: handle error result
		return nil, fmt.Errorf("failed to get messages")
	}
	witMessages := result.OK()

	// Convert context.Message to types.Message
	witMessagesSlice := witMessages.Slice()
	msgs := make([]types.Message, len(witMessagesSlice))
	for i, msg := range witMessagesSlice {
		msgs[i] = cm.Reinterpret[types.Message](msg)
	}

	return msgs, nil
}

func (c Context) Next() (*types.Message, error) {
	witContext := cm.Reinterpret[context.Context](c)
	result := witContext.Next()
	if result.IsErr() {
		// TODO : handle error result
		return nil, fmt.Errorf("failed to get next message")
	}

	msg := cm.Reinterpret[*types.Message](result.OK())
	return msg, nil
}

// Create the resource
func NewContext() Context {
	return Context(context.NewContext())
}
