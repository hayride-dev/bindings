package agent

import (
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/imports/ai/types"
	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	witAgent "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/agent"
	witTypes "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

func Invoke(message *types.Message, w io.Writer) error {
	if _, ok := w.(wasiio.WasiWriter); !ok {
		return fmt.Errorf("expected io.WasiWriter, got %T", w)
	}

	if message.Role != types.RoleUser {
		return fmt.Errorf("expected user role")
	}
	content := make([]witTypes.Content, 0)
	for _, c := range message.Content {
		switch c.Type() {
		case "text":
			textContent := c.(*types.TextContent)
			content = append(content, witTypes.ContentText(witTypes.TextContent{
				Text:        textContent.Text,
				ContentType: textContent.ContentType,
			}))
		}
	}

	witMsg := witTypes.Message{
		Role:    witTypes.Role(message.Role),
		Content: cm.ToList(content),
	}

	wasiStream := w.(wasiio.WasiWriter)
	ptr := wasiStream.Ptr()
	output := witAgent.OutputStream(cm.Rep(ptr))

	result := witAgent.Invoke(witMsg, output)

	if result.IsErr() {
		return fmt.Errorf("invoke error: %s", result.Err().Data())
	}

	return nil
}
