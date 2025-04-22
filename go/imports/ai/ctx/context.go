package ctx

/*
This file contains a ergonamic wrapper around the wit generated code
for interacting with a imported context resource.
*/
import (
	"fmt"

	witContext "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/context"
	"go.bytecodealliance.org/cm"
)

type Context cm.Resource

// Push take a list of messages, convert them to a list of wit Messages
// and call imported context push
func (c Context) Push(messages ...witContext.Message) error {
	witContext := cm.Reinterpret[witContext.Context](c)

	result := witContext.Push(cm.ToList(messages))
	if result.IsErr() {
		// TODO: handle error result
		return fmt.Errorf("failed to push message")
	}
	return nil
}

// Messages returns the list of messages in the context
func (c Context) Messages() ([]witContext.Message, error) {
	witContext := cm.Reinterpret[witContext.Context](c)
	result := witContext.Messages()
	if result.IsErr() {
		// TODO: handle error result
		return nil, fmt.Errorf("failed to get messages")
	}
	witMessages := result.OK()

	return witMessages.Slice(), nil
}

func (c Context) Next() (*witContext.Message, error) {
	witContext := cm.Reinterpret[witContext.Context](c)
	result := witContext.Next()
	if result.IsErr() {
		// TODO : handle error result
		return nil, fmt.Errorf("failed to get next message")
	}

	return result.OK(), nil
}

// Create the resource
func NewContext() Context {
	return Context(witContext.NewContext())
}
