package msg

type Role uint8

const (
	RoleUser Role = iota
	RoleAssistant
	RoleSystem
	RoleTool
	RoleUnknown
)

type Message struct {
	Role    Role      `json:"role"`
	Content []Content `json:"content"`
}

type Content interface {
	Type() string
}

type TextContent struct {
	Text        string `json:"text"`
	ContentType string `json:"content-type"`
}

func (t *TextContent) Type() string {
	return "text"
}

type ToolInput struct {
	ContentType string `json:"content-type"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Input       string `json:"input"`
}

func (t *ToolInput) Type() string {
	return "tool-input"
}

type ToolOutput struct {
	ContentType string `json:"content-type"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Output      string `json:"Output"`
}

func (t *ToolOutput) Type() string {
	return "tool-output"
}
