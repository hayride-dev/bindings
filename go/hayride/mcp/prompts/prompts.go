package prompts

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/mcp"
	witPrompts "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/mcp/prompts"

	"go.bytecodealliance.org/cm"
)

var _ Prompts = (*PromptResource)(nil)

type Prompts interface {
	Get(params mcp.GetPromptParams) (*mcp.GetPromptResult, error)
	List(cursor string) (*mcp.ListPromptsResult, error)
}

type PromptResource cm.Resource

func New() (PromptResource, error) {
	return PromptResource(witPrompts.NewPrompts()), nil
}

func (t PromptResource) Get(params mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	witPromptsResource := cm.Reinterpret[witPrompts.Prompts](t)

	result := witPromptsResource.GetPrompt(cm.Reinterpret[witPrompts.GetPromptParams](params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get prompt: %s", result.Err().Data())
	}

	return cm.Reinterpret[*mcp.GetPromptResult](result.OK()), nil
}
func (t PromptResource) List(cursor string) (*mcp.ListPromptsResult, error) {
	witPromptsToolbox := cm.Reinterpret[witPrompts.Prompts](t)

	result := witPromptsToolbox.ListPrompts(cursor)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to list prompts: %s", result.Err().Data())
	}

	return cm.Reinterpret[*mcp.ListPromptsResult](result.OK()), nil
}
