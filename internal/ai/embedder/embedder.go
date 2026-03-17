package embedder

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/gogf/gf/v2/frame/g"
)

func DoubaoEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	model, err := g.Cfg().Get(ctx, "doubao_embedding_model.model")
	if err != nil {
		return nil, err
	}
	api_key, err := g.Cfg().Get(ctx, "doubao_embedding_model.api_key")
	if err != nil {
		return nil, err
	}
	base_url, err := g.Cfg().Get(ctx, "doubao_embedding_model.base_url")
	if err != nil {
		return nil, err
	}

	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		Model:   model.String(),
		APIKey:  api_key.String(),
		BaseURL: base_url.String(),
	})
	if err != nil {
		log.Printf("new embedder error: %v\n", err)
		return nil, err
	}
	return embedder, nil
}
