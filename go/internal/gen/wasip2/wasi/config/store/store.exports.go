// Code generated by wit-bindgen-go. DO NOT EDIT.

package store

import (
	"go.bytecodealliance.org/cm"
)

// Exports represents the caller-defined exports from "wasi:config/store@0.2.0-draft".
var Exports struct {
	// Get represents the caller-defined, exported function "get".
	//
	// Gets a configuration value of type `string` associated with the `key`.
	//
	// The value is returned as an `option<string>`. If the key is not found,
	// `Ok(none)` is returned. If an error occurs, an `Err(error)` is returned.
	//
	//	get: func(key: string) -> result<option<string>, error>
	Get func(key string) (result cm.Result[OptionStringShape_, cm.Option[string], Error])

	// GetAll represents the caller-defined, exported function "get-all".
	//
	// Gets a list of configuration key-value pairs of type `string`.
	//
	// If an error occurs, an `Err(error)` is returned.
	//
	//	get-all: func() -> result<list<tuple<string, string>>, error>
	GetAll func() (result cm.Result[ErrorShape_, cm.List[[2]string], Error])
}
