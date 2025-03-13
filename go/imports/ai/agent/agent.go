package agent

import (
	"fmt"
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agent"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"github.com/hayride-dev/bindings/go/shared/ai/msg"
	"go.bytecodealliance.org/cm"
)

func Invoke(message *msg.Message, w io.Writer) error {
	if _, ok := w.(wasiio.WasiWriter); !ok {
		return fmt.Errorf("expected io.WasiWriter, got %T", w)
	}
	// wasi message
	if message.Role != msg.RoleUser {
		return fmt.Errorf("expected user role")
	}
	content := make([]types.Content, 0)
	for _, c := range message.Content {
		switch c.Type() {
		case "text":
			textContent := c.(*msg.TextContent)
			content = append(content, types.ContentText(types.TextContent{
				Text:        textContent.Text,
				ContentType: textContent.ContentType,
			}))
		}
	}

	witMsg := types.Message{
		Role:    types.Role(message.Role),
		Content: cm.ToList(content),
	}

	wasiStream := w.(wasiio.WasiWriter)
	ptr := wasiStream.Ptr()
	output := agent.OutputStream(cm.Rep(ptr))

	result := agent.Invoke(witMsg, output)

	if result.IsErr() {
		return fmt.Errorf("invoke error: %s", result.Err().Data())
	}

	return nil
}
