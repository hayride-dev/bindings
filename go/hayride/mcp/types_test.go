package mcp

import (
	"encoding/json"
	"reflect"
	"testing"

	"go.bytecodealliance.org/cm"
)

func TestMarshalToolSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *ToolSchema
		expected string // JSON string
	}{
		{
			name: "empty schema",
			schema: &ToolSchema{
				SchemaType: "object",
				Properties: cm.ToList([][2]string{}),
				Required:   cm.ToList([]string{}),
			},
			expected: `{"type":"object", "properties":{}}`,
		},
		{
			name: "valid JSON property",
			schema: &ToolSchema{
				SchemaType: "object",
				Properties: cm.ToList([][2]string{
					{"name", `{"type":"string","description":"Name","default":""}`},
				}),
			},
			expected: `{"type":"object","properties":{"name":{"type":"string","description":"Name","default":""}}}`,
		},
		{
			name: "non JSON fallback",
			schema: &ToolSchema{
				SchemaType: "object",
				Properties: cm.ToList([][2]string{
					{"name", "not valid json"},
				}),
			},
			expected: `{"type":"object","properties":{"name":"not valid json"}}`,
		},
		{
			name: "required fields",
			schema: &ToolSchema{
				SchemaType: "object",
				Properties: cm.ToList([][2]string{
					{"name", `{"type":"string"}`},
				}),
				Required: cm.ToList([]string{"name"}),
			},
			expected: `{"type":"object","properties":{"name":{"type":"string"}},"required":["name"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualBytes, err := json.Marshal(tt.schema)
			if err != nil {
				t.Fatalf("unexpected error marshaling schema: %v", err)
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
