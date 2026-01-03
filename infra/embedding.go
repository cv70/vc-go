package infra

import (
	"context"
	"vc-go/config"
	"vc-go/pkg/sdk"
)

func NewEmbeddingModel(ctx context.Context, c *config.EmbeddingConfig) (sdk.EmbeddingClient, error) {
	embeddingClient := &sdk.AnythingEmbeddingClient{
		URL: c.BaseURL,
	}
	return embeddingClient, nil
}
