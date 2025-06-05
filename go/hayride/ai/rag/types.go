package rag

import (
	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/rag"
	"github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/transformer"
)

type Transformer = rag.Transformer
type RagOption = rag.RagOption
type EmbeddingType = transformer.EmbeddingType

const EmbeddingTypeSentence = transformer.EmbeddingTypeSentence
