// Code generated by wit-bindgen-go. DO NOT EDIT.

package agents

import (
	"go.bytecodealliance.org/cm"
)

// Exports represents the caller-defined exports from "hayride:ai/agents@0.0.61".
var Exports struct {
	// Error represents the caller-defined exports for resource "hayride:ai/agents@0.0.61#error".
	Error struct {
		// Destructor represents the caller-defined, exported destructor for resource "error".
		//
		// Resource destructor.
		Destructor func(self cm.Rep)

		// Code represents the caller-defined, exported method "code".
		//
		// return the error code.
		//
		//	code: func() -> error-code
		Code func(self cm.Rep) (result ErrorCode)

		// Data represents the caller-defined, exported method "data".
		//
		// errors can propagated with backend specific status through a string value.
		//
		//	data: func() -> string
		Data func(self cm.Rep) (result string)
	}

	// Agent represents the caller-defined exports for resource "hayride:ai/agents@0.0.61#agent".
	Agent struct {
		// Destructor represents the caller-defined, exported destructor for resource "agent".
		//
		// Resource destructor.
		Destructor func(self cm.Rep)

		// Constructor represents the caller-defined, exported constructor for resource "agent".
		//
		//	constructor(name: string, instruction: string, format: format, graph: graph-execution-context-stream,
		//	tools: option<tools>, context: option<context>)
		Constructor func(name string, instruction string, format Format, graph GraphExecutionContextStream, tools_ cm.Option[Tools], context_ cm.Option[Context]) (result Agent)

		// Capabilities represents the caller-defined, exported method "capabilities".
		//
		//	capabilities: func() -> result<list<tool>, error>
		Capabilities func(self cm.Rep) (result cm.Result[cm.List[Tool], cm.List[Tool], Error])

		// Compute represents the caller-defined, exported method "compute".
		//
		//	compute: func(message: message) -> result<message, error>
		Compute func(self cm.Rep, message Message) (result cm.Result[MessageShape, Message, Error])

		// Context represents the caller-defined, exported method "context".
		//
		//	context: func() -> result<list<message>, error>
		Context func(self cm.Rep) (result cm.Result[cm.List[Message], cm.List[Message], Error])

		// Execute represents the caller-defined, exported method "execute".
		//
		//	execute: func(params: call-tool-params) -> result<call-tool-result, error>
		Execute func(self cm.Rep, params CallToolParams) (result cm.Result[CallToolResultShape, CallToolResult, Error])

		// Instruction represents the caller-defined, exported method "instruction".
		//
		//	instruction: func() -> string
		Instruction func(self cm.Rep) (result string)

		// Name represents the caller-defined, exported method "name".
		//
		//	name: func() -> string
		Name func(self cm.Rep) (result string)
	}
}
