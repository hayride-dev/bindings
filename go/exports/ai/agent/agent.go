package agent

import (
	"io"

	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	witAgent "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	"github.com/hayride-dev/bindings/go/shared/ai/msg"
	"go.bytecodealliance.org/cm"
)

type Agent interface {
	Invoke(msg *msg.Message, w io.Writer) error
}

func init() {
	witAgent.Exports.Invoke = wacInvoke
}

var a Agent

func Register(agent Agent) error {
	a = agent
	return nil
}

func wacInvoke(message witAgent.Message, output witAgent.OutputStream) cm.Result[witAgent.Error, struct{}, witAgent.Error] {
	if a == nil {
		agentErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeUnknown))
		return cm.Err[cm.Result[witAgent.Error, struct{}, witAgent.Error]](agentErr)
	}

	w := wasiio.Clone(uint32(output))

	content := make([]msg.Content, 0)
	for _, c := range message.Content.Slice() {
		switch c.String() {
		case "text":
			content = append(content, &msg.TextContent{
				Text:        c.Text().Text,
				ContentType: c.Text().ContentType,
			})
		}
	}

	m := &msg.Message{
		Role:    msg.Role(message.Role),
		Content: content,
	}

	err := a.Invoke(m, w)
	if err != nil {
		agentErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeUnknown))
		return cm.Err[cm.Result[witAgent.Error, struct{}, witAgent.Error]](agentErr)
	}

	return cm.OK[cm.Result[witAgent.Error, struct{}, witAgent.Error]](struct{}{})
}
