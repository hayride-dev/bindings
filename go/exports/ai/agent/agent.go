package agent

import (
	"io"

	"github.com/hayride-dev/bindings/go/exports/ai/types"
	wasiio "github.com/hayride-dev/bindings/go/imports/io"
	witAgent "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/agent"
	"go.bytecodealliance.org/cm"
)

type Agent interface {
	Invoke(msg *types.Message, w io.Writer) error
}

func init() {
	witAgent.Exports.Invoke = wacInvoke
}

var a Agent

func Register(agent Agent) error {
	a = agent
	return nil
}

func wacInvoke(message witAgent.Message, output cm.Rep) cm.Result[witAgent.Error, struct{}, witAgent.Error] {
	if a == nil {
		agentErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeUnknown))
		return cm.Err[cm.Result[witAgent.Error, struct{}, witAgent.Error]](agentErr)
	}

	w := wasiio.Clone(uint32(output))
	defer witAgent.OutputStream(cm.Rep(uint32(output))).ResourceDrop()

	content := make([]types.Content, 0)
	for _, c := range message.Content.Slice() {
		switch c.String() {
		case "text":
			content = append(content, &types.TextContent{
				Text:        c.Text().Text,
				ContentType: c.Text().ContentType,
			})
		}
	}

	m := &types.Message{
		Role:    types.Role(message.Role),
		Content: content,
	}

	err := a.Invoke(m, w)
	if err != nil {
		agentErr := witAgent.ErrorResourceNew(cm.Rep(witAgent.ErrorCodeUnknown))
		return cm.Err[cm.Result[witAgent.Error, struct{}, witAgent.Error]](agentErr)
	}

	return cm.OK[cm.Result[witAgent.Error, struct{}, witAgent.Error]](struct{}{})
}
