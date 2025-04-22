package model

import (
	"fmt"

	witGraph "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/graph-stream"
	witModel "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/model"
	witTypes "github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/types"
	"go.bytecodealliance.org/cm"
)

type Model cm.Resource

func New(options ...Option[*ModelOptions]) (Model, error) {
	opts := defaultModelOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return cm.ResourceNone, err
		}
	}

	result := witGraph.LoadByName(opts.name)
	if result.IsErr() {
		return cm.ResourceNone, fmt.Errorf("failed to load graph")
	}
	graph := result.OK()
	resultCtxStream := graph.InitExecutionContextStream()
	if result.IsErr() {
		return cm.ResourceNone, fmt.Errorf("failed to init execution graph context stream")
	}
	stream := *resultCtxStream.OK()

	format := witModel.NewFormat()

	return Model(witModel.NewModel(format, stream)), nil
}

func (m Model) Compute(messages []witModel.Message) (*witModel.Message, error) {
	wModel := cm.Reinterpret[witModel.Model](m)
	result := wModel.Compute(cm.ToList(messages))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compute")
	}

	// message should always be a model response ( aka assistant)
	witMsg := result.OK()
	if witMsg.Role != witTypes.RoleAssistant {
		return nil, fmt.Errorf("expected assistant role, got %v", witMsg.Role)
	}

	return witMsg, nil
}
