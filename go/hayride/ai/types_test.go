package ai

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hayride-dev/bindings/go/hayride/mcp"
	"go.bytecodealliance.org/cm"
)

func TestMarshalMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string // JSON string
	}{
		{
			name: "user message with text content",
			message: &Message{
				Role:    RoleUser,
				Content: cm.ToList([]MessageContent{NewMessageContent(Text("Hello, world!"))}),
				Final:   true,
			},
			expected: `{"role":"user","content":[{"text":"Hello, world!"}],"final":true}`,
		},
		{
			name: "assistant message with no content",
			message: &Message{
				Role:    RoleAssistant,
				Content: cm.ToList([]MessageContent{NewMessageContent(None{})}),
				Final:   false,
			},
			expected: `{"role":"assistant","content":[{"none":null}],"final":false}`,
		},
		{
			name: "system message with multiple content types",
			message: &Message{
				Role: RoleSystem,
				Content: cm.ToList([]MessageContent{
					NewMessageContent(Text("System prompt")),
					NewMessageContent(cm.ToList([]uint8{1, 2, 3, 4})),
				}),
				Final: true,
			},
			expected: `{"role":"system","content":[{"text":"System prompt"},{"blob":[1,2,3,4]}],"final":true}`,
		},
		{
			name: "tool message with tool input",
			message: &Message{
				Role: RoleTool,
				Content: cm.ToList([]MessageContent{
					NewMessageContent(mcp.CallToolParams{
						Name:      "test_tool",
						Arguments: cm.ToList([][2]string{{"param1", "value1"}}),
					}),
				}),
				Final: false,
			},
			expected: `{"role":"tool","content":[{"tool-input":{"name":"test_tool","arguments":[["param1","value1"]]}}],"final":false}`,
		},
		{
			name: "tool message with tool output",
			message: &Message{
				Role: RoleTool,
				Content: cm.ToList([]MessageContent{
					NewMessageContent(mcp.CallToolResult{
						Content:           cm.ToList([]mcp.Content{}),
						StructuredContent: cm.ToList([][2]string{{"result", "success"}}),
						IsError:           false,
						Meta:              cm.ToList([][2]string{}),
					}),
				}),
				Final: true,
			},
			expected: `{"role":"tool","content":[{"tool-output":{"content":[],"structured-content":[["result","success"]],"is-error":false,"meta":[]}}],"final":true}`,
		},
		{
			name: "unknown role message",
			message: &Message{
				Role:    RoleUnknown,
				Content: cm.ToList([]MessageContent{NewMessageContent(Text("Unknown role"))}),
				Final:   true,
			},
			expected: `{"role":"unknown","content":[{"text":"Unknown role"}],"final":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualBytes, err := json.Marshal(tt.message)
			if err != nil {
				t.Fatalf("unexpected error marshaling message: %v", err)
			}

			var actualObj, expectedObj any
			if err := json.Unmarshal(actualBytes, &actualObj); err != nil {
				t.Fatalf("failed to unmarshal actual JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expected), &expectedObj); err != nil {
				t.Fatalf("failed to unmarshal expected JSON: %v", err)
			}

			if !reflect.DeepEqual(actualObj, expectedObj) {
				t.Errorf("unexpected marshaled JSON.\nExpected: %s\nGot:      %s", tt.expected, string(actualBytes))
			}
		})
	}
}
