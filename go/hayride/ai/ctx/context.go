package ctx

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported context resource.
*/
import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/ai"
	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/context"
	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type Context interface {
	Push(messages ...ai.Message) error
	Messages() ([]ai.Message, error)
}

type Ctx cm.Resource

// Create the resource
func New() (Context, error) {
	return Ctx(context.NewContext()), nil
}

// Push take a list of messages, convert them to a list of wit Messages
// and call imported context push
func (c Ctx) Push(messages ...ai.Message) error {
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
func (c Ctx) Messages() ([]ai.Message, error) {
	witContext := cm.Reinterpret[context.Context](c)
	result := witContext.Messages()
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get messages: %s", result.Err().Data())
	}
	msgs := cm.Reinterpret[cm.List[types.Message]](result.OK())
	var aiMsgs []ai.Message
	for _, m := range msgs.Slice() {
		aiMsgs = append(aiMsgs, ai.Message{Message: m})
	}
	return aiMsgs, nil
}
