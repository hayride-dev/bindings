package rag

import (
	"github.com/hayride-dev/bindings/go/ai/gen/imports/hayride/ai/rag"
	"github.com/hayride-dev/bindings/go/ai/gen/imports/hayride/ai/transformer"
	"go.bytecodealliance.org/cm"
)

type ragConnection struct {
	conn *rag.Connection

	options *RagOptions
}

func NewConnection(dsn string, options ...Option[*RagOptions]) (*ragConnection, error) {
	opts := defaultModelOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	result := rag.Connect(dsn)
	if result.IsErr() {
		return nil, &ragErr{result.Err()}
	}

	conn := result.OK()

	return &ragConnection{
		conn:    conn,
		options: opts,
	}, nil
}

func NewTransformer(embeddingType EmbeddingType, model string, dataColumn string, vectorColumn string) Transformer {
	return transformer.NewTransformer(embeddingType, model, dataColumn, vectorColumn)
}

func (c *ragConnection) Register(transformer rag.Transformer) error {
	result := c.conn.Register(transformer)
	if result.IsErr() {
		return &ragErr{result.Err()}
	}

	return nil
}

func (c *ragConnection) Embed(table string, data string) error {
	result := c.conn.Embed(table, data)
	if result.IsErr() {
		return &ragErr{result.Err()}
	}

	return nil
}

func (c *ragConnection) Query(table string, data string, options []RagOption) ([]string, error) {
	list := cm.ToList(options)

	result := c.conn.Query(table, data, list)
	if result.IsErr() {
		return nil, &ragErr{result.Err()}
	}

	return result.OK().Slice(), nil
}
