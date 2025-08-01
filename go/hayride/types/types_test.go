package types

import (
	"encoding/json"
	"reflect"
	"testing"

	"go.bytecodealliance.org/cm"
)

func TestRequestData(t *testing.T) {
	t.Run("Cast", func(t *testing.T) {
		data := NewRequestData(Cast{})
		if data.Tag() != 1 {
			t.Errorf("expected tag 1, got %d", data.Tag())
		}
	})

	t.Run("SessionID", func(t *testing.T) {
		data := NewRequestData(SessionID("test-session"))
		if data.Tag() != 2 {
			t.Errorf("expected tag 2, got %d", data.Tag())
		}
	})

	t.Run("Generate", func(t *testing.T) {
		data := NewRequestData(Generate{})
		if data.Tag() != 3 {
			t.Errorf("expected tag 3, got %d", data.Tag())
		}
	})
}

func TestResponseData(t *testing.T) {
	t.Run("ThreadMetadata", func(t *testing.T) {
		data := NewResponseData(cm.ToList([]ThreadMetadata{{}}))
		if data.Tag() != 1 {
			t.Errorf("expected tag 1, got %d", data.Tag())
		}
	})

	t.Run("SessionID", func(t *testing.T) {
		data := NewResponseData(SessionID("test-session"))
		if data.Tag() != 2 {
			t.Errorf("expected tag 2, got %d", data.Tag())
		}
	})

	t.Run("Path", func(t *testing.T) {
		data := NewResponseData(Path("test/path"))
		if data.Tag() != 5 {
			t.Errorf("expected tag 5, got %d", data.Tag())
		}
	})

	t.Run("Messages", func(t *testing.T) {
		data := NewResponseData(cm.ToList([]Message{{}}))
		if data.Tag() != 4 {
			t.Errorf("expected tag 4, got %d", data.Tag())
		}
	})

	t.Run("Paths", func(t *testing.T) {
		data := NewResponseData(cm.ToList([]string{"path1", "path2"}))
		if data.Tag() != 6 {
			t.Errorf("expected tag 6, got %d", data.Tag())
		}
	})

	t.Run("ThreadStatus", func(t *testing.T) {
		data := NewResponseData(ThreadStatus(0))
		if data.Tag() != 3 {
			t.Errorf("expected tag 3, got %d", data.Tag())
		}
	})
}

func TestMarshalRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   Request
		expectErr bool
	}{
		{
			name: "User Request",
			request: Request{
				Data: NewRequestData(
					Generate{
						Model:    "test-model",
						System:   "This is a system message.",
						Messages: cm.ToList([]Message{{Role: RoleUser, Content: cm.ToList([]MessageContent{NewMessageContent(Text("Hello, world!"))})}}),
					},
				),
				Metadata: cm.ToList([][2]string{{"key", "value"}}),
			},
			expectErr: false,
		},
		{
			name: "System and User Request",
			request: Request{
				Data: NewRequestData(
					Generate{
						Model:  "test-model",
						System: "This is a system message.",
						Messages: cm.ToList([]Message{
							{Role: RoleSystem, Content: cm.ToList([]MessageContent{NewMessageContent(Text("System message."))})},
							{Role: RoleUser, Content: cm.ToList([]MessageContent{NewMessageContent(Text("User message."))})},
						}),
					}),
				Metadata: cm.ToList([][2]string{{"key1", "value1"}, {"key2", "value2"}}),
			},
			expectErr: false,
		},
		{
			name: "Tool Schema Content",
			request: Request{
				Data: NewRequestData(
					Generate{
						Model:  "test-model",
						System: "This is a system message.",
						Messages: cm.ToList([]Message{
							{Role: RoleUser, Content: cm.ToList([]MessageContent{NewMessageContent(cm.ToList([]Tool{
								{Name: "example-tool", Description: "An example tool", InputSchema: ToolSchema{SchemaType: "object", Properties: cm.ToList([][2]string{{"arg1", "string"}, {"arg2", "string"}}), Required: cm.ToList([]string{"arg1", "arg2"})}},
							}))})},
							{Role: RoleAssistant, Content: cm.ToList([]MessageContent{NewMessageContent(CallToolParams{
								Name:      "example-tool",
								Arguments: cm.ToList([][2]string{{"arg1", "value1"}, {"arg2", "value2"}}),
							})})},
							{Role: RoleTool, Content: cm.ToList([]MessageContent{NewMessageContent(CallToolResult{
								Content: cm.ToList([]Content{
									NewContent(TextContent{ContentType: "text", Text: "Tool output"}),
									NewContent(ImageContent{ContentType: "image", Data: cm.ToList([]byte{0x89, 0x50, 0x4E, 0x47})}),
									NewContent(AudioContent{ContentType: "audio", Data: cm.ToList([]byte("audio data"))}),
									NewContent(ResourceLinkContent{ContentType: "resource_link", URI: "https://example.com/resource"}),
									NewContent(EmbeddedResourceContent{ContentType: "resource", ResourceContents: NewResourceContents(
										TextResourceContents{
											URI:      "file:///example.txt",
											Name:     "example.txt",
											Title:    "Example Text File",
											MIMEType: "text/plain",
											Text:     "Resource content",
										},
									)}),
									NewContent(EmbeddedResourceContent{ContentType: "resource", ResourceContents: NewResourceContents(
										BlobResourceContents{
											URI:      "file:///example.bin",
											Name:     "example.bin",
											Title:    "Example Binary File",
											MIMEType: "application/octet-stream",
											Blob:     cm.ToList([]byte{0x01, 0x02, 0x03, 0x04}),
										},
									)}),
								}),
							})})},
						}),
					}),
				Metadata: cm.ToList([][2]string{{"key1", "value1"}, {"key2", "value2"}}),
			},
			expectErr: false,
		},
		{
			name:      "Empty Request",
			request:   Request{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeMarshal := tt.request

			data, err := tt.request.MarshalJSON()
			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error status: %v", err)
			}

			if !tt.expectErr {
				var unmarshalledRequest Request
				if err := unmarshalledRequest.UnmarshalJSON(data); err != nil {
					t.Fatalf("failed to unmarshal request: %v", err)
				}

				if tt.request.Data.String() != unmarshalledRequest.Data.String() {
					t.Fatalf("expected %s, got %s", tt.request.Data.String(), unmarshalledRequest.Data.String())
				}

				if tt.request.Metadata.Len() != unmarshalledRequest.Metadata.Len() {
					t.Fatalf("expected %d metadata items, got %d", tt.request.Metadata.Len(), unmarshalledRequest.Metadata.Len())
				}

				// Check that unmarshalled request matches the original
				switch beforeMarshal.Data.Tag() {
				case 0:
					if unmarshalledRequest.Data.Tag() != 0 {
						t.Fatal("expected tag 0 for Unknown data")
					}
				case 1: // Cast
					if beforeMarshal.Data.Cast() == nil || unmarshalledRequest.Data.Cast() == nil {
						t.Fatal("expected non-nil Cast data")
					}
					dataCast := beforeMarshal.Data.Cast()
					unmarshalledCast := unmarshalledRequest.Data.Cast()
					if dataCast.Name != unmarshalledCast.Name || dataCast.Args != unmarshalledCast.Args || dataCast.Function != unmarshalledCast.Function {
						t.Fatal("expected matching Cast data")
					}
				case 2: // SessionID
					if beforeMarshal.Data.SessionID() == nil || unmarshalledRequest.Data.SessionID() == nil {
						t.Fatal("expected non-nil SessionID data")
					}
					if beforeMarshal.Data.SessionID() != unmarshalledRequest.Data.SessionID() {
						t.Fatal("expected matching SessionID data")
					}
				case 3: // Generate
					if beforeMarshal.Data.Generate() == nil || unmarshalledRequest.Data.Generate() == nil {
						t.Fatal("expected non-nil Generate data")
					}
					dataGenerate := beforeMarshal.Data.Generate()
					unmarshalledGenerate := unmarshalledRequest.Data.Generate()
					if dataGenerate.Model != unmarshalledGenerate.Model || dataGenerate.System != unmarshalledGenerate.System {
						t.Fatal("expected matching Generate data")
					}
					if dataGenerate.Messages.Len() != unmarshalledGenerate.Messages.Len() {
						t.Fatalf("expected %d messages, got %d", dataGenerate.Messages.Len(), unmarshalledGenerate.Messages.Len())
					}
					for i, dataMessage := range dataGenerate.Messages.Slice() {
						unmarshalledMessage := unmarshalledGenerate.Messages.Slice()[i]
						if dataMessage.Role != unmarshalledMessage.Role {
							t.Fatalf("expected message role %s, got %s", dataMessage.Role, unmarshalledMessage.Role)
						}
						if dataMessage.Content.Len() != unmarshalledMessage.Content.Len() {
							t.Fatalf("expected %d content items, got %d", dataMessage.Content.Len(), unmarshalledMessage.Content.Len())
						}
						for j, dataContent := range dataMessage.Content.Slice() {
							unmarshalledContent := unmarshalledMessage.Content.Slice()[j]
							switch dataContent.Tag() {
							case 1: // TextContent
								if dataContent.Text() == nil || unmarshalledContent.Text() == nil {
									t.Fatal("expected non-nil TextContent")
								}
								if *dataContent.Text() != *unmarshalledContent.Text() {
									t.Fatalf("expected text content %s, got %s", *dataContent.Text(), *unmarshalledContent.Text())
								}
							case 2: // BlobContent
								if dataContent.Blob() == nil || unmarshalledContent.Blob() == nil {
									t.Fatal("expected non-nil BlobContent")
								}
								if dataContent.Blob().Len() != unmarshalledContent.Blob().Len() {
									t.Fatalf("expected %d bytes in BlobContent, got %d", dataContent.Blob().Len(), unmarshalledContent.Blob().Len())
								}
								for k, dataByte := range dataContent.Blob().Slice() {
									if dataByte != unmarshalledContent.Blob().Slice()[k] {
										t.Fatalf("expected byte %d in BlobContent, got %d", dataByte, unmarshalledContent.Blob().Slice()[k])
									}
								}
							case 3: // ToolsContent
								if dataContent.Tools() == nil || unmarshalledContent.Tools() == nil {
									t.Fatal("expected non-nil ToolsContent")
								}
								if dataContent.Tools().Len() != unmarshalledContent.Tools().Len() {
									t.Fatalf("expected %d tools, got %d", dataContent.Tools().Len(), unmarshalledContent.Tools().Len())
								}
								for k, dataTool := range dataContent.Tools().Slice() {
									unmarshalledTool := unmarshalledContent.Tools().Slice()[k]
									if dataTool.Name != unmarshalledTool.Name || dataTool.Description != unmarshalledTool.Description {
										t.Fatalf("expected tool %s with description %s, got %s with description %s", dataTool.Name, dataTool.Description, unmarshalledTool.Name, unmarshalledTool.Description)
									}
								}
							case 4: // ToolInputContent
								if dataContent.ToolInput() == nil || unmarshalledContent.ToolInput() == nil {
									t.Fatal("expected non-nil ToolInputContent")
								}
								dataInput := dataContent.ToolInput()
								unmarshalledInput := unmarshalledContent.ToolInput()
								if dataInput.Name != unmarshalledInput.Name || dataInput.Arguments.Len() !=
									unmarshalledInput.Arguments.Len() {
									t.Fatalf("expected tool input %s with %d arguments, got %s with %d arguments",
										dataInput.Name, dataInput.Arguments.Len(), unmarshalledInput.Name, unmarshalledInput.Arguments.Len())
								}
								for k, dataArg := range dataInput.Arguments.Slice() {
									unmarshalledArg := unmarshalledInput.Arguments.Slice()[k]
									if dataArg[0] != unmarshalledArg[0] || dataArg[1] != unmarshalledArg[1] {
										t.Fatalf("expected argument %s: %s, got %s: %s", dataArg[0], dataArg[1], unmarshalledArg[0], unmarshalledArg[1])
									}
								}
							case 5: // ToolOutputContent
								if dataContent.ToolOutput() == nil || unmarshalledContent.ToolOutput() == nil {
									t.Fatal("expected non-nil ToolOutputContent")
								}
								dataOutput := dataContent.ToolOutput()
								unmarshalledOutput := unmarshalledContent.ToolOutput()
								if dataOutput.Content.Len() != unmarshalledOutput.Content.Len() {
									t.Fatalf("expected %d content items in ToolOutputContent, got %d",
										dataOutput.Content.Len(), unmarshalledOutput.Content.Len())
								}
								for k, dataContentItem := range dataOutput.Content.Slice() {
									unmarshalledContentItem := unmarshalledOutput.Content.Slice()[k]
									switch dataContentItem.Tag() {
									case 1: // TextContent
										if dataContentItem.Text() == nil || unmarshalledContentItem.Text() == nil {
											t.Fatal("expected non-nil TextContent in ToolOutputContent")
										}
										if *dataContentItem.Text() != *unmarshalledContentItem.Text() {
											t.Fatalf("expected text content %s, got %s", *dataContentItem.Text(), *unmarshalledContentItem.Text())
										}
									case 2: // ImageContent
										if dataContentItem.Image() == nil || unmarshalledContentItem.Image() == nil {
											t.Fatal("expected non-nil ImageContent in ToolOutputContent")
										}
									case 3: // AudioContent
										if dataContentItem.Audio() == nil || unmarshalledContentItem.Audio() == nil {
											t.Fatal("expected non-nil AudioContent in ToolOutputContent")
										}
										if dataContentItem.Audio().Data.Len() != unmarshalledContentItem.Audio().Data.Len() {
											t.Fatalf("expected %d bytes in AudioContent, got %d",
												dataContentItem.Audio().Data.Len(), unmarshalledContentItem.Audio().Data.Len())
										}
										for l, dataByte := range dataContentItem.Audio().Data.Slice() {
											if dataByte != unmarshalledContentItem.Audio().Data.Slice()[l] {
												t.Fatalf("expected audio data %v, got %v",
													dataByte, unmarshalledContentItem.Audio().Data.Slice()[l])
											}
										}
									case 4: // ResourceLinkContent
										if dataContentItem.ResourceLink() == nil || unmarshalledContentItem.ResourceLink() == nil {
											t.Fatal("expected non-nil ResourceLinkContent in ToolOutputContent")
										}
									case 5: // ResourceContent
										if dataContentItem.ResourceContent() == nil || unmarshalledContentItem.ResourceContent() == nil {
											t.Fatal("expected non-nil ResourceContent in ToolOutputContent")
										}
										EmbeddedResourceContent := dataContentItem.ResourceContent().ResourceContents
										unmarshalledEmbeddedResourceContent := unmarshalledContentItem.ResourceContent().ResourceContents
										switch EmbeddedResourceContent.Tag() {
										case 1: // TextResourceContents
											if EmbeddedResourceContent.Text() == nil || unmarshalledEmbeddedResourceContent.Text() == nil {
												t.Fatal("expected non-nil TextResourceContents in ToolOutputContent")
											}
											if EmbeddedResourceContent.Text().URI != unmarshalledEmbeddedResourceContent.Text().URI ||
												EmbeddedResourceContent.Text().Name != unmarshalledEmbeddedResourceContent.Text().Name ||
												EmbeddedResourceContent.Text().Title != unmarshalledEmbeddedResourceContent.Text().Title ||
												EmbeddedResourceContent.Text().MIMEType != unmarshalledEmbeddedResourceContent.Text().MIMEType ||
												EmbeddedResourceContent.Text().Text != unmarshalledEmbeddedResourceContent.Text().Text {
												t.Fatalf("expected TextResourceContents to match, got %v vs %v",
													EmbeddedResourceContent.Text(), unmarshalledEmbeddedResourceContent.Text())
											}
										case 2: // BlobResourceContents
											if EmbeddedResourceContent.Blob() == nil || unmarshalledEmbeddedResourceContent.Blob() == nil {
												t.Fatal("expected non-nil BlobResourceContents in ToolOutputContent")
											}
											if EmbeddedResourceContent.Blob().URI != unmarshalledEmbeddedResourceContent.Blob().URI ||
												EmbeddedResourceContent.Blob().Name != unmarshalledEmbeddedResourceContent.Blob().Name ||
												EmbeddedResourceContent.Blob().Title != unmarshalledEmbeddedResourceContent.Blob().Title ||
												EmbeddedResourceContent.Blob().MIMEType != unmarshalledEmbeddedResourceContent.Blob().MIMEType ||
												EmbeddedResourceContent.Blob().Blob.Len() != unmarshalledEmbeddedResourceContent.Blob().Blob.Len() {
												t.Fatalf("expected BlobResourceContents to match, got %v vs %v",
													EmbeddedResourceContent.Blob(), unmarshalledEmbeddedResourceContent.Blob())
											}
											for m, dataByte := range EmbeddedResourceContent.Blob().Blob.Slice() {
												if dataByte != unmarshalledEmbeddedResourceContent.Blob().Blob.Slice()[m] {
													t.Fatalf("expected blob data %v, got %v",
														dataByte, unmarshalledEmbeddedResourceContent.Blob().Blob.Slice()[m])
												}
											}
										default:
											t.Fatalf("unexpected content tag %d in ToolOutputContent", dataContentItem.Tag())
										}
									}
								}
							}
						}
					}
				default:
					t.Fatal("unexpected data tag")
				}
			}
		})
	}
}

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
