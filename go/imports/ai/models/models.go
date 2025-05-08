package models

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/gen/types/hayride/ai/types"
	graphStream "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/graph-stream"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model"
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
	// host provides a graph stream
	result := graphStream.LoadByName(opts.name)
	if result.IsErr() {
		return cm.ResourceNone, fmt.Errorf("failed to load graph")
	}
	graph := result.OK()
	resultCtxStream := graph.InitExecutionContextStream()
	if result.IsErr() {
		return cm.ResourceNone, fmt.Errorf("failed to init execution graph context stream")
	}
	stream := *resultCtxStream.OK()

	// assumed a model is wac'd or host provides a format
	format := model.NewFormat()

	return Model(model.NewModel(format, stream)), nil
}

func (m Model) Compute(messages []types.Message) (*types.Message, error) {
	wModel := cm.Reinterpret[model.Model](m)
	cmMessages := make([]model.Message, len(messages))
	for i, msg := range messages {
		cmMessages[i] = cm.Reinterpret[model.Message](msg)
	}

	result := wModel.Compute(cm.ToList(cmMessages))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to compute")
	}

	witMsg := result.OK()
	msg := cm.Reinterpret[*types.Message](witMsg)

	return msg, nil
}
