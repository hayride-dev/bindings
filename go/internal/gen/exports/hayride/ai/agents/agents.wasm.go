// Code generated by wit-bindgen-go. DO NOT EDIT.

package agents

import (
	"github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

// This file contains wasmimport and wasmexport declarations for "hayride:ai@0.0.61".

//go:wasmimport [export]hayride:ai/agents@0.0.61 [resource-new]error
//go:noescape
func wasmimport_ErrorResourceNew(rep0 uint32) (result0 uint32)

//go:wasmimport [export]hayride:ai/agents@0.0.61 [resource-rep]error
//go:noescape
func wasmimport_ErrorResourceRep(self0 uint32) (result0 uint32)

//go:wasmimport [export]hayride:ai/agents@0.0.61 [resource-drop]error
//go:noescape
func wasmimport_ErrorResourceDrop(self0 uint32)

//go:wasmexport hayride:ai/agents@0.0.61#[dtor]error
//export hayride:ai/agents@0.0.61#[dtor]error
func wasmexport_ErrorDestructor(self0 uint32) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	Exports.Error.Destructor(self)
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]error.code
//export hayride:ai/agents@0.0.61#[method]error.code
func wasmexport_ErrorCode(self0 uint32) (result0 uint32) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result := Exports.Error.Code(self)
	result0 = (uint32)(result)
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]error.data
//export hayride:ai/agents@0.0.61#[method]error.data
func wasmexport_ErrorData(self0 uint32) (result *string) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result_ := Exports.Error.Data(self)
	result = &result_
	return
}

//go:wasmimport [export]hayride:ai/agents@0.0.61 [resource-new]agent
//go:noescape
func wasmimport_AgentResourceNew(rep0 uint32) (result0 uint32)

//go:wasmimport [export]hayride:ai/agents@0.0.61 [resource-rep]agent
//go:noescape
func wasmimport_AgentResourceRep(self0 uint32) (result0 uint32)

//go:wasmimport [export]hayride:ai/agents@0.0.61 [resource-drop]agent
//go:noescape
func wasmimport_AgentResourceDrop(self0 uint32)

//go:wasmexport hayride:ai/agents@0.0.61#[dtor]agent
//export hayride:ai/agents@0.0.61#[dtor]agent
func wasmexport_AgentDestructor(self0 uint32) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	Exports.Agent.Destructor(self)
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[constructor]agent
//export hayride:ai/agents@0.0.61#[constructor]agent
func wasmexport_Constructor(name0 *uint8, name1 uint32, instruction0 *uint8, instruction1 uint32, format0 uint32, graph0 uint32, tools0 uint32, tools1 uint32, context0 uint32, context1 uint32) (result0 uint32) {
	name := cm.LiftString[string]((*uint8)(name0), (uint32)(name1))
	instruction := cm.LiftString[string]((*uint8)(instruction0), (uint32)(instruction1))
	format := cm.Reinterpret[Format]((uint32)(format0))
	graph := cm.Reinterpret[GraphExecutionContextStream]((uint32)(graph0))
	tools_ := lift_OptionTools((uint32)(tools0), (uint32)(tools1))
	context_ := lift_OptionContext((uint32)(context0), (uint32)(context1))
	result := Exports.Agent.Constructor(name, instruction, format, graph, tools_, context_)
	result0 = cm.Reinterpret[uint32](result)
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]agent.capabilities
//export hayride:ai/agents@0.0.61#[method]agent.capabilities
func wasmexport_AgentCapabilities(self0 uint32) (result *cm.Result[cm.List[Tool], cm.List[Tool], Error]) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result_ := Exports.Agent.Capabilities(self)
	result = &result_
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]agent.compute
//export hayride:ai/agents@0.0.61#[method]agent.compute
func wasmexport_AgentCompute(self0 uint32, message0 uint32, message1 *types.MessageContent, message2 uint32) (result *cm.Result[MessageShape, Message, Error]) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	message := lift_Message((uint32)(message0), (*types.MessageContent)(message1), (uint32)(message2))
	result_ := Exports.Agent.Compute(self, message)
	result = &result_
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]agent.context
//export hayride:ai/agents@0.0.61#[method]agent.context
func wasmexport_AgentContext(self0 uint32) (result *cm.Result[cm.List[Message], cm.List[Message], Error]) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result_ := Exports.Agent.Context(self)
	result = &result_
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]agent.execute
//export hayride:ai/agents@0.0.61#[method]agent.execute
func wasmexport_AgentExecute(self0 uint32, params0 *uint8, params1 uint32, params2 *[2]string, params3 uint32) (result *cm.Result[CallToolResultShape, CallToolResult, Error]) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	params := lift_CallToolParams((*uint8)(params0), (uint32)(params1), (*[2]string)(params2), (uint32)(params3))
	result_ := Exports.Agent.Execute(self, params)
	result = &result_
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]agent.instruction
//export hayride:ai/agents@0.0.61#[method]agent.instruction
func wasmexport_AgentInstruction(self0 uint32) (result *string) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result_ := Exports.Agent.Instruction(self)
	result = &result_
	return
}

//go:wasmexport hayride:ai/agents@0.0.61#[method]agent.name
//export hayride:ai/agents@0.0.61#[method]agent.name
func wasmexport_AgentName(self0 uint32) (result *string) {
	self := cm.Reinterpret[cm.Rep]((uint32)(self0))
	result_ := Exports.Agent.Name(self)
	result = &result_
	return
}
