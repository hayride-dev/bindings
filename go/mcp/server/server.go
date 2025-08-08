package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hayride-dev/bindings/go/hayride/mcp/auth"
	"github.com/hayride-dev/bindings/go/hayride/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.bytecodealliance.org/cm"
)

type httpServer struct {
	server *server.StreamableHTTPServer

	authProvider auth.Provider
}

// NewMCPRouter creates a new Streamable HTTP MCP server with the given options.
// Notifications are not supported
func NewMCPRouter(options ...Option[*MCPServerOptions]) (*http.ServeMux, error) {
	opts := defaultMCPServerOptions()
	for _, o := range options {
		if err := o.Apply(opts); err != nil {
			return nil, fmt.Errorf("failed to apply option: %v", err)
		}
	}

	s, err := createServer(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %v", err)
	}

	server := &httpServer{
		server:       server.NewStreamableHTTPServer(s),
		authProvider: opts.authProvider,
	}

	mux := http.NewServeMux()

	// Add CORS middleware wrapper
	corsHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	mux.HandleFunc("/mcp", corsHandler(server.ServeHTTP))

	// If auth provider is set, handle the default auth routes
	if opts.authProvider != nil {
		mux.HandleFunc("/authorize", corsHandler(func(w http.ResponseWriter, r *http.Request) {
			baseURL, err := opts.authProvider.AuthURL()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to get auth URL: %v", err), http.StatusInternalServerError)
				return
			}

			// TODO: The client should provide and validate this, added to work with mcp inspector client for now
			// Generate a random 16-digit state parameter
			stateBytes := make([]byte, 8)
			if _, err := rand.Read(stateBytes); err != nil {
				http.Error(w, fmt.Sprintf("failed to generate state parameter: %v", err), http.StatusInternalServerError)
				return
			}
			state := hex.EncodeToString(stateBytes)

			// Get existing query parameters and add the state parameter
			queryParams := r.URL.Query()
			queryParams.Set("state", state)
			queryString := queryParams.Encode()

			// Append query parameters to the base URL
			separator := "?"
			if strings.Contains(baseURL, "?") {
				separator = "&"
			}
			baseURL = baseURL + separator + queryString

			http.Redirect(w, r, baseURL, http.StatusFound)
		}))

		mux.HandleFunc("/token", corsHandler(func(w http.ResponseWriter, r *http.Request) {
			// Send body to auth provider for code exchange to get a token
			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to read body: %v", err), http.StatusInternalServerError)
				return
			}

			response, err := opts.authProvider.ExchangeCode(data)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to exchange code: %v", err), http.StatusInternalServerError)
				return
			}

			w.Write(response)
		}))

		mux.HandleFunc("/register", corsHandler(func(w http.ResponseWriter, r *http.Request) {
			// Read body to pass to auth provider for registration
			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to read body: %v", err), http.StatusInternalServerError)
				return
			}

			response, err := opts.authProvider.Registration(data)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed Registration: %v", err), http.StatusInternalServerError)
				return
			}

			// Respond with 201 Created and the response body
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(response); err != nil {
				http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
				return
			}
		}))
	}

	return mux, nil
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If authProvider is set, require bearer authentication
	if s.authProvider != nil {
		// Check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expect Bearer token
		var bearer string
		n, err := fmt.Sscanf(authHeader, "Bearer %s", &bearer)
		if err != nil || n != 1 {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Validate the token
		valid, err := s.authProvider.Validate(bearer)
		if err != nil || !valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
	}

	s.server.ServeHTTP(w, r)
}

/*
func StartStdioServer(stdin io.Reader, stdout io.Writer, options ...Option[*MCPServerOptions]) error {
	s, err := createServer(options...)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	// Continuously process messages from stdio
	reader := bufio.NewReader(stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// Parse the message as raw JSON
		var rawMessage json.RawMessage
		if err := json.Unmarshal([]byte(line), &rawMessage); err != nil {
			// TODO: Handle parse error
			// response := createErrorResponse(nil, mcp.PARSE_ERROR, "Parse error")
			// return s.writeResponse(response, writer)
			return fmt.Errorf("failed to parse message: %v", err)
		}

		response := s.HandleMessage(context.TODO(), rawMessage)

		// Write the response as JSON
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		// Write response followed by newline
		if _, err := fmt.Fprintf(stdout, "%s\n", responseBytes); err != nil {
			return err
		}
	}
}
*/

func createServer(options ...Option[*MCPServerOptions]) (*server.MCPServer, error) {
	opts := defaultMCPServerOptions()
	for _, o := range options {
		if err := o.Apply(opts); err != nil {
			return nil, fmt.Errorf("failed to apply option: %v", err)
		}
	}

	mcpOptions := []server.ServerOption{}

	if opts.toolbox != nil {
		// Add tool capabilities with list changed notifications disabled
		mcpOptions = append(mcpOptions, server.WithToolCapabilities(false))
	}

	s := server.NewMCPServer(
		opts.name,
		opts.version,
		mcpOptions...,
	)

	if opts.toolbox != nil {
		cursor := ""
		toolResult, err := opts.toolbox.List(cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to get tool capabilities: %v", err)
		}

		for _, t := range toolResult.Tools.Slice() {
			var tool mcp.Tool
			if len(t.InputSchema.Properties.Slice()) > 0 {
				// Use the Input Schema as the Raw Schema for the MCP tool if set
				data, err := t.InputSchema.MarshalJSON()
				if err != nil {
					return nil, fmt.Errorf("failed to marshal input schema for %s: %v", t.Name, err)
				}

				tool = mcp.NewToolWithRawSchema(t.Name, t.Description, data)
			} else {
				// Use the basic tool definition without input schema
				tool = mcp.NewTool(t.Name, mcp.WithDescription(t.Description))
			}

			// Create a handler for the tool
			handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				args := request.GetArguments()

				// Convert map to [][2]string
				input := make([][2]string, 0, len(args))
				for k, v := range args {
					value := ""
					switch v := v.(type) {
					case string:
						// String argument
						value = v
					case int, int8, int16, int32, int64:
						// Number argument
						value = fmt.Sprintf("%d", v)
					case float32, float64:
						// Float argument
						value = fmt.Sprintf("%f", v)
					case bool:
						// Boolean argument
						value = fmt.Sprintf("%t", v)
					default:
						// Unsupported type, skip
						return nil, fmt.Errorf("unsupported argument type for key %s: %T", k, v)
					}

					input = append(input, [2]string{k, value})
				}

				toolInput := types.CallToolParams{
					Name:      tool.Name,
					Arguments: cm.ToList(input),
				}

				// Call the tool
				result, err := opts.toolbox.Call(toolInput)
				if err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}

				// Convert the result to MCP content
				callToolResult := &mcp.CallToolResult{}
				for _, out := range result.Content.Slice() {
					switch out.String() {
					case "text":
						content := out.Text()
						callToolResult.Content = append(callToolResult.Content, mcp.TextContent{
							Type: content.ContentType,
							Text: content.Text,
						})
					case "image":
						content := out.Image()
						callToolResult.Content = append(callToolResult.Content, mcp.ImageContent{
							Type:     content.ContentType,
							Data:     string(content.Data.Slice()),
							MIMEType: content.MIMEType,
						})
					case "audio":
						content := out.Audio()
						callToolResult.Content = append(callToolResult.Content, mcp.AudioContent{
							Type:     content.ContentType,
							Data:     string(content.Data.Slice()),
							MIMEType: content.MIMEType,
						})
					case "resource-link":
						content := out.ResourceLink()
						callToolResult.Content = append(callToolResult.Content, mcp.ResourceLink{
							Type:        content.ContentType,
							URI:         content.URI,
							Name:        content.Name,
							Description: content.Description,
							MIMEType:    content.MIMEType,
						})
					case "resource-content":
						// Either blob or text resource
						content := out.ResourceContent()
						switch content.ResourceContents.String() {
						case "text":
							text := content.ResourceContents.Text()
							callToolResult.Content = append(callToolResult.Content, mcp.EmbeddedResource{
								Type: content.ContentType,
								Resource: mcp.TextResourceContents{
									URI:      text.URI,
									MIMEType: text.MIMEType,
									Text:     text.Text,
								},
							})
						case "blob":
							blob := content.ResourceContents.Blob()
							callToolResult.Content = append(callToolResult.Content, mcp.EmbeddedResource{
								Type: content.ContentType,
								Resource: mcp.BlobResourceContents{
									URI:      blob.URI,
									MIMEType: blob.MIMEType,
									Blob:     string(blob.Blob.Slice()),
								},
							})
						default:
							return nil, fmt.Errorf("unsupported resource content type: %s", content.ResourceContents.String())
						}
					default:
						return nil, fmt.Errorf("unsupported output type: %s", out.String())
					}
				}

				return callToolResult, nil
			}

			// Add tool handler
			s.AddTool(tool, handler)
		}
	}

	return s, nil
}
