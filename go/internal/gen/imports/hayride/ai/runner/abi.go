// Code generated by wit-bindgen-go. DO NOT EDIT.

package runner

import (
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

func lower_Message(v types.Message) (f0 uint32, f1 *types.MessageContent, f2 uint32) {
	f0 = (uint32)(v.Role)
	f1, f2 = cm.LowerList(v.Content)
	return
}
