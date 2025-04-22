package rag

import (
	"github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/rag"
	"github.com/hayride-dev/bindings/go/gen/imports/hayride/ai/transformer"
)

type Transformer = rag.Transformer
type RagOption = rag.RagOption
type EmbeddingType = transformer.EmbeddingType

const EmbeddingTypeSentence = transformer.EmbeddingTypeSentence
